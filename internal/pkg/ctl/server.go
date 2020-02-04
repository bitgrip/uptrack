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
	"strconv"

	"bitbucket.org/bitgrip/uptrack/internal/pkg/config"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
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
	uri := fmt.Sprintf(":" + port + "/" + endpoint)
	registry := metrics.NewCombinedRegistry(
		metrics.NewPrometheusRegistry(uri),
		metrics.NewDatadogRegistry(config.DatadogCredentials()),
	)
	logrus.Info(fmt.Sprintf("Prometheus metrics available at  %v", uri))
	logrus.Info(registry)
	http.Handle(endpoint, promhttp.Handler())
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	return nil
}
