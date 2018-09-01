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
	"log"
	"time"

	"github.com/mpppk/tbf/crawl"
	"github.com/mpppk/tbf/util"
	"github.com/spf13/cobra"
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		csvFilePath := "circles.csv"
		crawler, err := crawl.NewTBFCrawler(context.Background())
		if err != nil {
			panic(err)
		}

		circleCSV, err := util.NewCircleCSV(csvFilePath)
		if err != nil {
			panic(err)
		}

		circleDetailMap, err := circleCSV.ToCircleDetailMap()
		if err != nil {
			panic(err)
		}

		circles, err := crawler.FetchCircles(context.Background())
		if err != nil {
			panic(err)
		}

		var filteredCircles []*crawl.Circle
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
				panic(err)
			}
			fmt.Printf("%#v\n", circleDetail)
			if err := circleCSV.AppendCircleDetail(circleDetail); err != nil {
				panic(err)
			}
			time.Sleep(10 * time.Second)
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

func init() {
	rootCmd.AddCommand(crawlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crawlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crawlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
