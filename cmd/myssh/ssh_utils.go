package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// loadPrivateKey lit la clé privée SSH depuis un fichier et retourne un ssh.Signer
func loadPrivateKey(path string) (ssh.Signer, error) {
	keyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return signer, nil
}

// buildSSHConfig prépare la configuration SSH avec clé privée ou mot de passe
func buildSSHConfig(user, keyPath, password string) (*ssh.ClientConfig, error) {
	authMethods := []ssh.AuthMethod{}

	if keyPath != "" {
		signer, err := loadPrivateKey(keyPath)
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // temporaire
	}

	return config, nil
}

// connectSSH établit la connexion SSH
func connectSSH(host string, port int, config *ssh.ClientConfig) (*ssh.Client, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}
	return client, nil
}

// startInteractiveSession lance un shell interactif sur le serveur
func startInteractiveSession(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // activer l'écho
		ssh.TTY_OP_ISPEED: 14400, // vitesse input
		ssh.TTY_OP_OSPEED: 14400, // vitesse output
	}

	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %v", err)
	}

	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %v", err)
	}

	return session.Wait()
}
