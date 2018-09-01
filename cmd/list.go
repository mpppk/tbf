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

	"github.com/mpppk/tbf/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list circle information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		csvFilePath := viper.GetString("file")
		circleCSV, err := util.NewCircleCSV(csvFilePath)
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

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "file", "", "circle csv file")

	rootCmd.PersistentFlags().String("file", "circles.csv", "circle csv file")
	viper.BindPFlag("file", rootCmd.PersistentFlags().Lookup("file"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}