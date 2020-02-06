package metrics

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// prometheusRegistry is a wrapper to forward Registry actions
// to a collection of Registries
type prometheusRegistry struct {
	metricsForChecks map[string]metrics
}
type metrics struct {
	Execution     prometheus.Counter //Execution Counter
	CanConnect    prometheus.Counter
	CannotConnect prometheus.Counter
	TTFB          prometheus.Gauge
}

func NewPrometheusRegistry(descriptor job.Descriptor) Registry {

	metricsMap := make(map[string]metrics, 5)
	for name, upCheck := range descriptor.UpChecks {
		var (
			execution = promauto.NewCounter(prometheus.CounterOpts{
				Name: "upcheck_" + name + "_total_count",
				Help: "The total number of check Iterations",
			})

			canConnect = promauto.NewCounter(prometheus.CounterOpts{
				Name: "upcheck_" + name + "_can_connect",
				Help: "TBD",
			})
			canNotConnect = promauto.NewCounter(prometheus.CounterOpts{
				Name: "upcheck_" + name + "_can_not_connect",
				Help: "TBD",
			})

			ttfb = promauto.NewGauge(prometheus.GaugeOpts{
				Name: "upcheck_" + name + "_TTFB",
				Help: "TBD",
			})
		)
		localMetrics := metrics{
			Execution:     execution,
			CanConnect:    canConnect,
			CannotConnect: canNotConnect,
			TTFB:          ttfb,
		}
		metricsMap[upCheck.Name] = localMetrics
	}

	return &prometheusRegistry{metricsForChecks: metricsMap}
}

func (r *prometheusRegistry) IncExecution(name string) {
	r.metricsForChecks[name].Execution.Inc()
}

func (r *prometheusRegistry) IncCanConnect(name string, uri string) {
	r.metricsForChecks[name].CanConnect.Inc()

}

func (r *prometheusRegistry) IncCanNotConnect(name string, uri string) {
	r.metricsForChecks[name].CannotConnect.Inc()
}

func (r *prometheusRegistry) SetSSLDaysLeft(name string, uri string, daysLeft int64) {
}

func (r *prometheusRegistry) SetConnectTime(name string, uri string, millis int64) {
}

func (r *prometheusRegistry) SetTTFB(name string, uri string, millis int64) {
	r.metricsForChecks[name].TTFB.Set(float64(millis))

}

func (r *prometheusRegistry) SetRequestTime(name string, uri string, millis int64) {
}

func (r *prometheusRegistry) SetBytesReceived(name string, uri string, bytes int64) {
}
