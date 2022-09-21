package middleware

import "net/http"

func Middleware(x http.Handler) http.Handler {
	return x
}
