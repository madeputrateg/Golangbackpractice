package main

import (
	"fmt"
	"os"
	"practice/pkg"

	"github.com/gorilla/mux"
)

func main() {
	database := make(map[string]string)
	database["DB_NAME"] = os.Getenv("DB_NAME")
	database["DB_PASSWORD"] = os.Getenv("DB_PASSWORD")
	database["DB_USERNAME"] = os.Getenv("DB_USERNAME")
	database["DB_SERVER"] = os.Getenv("DB_SERVER")
	db, err := pkg.DatabaseStart(database["DB_USERNAME"], database["DB_PASSWORD"], database["DB_NAME"], database["DB_SERVER"])
	if err != nil {
		fmt.Println("database err")
		return
	}
	r := mux.NewRouter()
	pkg.InitializeGlobalController(db, r)
}
