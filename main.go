package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type Haiku struct {
	Text string    `json:"text" bson:"text"`
	User string    `json:"user" bson:"user"`
	When time.Time `json:"when" bson:"when"`
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	http.HandleFunc("/", index)
	http.HandleFunc("/create", create(session))

	http.ListenAndServe(":3000", nil)
}

func create(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Println("GET create ")
			tpl.ExecuteTemplate(w, "create.html", nil)
		} else if r.Method == "POST" {
			fmt.Println("POST create")

			err := r.ParseForm()
			if err != nil {
				fmt.Println(err)
				return
			}

			var haiku Haiku
			haiku.Text = r.PostFormValue("text")
			haiku.User = r.PostFormValue("user")
			haiku.When = time.Now()

			session := s.Copy()
			defer session.Close()

			if err := session.DB("canaryhaiku").C("verses").Insert(&haiku); err != nil {
				log.Println("Failed insert")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
