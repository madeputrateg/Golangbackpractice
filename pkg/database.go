package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func DatabaseStart(username string, password string, database string) (*sql.DB, error) {
	hsl := fmt.Sprint("postgresql//", username, ":", password, "@", database, "/todos?sslmode=disable")
	db, err := sql.Open("postgres", hsl)
	if err != nil {
		fmt.Println("error database")
		return nil, err
	}
	return db, nil
}
