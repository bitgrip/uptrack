package dd

import (
	cons "bitbucket.org/bitgrip/uptrack/internal/pkg"
	"github.com/sirupsen/logrus"
	"sort"
	"sync"
	"time"
)

// todo: make proper type
const (
	TypeGauge   = "gauge"
	TypeRate    = "rate"
	TypeCounter = "counter"
)

// Metric is a data structure that represents the JSON that Datadog
// expects when posting to the API
type Metric struct {
	Name     string        `json:"metric"`
	Value    [1][2]float64 `json:"points"`
	Type     string        `json:"type"`
	Hostname string        `json:"host,omitempty"`
	Tags     []string      `json:"tags,omitempty"`
	Interval int           `json:"ratio,omitempty"`
}

func now() float64 {
	return float64(time.Now().Unix())
}

// NewMetric creates a new metric
func NewMetric(name string, host string, mtype string, tags []string) *Metric {
	return &Metric{
		Name:     name,
		Hostname: host,
		Type:     mtype,
		Tags:     tags,
	}
}

// Client is the main datastructure of metricsMap to upload
type Client struct {
	Series     []*Metric `json:"series"` // raw data
	histograms map[string]*Histogram
	tags       []string           // global tags, if any
	metricsMap map[string]*Metric // map of name to metric for fast lookup
	now        func() float64     // for testing
	api        API                // where output goes
	lastFlush  float64            // unix epoch as float64(t.Now().Unix())
	Stop       chan string
	sync.Mutex
}

// NewClient creates a new datadog metricsMap client
func NewClient(api API, in float64) *Client {
	client := &Client{
		now:        now,
		histograms: make(map[string]*Histogram),
		metricsMap: make(map[string]*Metric),
		api:        api,
		lastFlush:  now(),
		Stop:       make(chan string),
	}
	return client
}

// Gauge represents an observation
func (c *Client) Gauge(jobName string, checkName string, value float64, tags DDTags) error {
	c.Lock()
	m, ok := c.metricsMap[jobName+checkName]
	if !ok {
		m = NewMetric(checkName, tags[cons.Host], TypeGauge, tags.ToTagList())
		c.Series = append(c.Series, m)
		c.metricsMap[checkName] = m
	}
	m.Value[0][1] = value
	c.Unlock()
	return nil
}

// Rate represents a count of events
func (c *Client) Rate(jobName string, checkName string, value float64, tags DDTags) error {
	c.Lock()
	m, ok := c.metricsMap[jobName+checkName]
	if !ok {
		m = NewMetric(checkName, tags[cons.Host], TypeRate, tags.ToTagList())
		c.Series = append(c.Series, m)
		c.metricsMap[checkName] = m
	}
	m.Value[0][1] += value
	c.Unlock()
	return nil
}

func (c *Client) Timing(jobName string, checkName string, val float64, tags DDTags) error {
	return c.Hist(jobName, checkName, val, tags)
}

func (c *Client) Hist(jobName string, checkName string, val float64, tags DDTags) error {
	c.Lock()
	h := c.histograms[jobName+checkName]
	if h == nil {
		h = NewHistogram(1000, tags)
		c.histograms[jobName+checkName] = h
	}
	c.Unlock()
	h.Add(c, jobName, checkName, val)

	return nil
}

func (c *Client) Incr(jobName string, checkName string, tags DDTags) error {

	return c.Rate(jobName, checkName, 1.0, tags)
}

func (c *Client) Decr(jobName string, checkName string, tags DDTags) error {
	return c.Rate(jobName, checkName, -1.0, tags)
}

func (c *Client) Snapshot() *Client {
	c.Lock()
	defer func() {
		c.lastFlush = c.now()
		c.Unlock()
	}()

	if len(c.Series) == 0 {
		return nil
	}
	snap := Client{
		Series:     c.Series,
		metricsMap: c.metricsMap,
		histograms: c.histograms,
		lastFlush:  c.lastFlush,
	}
	c.metricsMap = make(map[string]*Metric)
	c.Series = nil
	return &snap
}

// not locked.. for use locally with snapshots
func (c *Client) finalize(nowUnix float64) {
	interval := nowUnix - c.lastFlush

	for i, m := range c.Series {
		c.Series[i].Value[0][0] = nowUnix
		c.Series[i].Hostname = m.Hostname
		c.Series[i].Interval = int(interval)
		if c.Series[i].Type == TypeRate {
			logrus.Warn("XXX: ", c.Series[i].Value[0][1])
			c.Series[i].Value[0][1] /= interval
		}
	}

}

func (c *Client) Flush() error {
	if c == nil {
		return nil
	}
	snap := c.Snapshot()
	if snap == nil {
		return nil
	}
	snap.finalize(c.lastFlush)

	return c.api.postSeries(snap.Series)
}

func unique(s []string) []string {
	if len(s) < 2 {
		return s
	}
	sort.Strings(s)
	j := 1
	for i := 1; i < len(s); i++ {
		if s[j-1] != s[i] {
			s[j] = s[i]
			j++
		}
	}
	return s[:j]
}
