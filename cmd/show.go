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
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fanzy618/easycert/pkg/cert"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show certificate information",
	Long:  "Load a certificate from disk and print its details. The certificate location is determined by --name and --dir flags.",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, _, err := cert.LoadFromDisk(globalCfg.Dir, globalCfg.Name)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Subject: %s\n", c.Subject.String())
		fmt.Fprintf(cmd.OutOrStdout(), "DNSNames: %s\n", strings.Join(c.DNSNames, ", "))
		var ips []string
		for _, ip := range c.IPAddresses {
			ips = append(ips, ip.String())
		}
		fmt.Fprintf(cmd.OutOrStdout(), "IPAddresses: %s\n", strings.Join(ips, ", "))
		fmt.Fprintf(cmd.OutOrStdout(), "NotBefore: %s\n", c.NotBefore.Format("2006-01-02T15:04:05Z07:00"))
		fmt.Fprintf(cmd.OutOrStdout(), "NotAfter: %s\n", c.NotAfter.Format("2006-01-02T15:04:05Z07:00"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
