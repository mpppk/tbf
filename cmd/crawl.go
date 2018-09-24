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

	"context"
	"time"

	"os"

	"github.com/mpppk/tbf/crawl"
	"github.com/mpppk/tbf/csv"
	"github.com/mpppk/tbf/tbf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileKey = "file"
var urlKey = "url"
var sleepKey = "sleep"

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "技術書典のウェブサイトをスクレイピングしてcsvとして保存",
	Long: `技術書典のウェブサイトをスクレイピングし、サークル情報を--fileで指定した名前のcsvとして書き込みます。
スクレイピングにはchromeを利用するため、実行する環境にあらかじめインストールしておく必要があります。`,

	Run: func(cmd *cobra.Command, args []string) {
		csvFilePath := viper.GetString(fileKey)
		circlesURL := viper.GetString(urlKey)
		sleep := time.Duration(viper.GetInt(sleepKey)) * time.Second
		crawler, err := crawl.NewTBFCrawler(context.Background(), tbf.BaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wait error: %v", err)
			os.Exit(1)
		}

		circleCSV, err := csv.NewCircleCSV(csvFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to prepare csv: %v", err)
			os.Exit(1)
		}

		circleDetailMap, err := circleCSV.ToCircleDetailMap()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse csv: %v", err)
			os.Exit(1)
		}

		// TODO: Add timeout
		circles, err := crawler.FetchCircles(context.Background(), circlesURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to fetch circle information: %v", err)
			os.Exit(1)
		}

		var filteredCircles []*tbf.Circle
		for _, c := range circles {
			if _, ok := circleDetailMap[c.Space]; !ok {
				filteredCircles = append(filteredCircles, c)
			}
		}

		for i, circle := range filteredCircles {
			fmt.Printf(
				"all: %d, saved: %d, new: %d\n",
				len(circles),
				len(circleDetailMap)+i,
				len(filteredCircles)-i,
			)

			circleDetail, err := crawler.FetchCircleDetail(context.Background(), circle)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to fetch circle detail information: %v", err)
				os.Exit(1)
			}
			fmt.Printf("%#v\n", circleDetail)
			if err := circleCSV.AppendCircleDetail(circleDetail); err != nil {
				panic(err)
			}
			time.Sleep(sleep)
		}

		// shutdown chrome
		if crawler.Shutdown(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "shutdown error: %v", err)
			os.Exit(1)
		}

		// wait for chrome to finish
		if err := crawler.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "wait error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)

	crawlCmd.Flags().StringP(fileKey, "f", "circles.csv", "サークル情報を書き出すcsvファイル名")
	viper.BindPFlag(fileKey, crawlCmd.Flags().Lookup(fileKey))

	crawlCmd.Flags().StringP(urlKey, "u", "https://techbookfest.org/event/tbf05/circle", "サークル情報を取得するURL")
	viper.BindPFlag(urlKey, crawlCmd.Flags().Lookup(urlKey))

	crawlCmd.Flags().Int(sleepKey, 10, "スクレイピングのためにHTTPリクエストを送る際のインターバル(秒)")
	viper.BindPFlag(sleepKey, crawlCmd.Flags().Lookup(sleepKey))
}
