package structs

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	Email string
	jwt.StandardClaims
}
