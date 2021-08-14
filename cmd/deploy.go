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

var (
	app_name  string
	image_url string
	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "A brief description of your command",
		Long:  ``,
		Run:   deploy,
		Args:  checkArgs,
	}
)

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployCmd.Flags().StringVarP(&app_name, "name", "n", "", "The application name")
	deployCmd.Flags().StringVarP(&image_url, "image", "r", "", "The url to docker image")
}

func deploy(cmd *cobra.Command, args []string) {

	cmds, err := functions.GetCommands(Backend_name, image_url)

	if err != nil {
		functions.HandleErr(err, fmt.Sprintf("invalid backend: %s", Backend_name))
	}

	if Ssh_key_path != "" {
		functions.RunServerCommands(Server_ip, cmds, &Ssh_key_path, Base64ssh_key, Ssh_key_string)
		return
	}
	functions.RunServerCommands(Server_ip, cmds, nil, Base64ssh_key, Ssh_key_string)

}

func checkArgs(cmd *cobra.Command, args []string) error {
	err := functions.CheckIPAddress(Server_ip)
	if err != nil {
		functions.HandleErr(err, "You have entered an Invalid ip")
	}

	if Ssh_key_string == "none" && Ssh_key_path == "" {
		functions.HandleErr(nil, "Please enter ssh string or ssh path")
	}
	return nil
}
