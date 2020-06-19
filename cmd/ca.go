/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"crypto/x509"

	"github.com/spf13/cobra"

	certutil "k8s.io/client-go/util/cert"

	"github.com/fanzy618/easycert/pkg/cert"
)

// caCmd represents the ca command
var caCmd = &cobra.Command{
	Use:   "ca",
	Short: "Generate a new CA certificate",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("ca called")
	// 	_ = &certutil.Config{}
	// },
	RunE: func(cmd *cobra.Command, args []string) error {
		return createCA(globalCfg)
	},
}

func init() {
	rootCmd.AddCommand(caCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// caCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// caCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createCA(cfg *cert.Config) error {
	key, err := cert.NewKey("", cfg.Bits)
	if err != nil {
		return err
	}

	ca, err := certutil.NewSelfSignedCACert(
		certutil.Config{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		key,
	)
	if err != nil {
		return err
	}
	return cert.WriteToDisk(cfg.Dir, cfg.Name, []*x509.Certificate{ca}, key)
}
