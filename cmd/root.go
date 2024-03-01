/*
Copyright © 2023 Agustin Larreinegabe

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var (
	URLGRAPHQL       string
	URLREST          string
	organizationList []string
	aacFormatType    string
	aacFilePath      string
	endpoint         string
	org              string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-aac",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gh-aac.yaml)")

	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "GitHub Enterprise endpoint. By default is https://github.com.")
	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))

	rootCmd.PersistentFlags().StringVarP(&org, "organization", "o", "", "Slug organization name. By default from conf file.")
	viper.BindPFlag("organization", rootCmd.PersistentFlags().Lookup("organization"))

	rootCmd.PersistentFlags().StringVarP(&aacFormatType, "aac-format", "f", "", "It defines the format of the access as code (aac) output file. By default is YAML. It can be YAML or JSON.")
	viper.BindPFlag("aac-format", rootCmd.PersistentFlags().Lookup("aac-format"))
	viper.SetDefault("aac-format", "yaml")

	rootCmd.PersistentFlags().StringVarP(&aacFilePath, "aac-path", "p", "", "Path to the output YAML file")
	viper.BindPFlag("aac-path", rootCmd.PersistentFlags().Lookup("aac-path"))
	viper.SetDefault("aac-path", "access-config.yaml")

	rootCmd.MarkFlagFilename("aac-path", "yaml")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gh-aac" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gh-aac")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if viper.GetString("endpoint") == "" {
		// By default http://github.com
		viper.Set("endpoint", "https://github.com")

		URLGRAPHQL = "https://api.github.com/graphql"
		URLREST = "https://api.github.com"
	} else {
		URLREST = fmt.Sprintf("%s/api/v3", viper.GetString("endpoint"))
		URLGRAPHQL = fmt.Sprintf("https://api.%s/graphql", strings.Replace(viper.GetString("endpoint"), "https://", "", 1))
	}

	if viper.GetString("organization") != "" {
		organizationList = append(organizationList, viper.GetString("organization"))
	} else {
		organizationList = viper.GetStringSlice("organizations")
	}

	// print endpoint
	log.Printf("Endpoint: %s", viper.GetString("endpoint"))
}
