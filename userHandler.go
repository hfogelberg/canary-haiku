package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

func (connection *Connection) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl, err := template.New("").ParseFiles("templates/signup.html", "templates/base.html")
		err = tpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			return
		}

		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		_, err = connection.UsernameIsInDb(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		token, err := connection.SaveUser(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Create cookie with session token
		expiration := time.Now().Add(60 * time.Second)
		cookie := http.Cookie{Name: "token", Value: token, Expires: expiration}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/admin", http.StatusFound)
		tpl, err := template.New("").ParseFiles("templates/admin.html", "templates/base.html")
		err = tpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
