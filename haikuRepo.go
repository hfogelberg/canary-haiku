package main

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

func (connection *Connection) GetHaikus() (h []*Haiku, e error) {
	var haikus []*Haiku
	err := connection.Db.C("verses").Find(nil).All(&haikus)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return haikus, nil
}

func (connection *Connection) GetDailyHaiku() (h *Haiku, e error) {
	var haiku *Haiku
	// d := time.Now()
	//err := session.DB("canaryhaiku").C("verses").Find(bson.M{"when": bson.M{"$gte": d}}).One(&Verse)
	err := connection.Db.C("verses").Find(bson.M{"user": "bosse"}).One(&haiku)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return haiku, nil
}

func (connection *Connection) CreateHaiku(haiku *Haiku) (e error) {
	err := connection.Db.C("verses").Insert(&haiku)
	if err != nil {
		log.Println("Failed insert")
		return err
	}
	return nil
}
