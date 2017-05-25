package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (connection *Connection) GetLogin(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/login.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (connection *Connection) PostLogin(w http.ResponseWriter, r *http.Request) {

	var user User

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	user, err2 := connection.UsernameIsInDb(username)
	if err2 != nil {
		log.Println("Error checking user in Db")
		log.Print(err)
		http.Error(w, err2.Error(), http.StatusInternalServerError)
	}

	log.Println("User")
	log.Print(user)

	if user.Username != "" {
		// We have a user. Verify password
		pwd := []byte(password)
		if err := bcrypt.CompareHashAndPassword(user.Password, pwd); err != nil {
			log.Println("Wrong password or error after checking password")
			log.Print(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		// Password OK. Generate token, create cookie and redirect to admin
		token := CreateToken(username)
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

	} else {
		// Unknown user
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}

func (connection *Connection) GetSignup(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/signup.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (connection *Connection) PostSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
