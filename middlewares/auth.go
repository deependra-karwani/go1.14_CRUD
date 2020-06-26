package middlewares

import (
	"CRUD/config"
	"CRUD/structs"
	"context"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
			if value == "/user"+requestPath {
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
			if requestPath == "/user/refresh" {
				ctx := context.WithValue(r.Context(), "email", tk.Email)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			config.SendUnauthorizedResponse(w, `{"message": "Malformed Headers"}`)
			return
		}

		if !token.Valid {
			if requestPath == "/user/refresh" {
				ctx := context.WithValue(r.Context(), "email", tk.Email)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
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

func UserAuthGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		unauth := []string{"/register", "/login", "/forgot", "/logout"}
		requestPath := c.Request.URL.Path

		for _, value := range unauth {
			if value == requestPath {
				c.Next()
				return
			}
		}

		tokenHeader := c.GetHeader("token")

		if tokenHeader == "" {
			c.Header("Content-Type", "application/json")
			config.SendUnauthorizedResponse(c.Writer, `{"message": "Missing Headers"}`)
			return
		}

		tk := &structs.Token{}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("auth_pass")), nil
		})

		if err != nil {
			if requestPath == "/user/refresh" {
				r := c.Request
				ctx := context.WithValue(r.Context(), "email", tk.Email)
				r = r.WithContext(ctx)
				c.Next()
				return
			}
			c.Header("Content-Type", "application/json")
			config.SendUnauthorizedResponse(c.Writer, `{"message": "Malformed Headers"}`)
			return
		}

		if !token.Valid {
			if requestPath == "/user/refresh" {
				r := c.Request
				ctx := context.WithValue(r.Context(), "email", tk.Email)
				r = r.WithContext(ctx)
				c.Next()
				return
			}
			c.Header("Content-Type", "application/json")
			config.SendUnauthorizedResponse(c.Writer, `{"message": "Invalid Headers"}`)
			return
		}

		var forCountScan int
		if err := db.QueryRow("SELECT id FROM users WHERE token = $1", tokenHeader).Scan(&forCountScan); err != nil {
			c.Header("Content-Type", "application/json")
			config.SendUnauthorizedResponse(c.Writer, `{"message": "Invalid Session"}`)
			return
		}

		c.Next()
	}
}
