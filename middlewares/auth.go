package middlewares

import (
	"CRUD/config"
	"CRUD/structs"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	// "strings"
)

var db = config.GetDB()

var UserAuthMux = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unauth := []string{"/register", "/login", "/forgot", "/logout"}
		requestPath := r.URL.Path

		// if strings.HasPrefix(requestPath, "/static/") {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// fmt.Printf("%#v", r)

		for _, value := range unauth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("token")

		if tokenHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			config.SendUnauthorizedResponse(w, `{"message": "Missing Headers"}`)
			return
		}

		tk := &structs.Token{}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("auth_pass")), nil
		})

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			config.SendUnauthorizedResponse(w, `{"message": "Malformed Headers"}`)
			return
		}

		if !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			config.SendUnauthorizedResponse(w, `{"message": "Invalid Headers"}`)
			return
		}

		var forCountScan int
		if err := db.QueryRow("SELECT id FROM users WHERE token = $1", tokenHeader).Scan(&forCountScan); err != nil {
			w.Header().Add("Content-Type", "application/json")
			config.SendUnauthorizedResponse(w, `{"message": "Invalid Session"}`)
			return
		}

		// Create Context for Request
		// ctx := context.WithValue(r.Context(), "email", tk.Email)
		// r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Gin Middleware
