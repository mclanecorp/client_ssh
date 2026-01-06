package main

import (
	"fmt"
	"myssh/internal/profile"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Gestion des profils SSH",
}

var (
    hostFlag     string
    portFlag     int
    userFlag     string
    passwordFlag string
    keyFlag      string
)

var profileAddCmd = &cobra.Command{
	Use:   "add [nom]",
	Short: "Ajouter un profil SSH",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		p := profile.Profile{
			Name:     name,
			Host:     hostFlag,
			Port:     portFlag,
			User:     userFlag,
			Password: passwordFlag,
			KeyPath:  keyFlag,
		}
		if err := profile.Create(p); err != nil {
			fmt.Println("Erreur création profil:", err)
			return
		}
		fmt.Println("Profil créé avec succès:", name)
	},
}

var profileEditCmd = &cobra.Command{
	Use:   "edit [nom]",
	Short: "Modifier un profil SSH",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		p, err := profile.GetByName(name)
		if err != nil {
			fmt.Println("Profil introuvable:", err)
			return
		}

		updates := make(map[string]interface{})

		if cmd.Flags().Changed("host") {
			updates["host"] = hostFlag
		}
		if cmd.Flags().Changed("port") {
			updates["port"] = portFlag
		}
		if cmd.Flags().Changed("user") {
			updates["user"] = userFlag
		}
		if cmd.Flags().Changed("password") {
			updates["password"] = passwordFlag
		}
		if cmd.Flags().Changed("key") {
			updates["key"] = keyFlag
		}

		if err := profile.Update(p, updates); err != nil {
			fmt.Println("Erreur mise à jour profil:", err)
			return
		}

		fmt.Println("Profil mis à jour:", name)
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lister tous les profils",
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := profile.List()
		if err != nil {
			fmt.Println("Erreur listing profils:", err)
			return
		}
		if len(profiles) == 0 {
			fmt.Println("Aucun profil enregistré.")
			return
		}
		fmt.Printf("%-10s %-16s %-12s %-5s\n", "NAME", "HOST", "USER", "PORT")
		for _, p := range profiles {
			fmt.Printf("%-10s %-16s %-12s %-5d\n", p.Name, p.Host, p.User, p.Port)
		}
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete [nom]",
	Short: "Supprimer un profil",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := profile.Delete(name); err != nil {
			fmt.Println("Erreur suppression profil:", err)
			return
		}
		fmt.Println("Profil supprimé:", name)
	},
}

func init() {
    rootCmd.AddCommand(profileCmd)

    profileCmd.AddCommand(profileAddCmd)
    profileCmd.AddCommand(profileEditCmd)
    profileCmd.AddCommand(profileListCmd)
    profileCmd.AddCommand(profileDeleteCmd)

    profileAddCmd.Flags().StringVar(&hostFlag, "host", "", "Adresse SSH")
    profileAddCmd.Flags().IntVar(&portFlag, "port", 22, "Port SSH")
    profileAddCmd.Flags().StringVar(&userFlag, "user", "", "Utilisateur SSH")
    profileAddCmd.Flags().StringVar(&passwordFlag, "password", "", "Mot de passe SSH")
    profileAddCmd.Flags().StringVar(&keyFlag, "key", "", "Chemin clé privée SSH")

    profileEditCmd.Flags().StringVar(&hostFlag, "host", "", "Adresse SSH")
    profileEditCmd.Flags().IntVar(&portFlag, "port", 0, "Port SSH")
    profileEditCmd.Flags().StringVar(&userFlag, "user", "", "Utilisateur SSH")
    profileEditCmd.Flags().StringVar(&passwordFlag, "password", "", "Mot de passe SSH")
    profileEditCmd.Flags().StringVar(&keyFlag, "key", "", "Chemin clé privée SSH")
}

