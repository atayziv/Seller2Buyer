package middleware

import (
	"log"
	"net/http"
)

func Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("Middleware: Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		log.Println("Middleware: Method is valid")
		next.ServeHTTP(w, r)
	})
}
