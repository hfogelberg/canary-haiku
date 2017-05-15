package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Haiku struct {
	Text string `json:"text" bson:"text"`
	User string `json:"user" bson:"user"`
	// Display time.Time `json:"displayDate" bson:"displayDate"`
	When time.Time `json:"when" bson:"when"`
}

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Claims struct {
	Username string `json:"username`
	jwt.StandardClaims
}

var tpl *template.Template
var hmacSampleSecret = []byte("secret")

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// serve assets
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.HandleFunc("/admin", admin(session))
	http.HandleFunc("/create", create(session))
	http.HandleFunc("/about", about)
	http.HandleFunc("/archive", archive(session))
	http.HandleFunc("/signup", signup(session))
	http.HandleFunc("/", index(session))

	http.ListenAndServe(":3000", nil)
}

func about(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("").ParseFiles("templates/about.html", "templates/base.html")
	err = tpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func signup(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
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

			token := createToken(user.Username)
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

func createToken(username string) string {
	expireToken := time.Now().Add(time.Minute * 1).Unix()

	claims := Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:3000",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(hmacSampleSecret)

	if err != nil {
		panic(err)
	}

	return tokenString
}

func archive(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
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

func admin(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
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

func index(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
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
