package stats

import (
	"time"

	"github.com/montanaflynn/stats"
)

// Stats contains all the basic stats for a run of the command
type Stats struct {
	Queries int
	Min     time.Duration
	Max     time.Duration
	Total   time.Duration
	Median  time.Duration
	Average time.Duration
}

const nsToMs = 1000000

// ProcessDurations gets a list of Durations from the workers and calculates the basics stats
// measures cannot be an empty slice
func ProcessDurations(measures []time.Duration) Stats {

	//Create fists a slice of floats in ms
	ms := make([]float64, 0, len(measures))

	min := int64(measures[0])
	max := int64(0)
	total := int64(0)
	queries := len(measures)

	// copy to ms to get the median and at the same time find some other values
	for _, measure := range measures {
		ns := int64(measure)
		if ns < min {
			min = ns
		}
		if ns > max {
			max = ns
		}
		total = total + ns
		ms = append(ms, float64((measure))/nsToMs)
	}

	median, _ := stats.Median(ms)
	average := total / int64(queries)

	results := Stats{
		Queries: queries,
		Median:  time.Duration(median * nsToMs),
		Average: time.Duration(average),
		Max:     time.Duration(max),
		Min:     time.Duration(min),
		Total:   time.Duration(total),
	}

	return results
}
