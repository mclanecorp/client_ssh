package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var scpCmd = &cobra.Command{
	Use:   "scp",
	Short: "Secure Copy over SSH",
	Long:  "Transfer files to or from a remote host using SCP over SSH.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SCP command. Use upload or download.")
	},
}
var scpUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to a remote server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SCP upload not implemented yet")
	},
}
var scpDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from a remote server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SCP download not implemented yet")
	},
}

func init() {
	scpCmd.AddCommand(scpUploadCmd)
	scpCmd.AddCommand(scpDownloadCmd)

	rootCmd.AddCommand(scpCmd)
}
