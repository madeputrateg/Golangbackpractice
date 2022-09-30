package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"practice/internal/user"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

func MakeUserRepoInterface(Dr user.DataRepo) Service {
	return Service{Dr: Dr}
}

func GenerateToken(header string, payload map[string]string, secret string) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))
	header64 := base64.StdEncoding.EncodeToString([]byte(header))
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error -> payload marshal")
		return string(payloadstr), err
	}
	payload64 := base64.StdEncoding.EncodeToString([]byte(payloadstr))
	message := header64 + "." + payload64
	unsingnedStr := header + string(payloadstr)
	h.Write([]byte(unsingnedStr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	tokenstr := message + "." + signature
	return tokenstr, nil

}

func ValidateToken(token string, secret string) (bool, error) {
	splitToken := strings.Split(token, ".")
	if len(splitToken) != 3 {
		fmt.Println("error -> split token error")
		return false, nil
	}
	header, err := base64.StdEncoding.DecodeString(splitToken[0])
	if err != nil {
		fmt.Println("error -> header decode error")
		return false, nil
	}
	payload, err := base64.StdEncoding.DecodeString(splitToken[1])
	if err != nil {
		fmt.Println("error -> payload decode error")
		return false, nil
	}
	unsingnedStr := string(header) + string(payload)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsingnedStr))

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(signature)

	if signature != splitToken[2] {
		return false, nil
	}
	return true, nil
}

func (t Service) SignupHandler(rw http.ResponseWriter, r *http.Request) {
	var userdata user.User
	json.NewEncoder(rw).Encode(&userdata)
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
			rw.Write([]byte("Internal Server Error"))
			return
		}
	}
	back := user.Changebackdatatype(userdata)
	back.Password = Hashpassword(back.Password)
	tokenjwt, err := Generatejwttoken(back.Name, back.Email)
	if err != nil {
		rw.Write([]byte("Internal Server Error"))
		return
	}
	expire := time.Now().Add(11 * time.Minute)
	http.SetCookie(rw, &http.Cookie{
		Name:    "testtoken",
		Value:   tokenjwt,
		Expires: expire,
	})
	t.Dr.InsertDataUserRepo(ctx, user.Userdb{Name: back.Name, Password: back.Password, Email: back.Email})
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("User Created"))
}

func Generatejwttoken(user string, email string) (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(11 * time.Minute)
	claims["authorized"] = true
	claims["user"] = user
	claims["email"] = email
	lubang := os.Getenv("secretkey")
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
	headsedpass = base64.RawStdEncoding.EncodeToString([]byte(headsedpass))
	return headsedpass
}

type authser interface {
	SignupHandler(rw http.ResponseWriter, r *http.Request)
}
