/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	Backend_name   string
	cfgFile        string
	Server_ip      string
	Ssh_key_path   string
	Ssh_key_string string
	Base64ssh_key  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gogamic-ci-cli",
	Short: "gogamic-cli is tool which is used to deploy to dokku from a CI/CD server",
	Long:  `Use this tool for deploying your apps to dokku`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&Backend_name, "backend_type", "t", "dokku", "The type of backend default:dokku")
	rootCmd.PersistentFlags().StringVarP(&Server_ip, "server_ip", "i", "", "The Server IP addr where the app should be deployed")
	rootCmd.PersistentFlags().StringVarP(&Ssh_key_path, "ssh_key_file", "f", "", "Enter the path to SSH Private Key (default is $HOME/.ssh/id_rsa)")
	rootCmd.PersistentFlags().StringVarP(&Ssh_key_string, "ssh_key_string", "s", "null", "The Base64 encoded SSH Private Key string")
	rootCmd.PersistentFlags().BoolVarP(&Base64ssh_key, "base64_key", "b", false, "Weather the private key file is encoded in BASE64")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gogamic-ci-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gogamic-ci-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gogamic-ci-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
