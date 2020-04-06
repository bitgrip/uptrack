// Copyright © 2018 Bitgrip <berlin@bitgrip.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctl

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/config"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/metrics"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"time"
)

// StartUpTrackServer is initialising the UpTrack Server
//
// It is:
// * Timing all jobs
// * Initialising Prometheus
// * Start and Listening the UpTrack Server
func StartUpTrackServer(config config.Config) error {
	logrus.Info("Start UpTrack server")

	descriptor, _ := job.DescriptorFromFile(config.JobConfigDir())
	if descriptor.DNSJobs == nil && descriptor.UpJobs == nil {
		logrus.Warn("No Jobs defined or not properly parsed from", config.JobConfigDir())
		return nil
	}
	registry := metrics.NewCombinedRegistry(
		metrics.NewPrometheusRegistry(config, descriptor),
		metrics.NewDatadogRegistry(config, descriptor),
	)
	if !registry.Enabled() {
		logrus.Info(fmt.Sprintf("No enabled registries found. \n Exiting."))
		return nil
	}

	logrus.Info(fmt.Sprintf("Load Jobs from %s\n", config.JobConfigDir()))

	endpoint := config.PrometheusEndpoint()
	port := strconv.Itoa(config.PrometheusPort())
	promUri := fmt.Sprintf(":" + port + endpoint)

	if config.PrometheusEnabled() {
		logrus.Info(fmt.Sprintf("Prometheus metrics available at  %v", promUri))
	} else {
		logrus.Info(fmt.Sprintf("Prometheus Endpoint is disabled"))
	}
	if config.DDEnabled() {
		logrus.Info(fmt.Sprintf("Initialize DataDog Registry for endpoint '%s'", config.DDEndpoint()))
		logrus.Info(fmt.Sprintf("DataDog Interval: '%ds'", int(config.DDInterval().Seconds())))
	} else {
		logrus.Info(fmt.Sprintf("Sending Metrics to DataDog is disabled"))
	}

	logrus.Info(fmt.Sprintf("Use check frequency %v", config.CheckFrequency()))
	info := "Initialized UpChecks:\n"

	for _, upJob := range descriptor.UpJobs {
		info = info + fmt.Sprintf("Name: '%s', url:'%s' \n", upJob.Name, upJob.URL)
	}
	logrus.Info(info)

	info = "Initialized DNSChecks:\n"

	for _, dnsJob := range descriptor.DNSJobs {
		info = info + fmt.Sprintf("Name: '%s', FQDN:'%s' \n", dnsJob.Name, dnsJob.FQDN)
	}
	logrus.Info(info)

	go runJobs(descriptor, registry, config.CheckFrequency())

	// starting server with prometheus endpoint
	http.Handle(endpoint, promhttp.Handler())

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	return nil
}

func runJobs(descriptor job.Descriptor, registry metrics.Registry, interval time.Duration) {
	for {
		//having one upJob iteration per interval
		startJobs := time.Now()
		//count iterations
		for _, upJob := range descriptor.UpJobs {
			doUpChecks(registry, upJob)
		}

		for _, dnsJob := range descriptor.DNSJobs {
			doDnsChecks(registry, dnsJob)
		}
		duration := startJobs.Add(interval).Sub(time.Now())

		if duration.Milliseconds() < 0 {
			duration = 0
		}

		time.Sleep(duration)

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
