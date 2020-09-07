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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fanzy618/easycert/pkg/cert"
)

var cfgFile string
var globalCfg = &cert.Config{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "easycert",
	Short: "Help you create certificates easily",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var ipStrSlic []string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&globalCfg.Name, "name", "n", "ca", "CA file name")
	rootCmd.PersistentFlags().StringVarP(&globalCfg.Dir, "dir", "d", ".", "Dictory where save ca files")

	rootCmd.PersistentFlags().StringVar(&globalCfg.CommonName, "cn", "easycert", "Common Name")
	rootCmd.PersistentFlags().StringSliceVar(&globalCfg.Organization, "orgs", []string{"easycert"}, "Organizations")
	rootCmd.PersistentFlags().StringSliceVar(&globalCfg.AltNames.DNSNames, "dns", []string{}, "DNS")
	rootCmd.PersistentFlags().IPSliceVar(&globalCfg.AltNames.IPs, "ip", nil, "IPs")

	rootCmd.PersistentFlags().IntVar(&globalCfg.Bits, "b", 2048, "Bits of RSA key")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.easycert.yaml)")

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
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".easycert" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".easycert")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
