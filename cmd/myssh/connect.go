package main

import (
	"fmt"

	"myssh/internal/profile"
	"myssh/internal/ssh"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Se connecter en SSH et lancer un shell interactif",
	Run: func(cmd *cobra.Command, args []string) {

		// Résolution SSH via profil ou flags
		host, port, user, password, key, err := profile.ResolveSSHConfig(
			profileName,
			sshHost,
			sshPort,
			sshUser,
			sshPassword,
			sshKey,
		)
		if err != nil {
			fmt.Println("Erreur profil:", err)
			return
		}

		// Configuration SSH
		cfg := ssh.Config{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			KeyPath:  key,
		}

		// Connexion SSH
		client, err := ssh.Connect(cfg)
		if err != nil {
			fmt.Println("Erreur connexion SSH:", err)
			return
		}
		defer client.Close()

		// Session interactive
		if err := client.StartInteractiveSession(); err != nil {
			fmt.Println("Erreur session interactive:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
