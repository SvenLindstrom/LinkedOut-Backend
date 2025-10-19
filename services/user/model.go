package user

import (
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) getInfo(id string) (UserInfo, error) {

	var userInfo UserInfo

	err := m.DB.QueryRow(
		"SELECT id, name, bio FROM users WHERE id=$1",
		id,
	).Scan(&userInfo.Id, &userInfo.Name, &userInfo.Bio)

	return userInfo, err
}

func (m *UserModel) updateBio(id string, bio string) error {
	_, err := m.DB.Exec("UPDATE users SET bio=$1 WHERE id=$2", bio, id)
	return err
}
