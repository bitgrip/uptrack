package dd

import (
	"sort"
)

// HistogramResult returns some descriptive statistics
// Add other stats as needed
type HistogramResult struct {
	count  float64
	min    float64
	max    float64
	avg    float64
	median float64
	p95    float64
}

// Histogram is the dumbest way possible to compute various descriptive statistics
//  It keeps all data, does a sort, and the figures out various stats.
//  That said for 1000 elements, it takes under 1/20 of a millisecond to compute.
//
// Also the "sort" method is what datadog's agent does, so it can't be too painful.
//
type Histogram struct {
	samples []float64
	tags    DDTags
}

// NewHistogram creates a new object
func NewHistogram(points int, tags DDTags) *Histogram {
	if points == 0 {
		return &Histogram{}
	}
	return &Histogram{
		samples: make([]float64, 0, points),
		tags:    tags,
	}
}

// Add adds a data point
func (h *Histogram) Add(c *Client, jobName string, checkName string, val float64) {
	h.samples = append(h.samples, val)
	res := h.calc()
	c.Rate(jobName, checkName+".count", res.count, h.tags)
	c.Gauge(jobName, checkName+".max", res.max, h.tags)
	c.Gauge(jobName, checkName+".avg", res.avg, h.tags)
	c.Gauge(jobName, checkName+".median", res.median, h.tags)
	c.Gauge(jobName, checkName+".95percentile", res.p95, h.tags)
}

func (h *Histogram) calc() HistogramResult {
	if len(h.samples) == 0 {
		// caller can check to see if count = 0
		return HistogramResult{}
	}

	sort.Float64s(h.samples)
	count := len(h.samples)
	sum := 0.0
	for _, val := range h.samples {
		sum += val
	}

	return HistogramResult{
		count:  float64(count),
		min:    h.samples[0],
		max:    h.samples[count-1],
		avg:    sum / float64(count),
		median: h.samples[count/2],
		p95:    h.samples[(count*95)/100],
	}
}
