package location

import "database/sql"

type LocationModel struct {
	DB *sql.DB
}

func (m *LocationModel) UpdateLocation(id string, location Location) error {
	_, err := m.DB.Exec(
		"UPDATE users SET location=ST_SetSRID(ST_MakePoint($1, $2),4326) WHERE id=$3",
		location.Lon,
		location.Lat,
		id,
	)
	return err
}

func (m *LocationModel) UpdateStatus(id string, status bool) error {
	_, err := m.DB.Exec("UPDATE users SET connecting=$1  WHERE id=$2", status, id)
	return err
}

func (m *LocationModel) getProximity(id string, location Location) ([]UserProx, error) {
	rows, err := m.DB.Query(
		"SELECT name FROM users WHERE  connecting=true AND ST_DWithin(location, ST_SetSRID(ST_MakePoint($1, $2), 4326), 50) AND id!=$3",
		location.Lon,
		location.Lat,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]UserProx, 0)
	for rows.Next() {

		var u UserProx

		if err := rows.Scan(&u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, err
}
