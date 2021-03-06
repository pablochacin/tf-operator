/*
Copyright © 2020 Pablo Chacin <pablochacin@gmail.com>

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
package main

import (
	"fmt"

	"github.com/pablochacin/tf-operator/pkg/client"
	"github.com/spf13/cobra"
)

var (
	requiredArgs = []string{"stack"}
)

func newCreateCmd() *cobra.Command {

	var kubeconfig string

	opts := &createOpts{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a terraform operator stack",
		Long: `Create a terraform operator stack from a terrafrom configuration
and a tfvars file. The terraform configuration is obtained from a local
directory.`,
		Example: `
# Create stack from working directory. All .tf files will be used as config
# and the terraform.tfvars file will be used to provider the input variables
tfoctl -s MyStack`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.validateArgs(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.NewFromKubeconfig(kubeconfig)
			if err != nil {
				return err
			}
			opts.client = client
			return opts.run()
		},
	}

	cmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "path to kubeconfig for cluster. If not specified, default discovery rules will apply")
	cmd.Flags().StringVarP(&opts.stack, "stack", "s", "", "stack name")
	cmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "default", "namespace for stack")
	cmd.Flags().StringVarP(&opts.configDir, "config", "c", "./", "path to the terraform configuration directory. All .tf files will be used as the stack configuration. Default is current directory")
	cmd.Flags().StringVarP(&opts.tfvars, "vars", "v", "terraform.tfvars", "Path toterraform vars file.")

	return cmd
}

// validateArgs validates the arguments
func (opts *createOpts) validateArgs(cmd *cobra.Command) error {
	// check for required arguments
	for _, arg := range requiredArgs {
		if !cmd.Flags().Lookup(arg).Changed {
			return fmt.Errorf("argument %s must be specified", arg)
		}
	}

	return nil
}
