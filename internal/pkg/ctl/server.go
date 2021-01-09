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
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/bitgrip/uptrack/internal/pkg/config"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/job"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// StartUpTrackServer is initialising the UpTrack Server
//
// It is:
// * Timing all jobs
// * Initialising Prometheus
// * Start and Listening the UpTrack Server
func StartUpTrackServer(config config.Config) error {
	logrus.Info("Start UpTrack server")
	endpoint := config.PrometheusEndpoint()
	port := strconv.Itoa(config.PrometheusPort())
	promUri := fmt.Sprintf(":" + port + endpoint)
	descriptor, _ := job.DescriptorFromFile(config.JobConfigDir())
	if descriptor.DNSJobs == nil && descriptor.UpJobs == nil {
		logrus.Warn("No Jobs defined or not properly parsed from", config.JobConfigDir())
		return nil
	}
	var registries []metrics.Registry

	if config.PrometheusEnabled() {
		logrus.Info(fmt.Sprintf("Prometheus metrics available at  %v", promUri))
		registries = append(registries, metrics.NewPrometheusRegistry(config, descriptor))
		go startPrometheus(endpoint, port)

	} else {
		logrus.Info(fmt.Sprintf("Prometheus Endpoint is disabled"))
	}
	if config.DDEnabled() {
		logrus.Info(fmt.Sprintf("Initialize DataDog Registry for endpoint '%s'", config.DDEndpoint()))
		logrus.Info(fmt.Sprintf("DataDog Interval: '%ds'", int(config.DDInterval().Seconds())))
		registries = append(registries, metrics.NewDatadogRegistry(config, descriptor))

	} else {
		logrus.Info(fmt.Sprintf("Sending Metrics to DataDog is disabled"))
	}

	registry := metrics.NewCombinedRegistry(registries)
	if !registry.Enabled() {
		logrus.Info(fmt.Sprintf("No enabled registries found. \n Exiting."))
		return nil
	}

	logrus.Info(fmt.Sprintf("Load Jobs from %s\n", config.JobConfigDir()))

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

	go runJobs(&descriptor, registry, config.CheckFrequency())

	for {
		time.Sleep(config.CheckFrequency())
	}
}

func startPrometheus(endpoint string, port string) {
	http.Handle(endpoint, promhttp.Handler())

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Fatal(err)
	}
}

func runJobs(descriptor *job.Descriptor, registry metrics.Registry, interval time.Duration) {
	for {
		//having one upJob iteration per interval
		startJobs := time.Now()
		//count iterations
		jobs := &descriptor.UpJobs
		for key, _ := range *jobs {
			doUpChecks(registry, (*jobs)[key])
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
