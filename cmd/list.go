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

	"strings"

	"os"

	"path"

	"github.com/mpppk/tbf/csv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type listConfig struct {
	url      string
	fileName string
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list circle information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config := createConfigFromSource(viper.GetString("source"))
		csvFilePath := config.fileName

		csvURL, ok := getCSVURL(config.url)
		if !ok {
			fmt.Fprintln(os.Stderr, "unknown CSV URL: "+config.url)
			os.Exit(1)
		}

		csvMetaURL := strings.Replace(csvURL, ".csv", ".json", 1)

		_, err := csv.DownloadCSVIfChanged(csvURL, csvMetaURL, csvFilePath)
		if err != nil {
			panic(err)
		}

		circleCSV, err := csv.NewCircleCSV(csvFilePath)
		if err != nil {
			panic(err)
		}

		circleDetailMap, err := circleCSV.ToCircleDetailMap()
		if err != nil {
			panic(err)
		}

		for space, circleDetail := range circleDetailMap {
			fmt.Printf("%s %s by %s 【%s】 : %s\n",
				space,
				circleDetail.Name,
				circleDetail.Penname,
				circleDetail.Genre,
				circleDetail.GenreFreeFormat,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	rootCmd.PersistentFlags().String("source", "circles.csv", "circle data source")
	viper.BindPFlag("source", rootCmd.PersistentFlags().Lookup("source"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createConfigFromSource(source string) *listConfig {
	url, ok := getCSVURL(source)
	fileName := path.Base(source)
	if ok {
		fileName = path.Base(url)
	}

	return &listConfig{
		url:      url,
		fileName: fileName,
	}
}

func getCSVURL(source string) (string, bool) {
	if strings.Contains(source, "http") {
		return source, true
	}
	u, ok := csv.URLMap[source]
	if ok {
		return u, true
	}
	return "", false
}
