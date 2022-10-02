package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"practice/internal/user"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

type Service struct {
	Dr user.DataRepo
}

type Contoller struct {
	S Service
	X *mux.Router
}

func MakeUserRepoInterface(Dr user.DataRepo) Service {
	return Service{Dr: Dr}
}

func (t Service) SigninHandler(rw http.ResponseWriter, r *http.Request) {
	var userdata user.User
	json.NewDecoder(r.Body).Decode(&userdata)
	userdata.Password = Hashpassword(userdata.Password)
	ctx := r.Context()
	data, err := t.Dr.GetUserDataRepo(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, isi := range data {
		if isi.Email == userdata.Email {
			if isi.Password != userdata.Password {
				rw.Write([]byte("false password"))
				return
			}
			jwttoken, err := Generatejwttoken(isi.Name, userdata.Email)
			if err != nil {
				fmt.Println(err)
				return
			}
			expire := time.Now().Add(11 * time.Minute)
			http.SetCookie(rw, &http.Cookie{
				Name:    "TestToken",
				Value:   jwttoken,
				Expires: expire,
			})
			rw.Write([]byte("login succesful"))
			return
		}
	}
	rw.Write([]byte("email not found"))
}

func (t Service) SignupHandler(rw http.ResponseWriter, r *http.Request) {
	var userdata user.User
	json.NewDecoder(r.Body).Decode(&userdata)
	ctx := r.Context()
	data, err := t.Dr.GetUserDataRepo(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsondata := user.Changedatatype(data)
	for _, isi := range jsondata {
		if isi.Email == userdata.Email {
			rw.WriteHeader(http.StatusConflict)
			rw.Write([]byte("User is udah ada"))
			return
		}
	}
	back := user.Changebackdatatype(userdata)
	back.Password = Hashpassword(back.Password)
	tokenjwt, err := Generatejwttoken(back.Name, back.Email)
	if err != nil {
		rw.Write([]byte("Internal Server Error di buat token"))
		fmt.Println(err)
		return
	}
	expire := time.Now().Add(11 * time.Minute)
	http.SetCookie(rw, &http.Cookie{
		Name:    "TestToken",
		Value:   tokenjwt,
		Expires: expire,
	})
	err = t.Dr.InsertDataUserRepo(ctx, user.Userdb{Name: back.Name, Password: back.Password, Email: back.Email})
	if err != nil {
		fmt.Println(err)
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("User Created"))
}

func Generatejwttoken(user string, email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(11 * time.Minute)
	claims["authorized"] = true
	claims["user"] = user
	claims["email"] = email
	lubang := os.Getenv("KEY")
	tokenString, err := token.SignedString([]byte(lubang))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Hashpassword(pass string) string {
	var headsedpass string
	supersecret := os.Getenv("salt")
	headsedpass = supersecret + pass + supersecret
	h := sha256.New()
	h.Write([]byte(headsedpass))
	bs := h.Sum(nil)
	return string(bs)
}

func (c Contoller) InitializeController() {
	c.X.HandleFunc("/Signup", c.S.SignupHandler).Methods(http.MethodPost)
	c.X.HandleFunc("/Signin", c.S.SigninHandler).Methods(http.MethodPost)
}

func SetupaAuthContoller(X *mux.Router, S Service) Contoller {
	return Contoller{
		X: X,
		S: S,
	}
}

type Authser interface {
	SignupHandler(rw http.ResponseWriter, r *http.Request)
}
