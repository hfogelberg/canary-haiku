package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	mgo "gopkg.in/mgo.v2"
)

var tpl *template.Template

func main() {
	// Hook up Db
	session, err := mgo.Dial(MongoDBHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	connection := Connection{session.DB(MongoDb)}

	router := mux.NewRouter().StrictSlash(true)

	// serve assets
	static := http.StripPrefix("/public/", http.FileServer(http.Dir("public")))
	router.PathPrefix("/public/").Handler(static)

	// Routes
	router.HandleFunc("/admin", connection.Admin)
	router.HandleFunc("/create", connection.Create)
	router.HandleFunc("/archive", connection.Archive)
	router.HandleFunc("/signup", connection.Signup)
	router.HandleFunc("/", connection.Index)

	// Start server
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(Port)
}
