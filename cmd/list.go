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

	"github.com/mpppk/tbf/csv"
	"github.com/mpppk/tbf/tbf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "与えられたソースのサークル情報を表示",
	Long: `1行に１サークルの情報を以下のフォーマットで表示します。 
[スペース名] [サークル名] by [ペンネーム]【[ジャンル名]】 : [頒布物説明]
ソースにはローカルファイル, URL, エイリアスが使用可能です。
ローカルファイルの例: ./circles.csv
URLの例: https://raw.githubusercontent.com/mpppk/tbf/master/data/latest_circles.csv

次のエイリアスが利用可能です。
latest(最新の技術書典サークル情報) → https://raw.githubusercontent.com/mpppk/tbf/master/data/latest_circles.csv
tbf4(技術書典4サークル情報) → https://raw.githubusercontent.com/mpppk/tbf/master/data/tbf4_circles.csv
`,
	Run: func(cmd *cobra.Command, args []string) {
		source := tbf.NewSource(viper.GetString("source"))
		csvFilePath := source.FileName

		if source.Url != "" {
			csvURL, ok := tbf.GetCSVURL(source.Url)
			if !ok {
				fmt.Fprintln(os.Stderr, "unknown CSV URL: "+source.Url)
				os.Exit(1)
			}

			csvMetaURL := strings.Replace(csvURL, ".csv", ".json", 1)

			downloaded, err := csv.DownloadCSVIfChanged(csvURL, csvMetaURL, csvFilePath)
			if err != nil {
				panic(err)
			}

			if downloaded {
				fmt.Fprintf(os.Stderr, "new csv file is downloaded from %s to %s\n", csvURL, csvFilePath)
			}
		} else if !csv.IsExist(source.FileName) {
			fmt.Fprintf(os.Stderr, "csv file not found: %s\n", csvFilePath)
			os.Exit(1)
		}

		circleCSV, err := csv.NewCircleCSV(csvFilePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to load csv from: "+source.Url)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		circleDetailMap, err := circleCSV.ToCircleDetailMap()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to parse csv from: "+source.Url)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
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

	listCmd.Flags().StringP("source", "s", "latest", "表示するサークル情報のソース(ファイルパスorURLorエイリアス)")
	viper.BindPFlag("source", listCmd.Flags().Lookup("source"))
}
