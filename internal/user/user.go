package user

import (
	"context"
	"database/sql"
	"fmt"
)

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

type Users []*User
type Userdbs []*Userdb

type DataRepo struct {
	Db *sql.DB
}

func ProvideDB(Db *sql.DB) DataRepo {
	return DataRepo{Db: Db}
}

const (
	SELECT_DATA_USERS = `SELECT name,password,email FROM users`
	INSERT_DATA_USER  = `INSERT INTO users(name,password,email) VALUES ($1,$2,$3)`
)

func (t DataRepo) GetUserDataRepo(ctx context.Context) (Userdbs, error) {
	stmt, err := t.Db.PrepareContext(ctx, SELECT_DATA_USERS)
	if err != nil {
		return nil, err
	}
	var temp Userdbs
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var temp2 Userdb
		if err := rows.Scan(&temp2.Name, &temp2.Password, &temp2.Email); err != nil {
			return nil, err
		}
		temp = append(temp, &temp2)
	}
	return temp, nil
}

func (t DataRepo) InsertDataUserRepo(ctx context.Context, userdata Userdb) error {
	stmt, err := t.Db.PrepareContext(ctx, INSERT_DATA_USER)
	if err != nil {
		fmt.Println("prepare err")
		return err
	}
	_, err = stmt.QueryContext(ctx, userdata.Name, userdata.Password, userdata.Email)
	if err != nil {
		fmt.Println("Querry err")
		return err
	}
	return nil
}

type UserRepoInterface interface {
	GetUserDataRepo(ctx context.Context) (Userdbs, error)
	InsertDataUserRepo(ctx context.Context, userdata Userdb) error
}

func Changedatatype(t Userdbs) Users {
	var temp Users
	for _, data := range t {
		temp = append(temp, &User{Email: data.Email, Name: data.Name, Password: data.Password})
	}
	return temp
}

func Changebackdatatype(t User) Userdb {
	var data Userdb
	data.Email = t.Email
	data.Name = t.Name
	data.Password = t.Password
	return data
}
