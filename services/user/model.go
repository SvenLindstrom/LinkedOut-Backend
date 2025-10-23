package user

import (
	"database/sql"
	"encoding/json"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) getInfo(id string) (UserInfo, error) {

	var userInfo UserInfo
	var interestsJSON []byte

	err := m.DB.QueryRow(
		`SELECT
		users.id,
		users.name,
		users.bio,
		users.profession,
		COALESCE(
			json_agg(
				json_build_object('id', interests.id, 'name', interests.name)
			)
		) AS interests
		FROM users
		LEFT JOIN users_interests ON users.id = users_interests.user_id
		LEFT JOIN interests ON users_interests.interest_id = interests.id
		WHERE users.id=$1
		GROUP BY users.id, users.name, users.bio, users.profession`,
		id,
	).Scan(&userInfo.Id, &userInfo.Name, &userInfo.Bio, &userInfo.Profession, &interestsJSON)

	if err != nil {
		println(err.Error())
		return UserInfo{}, err
	}

	err = json.Unmarshal(interestsJSON, &userInfo.Interests)
	if err != nil {
		println(err.Error())
		return UserInfo{}, err
	}

	return userInfo, err
}

func (um *UserModel) SaveInfo(userID, profession, bio string, interests []Interest) error {
	tx, err := um.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE users SET profession = $1, bio = $2 WHERE id=$3",
		profession,
		bio,
		userID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM users_interests WHERE user_id=$1", userID)
	if err != nil {
		return err
	}

	for _, i := range interests {
		_, err = tx.Exec(
			"INSERT INTO users_interests (user_id, interest_id) VALUES ($1, $2)",
			userID,
			i.Id,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (um *UserModel) FindAllInterests() ([]*Interest, error) {
	rows, err := um.DB.Query("SELECT id, name FROM interests")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	interests := make([]*Interest, 0)
	for rows.Next() {
		var i Interest
		if err := rows.Scan(&i.Id, &i.Name); err != nil {
			return nil, err
		}
		interests = append(interests, &i)
	}
	return interests, err
}
