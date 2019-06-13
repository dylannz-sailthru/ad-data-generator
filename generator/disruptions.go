package generator

import (
	"math/rand"
	"sort"

	"github.com/dylannz-sailthru/ad-data-generator/config"
	"github.com/sirupsen/logrus"
)

type DisruptionType int

const (
	None DisruptionType = iota
	Node
	Query
	Metric
)

type Disruption struct {
	Type     DisruptionType
	Affected sort.IntSlice

	Length int
}

func (d Disruption) Contains(v int) bool {
	for _, i := range d.Affected {
		if v == i {
			return true
		}
	}
	return false
}

// dedupe removes duplicates from a pre-sorted slice of integers
func dedupe(ints sort.IntSlice) sort.IntSlice {
	newInts := []int{}
	prev := 0
	for k, i := range ints {
		if k != 0 && prev == i {
			continue
		}
		newInts = append(newInts, i)
		prev = i
	}

	return newInts
}

func generateDisruptions(config *config.Config) map[int]Disruption {
	m := map[int]Disruption{}

	for i := 0; i < config.Disruptions; i++ {
		// Disruptions start after the 48th hour, and can last 2-24 hours long
		start := rand.Intn(config.Hours-24-48) + 48
		length := rand.Intn(24) - 2

		disruption := Disruption{
			Length: length,
		}
		switch rand.Intn(4) {
		case 0: // Node disruption: all (metric,queries) on the node are disrupted
			id := rand.Intn(config.Nodes)
			disruption.Type = Node
			disruption.Affected = append(disruption.Affected, id)
			logrus.Debugf("Node Disruption: %v-%v [%v]", start, start+length, id)
		case 1: // Query disruption: all (metric, node) with the query(ies) are disrupted
			// random number of queries from 1 to one-tenth of the total queries
			numQueries := genRange(1, int(config.Queries/10))
			disruption.Type = Query
			disruption.Affected = make(sort.IntSlice, numQueries)
			for q := 0; q < numQueries; q++ {
				disruption.Affected[q] = rand.Intn(config.Queries)
			}
			logrus.Debugf("Query Disruption: %v-%v [%v]", start, start+length, disruption.Affected)
		default: // Metric disruption: all (node,query) with the metric are disrupted
			// one-to-all metrics disrupted
			numMetrics := genRange(1, config.Metrics)
			disruption.Type = Metric
			disruption.Affected = make(sort.IntSlice, numMetrics)
			for m := 0; m < numMetrics; m++ {
				disruption.Affected[m] = rand.Intn(config.Metrics)
			}
			logrus.Debugf("Metric Disruption: %v-%v [%v]", start, start+length, disruption.Affected)
		}
		sort.Sort(disruption.Affected)
		disruption.Affected = dedupe(disruption.Affected)
		m[start] = disruption
	}

	return m
}
