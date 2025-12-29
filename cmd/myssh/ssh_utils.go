package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

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

	return &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

// connectSSH établit une connexion SSH
func connectSSH(host string, port int, config *ssh.ClientConfig) (*ssh.Client, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}
	return client, nil
}

// startInteractiveSession lance un shell interactif
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

// SCP upload
func scpUpload(host string, port int, user, password, keyPath, localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("cannot open local file: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	config, err := buildSSHConfig(user, keyPath, password)
	if err != nil {
		return err
	}

	client, err := connectSSH(host, port, config)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	if err := session.Start("scp -t " + remotePath); err != nil {
		return err
	}

	fmt.Fprintf(stdin, "C0644 %d %s\n", info.Size(), filepath.Base(remotePath))
	if err := waitForAck(stdout); err != nil {
		return err
	}

	if _, err := io.Copy(stdin, file); err != nil {
		return err
	}

	fmt.Fprint(stdin, "\x00")
	if err := waitForAck(stdout); err != nil {
		return err
	}

	return session.Wait()
}

// scpDownload récupère un fichier depuis le serveur
func scpDownload(host string, port int, user, password, keyPath, remotePath, localPath string) error {
	// 1. Construire la config SSH
	config, err := buildSSHConfig(user, keyPath, password)
	if err != nil {
		return err
	}

	// 2. Connexion SSH
	client, err := connectSSH(host, port, config)
	if err != nil {
		return err
	}
	defer client.Close()

	// 3. Nouvelle session
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// 4. Préparer les pipes
	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	// 5. Lancer SCP côté serveur en mode 'from' (-f)
	if err := session.Start("scp -f " + remotePath); err != nil {
		return err
	}

	// 6. Envoyer l'ack initial
	if _, err := stdin.Write([]byte{0}); err != nil {
		return err
	}

	// 7. Lire l’en-tête du fichier envoyé par le serveur
	buf := make([]byte, 1024)
	n, err := stdout.Read(buf)
	if err != nil {
		return err
	}

	var mode string
	var size int64
	var filename string
	_, err = fmt.Sscanf(string(buf[:n]), "C%s %d %s", &mode, &size, &filename)
	if err != nil {
		return fmt.Errorf("failed to parse scp header: %v", err)
	}

	// 8. Envoyer l'ack
	if _, err := stdin.Write([]byte{0}); err != nil {
		return err
	}

	// 9. Créer le fichier local
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("cannot create local file: %v", err)
	}
	defer localFile.Close()

	// 10. Copier le contenu
	if _, err := io.CopyN(localFile, stdout, size); err != nil {
		return err
	}

	// 11. Lire et envoyer le dernier ack
	if _, err := stdout.Read(buf[:1]); err != nil {
		return err
	}
	if _, err := stdin.Write([]byte{0}); err != nil {
		return err
	}

	return session.Wait()
}


// waitForAck lit le retour SCP
func waitForAck(r io.Reader) error {
	buf := make([]byte, 1)
	if _, err := r.Read(buf); err != nil {
		return err
	}
	if buf[0] != 0 {
		return fmt.Errorf("scp error, ack=%d", buf[0])
	}
	return nil
}
