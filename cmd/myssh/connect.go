package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	host     string
	user     string
	port     int
	keyPath  string
	password string
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a remote SSH server",
	Long:  "Establish a connection to a remote SSH server. Currently a placeholder.",
	Run: func(cmd *cobra.Command, args []string) {
		if keyPath == "" {
			fmt.Println("No key provided, nothing to do for now.")
			return
		}

		config, err := buildSSHConfig(user, keyPath)
		if err != nil {
			fmt.Println("Error building SSH config:", err)
			return
		}

		fmt.Println("SSH config built successfully for user:", config.User)
	},
}

func init() {
	connectCmd.Flags().StringVar(&host, "host", "", "SSH server host (required)")
	connectCmd.Flags().StringVar(&user, "user", "", "SSH username (required)")
	connectCmd.Flags().IntVar(&port, "port", 22, "SSH port (default 22)")
	connectCmd.Flags().StringVar(&keyPath, "key", "", "Path to private key")
	connectCmd.Flags().StringVar(&password, "password", "", "Password (not recommended)")

	connectCmd.MarkFlagRequired("host")
	connectCmd.MarkFlagRequired("user")

	rootCmd.AddCommand(connectCmd)
}
