package dd

import (
	"sort"
	"sync"
	"time"
)

// todo: make proper type
const (
	TypeGauge = "gauge"
	TypeRate  = "rate"
)

// Metric is a data structure that represents the JSON that Datadog
// wants when posting to the API
type Metric struct {
	Name     string        `json:"metric"`
	Value    [1][2]float64 `json:"points"`
	Type     string        `json:"type"`
	Hostname string        `json:"host,omitempty"`
	Tags     []string      `json:"tags,omitempty"`
	Interval int           `json:"interval,omitempty"`
}

func now() float64 {
	return float64(time.Now().Unix())
}

// NewMetric creates a new metric
func NewMetric(name string, mtype string, tags []string) *Metric {
	return &Metric{
		Name: name,
		Type: mtype,
		Tags: tags,
	}
}

// Client is the main datastructure of metrics to upload
type Client struct {
	Series    []*Metric          `json:"series"` // raw data
	hostname  string             // hostname
	tags      []string           // global tags, if any
	metrics   map[string]*Metric // map of name to metric for fast lookup
	now       func() float64     // for testing
	api       API                // where output goes
	lastFlush float64            // unix epoch as float64(t.Now().Unix())
	Stop      chan string
	sync.Mutex
}

// NewClient creates a new datadog metrics client
func NewClient(hostname string, api API) *Client {
	client := &Client{
		now:       now,
		hostname:  hostname,
		metrics:   make(map[string]*Metric),
		api:       api,
		lastFlush: now(),
		Stop:      make(chan string),
	}
	return client
}

// Gauge represents an observation
func (c *Client) Gauge(name string, value float64, tags []string) error {
	c.Lock()
	m, ok := c.metrics[name]
	if !ok {
		m = NewMetric(name, TypeGauge, unique(tags))
		c.Series = append(c.Series, m)
		c.metrics[name] = m
	}
	m.Value[0][1] = value
	c.Unlock()
	return nil
}

// Count represents a count of events
func (c *Client) Count(name string, value float64, tags []string) error {
	c.Lock()
	m, ok := c.metrics[name]
	if !ok {
		m = NewMetric(name, TypeRate, unique(tags))
		c.Series = append(c.Series, m)
		c.metrics[name] = m
	}
	// note, this sum must be divided by the interval length
	//  before sending.
	m.Value[0][1] += value
	c.Unlock()
	return nil
}

// Incr adds one event count, same as Count(name, 1)
func (c *Client) Incr(name string, tags []string) error {
	return c.Count(name, 1.0, tags)
}

// Decr subtracts one event, same as Count(name, -1)
func (c *Client) Decr(name string, tags []string) error {
	return c.Count(name, -1.0, tags)
}

// Snapshot makes a copy of the data and resets everything locally
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
		hostname:  c.hostname,
		Series:    c.Series,
		metrics:   c.metrics,
		lastFlush: c.lastFlush,
	}
	c.metrics = make(map[string]*Metric)
	c.Series = nil
	return &snap
}

// not locked.. for use locally with snapshots
func (c *Client) finalize(nowUnix float64) {
	interval := nowUnix - c.lastFlush

	// histograms: convert to various descriptive statistic gauges
	for i := 0; i < len(c.Series); i++ {
		c.Series[i].Value[0][0] = nowUnix
		c.Series[i].Hostname = c.hostname
		c.Series[i].Interval = int(interval)
		if c.Series[i].Type == "rate" {
			c.Series[i].Value[0][1] /= interval
		}
	}
}

// Flush forces a flush of the pending commands in the buffer
func (c *Client) Flush() error {
	if c == nil {
		return nil
	}
	snap := c.Snapshot()
	if snap == nil {
		return nil
	}
	snap.finalize(c.lastFlush)

	return c.api.postMetrics(snap.Series)
}

// Close the client connection.

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
