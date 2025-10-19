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

func (m *LocationModel) getProximity(id string, prox Proximity) ([]UserProx, error) {
	rows, err := m.DB.Query(`
		SELECT id, name, bio, ST_Distance(
			location,
		ST_SetSRID(ST_MakePoint($1, $2), 4326)) AS distance FROM users WHERE  connecting=true AND ST_DWithin(location, ST_SetSRID(ST_MakePoint($1, $2), 4326), $3) AND id!=$4 ORDER BY distance ASC`,
		prox.Location.Lon,
		prox.Location.Lat,
		prox.Distance,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]UserProx, 0)
	for rows.Next() {

		var u UserProx

		if err := rows.Scan(&u.Id, &u.Name, &u.Bio, &u.Distance); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, err
}
