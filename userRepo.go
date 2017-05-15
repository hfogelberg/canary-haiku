package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func (connection *Connection) SaveUser(username string, pwd string) (token string, err error) {
	password := []byte(pwd)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "token", err
	}
	token := CreateToken(username)

	var user User
	user.Username = username
	user.Password = string(hashedPassword)
	user.Tokens[0] = token

	log.Println(user)

	err = connection.Db.C("users").Insert(&user)
	if err != nil {
		log.Println("Failed insert")
		return "", err
	}
	return token, nil
}

func (connection *Connection) UsernameIsInDb(username string) (u User, e error) {
	// TODO Error handling
	user := User{}
	err := connection.Db.C("users").Find(bson.M{"username": username}).One(&user)

	if err != nil {
		log.Println(err)
		return user, err
	}

	return user, nil
}
