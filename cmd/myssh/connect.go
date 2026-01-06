package main

import (
	"fmt"

	"myssh/internal/profile"

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

		// Construction config SSH
		config, err := buildSSHConfig(user, key, password)
		if err != nil {
			fmt.Println("Erreur config SSH:", err)
			return
		}

		// Connexion SSH
		client, err := connectSSH(host, port, config)
		if err != nil {
			fmt.Println("Erreur connexion SSH:", err)
			return
		}
		defer client.Close()

		// Session interactive
		if err := startInteractiveSession(client); err != nil {
			fmt.Println("Erreur session interactive:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
