// Copyright © 2018 mpppk <niboshiporipori@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"

	"encoding/json"
	"io/ioutil"

	"context"

	"github.com/mitchellh/go-homedir"
	"github.com/mpppk/tbf/crawl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "tbf",
	Short: "CLI for tech book festival",
	Long:  `CLI for tech book festival`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		crawler, err := crawl.NewTBFCrawler(context.Background())
		if err != nil {
			panic(err)
		}

		circles, err := crawler.FetchCircles(context.Background())
		if err != nil {
			panic(err)
		}

		fmt.Println("circles")
		for _, circle := range circles {
			fmt.Printf("%#v\n", circle)
		}

		encodedCircles, err := json.Marshal(circles)
		if err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile("circles.json", encodedCircles, 0755); err != nil {
			panic(err)
		}

		// shutdown chrome
		if crawler.Shutdown(context.Background()); err != nil {
			log.Fatal("shutdown error:", err)
		}

		// wait for chrome to finish
		if err := crawler.Wait(); err != nil {
			log.Fatal("wait error:", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tbf.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".tbf" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tbf")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
