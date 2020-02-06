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

package cmd

import (
	"github.com/spf13/viper"
)

var uptrackConfig = config{}

type config struct {
}

func (config) JobConfigDir() string {
	return viper.GetString("uptrack.jobs_config")
}

func (config) DefaultInterval() int {
	return viper.GetInt("uptrack.default_interval")
}

func (config) DatadogCredentials() string {
	return viper.GetString("uptrack.datadog.credentials")
}

func (config) PrometheusEndpoint() string {
	return viper.GetString("uptrack.prometheus.endpoint")
}

func (config) PrometheusPort() int {
	return viper.GetInt("uptrack.prometheus.port")
}
