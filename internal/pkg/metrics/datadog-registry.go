package metrics

import (
	"fmt"

	cons "bitbucket.org/bitgrip/uptrack/internal/pkg"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/api/dd"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/config"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
)

// datadogRegistry is a wrapper to forward Registry actions
// to a collection of Registries
type datadogRegistry struct {
	enabled              bool
	Client               *dd.Client
	ExecutionCounterTags dd.DDTags
	tagsForChecks        map[string]ddTags
	keysForChecks        map[string]metricKeys
}

func (r *datadogRegistry) Enabled() bool {
	return r.enabled
}

const metricsRootName = "uptrack"

type ddTags struct {
	Execution     dd.DDTags
	CanConnect    dd.DDTags
	CannotConnect dd.DDTags
	SSLDaysLeft   dd.DDTags
	ConnectTime   dd.DDTags
	TTFB          dd.DDTags
	RequestTime   dd.DDTags
	BytesReceived dd.DDTags
	DNSIpsRatio   dd.DDTags
}

type metricKeys struct {
	Execution     string
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
	api := dd.NewAPI(config.DDEndpoint(), config.DDApiKey(), config.DDAppKey())
	client := dd.NewClient(api, config.DDInterval().Seconds())
	client.Watch(config.DDInterval())
	localTagsForChecks := make(map[string]ddTags, 5)
	localKeysForChecks := make(map[string]metricKeys, 5)

	for name, upJob := range descriptor.UpJobs {
		localTagsForChecks[name] = ddTags{
			CanConnect:    upTags(descriptor, *upJob, cons.DDCanConnect),
			CannotConnect: upTags(descriptor, *upJob, cons.DDCannotConnect),
			SSLDaysLeft:   upTags(descriptor, *upJob, cons.DDSSLDaysLeft),
			ConnectTime:   upTags(descriptor, *upJob, cons.DDConnectTime),
			TTFB:          upTags(descriptor, *upJob, cons.DDTTFB),
			RequestTime:   upTags(descriptor, *upJob, cons.DDRequestTime),
			BytesReceived: upTags(descriptor, *upJob, cons.DDBytesReceived),
		}

		localKeysForChecks[name] = metricKeys{
			CanConnect:    keys(cons.DDCanConnect),
			CannotConnect: keys(cons.DDCannotConnect),
			SSLDaysLeft:   keys(cons.DDSSLDaysLeft),
			ConnectTime:   keys(cons.DDConnectTime),
			TTFB:          keys(cons.DDTTFB),
			RequestTime:   keys(cons.DDRequestTime),
			BytesReceived: keys(cons.DDBytesReceived),
		}

	}

	for name, dnsJob := range descriptor.DNSJobs {
		localTagsForChecks[name] = ddTags{
			DNSIpsRatio: dnsTags(descriptor, dnsJob, cons.DDFoundIps),
		}

		localKeysForChecks[name] = metricKeys{
			DNSIpsRatio: keys(cons.DDFoundIps),
		}
	}

	return &datadogRegistry{
		Client:        client,
		enabled:       config.DDEnabled(),
		keysForChecks: localKeysForChecks,
		tagsForChecks: localTagsForChecks,
	}
}

func (r *datadogRegistry) IncCanConnect(job string) {
	r.Client.Gauge(job+"_up", r.keysForChecks[job].CanConnect, 1, r.tagsForChecks[job].CanConnect)
	r.Client.Gauge(job+"_up", r.keysForChecks[job].CannotConnect, 0, r.tagsForChecks[job].CannotConnect)

}

func (r *datadogRegistry) IncCanNotConnect(job string) {
	r.Client.Gauge(job+"_up", r.keysForChecks[job].CanConnect, 0, r.tagsForChecks[job].CanConnect)
	r.Client.Gauge(job+"_up", r.keysForChecks[job].CannotConnect, 1, r.tagsForChecks[job].CannotConnect)

}

func (r *datadogRegistry) SetSSLDaysLeft(job string, daysLeft float64) {
	r.Client.Gauge(job+"_up", r.keysForChecks[job].SSLDaysLeft, daysLeft, r.tagsForChecks[job].SSLDaysLeft)
}

func (r *datadogRegistry) SetConnectTime(job string, millis float64) {
	r.Client.Gauge(job+"_up", r.keysForChecks[job].ConnectTime, millis, r.tagsForChecks[job].ConnectTime)

}

func (r *datadogRegistry) SetTTFB(job string, millis float64) {
	r.Client.Gauge(job+"_up", r.keysForChecks[job].TTFB, millis, r.tagsForChecks[job].TTFB)

}

func (r *datadogRegistry) SetRequestTime(job string, millis float64) {
	r.Client.Gauge(job+"_up", r.keysForChecks[job].RequestTime, millis, r.tagsForChecks[job].RequestTime)

}

func (r *datadogRegistry) SetBytesReceived(job string, bytes float64) {
	r.Client.Gauge(job+"_up", r.keysForChecks[job].BytesReceived, bytes, r.tagsForChecks[job].BytesReceived)

}

func (r *datadogRegistry) SetIpsRatio(job string, ratio float64) {
	r.Client.Gauge(job+"_dns", r.keysForChecks[job].DNSIpsRatio, ratio, r.tagsForChecks[job].DNSIpsRatio)
}

func dnsTags(descriptor job.Descriptor, dnsJob job.DnsJob, check string) dd.DDTags {
	tags := dd.DDTags{
		cons.ProjectName: descriptor.Name,
		cons.JobName:     dnsJob.Name,
		cons.Host:        dnsJob.Host,
		cons.CheckName:   check,
		cons.FQDN:        dnsJob.FQDN,
	}
	for k, v := range dnsJob.CustomTags {
		tags[k] = v
	}
	return tags
}

func upTags(descriptor job.Descriptor, upJob job.UpJob, name string) dd.DDTags {

	host, _ := upJob.HostString()
	tags := dd.DDTags{
		cons.ProjectName: descriptor.Name,
		cons.JobName:     upJob.Name,
		cons.Host:        host,
		cons.CheckName:   name,
		cons.UrlString:   upJob.URL,
		cons.ReqMethod:   string(upJob.Method),
	}

	for k, v := range upJob.CustomTags {
		tags[k] = v
	}
	return tags
}
func keys(check string) string {
	return fmt.Sprintf("%s.%s", metricsRootName, check)
}
