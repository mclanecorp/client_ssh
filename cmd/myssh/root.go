package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	sshHost     string
	sshUser     string
	sshPort     int
	sshPassword string
	sshKey      string
)

var rootCmd = &cobra.Command{
	Use:   "myssh",
	Short: "myssh est un client SSH minimaliste pour développeurs",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&sshHost, "host", "", "SSH host")
	rootCmd.PersistentFlags().StringVar(&sshUser, "user", "", "SSH user")
	rootCmd.PersistentFlags().IntVar(&sshPort, "port", 22, "SSH port")
	rootCmd.PersistentFlags().StringVar(&sshPassword, "password", "", "SSH password")
	rootCmd.PersistentFlags().StringVar(&sshKey, "key", "", "SSH private key path")

	// Ajouter les sous-commandes
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(scpCmd)
}
