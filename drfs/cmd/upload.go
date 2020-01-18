/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"io"
	"os"

	drfs "github.com/kaiserkarel/drfs/os"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload up a file to DRFS",
	Long: `Upload files in DRFS. If the file name already exists; backup throws an error (DRFS cannot update 
		files at the moment`,
	Run: func(cmd *cobra.Command, args []string) {
		backup(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func backup(cmd *cobra.Command, args []string) {
	var fileName = args[0]
	src, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("cannot open local file: %s", err)
		os.Exit(1)
	}

	fmt.Println("creating drfs file")
	dst, err := drfs.Open(fileName)
	if err != nil {
		fmt.Printf("cannot open drfs file: %s", err)
		os.Exit(1)
	}

	fmt.Println("starting transfer")
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Printf("cannot copy %s to drfs: %s", fileName, err)
		os.Exit(1)
	}
	fmt.Println("upload complete")
}
