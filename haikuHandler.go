package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func About(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/about.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (connection *Connection) Archive(w http.ResponseWriter, r *http.Request) {
	haikus, err := connection.GetHaikus()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tpl, err = template.New("").ParseFiles("templates/archive.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", haikus)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (connection *Connection) Admin(w http.ResponseWriter, r *http.Request) {
	haikus, err := connection.GetHaikus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tpl, err := template.New("").ParseFiles("templates/admin.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", haikus)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (connection *Connection) Index(w http.ResponseWriter, r *http.Request) {
	haiku, err := connection.GetDailyHaiku()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tpl, err := template.New("").ParseFiles("templates/index.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", haiku)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

		err = connection.CreateHaiku(&haiku)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/admin/", http.StatusInternalServerError)
	}
}
