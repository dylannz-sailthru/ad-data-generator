package generator

import (
	"github.com/dylannz-sailthru/ad-data-generator/config"
	"github.com/sirupsen/logrus"
)

type NodeQueryMetric [3]int

type NormalParams struct {
	Mean int
	Std  int
}

type MetricParams struct {
	Regular   NormalParams
	Disrupted NormalParams
}

// Generate the distributions for each (node, query, metric) tuple
func generateDistributions(config *config.Config) map[NodeQueryMetric]MetricParams {
	// Generate the distributions for each (node, query, metric) tuple
	logrus.Debug("generating distributions per (node,query,metric) tuple...")

	m := map[NodeQueryMetric]MetricParams{}
	for node := 0; node < config.Nodes; node++ {
		for query := 0; query < config.Queries; query++ {
			for metric := 0; metric < config.Metrics; metric++ {
				// Each tuple gets a "regular" and "disrupted"
				// distribution from the allowed bounds
				m[NodeQueryMetric{node, query, metric}] = MetricParams{
					Regular: NormalParams{
						Mean: genRange(config.RegularDistribution.MinMean, config.RegularDistribution.MaxMean),
						Std:  genRange(config.RegularDistribution.MinStd, config.RegularDistribution.MaxStd),
					},
					Disrupted: NormalParams{
						Mean: genRange(config.DisruptedDistribution.MinMean, config.DisruptedDistribution.MaxMean),
						Std:  genRange(config.DisruptedDistribution.MinStd, config.DisruptedDistribution.MaxStd),
					},
				}
			}
		}
	}

	return m
}

/*

// Generate the distributions for each (node, query, metric) tuple
fn generate_distributions(config: &Config, rng: &mut ThreadRng)
                            -> HashMap<(usize, usize, usize), (NormalParams, NormalParams)> {
    // Generate the distributions for each (node, query, metric) tuple
    debug!("generating distributions per (node,query,metric) tuple...");
    let mut distributions = HashMap::with_capacity(config.nodes * config.queries * config.metrics);
    for node in 0..config.nodes {
        for query in 0..config.queries {
            for metric in 0..config.metrics {

                // Each tuple gets a "regular" and "disrupted"
                // distribution from the allowed bounds
                let regular = NormalParams {
                    mean: rng.gen_range(config.regular_distribution.min_mean, config.regular_distribution.max_mean),
                    std: rng.gen_range(config.regular_distribution.min_std, config.regular_distribution.max_std)
                };

                let disrupted = NormalParams {
                    mean: rng.gen_range(config.disrupted_distribution.min_mean, config.disrupted_distribution.max_mean),
                    std: rng.gen_range(config.disrupted_distribution.min_std, config.disrupted_distribution.max_std)
                };

                distributions.insert((node, query, metric), (regular, disrupted));
            }
        }
    }

    distributions
}
*/
