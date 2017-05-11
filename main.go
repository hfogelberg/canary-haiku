package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Haiku struct {
	Text string `json:"text" bson:"text"`
	User string `json:"user" bson:"user"`
	// Display time.Time `json:"displayDate" bson:"displayDate"`
	When time.Time `json:"when" bson:"when"`
}

var tpl *template.Template

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// serve assets
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	// fs := http.FileServer(http.Dir("public"))
	// http.Handle("/public/styles/", http.StripPrefix("/public/", fs))

	// http.Handle("/canary-haiku/public/styles", http.StripPrefix("/canary-haiku/public/styles", http.FileServer(http.Dir("./canary-haiku/public/styles"))))

	http.HandleFunc("/", index(session))
	http.HandleFunc("/admin", admin(session))
	http.HandleFunc("/create", create(session))

	http.ListenAndServe(":3000", nil)
}

func admin(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var Haikus []*Haiku
		err := session.DB("canaryhaiku").C("verses").Find(nil).All(&Haikus)
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
}

func index(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var Verse *Haiku
		// d := time.Now()
		//err := session.DB("canaryhaiku").C("verses").Find(bson.M{"when": bson.M{"$gte": d}}).One(&Verse)
		err := session.DB("canaryhaiku").C("verses").Find(bson.M{"user": "bosse"}).One(&Verse)
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
}

func create(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

			session := s.Copy()
			defer session.Close()

			if err := session.DB("canaryhaiku").C("verses").Insert(&haiku); err != nil {
				log.Println("Failed insert")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
	}
}
