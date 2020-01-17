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
	"fmt"
	"os"

	"bitbucket.org/bitgrip/uptrack/cmd/server"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:               "uptrack",
		Short:             "track down your uptime",
		Long:              `uptrack is a service to automaticly check the uptime of your HTTP services`,
		DisableAutoGenTag: true,
	}
	cfgFile  string
	logJSON  bool
	logLevel int
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.uptrack.yaml)")
	rootCmd.PersistentFlags().IntVarP(&logLevel, "verbosity", "v", 0, "verbosity level to use")
	rootCmd.PersistentFlags().BoolVar(&logJSON, "log-json", false, "if to log using json format")

	rootCmd.AddCommand(server.BaseCommand(rancherConfig, initSubCommand))
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(genDocCmd)
	rootCmd.AddCommand(completionCmd)
}

func initSubCommand() {
	if logJSON {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
	}
	if logLevel > 0 {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if logLevel >= 2 {
		logrus.SetLevel(logrus.TraceLevel)
		logrus.SetReportCaller(true)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".uptrack")
	}

	viper.AutomaticEnv()

	viper.ReadInConfig()
}