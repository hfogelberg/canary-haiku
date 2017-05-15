package main

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

func (connection *Connection) SaveUser(user User) error {
	// TODO Error handling
	err := connection.Db.C("users").Insert(&user)
	if err != nil {
		log.Println("Failed insert")
		return err
	}
	return nil
}

func (connection *Connection) UsernameIsInDb(username string) User {
	// TODO Error handling
	user := User{}
	err := connection.Db.C("users").Find(bson.M{"username": username}).One(&user)

	if err != nil {
		log.Println(err)
	}

	return user
}
