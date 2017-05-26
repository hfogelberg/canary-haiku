package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (connection *Connection) SaveUser(username string, pwd string, token string) (err error) {
	password := []byte(pwd)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}

	var user User
	user.Username = username
	user.Password = hashedPassword
	user.Tokens = append(user.Tokens, token)

	log.Println(user)

	err = connection.Db.C("users").Insert(&user)
	if err != nil {
		log.Println("Failed insert")
		return err
	}
	return nil
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

func (connection *Connection) PushToken(username string, token string) (e error) {
	log.Println("Push token")
	user := User{}

	change := mgo.Change{
		Update:    bson.M{"$push": bson.M{"tokens": token}},
		ReturnNew: true,
	}

	_, err := connection.Db.C("users").Find(bson.M{"username": username}).Apply(change, &user)
	if err != nil {
		log.Println("Error updating user with token")
		log.Print(err)
		return err
	}

	log.Println("OK appending token")

	return nil
}
