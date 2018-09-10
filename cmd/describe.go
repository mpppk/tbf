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

	"os"

	"encoding/json"

	"github.com/mpppk/tbf/csv"
	"github.com/mpppk/tbf/tbf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "サークル情報を表示します",
	Long: `引数として与えられたスペース名のサークル情報をjsonで表示します
ex)
$ tbf describe あ01
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		spaces := args

		source := tbf.NewSource(viper.GetString("source"))
		circleCSV, err := csv.NewCircleCSV(source.FileName)
		if err != nil {
			panic(err)
		}

		circleDetailMap, err := circleCSV.ToCircleDetailMap()
		if err != nil {
			panic(err)
		}

		for _, space := range spaces {
			circleDetail, ok := circleDetailMap[space]
			if !ok {
				fmt.Fprintf(os.Stderr, "circle on %s not found\n", space)
				continue
			}

			marshaledCircleDetail, err := json.Marshal(circleDetail)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to convert to json from circle detail(%s): %s\n", circleDetail, err)
				continue
			}
			fmt.Println(string(marshaledCircleDetail))
		}
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)

	describeCmd.Flags().StringP("source", "s", "latest", "表示するサークル情報のソース(ファイルパスorURLorエイリアス)")
	viper.BindPFlag("source", listCmd.Flags().Lookup("source"))
}
