package ssh

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// Config contient les paramètres de connexion SSH
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	KeyPath  string
}

// Client encapsule un client SSH
type Client struct {
	*ssh.Client
}

// loadPrivateKey lit la clé privée SSH depuis un fichier
func loadPrivateKey(path string) (ssh.Signer, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return signer, nil
}

// buildSSHConfig construit la config SSH avec clé ou mot de passe
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

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method provided (need password or key)")
	}

	return &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

// Connect établit une connexion SSH
func Connect(cfg Config) (*Client, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("host is required")
	}
	if cfg.User == "" {
		return nil, fmt.Errorf("user is required")
	}
	if cfg.Port == 0 {
		cfg.Port = 22
	}

	sshConfig, err := buildSSHConfig(cfg.User, cfg.KeyPath, cfg.Password)
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}

	return &Client{Client: client}, nil
}

// StartInteractiveSession lance un shell interactif
func (c *Client) StartInteractiveSession() error {
	session, err := c.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %v", err)
	}

	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %v", err)
	}

	return session.Wait()
}
