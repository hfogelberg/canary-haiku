package main

import (
	"html/template"
	"log"
	"net/http"
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
