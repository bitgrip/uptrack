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
	Execution           prometheus.Counter //Execution Counter
	metricsForUpChecks  map[string]metrics
	metricsForDnsChecks map[string]metrics
}

const jobName string = "job"
const checkName string = "check"
const url string = "url"
const FQDN string = "FQDN"
const Project = "project"

type metrics struct {
	CanConnect    prometheus.Counter
	CannotConnect prometheus.Counter
	SSLDaysLeft   prometheus.Gauge
	ConnectTime   prometheus.Gauge
	TTFB          prometheus.Gauge
	RequestTime   prometheus.Gauge
	BytesReceived prometheus.Gauge
	DNSIpsRatio   prometheus.Gauge
}

func NewPrometheusRegistry(descriptor job.Descriptor) Registry {
	projectName := descriptor.Name
	exec := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "uptrack",
		Name:      "counter",
		ConstLabels: prometheus.Labels{
			Project: projectName,
		},
	})

	localMetricsForChecks := make(map[string]metrics, 8)
	for name, upJob := range descriptor.UpJobs {
		localMetricsForChecks[name] = metrics{
			CanConnect:    checkCounter(projectName, "can_connect", upJob),
			CannotConnect: checkCounter(projectName, "can_not_connect", upJob),
			SSLDaysLeft:   upGauge(projectName, "ssl_days_left", upJob),
			ConnectTime:   upGauge(projectName, "connect_time", upJob),
			TTFB:          upGauge(projectName, "TTFB", upJob),
			RequestTime:   upGauge(projectName, "request_time", upJob),
			BytesReceived: upGauge(projectName, "bytes_received", upJob),
		}
	}
	localMetricsForDns := make(map[string]metrics, 8)

	for name, dnsJob := range descriptor.DNSJobs {
		localMetricsForDns[name] = metrics{
			DNSIpsRatio: dnsGauge(projectName, "found_ips_ratio", dnsJob),
		}
	}

	return &prometheusRegistry{Execution: exec, metricsForUpChecks: localMetricsForChecks, metricsForDnsChecks: localMetricsForDns}
}

func checkCounter(project string, name string, job job.UpJob) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "uptrack",
		Name:      "upcheck_counter",
		ConstLabels: prometheus.Labels{
			Project:   project,
			jobName:   job.Name,
			checkName: name,
			url:       job.URL,
		},
	})
}

func dnsCounter(project string, name string, job job.DnsJob) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "uptrack",
		Name:      "dns_counter",
		ConstLabels: prometheus.Labels{
			Project:   project,
			jobName:   job.Name,
			checkName: name,
			FQDN:      job.FQDN,
		},
	})
}

func upGauge(project string, name string, job job.UpJob) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "uptrack",
		Name:      "upcheck_gauge",
		ConstLabels: prometheus.Labels{
			Project:   project,
			jobName:   job.Name,
			checkName: name,
			url:       job.URL,
		},
	})
}

func dnsGauge(project string, name string, job job.DnsJob) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "uptrack",
		Name:      "dns_gauge",
		ConstLabels: prometheus.Labels{
			Project:   project,
			jobName:   job.Name,
			checkName: name,
			FQDN:      job.FQDN,
		},
	})
}
func (r *prometheusRegistry) IncExecution(name string) {
	r.Execution.Inc()
}

func (r *prometheusRegistry) IncCanConnect(name string) {
	r.metricsForUpChecks[name].CanConnect.Inc()

}

func (r *prometheusRegistry) IncCanNotConnect(name string) {
	r.metricsForUpChecks[name].CannotConnect.Inc()
}

func (r *prometheusRegistry) SetSSLDaysLeft(name string, daysLeft float64) {
	r.metricsForUpChecks[name].SSLDaysLeft.Set(math.Round(daysLeft))
}

func (r *prometheusRegistry) SetConnectTime(name string, millis int64) {
	r.metricsForUpChecks[name].ConnectTime.Set(float64(millis))

}

func (r *prometheusRegistry) SetTTFB(name string, millis int64) {
	r.metricsForUpChecks[name].TTFB.Set(float64(millis))

}

func (r *prometheusRegistry) SetRequestTime(name string, millis int64) {
	r.metricsForUpChecks[name].RequestTime.Set(float64(millis))

}

func (r *prometheusRegistry) SetBytesReceived(name string, bytes int64) {
	r.metricsForUpChecks[name].BytesReceived.Set(float64(bytes))
}
func (r *prometheusRegistry) SetIpsRatio(job string, ratio float64) {
	r.metricsForDnsChecks[job].DNSIpsRatio.Set(ratio)

}
