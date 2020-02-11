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

package server

import (
	"bitbucket.org/bitgrip/uptrack/internal/pkg/config"
	"bitbucket.org/bitgrip/uptrack/internal/pkg/ctl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Interact with the uptrack server",
		Long:  serverLongDescription,
	}
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the uptrack server",
		Long:  startLongDescription,
		Run:   start,
	}
	rootConfig  config.Config
	initCommand func()
)

// BaseCommand initialises the server base command
func BaseCommand(config config.Config, init func()) *cobra.Command {
	rootConfig = config
	initCommand = init
	return serverCmd
}

func start(cmd *cobra.Command, args []string) {
	initCommand()

	err := ctl.StartUpTrackServer(rootConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	serverCmd.PersistentFlags().String("jobs-config", "./config/jobs.yaml", "Descriptor file defining all checks")
	serverCmd.PersistentFlags().Int("default-interval", 10, "Default interval to execute job")

	serverCmd.PersistentFlags().String("prometheus-port", "9001", "Port exposed by prometheus")
	serverCmd.PersistentFlags().String("prometheus-endpoint", "/metrics", "Prometheus Endpoint")

	serverCmd.PersistentFlags().String("datadog-endpoint", "https://api.datadoghq.com/api/v1", "Datadog Endpoint")
	serverCmd.PersistentFlags().String("datadog-appKey", "", "Datadog APP-key")
	serverCmd.PersistentFlags().String("datadog-apiKey", "", "Datadog API-Key")

	serverCmd.PersistentFlags().Int("datadog-interval", 5, "Interval for sending metrics to Datadog")

	viper.BindPFlag("uptrack.jobs_config", serverCmd.PersistentFlags().Lookup("jobs-config"))
	viper.BindEnv("uptrack.jobs_config", "UPTRACK_JOBS_CONFIG")

	viper.BindPFlag("uptrack.default_interval", serverCmd.PersistentFlags().Lookup("default-interval"))
	viper.BindEnv("uptrack.default_interval", "UPTRACK_DEFAULT_INTERVAL")

	viper.BindPFlag("uptrack.prometheus.port", serverCmd.PersistentFlags().Lookup("prometheus-port"))
	viper.BindEnv("uptrack.prometheus.port", "UPTRACK_PROMETHEUS_PORT")

	viper.BindPFlag("uptrack.prometheus.endpoint", serverCmd.PersistentFlags().Lookup("prometheus-endpoint"))
	viper.BindEnv("uptrack.prometheus.endpoint", "UPTRACK_PROMETHEUS_ENDPOINT")

	viper.BindPFlag("uptrack.datadog.endpoint", serverCmd.PersistentFlags().Lookup("datadog-endpoint"))
	viper.BindEnv("uptrack.datadog.endpoint", "UPTRACK_DATADOG_ENDPOINT")

	viper.BindPFlag("uptrack.datadog.appKey", serverCmd.PersistentFlags().Lookup("datadog-appKey"))
	viper.BindEnv("uptrack.datadog.appKey", "UPTRACK_DATADOG_APPKEY")

	viper.BindPFlag("uptrack.datadog.apiKey", serverCmd.PersistentFlags().Lookup("datadog-apiKey"))
	viper.BindEnv("uptrack.datadog.apiKey", "UPTRACK_DATADOG_APIKEY")

	viper.BindPFlag("uptrack.datadog.interval", serverCmd.PersistentFlags().Lookup("datadog-interval"))
	viper.BindEnv("uptrack.datadog.interval", "UPTRACK_DATADOG_INTERVAL")

	serverCmd.AddCommand(startCmd)
}
