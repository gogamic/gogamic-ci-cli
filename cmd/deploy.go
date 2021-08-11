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

	"github.com/gogamic/gogamic-ci-cli/functions"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

type ISshClient interface {
	GetConnection(host string, username string, password string) ssh.Client
}

var ssh_key_path string
var config_path string
var base64ssh_key bool

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   deploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployCmd.Flags().StringVarP(&config_path, "config", "c", "", "Please enter the path to YAML file")
	deployCmd.Flags().StringVarP(&ssh_key_path, "ssh_key", "f", "", "Enter the path to SSH Private Key (default is $HOME/.ssh/id_rsa) ")
	deployCmd.Flags().BoolVarP(&base64ssh_key, "base64_key", "b", false, "Weather the private key file is encoded in BASE64")
}

func deploy(cmd *cobra.Command, args []string) {
	data, err := functions.ParseYAMLFile(config_path)
	if err != nil {
		panic(fmt.Sprintf("unable to parse file: %s", err.Error()))
	}

	cmds, err := functions.GetCommands(data)

	if err != nil {
		panic(fmt.Sprintf("invalid backend: %s", data.Config.Backend))
	}
	if ssh_key_path == "" {
		err = functions.RunServerCommands(data.Config.IP, cmds, nil, base64ssh_key)
	} else {
		err = functions.RunServerCommands(data.Config.IP, cmds, &ssh_key_path, base64ssh_key)
	}

	if err != nil {
		panic(err.Error())
	}

}

/*
func checkArgs(cmd *cobra.Command, args []string) error {
	err := functions.ValidateYAMLFile(config_path)
	if err != nil {
		functions.HandleErr(err, "config file err")
	}
	return nil
}
*/
