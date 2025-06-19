package middleware

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok && username == "admin" && password == "latuerts" {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Nicht autorisiert"))
	})
}
