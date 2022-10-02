package pkg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func DatabaseStart(username string, password string, database string, server string) (*sql.DB, error) {
	hsl := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		server, 5500, username, password, database)
	fmt.Println(hsl)
	db, err := sql.Open("postgres", hsl)
	if err != nil {
		fmt.Println("error database")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("database initialize")
	return db, nil
}
