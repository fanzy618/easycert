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
	"crypto"
	cryptorand "crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/spf13/cobra"

	"github.com/fanzy618/easycert/pkg/cert"
)

// certCmd represents the cert command
var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Generate a new certificate",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("cert called")
	// },
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cert called")
		if caName != "" {
			return createFromCA(globalCfg, caName)
		}
		return createFromCA(globalCfg, "")
	},
}

var caName string

func init() {
	rootCmd.AddCommand(certCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// certCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// certCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	certCmd.Flags().StringVar(&caName, "ca", "", "Name of ca.")
}

func createFromCA(cfg *cert.Config, caName string) error {
	var caCert *x509.Certificate
	var caKey crypto.Signer
	var err error
	if caName != "" {
		caCert, caKey, err = cert.LoadFromDisk(cfg.Dir, caName)
		if err != nil {
			return err
		}
	} else {
		caKey, err = cert.NewKey("", cfg.Bits)
		if err == nil {
			caCert, err = cert.NewSelfSignedCACert(cfg, caKey)
		}
		if err != nil {
			return err
		}
	}
	certificate, key, err := newCertAndKey(caCert, caKey, cfg)
	if err != nil {
		return err
	}
	return cert.WriteToDisk(cfg.Dir, cfg.Name, []*x509.Certificate{certificate, caCert}, key)
}

func newCertAndKey(caCert *x509.Certificate, caKey crypto.Signer, cfg *cert.Config) (*x509.Certificate, crypto.Signer, error) {
	key, err := cert.NewKey("", cfg.Bits)
	if err != nil {
		return nil, nil, err
	}

	serial, err := cryptorand.Int(cryptorand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, nil, err
	}
	if len(cfg.CommonName) == 0 {
		return nil, nil, errors.New("must specify a CommonName")
	}

	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		DNSNames:     cfg.AltNames.DNSNames,
		IPAddresses:  cfg.AltNames.IPs,
		SerialNumber: serial,
		NotBefore:    caCert.NotBefore,
		NotAfter:     time.Now().Add(cfg.ValidFor).UTC(),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	certDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, nil, err
	}
	certificate, err := x509.ParseCertificate(certDERBytes)
	if err != nil {
		return nil, nil, err
	}
	return certificate, key, nil
}
