package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func About(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/about.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func (connection *Connection) Archive(w http.ResponseWriter, r *http.Request) {
	var Haikus []*Haiku
	err := connection.Db.C("verses").Find(nil).All(&Haikus)
	if err != nil {
		log.Fatalln(err)
		return
	}

	tpl, err := template.New("").ParseFiles("templates/archive.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", Haikus)
	if err != nil {
		log.Println(err)
		return
	}
}

func (connection *Connection) Admin(w http.ResponseWriter, r *http.Request) {
	log.Println("GET ADMIN")

	var Haikus []*Haiku
	err := connection.Db.C("verses").Find(nil).All(&Haikus)
	if err != nil {
		log.Fatalln(err)
		return
	}

	tpl, err := template.New("").ParseFiles("templates/admin.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", Haikus)
	if err != nil {
		log.Println(err)
		return
	}
}

func (connection *Connection) Index(w http.ResponseWriter, r *http.Request) {
	log.Println("GET INDEX")

	var Verse *Haiku
	// d := time.Now()
	//err := session.DB("canaryhaiku").C("verses").Find(bson.M{"when": bson.M{"$gte": d}}).One(&Verse)
	err := connection.Db.C("verses").Find(bson.M{"user": "bosse"}).One(&Verse)
	if err != nil {
		log.Println(err)
		return
	}

	tpl, err := template.New("").ParseFiles("templates/index.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", Verse)
	if err != nil {
		log.Println(err)
		return
	}
}

func (connection *Connection) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl, err := template.New("").ParseFiles("templates/create.html", "templates/base.html")
		err = tpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println(err)
			return
		}
	} else if r.Method == "POST" {
		fmt.Println("POST create")

		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			return
		}

		var haiku Haiku
		haiku.Text = r.PostFormValue("text")
		haiku.User = r.PostFormValue("createdBy")
		// dt := r.PostFormValue("displayDate")
		// haiku.Display = time.Parse("2006-01-02", dt)
		haiku.When = time.Now()

		err = connection.Db.C("verses").Insert(&haiku)
		if err != nil {
			log.Println("Failed insert")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/admin/", http.StatusSeeOther)
	}
}
