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
	"context"
	"fmt"
	"log"
	"os"

	"encoding/json"
	"io/ioutil"

	"github.com/chromedp/chromedp"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

type Circle struct {
	Href    string
	Space   string
	Name    string
	Penname string
	Genre   string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tbf",
	Short: "CLI for tech book festival",
	Long:  `CLI for tech book festival`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// create context
		ctxt, cancel := context.WithCancel(context.Background())
		defer cancel()

		// create chrome instance
		c, err := chromedp.New(ctxt)
		//c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
		if err != nil {
			log.Fatal("chromedep new error:", err)
		}

		var res []Circle
		err = c.Run(ctxt, chromedp.Tasks{
			chromedp.Navigate(`https://techbookfest.org/event/tbf04/circle`),
			chromedp.WaitVisible(`li.circle-list-item`),
			chromedp.Evaluate(
				`Array.from(document.querySelectorAll('li.circle-list-item')).map((l) => ({href: l.querySelector('a.circle-list-item-link').getAttribute('href'), space: l.querySelector('span.circle-space-label').textContent, name: l.querySelector('span.circle-name').textContent, penname: l.querySelector('p.circle-list-item-penname').textContent, genre: l.querySelector('p.circle-list-item-genre').textContent}))`,
				&res,
			),
		})

		fmt.Println("res")
		for _, r := range res {
			fmt.Printf("%#v\n", r)
		}

		encodedCircles, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile("circles.json", encodedCircles, 0755); err != nil {
			panic(err)
		}

		if err != nil {
			fmt.Println("click error")
			log.Fatal("failed to click:", err)
		}
		fmt.Println("shutdown start")
		// shutdown chrome
		err = c.Shutdown(ctxt)
		if err != nil {
			log.Fatal("shutdown error:", err)
		}

		// wait for chrome to finish
		err = c.Wait()
		if err != nil {
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
