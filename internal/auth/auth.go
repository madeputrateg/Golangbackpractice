package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DataRepo struct {
	Db *sql.DB
}

func ProvideDB(Db *sql.DB) DataRepo {
	return DataRepo{Db: Db}
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

func SignupHandler(rw http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Email"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Email Missing"))
		return
	}
	if _, ok := r.Header["Passwordhash"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Passwordhash Missing"))
		return
	}
	if _, ok := r.Header["Fullname"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Fullname Missing"))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("User Created"))
}