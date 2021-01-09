package ctl

import (
	"encoding/base64"
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
	//Perform authentication via oauth client credentials flow
	var bearerToken string
	if upJob.BearerToken == "" {
		bearerToken = getBearerToken(upJob)
		upJob.BearerToken = bearerToken
	} else {
		bearerToken = upJob.BearerToken
	}
	if upJob.OauthClientCredentials != nil {
	}

	//prepare request
	req, _ := http.NewRequest(string(upJob.Method), url, upJob.Body())
	clientTrace := trace(registry, jobName)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &clientTrace))
	//Add headers
	//if bearerToken != "" {
	//	req.Header.Add("Authorization", "Bearer "+bearerToken)
	//}
	for k, v := range upJob.Headers {
		req.Header.Add(k, v)
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

	up := resp.StatusCode == upJob.Expected
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
		registry.IncCanNotConnect(jobName)
	}

}

func getBearerToken(upJob *job.UpJob) string {
	authReq, _ := http.NewRequest("GET", upJob.OauthClientCredentials["auth_url"], upJob.Body())
	//Add basic auth header
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(upJob.OauthClientCredentials["client_id"]+":"+upJob.OauthClientCredentials["client_secret"]))
	authReq.Header.Add("Authorization", authHeader)

	client := &http.Client{}
	authResp, err := client.Do(authReq)
	if err != nil {
		logrus.Fatal(err)

	}
	bytes, _ := ioutil.ReadAll(authResp.Body)
	return string(bytes)
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
