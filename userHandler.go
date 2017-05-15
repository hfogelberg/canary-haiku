package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (connection *Connection) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl, err := template.New("").ParseFiles("templates/signup.html", "templates/base.html")
		err = tpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println("POST SIGNUP")

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			return
		}

		var user User
		user.Username = r.PostFormValue("username")
		password := []byte(r.PostFormValue("password"))
		hashedPassword, err2 := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err2 != nil {
			log.Println(err2)
			return
		}
		user.Password = string(hashedPassword)

		err = connection.SaveUser(user)
		if err != nil {
			//TODO Send error response
			return
		}

		token := CreateToken(user.Username)
		expiration := time.Now().Add(60 * time.Second)
		cookie := http.Cookie{Name: "token", Value: token, Expires: expiration}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/admin", http.StatusFound)
		tpl, err := template.New("").ParseFiles("templates/admin.html", "templates/base.html")
		err = tpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
