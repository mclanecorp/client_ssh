package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Se connecter en SSH et lancer un shell interactif",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := buildSSHConfig(sshUser, sshKey, sshPassword)
		if err != nil {
			fmt.Println("Erreur config SSH:", err)
			return
		}

		client, err := connectSSH(sshHost, sshPort, config)
		if err != nil {
			fmt.Println("Erreur connexion SSH:", err)
			return
		}
		defer client.Close()

		if err := startInteractiveSession(client); err != nil {
			fmt.Println("Erreur session interactive:", err)
		}
	},
}
