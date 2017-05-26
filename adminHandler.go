package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func (connection *Connection) GetAdmin(w http.ResponseWriter, r *http.Request) {
	log.Println("GET ADMIN")

	haikus, err := connection.GetHaikus()
	if err != nil {
		log.Println("Get haikus failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tpl, err := template.New("").ParseFiles("templates/admin.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", haikus)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (connection *Connection) PostCreateHaiku(w http.ResponseWriter, r *http.Request) {
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

func (connection *Connection) GetCreateHaiku(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/create.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
