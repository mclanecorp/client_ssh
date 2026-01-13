package scp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"myssh/internal/ssh"
)

// Upload transfère un fichier local vers un serveur distant
func Upload(cfg ssh.Config, localPath, remotePath string) error {
	// Ouvrir le fichier local
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("cannot open local file: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Connexion SSH
	client, err := ssh.Connect(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	// Nouvelle session pour SCP
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

	// Lancer SCP côté serveur en mode réception (-t)
	if err := session.Start("scp -t " + remotePath); err != nil {
		return err
	}

	// Envoyer l'en-tête du fichier
	fmt.Fprintf(stdin, "C0644 %d %s\n", info.Size(), filepath.Base(remotePath))
	if err := waitForAck(stdout); err != nil {
		return err
	}

	// Copier le contenu du fichier
	if _, err := io.Copy(stdin, file); err != nil {
		return err
	}

	// Envoyer le signal de fin
	fmt.Fprint(stdin, "\x00")
	if err := waitForAck(stdout); err != nil {
		return err
	}

	return session.Wait()
}

// Download récupère un fichier depuis un serveur distant
func Download(cfg ssh.Config, remotePath, localPath string) error {
	// Connexion SSH
	client, err := ssh.Connect(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	// Nouvelle session pour SCP
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Préparer les pipes
	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	// Lancer SCP côté serveur en mode envoi (-f)
	if err := session.Start("scp -f " + remotePath); err != nil {
		return err
	}

	// Envoyer l'ack initial
	if _, err := stdin.Write([]byte{0}); err != nil {
		return err
	}

	// Lire l'en-tête du fichier envoyé par le serveur
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

	// Envoyer l'ack
	if _, err := stdin.Write([]byte{0}); err != nil {
		return err
	}

	// Créer le fichier local
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("cannot create local file: %v", err)
	}
	defer localFile.Close()

	// Copier le contenu
	if _, err := io.CopyN(localFile, stdout, size); err != nil {
		return err
	}

	// Lire et envoyer le dernier ack
	if _, err := stdout.Read(buf[:1]); err != nil {
		return err
	}
	if _, err := stdin.Write([]byte{0}); err != nil {
		return err
	}

	return session.Wait()
}

// waitForAck lit le retour SCP et vérifie que c'est un succès
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
