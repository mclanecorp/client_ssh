package main

import (
	"fmt"

	"myssh/internal/profile"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage SSH profiles",
	Long:  "Create, list and delete SSH connection profiles.",
}

var (
	profileName     string
	profileHost     string
	profileUser     string
	profilePort     int
	profilePassword string
	profileKeyPath  string
)

var profileAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new SSH profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileName = args[0]

		p := profile.Profile{
			Name:     profileName,
			Host:     profileHost,
			Port:     profilePort,
			User:     profileUser,
			Password: profilePassword,
			KeyPath:  profileKeyPath,
		}

		if err := profile.Create(p); err != nil {
			fmt.Println("Erreur création profil:", err)
			return
		}

		fmt.Println("Profil créé avec succès:", profileName)
	},
}

func init() {
	// Commande profile
	rootCmd.AddCommand(profileCmd)

	// Sous-commande add
	profileCmd.AddCommand(profileAddCmd)

	profileAddCmd.Flags().StringVar(&profileHost, "host", "", "SSH host")
	profileAddCmd.Flags().StringVar(&profileUser, "user", "", "SSH user")
	profileAddCmd.Flags().IntVar(&profilePort, "port", 22, "SSH port")
	profileAddCmd.Flags().StringVar(&profilePassword, "password", "", "SSH password")
	profileAddCmd.Flags().StringVar(&profileKeyPath, "key", "", "SSH private key path")

	// Flags obligatoires
	profileAddCmd.MarkFlagRequired("host")
	profileAddCmd.MarkFlagRequired("user")
}
