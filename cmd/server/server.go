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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Interact to the server",
		Long:  serverLongDescription,
	}
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the server",
		Long:  startLongDescription,
		Run:   start,
	}
	showFile    string
	valuesFiles []string
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
	serverCmd.PersistentFlags().String("config-file", ".", "Path of configuration file")
	serverCmd.PersistentFlags().String("jobs-dir", "/etc/uptrack/jobs.d", "Directory to find job descriptors")
	serverCmd.PersistentFlags().Int("default-interval", 10, "Default interval to execute job")
	serverCmd.PersistentFlags().String("datadog-credentials", "/etc/uptrack/datadog/credentials", "File containing datadog credentials")
	serverCmd.PersistentFlags().String("prometheus-port", "9001", "Port exposed by prometheus")
	serverCmd.PersistentFlags().String("prometheus-endpoint", "metrics", "Prometheus Endpoint")

	//Server Configuration
	viper.BindPFlag("uptrack.config_file", serverCmd.PersistentFlags().Lookup("config-file"))
	viper.BindEnv("uptrack.config_file", "UPTRACK_CONFIG_FILE")

	viper.BindPFlag("uptrack.jobs_dir", serverCmd.PersistentFlags().Lookup("jobs-dir"))
	viper.BindEnv("uptrack.jobs_dir", "UPTRACK_JOBS_DIR")

	viper.BindPFlag("uptrack.default_interval", serverCmd.PersistentFlags().Lookup("default-interval"))
	viper.BindEnv("uptrack.default_interval", "UPTRACK_DEFAULT_INTERVAL")

	viper.BindPFlag("uptrack.datadog.credentials", serverCmd.PersistentFlags().Lookup("datadog-credentials"))
	viper.BindEnv("uptrack.datadog.credentials", "UPTRACK_DATADOG_CREDENTIALS")

	viper.BindPFlag("uptrack.prometheus.port", serverCmd.PersistentFlags().Lookup("prometheus-port"))
	viper.BindEnv("uptrack.prometheus.port", "UPTRACK_PROMETHEUS_PORT")

	viper.BindPFlag("uptrack.prometheus.endpoint", serverCmd.PersistentFlags().Lookup("prometheus-endpoint"))
	viper.BindEnv("uptrack.prometheus.endpoint", "UPTRACK_PROMETHEUS_ENDPOINT")

	viper.SetConfigFile(viper.GetString("uptrack.config_file"))

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Panic(fmt.Sprintf("Fatal error config file: %s \n", err))
	}
	serverCmd.AddCommand(startCmd)
}
