package profile

import (
	"database/sql"
	"fmt"

	"myssh/internal/config"
)

// Profile représente un profil SSH stocké en base
type Profile struct {
	ID       int
	Name     string
	Host     string
	Port     int
	User     string
	Password string
	KeyPath  string
}

// Create ajoute un nouveau profil en base
func Create(p Profile) error {
	query := `
	INSERT INTO profiles (name, host, port, user, password, key_path)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := config.DB.Exec(
		query,
		p.Name,
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.KeyPath,
	)

	if err != nil {
		return fmt.Errorf("cannot create profile: %w", err)
	}

	return nil
}

// GetByName récupère un profil par son nom
func GetByName(name string) (*Profile, error) {
	query := `
	SELECT id, name, host, port, user, password, key_path
	FROM profiles
	WHERE name = ?
	`

	row := config.DB.QueryRow(query, name)

	var p Profile
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Host,
		&p.Port,
		&p.User,
		&p.Password,
		&p.KeyPath,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("profile '%s' not found", name)
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// List retourne tous les profils existants
func List() ([]Profile, error) {
	query := `
	SELECT id, name, host, port, user, password, key_path
	FROM profiles
	ORDER BY name
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile

	for rows.Next() {
		var p Profile
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Host,
			&p.Port,
			&p.User,
			&p.Password,
			&p.KeyPath,
		); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}

	return profiles, nil
}

// Delete supprime un profil par son nom
func Delete(name string) error {
	result, err := config.DB.Exec(
		`DELETE FROM profiles WHERE name = ?`,
		name,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("profile '%s' not found", name)
	}

	return nil
}

func ResolveSSHConfig(
	profileName string,
	host string,
	port int,
	user string,
	password string,
	key string,
) (string, int, string, string, string, error) {

	if profileName == "" {
		return host, port, user, password, key, nil
	}

	p, err := Get(profileName)
	if err != nil {
		return "", 0, "", "", "", err
	}

	return p.Host, p.Port, p.User, p.Password, p.KeyPath, nil
}

// Get est un alias de GetByName pour simplifier ResolveSSHConfig
func Get(name string) (*Profile, error) {
	return GetByName(name)
}