package databases

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func Pg_init() *sql.DB {

	pwd := os.Getenv("POSTGRES_PASSWORD")
	name := os.Getenv("POSTGRES_DB")
	addr := os.Getenv("POSTGRES_ADDR")

	url := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", pwd, addr, name)
	db, err := sql.Open("pgx", url)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
