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
	"strconv"

	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func NewCmdList() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			insecure, _ := strconv.ParseBool(os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CLUSTER_INSECURE_SKIP_TLS_VERIFY"))
			config, err := clientcmd.BuildConfigFromKubeconfigGetter(os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CLUSTER_SERVER"), func() (*clientcmdapi.Config, error) {
				return &clientcmdapi.Config{
					CurrentContext: os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CURRENT_CONTEXT"),
					Clusters: map[string]*clientcmdapi.Cluster{
						os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CLUSTER_NAME"): &clientcmdapi.Cluster{
							Server:                   os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CLUSTER_SERVER"),
							APIVersion:               os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CLUSTER_API_VERSION"),
							InsecureSkipTLSVerify:    insecure,
							CertificateAuthority:     os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CLUSTER_CERTIFICATE_AUTHORITY"),
							CertificateAuthorityData: []byte(os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CLUSTER_CERTIFICATE_AUTHORITY_DATA")),
						},
					},
					AuthInfos: map[string]*clientcmdapi.AuthInfo{
						os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_AUTH_INFO_NAME"): &clientcmdapi.AuthInfo{
							ClientCertificate:     os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_AUTH_INFO_CLIENT_CERTIFICATE"),
							ClientCertificateData: []byte(os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_AUTH_INFO_CLIENT_CERTIFICATE_DATA")),
							ClientKey:             os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_AUTH_INFO_CLIENT_KEY"),
							ClientKeyData:         []byte(os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_AUTH_INFO_CLIENT_KEY_DATA")),
							Token:                 os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_AUTH_INFO_TOKEN"),
							TokenFile:             os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_AUTH_INFO_TOKEN_FILE"),
							Impersonate:           os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_AUTH_INFO_IMPERSONATE"),
							Username:              os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_AUTH_INFO_USERNAME"),
							Password:              os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_AUTH_INFO_PASSWORD"),
						},
					},
					Contexts: map[string]*clientcmdapi.Context{
						os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CONTEXT_NAME"): &clientcmdapi.Context{
							Cluster:   os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CONTEXT_CLUSTER"),
							AuthInfo:  os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CONTEXT_AUTH_INFO"),
							Namespace: os.Getenv("KUBECTL_PLUGINS_CLIENT_CONFIG_CONTEXT_NAMESPACE"),
						},
					},
				}, nil
			})
			if err != nil {
				panic(err.Error())
			}

			cs, err := clientset.NewForConfig(config)
			if err != nil {
				panic(err.Error())
			}

			serviceclasses, err := cs.ServicecatalogV1alpha1().ServiceClasses().List(metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}

			for _, serviceclass := range serviceclasses.Items {
				fmt.Printf("Name: %s\n", serviceclass.Name)
				fmt.Printf("From Broker: %s\n", serviceclass.BrokerName)
				fmt.Printf("Bindable? %v\n", serviceclass.Bindable)
				fmt.Printf("Description: %v\n", *serviceclass.Description)
			}
		},
	}
}
