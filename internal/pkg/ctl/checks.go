package ctl

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"regexp"
	"time"

	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/metrics"
	"github.com/sirupsen/logrus"
)

func doUpChecks(registry metrics.Registry, upJob *job.UpJob) {
	jobName := upJob.Name
	url := upJob.URL

	//prepare request
	req, _ := http.NewRequest(upJob.Method, url, upJob.Body())
	clientTrace := trace(registry, jobName)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &clientTrace))

	for k, v := range upJob.Headers {
		req.Header.Add(k, v)
	}

	t := transport(upJob)
	//measure request time
	//start here
	startReq := time.Now()
	resp, err := t.RoundTrip(req)
	if err != nil {
		logrus.Warn(fmt.Sprintf("Failed Request for job '%s'. msg: '%s' ", jobName, err))
		registry.IncCanNotConnect(jobName)
		return
	}
	//end measure request time
	endReq := time.Since(startReq)
	up := resp.StatusCode == upJob.Expected

	// if unexpected status code, check if bearer token expired
	if !up {
		if (upJob.Oauth.AuthUrl != "") && (resp.StatusCode == http.StatusUnauthorized) {
			logrus.Warn(fmt.Sprintf("Invalid Token for Job: '%s', url: '%s' ", upJob.Name, upJob.URL))
			upJob.Context.OauthResponse.Refresh = true
			return
		}
	}

	registry.SetRequestTime(jobName, float64(endReq.Milliseconds()))

	//time until expiry of ssl certs
	if upJob.CheckSSL {
		hours := time.Until(resp.TLS.PeerCertificates[0].NotAfter).Hours()
		registry.SetSSLDaysLeft(jobName, hours/24)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	} else {
		registry.SetBytesReceived(jobName, float64(len(bodyBytes)+len(resp.Header)))
	}
	if upJob.ContentMatch != "" {
		up = up && matchResponse(upJob.ContentMatch, bodyBytes, upJob.ReverseContentMatch)
	}
	if up {
		//count successful connections
		registry.IncCanConnect(jobName)

	} else {
		//count failed connections
		registry.IncCanNotConnect(jobName)
	}
}

func matchResponse(pattern string, bytes []byte, reverseMatch bool) bool {
	regex, _ := regexp.Compile(pattern)
	doesMatch := regex.Match(bytes)
	if reverseMatch {
		return !doesMatch
	}
	return doesMatch
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
