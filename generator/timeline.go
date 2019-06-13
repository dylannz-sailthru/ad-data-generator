package generator

import (
	"context"
	"os"
	"strconv"
	"time"

	"gopkg.in/olivere/elastic.v7"

	"github.com/dylannz-sailthru/ad-data-generator/config"
	"github.com/sirupsen/logrus"
)

func GenerateTimeline(ctx context.Context, config *config.Config) {
	// We generate gaussians on another thread and cache them in a channel
	// to be grabbed when necessary, since gaussian generation is a bit slow
	rx, cancelRx := StartNormalGenerator(ctx, 1000)
	defer cancelRx()

	results := make([]TupleResult, 100000)

	// Generate a list of disruptions that will be seeded into our timeline
	disruptions := generateDisruptions(config)

	// Generate the "regular" and "disrupted" distributions for
	// every (node,metric,query) tupe
	distributions := generateDistributions(config)

	logrus.Debug("Generating timeline...")

	disruption := Disruption{Type: None}
	counter := 0

	logrus.Info("Elasticsearch host: ", os.Getenv("ELASTICSEARCH_HOST"))

	c, err := elastic.NewClient(
		elastic.SetURL(os.Getenv("ELASTICSEARCH_HOST")),
		elastic.SetSniff(false),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	bulk := elastic.NewBulkService(c)

	for hour := 0; hour < config.Hours; hour++ {
		logrus.Debugf("%v -- %v", hour, bulk.NumberOfActions())

		// if we are not currently in a disruption, check to see
		// if there is one this hour
		if counter <= 0 {
			disruption = disruptions[hour]
			counter = disruption.Length
		}

		for node := 0; node < config.Nodes; node++ {
			for query := 0; query < config.Queries; query++ {
				for metric := 0; metric < config.Metrics; metric++ {
					d, ok := distributions[NodeQueryMetric{node, query, metric}]
					if !ok {
						logrus.Panicf("no distribution for node/query/metric: %v/%v/%v", node, query, metric)
					}

					// Set the disruption flag for this tuple at this hour, based on the
					// disruption type
					isDisrupted := false
					if counter > 0 {
						switch disruption.Type {
						case None:
						case Query:
							isDisrupted = disruption.Contains(query)
						case Metric:
							isDisrupted = disruption.Contains(metric)
						case Node:
							isDisrupted = disruption.Contains(node)
						}
					}

					// Grab a gaussian from the normalGenerator thread and use the "regular"
					// or "disrupted" distributions to find the final value
					value := <-rx
					if isDisrupted {
						value = (value * float64(d.Disrupted.Std)) + float64(d.Disrupted.Mean)
					} else {
						value = (value * float64(d.Regular.Std)) + float64(d.Regular.Mean)
					}

					h := -config.Hours + hour
					timestamp := time.Now().Add(time.Hour * time.Duration(h)).UTC()

					bulk.Add(
						elastic.NewBulkIndexRequest().
							Index("data").
							Doc(TupleResult{
								Node:          node,
								NodeKeyword:   strconv.Itoa(node),
								Metric:        metric,
								MetricKeyword: strconv.Itoa(metric),
								Query:         query,
								QueryKeyword:  strconv.Itoa(query),
								Hour:          timestamp,
								Value:         value,
								Disruption:    int(disruption.Type),
							}),
					)

					// default max bulk request size in ES is 10MB
					const miB = 1 << 20
					if bulk.EstimatedSizeInBytes() > 9*miB {
						err := doBulkRequest(ctx, bulk)
						if err != nil {
							logrus.Fatal(err)
						}
					}
				}
			}
		}

		if bulk.NumberOfActions() > 0 {
			err := doBulkRequest(ctx, bulk)
			if err != nil {
				logrus.Fatal(err)
			}
		}

		/*
				index := elastic.NewBulkIndexRequest().
			Index("data").
			Type("data")*/
		counter--
	}

	t(
		rx,
		results,
		disruptions,
		distributions,
		disruption,
		counter,
		bulk,
	)
}

func t(a ...interface{}) {
	return
}

func doBulkRequest(ctx context.Context, bulk *elastic.BulkService) error {
	res, err := bulk.Do(ctx)
	if err != nil {
		return err
	}
	if res.Errors {
		for _, item := range res.Failed() {
			return &elastic.Error{
				Details: item.Error,
				Status:  item.Status,
			}
		}
	}

	return nil
}
