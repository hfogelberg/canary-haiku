package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	mgo "gopkg.in/mgo.v2"
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

func Signup(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

			session := s.Copy()
			defer session.Close()

			err = session.DB("canaryhaiku").C("users").Insert(&user)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
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
}

func Archive(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var Haikus []*Haiku
		err := session.DB("canaryhaiku").C("verses").Find(nil).All(&Haikus)
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
}

func Admin(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET ADMIN")
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

func Index(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET INDEX")
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

func Create(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
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

			err = session.DB("canaryhaiku").C("verses").Insert(&haiku)
			if err != nil {
				log.Println("Failed insert")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Redirect(w, r, "/admin/", http.StatusSeeOther)
		}
	}
}
