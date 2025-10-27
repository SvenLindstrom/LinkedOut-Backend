package location

import (
	"database/sql"
	"encoding/json"
)

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
	var matchingInterestsJSON []byte
	rows, err := m.DB.Query(`
			SELECT u.id, u.name, u.bio, u.profession,
			ST_Distance(
				u.location,
				ST_SetSRID(ST_MakePoint($1, $2), 4326)
			) AS distance,
			COALESCE( json_agg(
				json_build_object(
					'id', i.id,
					'name', i.name
				)),
				'[]'::json
			) AS matching_interests
			FROM users u
			LEFT JOIN users_interests ui ON u.id = ui.user_id
			LEFT JOIN interests i ON ui.interest_id = i.id
			WHERE 
				u.connecting = TRUE
				AND u.id != $4
				AND ST_DWithin(
					u.location,
					ST_SetSRID(ST_MakePoint($1, $2), 4326),
					$3
				)
				AND ui.interest_id IN (
					SELECT interest_id
					FROM users_interests
					WHERE user_id = $4
				)
			GROUP BY u.id, u.name, u.bio, u.profession, u.location
			ORDER BY distance ASC;
		`,
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

		if err := rows.Scan(&u.Id, &u.Name, &u.Bio, &u.Profession, &u.Distance, &matchingInterestsJSON); err != nil {
			return nil, err
		}

		err = json.Unmarshal(matchingInterestsJSON, &u.Tags)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, err
}
