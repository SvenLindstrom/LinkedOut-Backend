package auth

import (
	"database/sql"
)

type AuthModel struct {
	DB *sql.DB
}

func (m *AuthModel) userExists(id string) (string, error) {

	var user_id string
	err := m.DB.QueryRow("SELECT id FROM users WHERE user_id=$1", id).
		Scan(&user_id)

	return user_id, err
}

func (m *AuthModel) creatUser(id string, name string) (string, error) {

	var user_id string
	err := m.DB.QueryRow("INSERT INTO users (user_id, name) VALUES ($1, $2) RETURNING id", id, name).
		Scan(&user_id)

	return user_id, err
}

func (m *AuthModel) setDeviceCode(id string, deviceCode string) error {
	_, err := m.DB.Exec("UPDATE users SET deviceCode=$1 WHERE id=$2", deviceCode, id)
	return err
}
