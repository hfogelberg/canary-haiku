package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
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
		if err := connection.PushToken(username, token); err != nil {
			log.Println("Error pushing token to Db")
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		expiration := time.Now().Add(5 * time.Minute)
		cookie := http.Cookie{Name: "haikuToken", Value: token, Expires: expiration}
		http.SetCookie(w, &cookie)
		log.Println("Cookie set")

		log.Println("Redirecting to admin")
		connection.GetAdmin(w, r)

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

	token := CreateToken(username)
	if err := connection.SaveUser(username, password, token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	saveSessionToken(token)

	log.Println("All done. Redirecting to admin")
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func saveSessionToken(token string) {
	log.Println("Saving to store")
	store := sessions.NewCookieStore([]byte(SessionsSecret))
	session, _ := store.Get(r, SessionName)
	session.Values["haikuToken"] = token
	session.Save(r, w)
}
