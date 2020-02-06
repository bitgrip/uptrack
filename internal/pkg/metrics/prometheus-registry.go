package metrics

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"math"
)

// prometheusRegistry is a wrapper to forward Registry actions
// to a collection of Registries
type prometheusRegistry struct {
	metricsForChecks map[string]metrics
}

const jobName string = "job"
const checkName string = "check"
const url string = "url"

type metrics struct {
	Execution     prometheus.Counter //Execution Counter
	CanConnect    prometheus.Counter
	CannotConnect prometheus.Counter
	SSLDaysLeft   prometheus.Gauge
	ConnectTime   prometheus.Gauge
	TTFB          prometheus.Gauge
	RequestTime   prometheus.Gauge
	BytesReceived prometheus.Gauge
}

func NewPrometheusRegistry(descriptor job.Descriptor) Registry {
	localMetricsForChecks := make(map[string]metrics, 8)
	for name, upJob := range descriptor.UpJobs {
		localMetricsForChecks[name] = metrics{
			Execution:     counter("iterations", upJob),
			CanConnect:    counter("can_connect", upJob),
			CannotConnect: counter("can_not_connect", upJob),
			SSLDaysLeft:   gauge("ssl_days_left", upJob),
			ConnectTime:   gauge("connect_time", upJob),
			TTFB:          gauge("TTFB", upJob),
			RequestTime:   gauge("request_time", upJob),
			BytesReceived: gauge("bytes_received", upJob),
		}
	}

	return &prometheusRegistry{metricsForChecks: localMetricsForChecks}
}

func counter(name string, job job.UpJob) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "uptrack",
		Name:      "counter",
		ConstLabels: prometheus.Labels{
			jobName:   job.Name,
			checkName: name,
			url:       job.URL,
		},
	})
}

func gauge(name string, job job.UpJob) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "uptrack",
		Name:      "gauge",
		ConstLabels: prometheus.Labels{
			jobName:   job.Name,
			checkName: name,
			url:       job.URL,
		},
	})
}

func (r *prometheusRegistry) IncExecution(name string) {
	r.metricsForChecks[name].Execution.Inc()
}

func (r *prometheusRegistry) IncCanConnect(name string) {
	r.metricsForChecks[name].CanConnect.Inc()

}

func (r *prometheusRegistry) IncCanNotConnect(name string) {
	r.metricsForChecks[name].CannotConnect.Inc()
}

func (r *prometheusRegistry) SetSSLDaysLeft(name string, daysLeft float64) {
	r.metricsForChecks[name].SSLDaysLeft.Set(math.Round(daysLeft))
}

func (r *prometheusRegistry) SetConnectTime(name string, millis int64) {
	r.metricsForChecks[name].ConnectTime.Set(float64(millis))

}

func (r *prometheusRegistry) SetTTFB(name string, millis int64) {
	r.metricsForChecks[name].TTFB.Set(float64(millis))

}

func (r *prometheusRegistry) SetRequestTime(name string, millis int64) {
	r.metricsForChecks[name].RequestTime.Set(float64(millis))

}

func (r *prometheusRegistry) SetBytesReceived(name string, bytes int64) {
	r.metricsForChecks[name].BytesReceived.Set(float64(bytes))
}
