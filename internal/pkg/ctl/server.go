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
	"crypto/tls"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
		//having one check per interval
		start := time.Now()
		for _, check := range descriptor.UpChecks {
			doChecks(descriptor, registry, check)
		}

		duration := (time.Duration(interval) * time.Second) - time.Since(start)
		if duration.Seconds() <= 0 {
			duration = time.Duration(0)
		}

		time.Sleep(duration)
	}

}

func doChecks(descriptor job.Descriptor, registry metrics.Registry, check job.UpCheck) {
	name := check.Name
	registry.IncExecution(name)
	url := descriptor.BaseURL
	req, _ := http.NewRequest(string(check.Method), url, nil)

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace(registry, name)))
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
	}
	if resp.StatusCode == 200 {
		registry.IncCanConnect(name, "")
	} else {
		registry.IncCanNotConnect(name, "")
	}
}

//TODO errorhandling!!!
func descriptor(path string) (job.Descriptor, error) {
	d := job.Descriptor{}
	data, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(data, &d)
	return d, err
}

func trace(registry metrics.Registry, name string) *httptrace.ClientTrace {
	var connect, dns, tlsHandshake time.Time
	start := time.Now()
	return &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
		},

		GotFirstResponseByte: func() {
			registry.SetTTFB(name, "", time.Since(start).Milliseconds())

		},
	}
}
