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
	"github.com/spf13/cobra"
)

var recursive bool

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Back up a file to DRFS",
	Long: `Store files in DRFS. If the file name already exists; backup throws an error (DRFS cannot update 
		files at the moment`,
	Run: func(cmd *cobra.Command, args []string) {
		backup(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	backupCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "-r", false, "recurse directory and copy all files")

}

func recursiveBackup(cmd *cobra.Command, args []string) {

}

func backup(cmd *cobra.Command, args []string) {

}
