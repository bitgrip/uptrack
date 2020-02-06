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
	"gopkg.in/yaml.v2"
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
	logrus.Info(fmt.Sprintf("Use default interval %v", config.DefaultInterval()))
	logrus.Info(fmt.Sprintf("Load Jobs from %s\n", config.JobConfigDir()))

	endpoint := config.PrometheusEndpoint()
	port := strconv.Itoa(config.PrometheusPort())
	uri := fmt.Sprintf(":" + port + endpoint)

	descriptor, _ := descriptor(config.JobConfigDir())
	registry := metrics.NewCombinedRegistry(
		metrics.NewPrometheusRegistry(descriptor),
		metrics.NewDatadogRegistry(config.DatadogCredentials()),
	)
	logrus.Info(fmt.Sprintf("Prometheus metrics available at  %v", uri))

	go runJobs(descriptor, registry, config.DefaultInterval())

	// starting server with prometheus endpoint
	http.Handle(endpoint, promhttp.Handler())

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	return nil
}

func runJobs(descriptor job.Descriptor, registry metrics.Registry, interval int) {
	for {
		//having one upJob iteration per interval
		startJobs := time.Now()
		for _, upJob := range descriptor.UpJobs {
			doChecks(registry, upJob)
		}

		duration := (time.Duration(interval) * time.Second) - time.Since(startJobs)
		if duration.Seconds() <= 0 {
			duration = time.Duration(0)
		}
		time.Sleep(duration)
	}

}

func doChecks(registry metrics.Registry, upJob job.UpJob) {
	jobName := upJob.Name
	url := upJob.URL
	//count iterations
	registry.IncExecution(jobName)

	//prepare request
	req, _ := http.NewRequest(string(upJob.Method), url, upJob.Body())
	clientTrace := trace(registry, jobName)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &clientTrace))
	t := transport()

	//measure request time
	startReq := time.Now()
	resp, err := t.RoundTrip(req)
	endReq := time.Since(startReq)
	registry.SetRequestTime(jobName, endReq.Milliseconds())

	//time until expiry of ssl certs
	if upJob.CheckSSL {
		hours := time.Until(resp.TLS.PeerCertificates[0].NotAfter).Hours()
		registry.SetSSLDaysLeft(jobName, hours/24)
	}

	if err != nil {
		registry.IncCanNotConnect(jobName)

	} else {
		if resp.StatusCode == upJob.Expected {
			//count successful connection
			registry.IncCanConnect(jobName)

			//measure received bytes, body+headers
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			registry.SetBytesReceived(jobName, int64(len(bodyBytes)+len(resp.Header)))

		} else {
			registry.IncCanNotConnect(jobName)
		}
	}
}

func transport() http.Transport {
	return http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func descriptor(path string) (job.Descriptor, error) {
	d := job.Descriptor{}
	data, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(data, &d)

	for name, upJob := range d.UpJobs {
		upJob.Name = name
		upJob.ConcatUrl(d.BaseURL)
		d.UpJobs[name] = upJob
	}
	return d, err
}

func trace(registry metrics.Registry, name string) httptrace.ClientTrace {
	var connect time.Time
	start := time.Now()
	return httptrace.ClientTrace{

		ConnectStart: func(network, addr string) {
			connect = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			registry.SetConnectTime(name, time.Since(connect).Milliseconds())
		},

		GotFirstResponseByte: func() {
			registry.SetTTFB(name, time.Since(start).Milliseconds())

		},
	}
}
