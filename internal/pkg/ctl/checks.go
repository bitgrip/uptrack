package ctl

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"

	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/metrics"
	"github.com/sirupsen/logrus"
)

func doUpChecks(registry metrics.Registry, upJob job.UpJob) {
	jobName := upJob.Name
	url := upJob.URL
	//prepare request
	req, _ := http.NewRequest(string(upJob.Method), url, upJob.Body())
	clientTrace := trace(registry, jobName)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &clientTrace))
	//Add headers
	for k, v := range upJob.Header {
		req.Header.Add(k, v[0])
	}

	t := transport()
	//measure request time
	startReq := time.Now()
	resp, err := t.RoundTrip(req)
	if err != nil {
		logrus.Warn(fmt.Sprintf("Failed Request for job '%s'. msg: '%s' ", jobName, err))
		registry.IncCanNotConnect(jobName)
		return
	}
	endReq := time.Since(startReq)
	registry.SetRequestTime(jobName, float64(endReq.Milliseconds()))

	//time until expiry of ssl certs
	if upJob.CheckSSL {
		hours := time.Until(resp.TLS.PeerCertificates[0].NotAfter).Hours()
		registry.SetSSLDaysLeft(jobName, hours/24)
	}

	if resp.StatusCode == upJob.Expected {
		//count successful connections
		registry.IncCanConnect(jobName)

		//measure received bytes, body+headers
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		registry.SetBytesReceived(jobName, float64(len(bodyBytes)+len(resp.Header)))

	} else {
		registry.IncCanNotConnect(jobName)
	}
}

func doDnsChecks(registry metrics.Registry, dnsJob job.DnsJob) {

	jobName := dnsJob.Name
	actIps, err := net.LookupIP(dnsJob.FQDN)

	if err != nil {
		logrus.Warn(fmt.Sprintf("Failed Request for job '%s'. msg: '%s' ", jobName, err))

		return
	}

	actIpsS := make([]string, 0)
	for _, v := range actIps {
		actIpsS = append(actIpsS, v.String())
	}
	expIps := dnsJob.IPs

	intersecIps := GetIntersecting(expIps, actIpsS)
	if len(expIps) == 0 {
		registry.SetIpsRatio(jobName, 0)
	}
	registry.SetIpsRatio(jobName, float64(len(intersecIps)/len(expIps)))
}

func GetIntersecting(expIps []string, actIps []string) []string {
	intersecIps := make([]string, 0)
	for _, expIp := range expIps {
		for _, actIp := range actIps {
			if actIp == expIp {
				intersecIps = append(intersecIps, expIp)
			}
		}

	}
	return intersecIps
}
