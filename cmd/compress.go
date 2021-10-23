/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:     "compress ",
	Example: "  melt compress --zip FILE",
	Short:   "Compresses a file in the specified format.",
	Long: `
Compresses the file in the format specified by the flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("compress called")
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// compressCmd.PersistentFlags().String("foo", "", "A help for foo")
	// compressCmd.PersistentFlags().Int("fooo", 0, "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	compressCmd.Flags().BoolP("zip", "z", false, "Compress in zip format")
}
