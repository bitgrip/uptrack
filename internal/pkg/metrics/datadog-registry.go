package metrics

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/config"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"fmt"
	"github.com/signalsciences/dogdirect"
	"net/url"
)

// datadogRegistry is a wrapper to forward Registry actions
// to a collection of Registries
type datadogRegistry struct {
	Client               *dogdirect.Client
	Periodic             *dogdirect.Periodic
	ExecutionCounterTags ddTags
	tagsForChecks        map[string]ddTagStruct
	keysForChecks        map[string]metricKeys
}

type ddTags map[string]string

const metricsRootName = "uptrack"

const (
	ddCanConnect    string = "connection.successful"
	ddCannotConnect string = "connection.failed"
	ddSSLDaysLeft   string = "ssl_days_left"
	ddConnectTime   string = "connection.time"
	ddTTFB          string = "TTFB"
	ddRequestTime   string = "request_time"
	ddBytesReceived string = "bytes_received"
	ddFoundIps      string = "found_ips_ratio"
)

type ddTagStruct struct {
	CanConnect    ddTags
	CannotConnect ddTags
	SSLDaysLeft   ddTags
	ConnectTime   ddTags
	TTFB          ddTags
	RequestTime   ddTags
	BytesReceived ddTags
	DNSIpsRatio   ddTags
}

type metricKeys struct {
	CanConnect    string
	CannotConnect string
	SSLDaysLeft   string
	ConnectTime   string
	TTFB          string
	RequestTime   string
	BytesReceived string
	DNSIpsRatio   string
}

func NewDatadogRegistry(config config.Config, descriptor job.Descriptor) Registry {
	api := dogdirect.NewAPI(config.DDApiKey(), config.DDAppKey(), config.DDInterval())
	client := dogdirect.New(replaceAll(descriptor.Name, " "), api)
	periodicClient := dogdirect.NewPeriodic(client, config.DDInterval())

	localTagsForChecks := make(map[string]ddTagStruct, 5)

	localKeysForChecks := make(map[string]metricKeys, 5)

	for name, upJob := range descriptor.UpJobs {
		localTagsForChecks[name] = ddTagStruct{
			CanConnect:    tags(descriptor, upJob, ddCanConnect),
			CannotConnect: tags(descriptor, upJob, ddCannotConnect),
			SSLDaysLeft:   tags(descriptor, upJob, ddSSLDaysLeft),
			ConnectTime:   tags(descriptor, upJob, ddConnectTime),
			TTFB:          tags(descriptor, upJob, ddTTFB),
			RequestTime:   tags(descriptor, upJob, ddRequestTime),
			BytesReceived: tags(descriptor, upJob, ddBytesReceived),
		}

		localKeysForChecks[name] = metricKeys{
			CanConnect:    keys(descriptor.Name, upJob.Name, ddCanConnect),
			CannotConnect: keys(descriptor.Name, upJob.Name, ddCannotConnect),
			SSLDaysLeft:   keys(descriptor.Name, upJob.Name, ddSSLDaysLeft),
			ConnectTime:   keys(descriptor.Name, upJob.Name, ddConnectTime),
			TTFB:          keys(descriptor.Name, upJob.Name, ddTTFB),
			RequestTime:   keys(descriptor.Name, upJob.Name, ddRequestTime),
			BytesReceived: keys(descriptor.Name, upJob.Name, ddBytesReceived),
		}

	}

	for name, dnsJob := range descriptor.DNSJobs {
		localTagsForChecks[name] = ddTagStruct{
			DNSIpsRatio: setDnsJobTags(descriptor, dnsJob, ddFoundIps),
		}

		localKeysForChecks[name] = metricKeys{
			DNSIpsRatio: keys(descriptor.Name, dnsJob.Name, ddFoundIps),
		}
	}
	executionCounterTags := ddTags{
		projectName: descriptor.Name,
		checkName:   "uptrack_counter",
	}

	return &datadogRegistry{Client: client, Periodic: periodicClient,
		keysForChecks:        localKeysForChecks,
		tagsForChecks:        localTagsForChecks,
		ExecutionCounterTags: executionCounterTags}
}

func (r *datadogRegistry) IncExecution(job string) {
	r.Client.Incr(metricsRootName+"."+job+"."+"counter", r.ExecutionCounterTags.toTagList())
}

func (r *datadogRegistry) IncCanConnect(job string) {
	r.Client.Incr(r.keysForChecks[job].CanConnect, r.tagsForChecks[job].CanConnect.toTagList())

}

func (r *datadogRegistry) IncCanNotConnect(job string) {
	r.Client.Incr(r.keysForChecks[job].CannotConnect, r.tagsForChecks[job].CannotConnect.toTagList())

}

func (r *datadogRegistry) SetSSLDaysLeft(job string, daysLeft float64) {
	r.Client.Gauge(r.keysForChecks[job].SSLDaysLeft, daysLeft, r.tagsForChecks[job].SSLDaysLeft.toTagList())
}

func (r *datadogRegistry) SetConnectTime(job string, millis float64) {
	r.Client.Gauge(r.keysForChecks[job].ConnectTime, millis, r.tagsForChecks[job].ConnectTime.toTagList())

}

func (r *datadogRegistry) SetTTFB(job string, millis float64) {
	r.Client.Gauge(r.keysForChecks[job].TTFB, millis, r.tagsForChecks[job].TTFB.toTagList())

}

func (r *datadogRegistry) SetRequestTime(job string, millis float64) {
	r.Client.Gauge(r.keysForChecks[job].RequestTime, millis, r.tagsForChecks[job].RequestTime.toTagList())

}

func (r *datadogRegistry) SetBytesReceived(job string, bytes float64) {
	r.Client.Gauge(r.keysForChecks[job].BytesReceived, bytes, r.tagsForChecks[job].BytesReceived.toTagList())

}

func (r *datadogRegistry) SetIpsRatio(job string, ratio float64) {
	r.Client.Gauge(r.keysForChecks[job].DNSIpsRatio, ratio, r.tagsForChecks[job].DNSIpsRatio.toTagList())

}

func (t ddTags) toTagList() []string {
	out := make([]string, 0)
	for k, v := range t {
		out = append(out, k+":"+v)
	}
	return out
}
func keys(project string, job string, check string) string {

	project = replaceAll(project, " +")

	return fmt.Sprintf("%s.%s.%s.%s", metricsRootName, project, job, check) // metricsRootName
}

func setDnsJobTags(descriptor job.Descriptor, dnsJob job.DnsJob, name string) ddTags {
	return ddTags{
		projectName: descriptor.Name,
		jobName:     dnsJob.Name,
		host:        dnsJob.FQDN,
		checkName:   name,
		FQDN:        dnsJob.FQDN,
	}
}
func tags(descriptor job.Descriptor, upJob job.UpJob, name string) ddTags {
	u, _ := url.Parse(upJob.URL)

	return ddTags{
		projectName: descriptor.Name,
		jobName:     upJob.Name,
		host:        u.Host,
		checkName:   name,
		URL:         upJob.URL,
	}
}
