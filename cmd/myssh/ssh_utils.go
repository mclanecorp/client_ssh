package main

import (
	"fmt"
	"io/ioutil"

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

// buildSSHConfig prépare la configuration SSH pour l'utilisateur et la clé privée
func buildSSHConfig(user string, keyPath string) (*ssh.ClientConfig, error) {
	authMethods := []ssh.AuthMethod{}

	if keyPath != "" {
		signer, err := loadPrivateKey(keyPath)
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // à remplacer plus tard par vérification
	}

	return config, nil
}
