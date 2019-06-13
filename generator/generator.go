package generator

import (
	"math/rand"
	"time"
)

// The generated data-point for a particular (node,metric,query) tuple.
// value contains the generated gaussian, disruption represents if/what
// disruption was active
type TupleResult struct {
	Node          int       `json:"node"`
	NodeKeyword   string    `json:"node_keyword"`
	Metric        int       `json:"metric"`
	MetricKeyword string    `json:"metric_keyword"`
	Query         int       `json:"query"`
	QueryKeyword  string    `json:"query_keyword"`
	Hour          time.Time `json:"hour"`
	Value         float64   `json:"value"`
	Disruption    int       `json:"disruption"`
}

func genRange(min, max int) int {
	return rand.Intn(max-min) + min
}
