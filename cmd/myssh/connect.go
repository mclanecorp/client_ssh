package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// connectCmd représente la commande "connect"
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a remote SSH server",
	Long:  "Establish a connection to a remote SSH server. Currently a placeholder.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("connect command executed. SSH logic will go here.")
	},
}

func init() {
	// On ajoute la commande connect à la commande racine rootCmd
	rootCmd.AddCommand(connectCmd)
}
