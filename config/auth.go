package config

import (
	"CRUD/structs"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

func GenToken(email string) (string, error) {
	tk := &structs.Token{Email: email}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("auth_pass")))
	if err != nil {
		return "No Token", err
	}
	return tokenString, nil
}
