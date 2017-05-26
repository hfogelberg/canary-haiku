package main

import (
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/dgrijalva/jwt-go"
)

type Haiku struct {
	Text string `json:"text" bson:"text"`
	User string `json:"user" bson:"user"`
	When time.Time `json:"when" bson:"when"`
}

type User struct {
	Username string   `json:"username" bson:"username"`
	Password []byte   `json:"password"`
	Tokens   []string `json:"tokens" bson:"tokens"`
}

type JwtClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Connection struct {
	Db *mgo.Database
}
