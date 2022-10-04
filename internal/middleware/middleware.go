package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func Middlewareauth(x http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("TestToken")
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("error cookie gak ada kue"))
			return
		}
		secret := os.Getenv("KEY")
		tknstr := c.Value
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tknstr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		fmt.Println(token)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("error cookie kue jelek"))
			return
		}
		if !token.Valid {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		x.ServeHTTP(rw, r)
	})
}
