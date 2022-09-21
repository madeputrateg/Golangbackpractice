package user

import "database/sql"

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Userdb struct {
	Name     string `db:"name"`
	Password string `db:"password"`
	Email    string `db:"Email"`
}

type DataRepo struct {
	Db *sql.DB
}

func ProvideDB(Db *sql.DB) DataRepo {
	return DataRepo{Db: Db}
}
