package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Haiku struct {
	Text string `json:"text" bson:"text"`
	User string `json:"user" bson:"user"`
	// Display time.Time `json:"displayDate" bson:"displayDate"`
	When time.Time `json:"when" bson:"when"`
}

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Claims struct {
	Username string `json:"username`
	jwt.StandardClaims
}
