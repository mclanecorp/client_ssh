package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	// 1. Déterminer le dossier config utilisateur
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".config", "myssh")

	// 2. Créer le dossier s’il n’existe pas
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	// 3. Chemin vers la base SQLite
	dbPath := filepath.Join(configDir, "myssh.db")

	// 4. Ouvrir la base
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// 5. Créer la table profiles si absente
	schema := `
	CREATE TABLE IF NOT EXISTS profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		host TEXT NOT NULL,
		port INTEGER NOT NULL,
		user TEXT NOT NULL,
		password TEXT,
		key_path TEXT
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	DB = db
	return nil
}
