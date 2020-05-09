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

)

type createOpts struct {
	kubeconfig string
	stack      string
	namespace  string
	configDir  string
	configMap  string
	tfvars     string
	tfstate    string
}

// run executes the create stack command
func (opts *createOpts) run() error {
	fmt.Printf("create stack %s in namepace %s from config in %s with vars in %s and state %s",
		opts.stack,
		opts.namespace,
		opts.configDir,
		opts.tfvars,
		opts.tfstate)

	return nil
}
