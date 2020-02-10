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
const URL string = "URL"
const host = "host"
const FQDN string = "FQDN"
const projectName = "project"

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

const (
	prCanConnect         string = "connection_successful"
	prCannotConnect      string = "connection_failed"
	prSSLDaysLeft        string = "ssl_days_left"
	prConnectTime        string = "connection_time"
	prTTFB               string = "TTFB"
	prRequestTime        string = "request_time"
	prBytesReceived      string = "bytes_received"
	prFoundIps           string = "found_ips_ratio"
	prNamespace          string = "uptrack"
	prNameUpcheckCounter string = "upcheck_counter"
	prNameUpCheckGauge   string = "upcheck_gauge"

	prNameDnsCheckCounter string = "upcheck_counter"

	prNameDnsCheckGauge string = "dnscheck_gauge"
)

func NewPrometheusRegistry(descriptor job.Descriptor) Registry {
	projectName := descriptor.Name
	exec := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prNamespace,
		Name:      "counter",
		ConstLabels: prometheus.Labels{
			projectName: projectName,
		},
	})

	localMetricsForChecks := make(map[string]metrics, 5)
	for name, upJob := range descriptor.UpJobs {
		localMetricsForChecks[name] = metrics{
			CanConnect:    checkCounter(projectName, prCanConnect, upJob),
			CannotConnect: checkCounter(projectName, prCannotConnect, upJob),
			SSLDaysLeft:   upGauge(projectName, prSSLDaysLeft, upJob),
			ConnectTime:   upGauge(projectName, prConnectTime, upJob),
			TTFB:          upGauge(projectName, prTTFB, upJob),
			RequestTime:   upGauge(projectName, prRequestTime, upJob),
			BytesReceived: upGauge(projectName, prBytesReceived, upJob),
		}
	}
	localMetricsForDns := make(map[string]metrics, 5)

	for name, dnsJob := range descriptor.DNSJobs {
		localMetricsForDns[name] = metrics{
			DNSIpsRatio: dnsGauge(projectName, prFoundIps, dnsJob),
		}
	}

	return &prometheusRegistry{Execution: exec, metricsForUpChecks: localMetricsForChecks, metricsForDnsChecks: localMetricsForDns}
}

func checkCounter(project string, name string, job job.UpJob) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prNamespace,
		Name:      prNameUpcheckCounter,
		ConstLabels: prometheus.Labels{
			projectName: project,
			jobName:     job.Name,
			checkName:   name,
			URL:         job.URL,
		},
	})
}

func dnsCounter(project string, name string, job job.DnsJob) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prNamespace,
		Name:      prNameDnsCheckCounter,
		ConstLabels: prometheus.Labels{
			projectName: project,
			jobName:     job.Name,
			checkName:   name,
			FQDN:        job.FQDN,
		},
	})
}

func upGauge(project string, name string, job job.UpJob) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: prNamespace,
		Name:      prNameUpCheckGauge,
		ConstLabels: prometheus.Labels{
			projectName: project,
			jobName:     job.Name,
			checkName:   name,
			URL:         job.URL,
		},
	})
}

func dnsGauge(project string, name string, job job.DnsJob) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: prNamespace,
		Name:      prNameDnsCheckGauge,
		ConstLabels: prometheus.Labels{
			projectName: project,
			jobName:     job.Name,
			checkName:   name,
			FQDN:        job.FQDN,
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

func (r *prometheusRegistry) SetConnectTime(name string, millis float64) {
	r.metricsForUpChecks[name].ConnectTime.Set(millis)

}

func (r *prometheusRegistry) SetTTFB(name string, millis float64) {
	r.metricsForUpChecks[name].TTFB.Set(millis)

}

func (r *prometheusRegistry) SetRequestTime(name string, millis float64) {
	r.metricsForUpChecks[name].RequestTime.Set(millis)

}

func (r *prometheusRegistry) SetBytesReceived(name string, bytes float64) {
	r.metricsForUpChecks[name].BytesReceived.Set(bytes)
}
func (r *prometheusRegistry) SetIpsRatio(job string, ratio float64) {
	r.metricsForDnsChecks[job].DNSIpsRatio.Set(ratio)

}
