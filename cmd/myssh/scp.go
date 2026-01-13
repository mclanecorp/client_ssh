package main

import (
	"fmt"

	"myssh/internal/profile"
	"myssh/internal/scp"
	"myssh/internal/ssh"

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

var (
	scpLocalPath  string
	scpRemotePath string
)

var scpUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to a remote server",
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

		// Upload SCP
		err = scp.Upload(cfg, scpLocalPath, scpRemotePath)
		if err != nil {
			fmt.Println("Erreur SCP upload:", err)
		} else {
			fmt.Println("Upload réussi !")
		}
	},
}

var scpDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from a remote server",
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

		// Download SCP
		err = scp.Download(cfg, scpRemotePath, scpLocalPath)
		if err != nil {
			fmt.Println("Erreur SCP download:", err)
		} else {
			fmt.Println("Download réussi !")
		}
	},
}


func init() {
	// Ajout des sous-commandes
	scpCmd.AddCommand(scpUploadCmd)
	scpCmd.AddCommand(scpDownloadCmd)

	rootCmd.AddCommand(scpCmd)

	// Flags upload
	scpUploadCmd.Flags().StringVar(&scpLocalPath, "local", "", "Local file path (required)")
	scpUploadCmd.Flags().StringVar(&scpRemotePath, "remote", "", "Remote file path (required)")
	scpUploadCmd.MarkFlagRequired("local")
	scpUploadCmd.MarkFlagRequired("remote")

	// Flags download
	scpDownloadCmd.Flags().StringVar(&scpLocalPath, "local", "", "Local file path (required)")
	scpDownloadCmd.Flags().StringVar(&scpRemotePath, "remote", "", "Remote file path (required)")
	scpDownloadCmd.MarkFlagRequired("local")
	scpDownloadCmd.MarkFlagRequired("remote")
}
