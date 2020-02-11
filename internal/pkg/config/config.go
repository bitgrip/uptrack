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

package config

import "time"

// Config provides information for DD and Prometheus configurations
type Config interface {
	JobConfigDir() string
	DefaultInterval() time.Duration
	PrometheusEndpoint() string
	PrometheusPort() int
	DatadogCredentials() string
	DDEndpoint() string
	DDApiKey() string
	DDAppKey() string
	DDInterval() time.Duration
}
