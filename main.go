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

	// Public routes
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/", connection.Index)
	router.HandleFunc("/archive", connection.Archive).Methods("GET")
	router.HandleFunc("/signup", connection.GetSignup).Methods("GET")
	router.HandleFunc("/signup", connection.PostSignup).Methods("POST")
	router.HandleFunc("/login", connection.GetLogin).Methods("GET")
	router.HandleFunc("/login", connection.PostLogin).Methods("POST")

	// Protected routes
	adm := router.PathPrefix("/admin").Subrouter()
	adm.HandleFunc("/", connection.GetAdmin).Methods("GET")
	adm.HandleFunc("/create", connection.GetCreateHaiku).Methods("GET")
	adm.HandleFunc("/create", connection.PostCreateHaiku).Methods("POST")

	// Mux
	mux := http.NewServeMux()
	mux.Handle("/", router)
	mux.Handle("/admin/", negroni.New(
		negroni.HandlerFunc(AuthMiddleware),
		negroni.Wrap(router),
	))

	// Serve assets
	static := http.StripPrefix("/public/", http.FileServer(http.Dir("public")))
	router.PathPrefix("/public/").Handler(static)

	// Kick off server
	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)

	// // Hook up Db
	// session, err := mgo.Dial(MongoDBHost)
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()
	// session.SetMode(mgo.Monotonic, true)
	// connection := Connection{session.DB(MongoDb)}

	// router := mux.NewRouter().StrictSlash(true)

	// // serve assets
	// static := http.StripPrefix("/public/", http.FileServer(http.Dir("public")))
	// router.PathPrefix("/public/").Handler(static)

	// // Routes
	// router.HandleFunc("/admin", negroni.New(
	// 	negroni.HandlerFunc(connection.ValidateUser),
	// 	negroni.Wrap(router),
	// ))
	// router.HandleFunc("/create", negroni.New(
	// 	negroni.HandlerFunc(connection.ValidateUser),
	// 	negroni.Wrap(router),
	// ))
	// router.HandleFunc("/archive", connection.Archive)
	// router.HandleFunc("/signup", connection.Signup)
	// router.HandleFunc("/login", connection.Login)
	// router.HandleFunc("/", connection.Index)

	// // Start server
	// n := negroni.Classic()
	// n.UseHandler(router)
	// n.Run(Port)
}
