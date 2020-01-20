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

	"bitbucket.org/bitgrip/uptrack/internal/pkg/config"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/metrics"
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
	logrus.Info(fmt.Sprintf("Use default interval %v", config.DefaultInterval()))
	logrus.Info(fmt.Sprintf("Load configs from %s\n", config.JobConfigDir()))
	registry := metrics.NewCombinedRegistry(
		metrics.NewPrometheusRegistry(config.PrometheusScrape()),
		metrics.NewDatadogRegistry(config.DatadogCredentials()),
	)
	logrus.Info(registry)
	return nil
}
