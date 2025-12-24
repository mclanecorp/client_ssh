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
	Long:  "Establish a connection to a remote SSH server and start an interactive session.",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := buildSSHConfig(user, keyPath, password)
		if err != nil {
			fmt.Println("Erreur configuration SSH:", err)
			return
		}

		client, err := connectSSH(host, port, config)
		if err != nil {
			fmt.Println("Erreur connexion SSH:", err)
			return
		}
		defer client.Close()

		fmt.Println("Connexion SSH établie avec succès sur", host)

		if err := startInteractiveSession(client); err != nil {
			fmt.Println("Erreur session interactive:", err)
		}
	},
}

func init() {
	connectCmd.Flags().StringVar(&host, "host", "", "SSH server host (required)")
	connectCmd.Flags().StringVar(&user, "user", "", "SSH username (required)")
	connectCmd.Flags().IntVar(&port, "port", 22, "SSH port (default 22)")
	connectCmd.Flags().StringVar(&keyPath, "key", "", "Path to private key")
	connectCmd.Flags().StringVar(&password, "password", "", "Password (if no key)")

	connectCmd.MarkFlagRequired("host")
	connectCmd.MarkFlagRequired("user")

	rootCmd.AddCommand(connectCmd)
}
