package main

import (
	"html/template"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

var tpl *template.Template

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// serve assets
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.HandleFunc("/admin", Admin(session))
	http.HandleFunc("/create", Create(session))
	http.HandleFunc("/about", About)
	http.HandleFunc("/archive", Archive(session))
	http.HandleFunc("/signup", Signup(session))
	http.HandleFunc("/", Index(session))

	http.ListenAndServe(":3000", nil)
}
