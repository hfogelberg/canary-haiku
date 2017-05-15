package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(username string) string {
	expireToken := time.Now().Add(time.Minute * 1).Unix()

	claims := Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:3000",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(hmacSampleSecret)

	if err != nil {
		panic(err)
	}

	return tokenString
}
