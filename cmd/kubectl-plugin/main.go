/*
Copyright 2017 The Kubernetes Authors.

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
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:     os.Getenv("KUBECTL_PLUGINS_DESCRIPTOR_NAME"),
		Short:   os.Getenv("KUBECTL_PLUGINS_DESCRIPTOR_SHORT_DESC"),
		Long:    os.Getenv("KUBECTL_PLUGINS_DESCRIPTOR_LONG_DESC"),
		Example: os.Getenv("KUBECTL_PLUGINS_DESCRIPTOR_EXAMPLE"),
		Run: func(cmd *cobra.Command, args []string) {
			for _, env := range os.Environ() {
				fmt.Println(env)
			}
			fmt.Println("SERVICE CATALOG")
		},
	}

	cmd.AddCommand(NewCmdList())
	cmd.AddCommand(NewCmdBind())
	cmd.AddCommand(NewCmdUnbind())

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
