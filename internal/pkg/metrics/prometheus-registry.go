package metrics

import (
	"math"

	cons "bitbucket.org/bitgrip/uptrack/internal/pkg"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/config"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// prometheusRegistry is a wrapper to forward Registry actions
// to a collection of Registries
type prometheusRegistry struct {
	enabled             bool
	Execution           prometheus.Counter //Execution Counter
	metricsForUpChecks  map[string]metrics
	metricsForDnsChecks map[string]metrics
}

func (r *prometheusRegistry) Enabled() bool {
	return r.enabled
}

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

func NewPrometheusRegistry(config config.Config, descriptor job.Descriptor) Registry {
	projectName := replaceAll(descriptor.Name, " +")
	exec := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: cons.PromNamespace,
		Name:      "counter",
		ConstLabels: prometheus.Labels{
			projectName: projectName,
		},
	})

	localMetricsForChecks := make(map[string]metrics, 5)
	for name, upJob := range descriptor.UpJobs {
		localMetricsForChecks[name] = metrics{
			CanConnect:    checkCounter(projectName, cons.PromCanConnect, *upJob),
			CannotConnect: checkCounter(projectName, cons.PromCannotConnect, *upJob),
			SSLDaysLeft:   upGauge(projectName, cons.PromSSLDaysLeft, *upJob),
			ConnectTime:   upGauge(projectName, cons.PromConnectTime, *upJob),
			TTFB:          upGauge(projectName, cons.PromTTFB, *upJob),
			RequestTime:   upGauge(projectName, cons.PromRequestTime, *upJob),
			BytesReceived: upGauge(projectName, cons.PromBytesReceived, *upJob),
		}
	}
	localMetricsForDns := make(map[string]metrics, 5)

	for name, dnsJob := range descriptor.DNSJobs {
		localMetricsForDns[name] = metrics{
			DNSIpsRatio: dnsGauge(projectName, cons.PromFoundIps, dnsJob),
		}
	}

	return &prometheusRegistry{
		Execution:           exec,
		enabled:             config.PrometheusEnabled(),
		metricsForUpChecks:  localMetricsForChecks,
		metricsForDnsChecks: localMetricsForDns,
	}
}

func checkCounter(project string, check string, upJob job.UpJob) prometheus.Counter {
	host, _ := upJob.HostString()

	labels := prometheus.Labels{
		cons.ProjectName: project,
		cons.JobName:     upJob.Name,
		cons.Host:        host,
		cons.CheckName:   check,
		cons.UrlString:   upJob.URL,
	}
	for k, v := range upJob.CustomTags {
		labels[k] = v
	}
	return promauto.NewCounter(prometheus.CounterOpts{
		Namespace:   cons.PromNamespace,
		Name:        cons.PromNameUpcheckCounter,
		ConstLabels: labels,
	})
}

//func dnsCounter(project string, check string, dnsJob job.DnsJob) prometheus.Counter {
//	return promauto.NewCounter(prometheus.CounterOpts{
//		Namespace: cons.PromNamespace,
//		Name:      cons.PromNameDnsCheckCounter,
//		ConstLabels: prometheus.Labels{
//			cons.ProjectName: project,
//			cons.JobName:     dnsJob.Name,
//			cons.Host:        dnsJob.Host,
//			cons.CheckName:   check,
//			cons.FQDN:        dnsJob.FQDN,
//		},
//	})
//}

func upGauge(project string, check string, upJob job.UpJob) prometheus.Gauge {
	host, _ := upJob.HostString()

	labels := prometheus.Labels{
		cons.ProjectName: project,
		cons.JobName:     upJob.Name,
		cons.Host:        host,
		cons.CheckName:   check,
		cons.UrlString:   upJob.URL,
		cons.ReqMethod:   string(upJob.Method),
	}
	for k, v := range upJob.CustomTags {
		labels[k] = v
	}
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   cons.PromNamespace,
		Name:        cons.PromNameUpCheckGauge,
		ConstLabels: labels,
	})
}

func dnsGauge(project string, check string, dnsJob job.DnsJob) prometheus.Gauge {
	labels := prometheus.Labels{
		cons.ProjectName: project,
		cons.JobName:     dnsJob.Name,
		cons.Host:        dnsJob.Host,
		cons.CheckName:   check,
		cons.FQDN:        dnsJob.FQDN,
	}
	for k, v := range dnsJob.CustomTags {
		labels[k] = v
	}
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   cons.PromNamespace,
		Name:        cons.PromNameDnsCheckGauge,
		ConstLabels: labels,
	})
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
