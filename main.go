package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
)

type Haiku struct {
	Text string    `json:"text" bson:"text"`
	User string    `json:"user" bson:"user"`
	When time.Time `json:"when" bson:"when"`
}

type Adapter func(http.Handler) http.Handler

func main() {
	db, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Cannot connect to Db", err)
		return
	}
	defer db.Close()

	h := Adapt(http.HandlerFunc(handle), withDB(db))
	http.Handle("/haiku", context.ClearHandler(h))

	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHaikus(w, r)
	case "POST":
		createHaiku(w, r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

func createHaiku(w http.ResponseWriter, r *http.Request) {
	log.Println("Create Haiku")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	var haiku Haiku
	haiku.Text = r.PostFormValue("text")
	haiku.User = r.PostFormValue("user")
	haiku.When = time.Now()

	log.Println("Haiku", haiku)

	db := context.Get(r, "database").(*mgo.Session)
	if err := db.DB("canryhaiku").C("verses").Insert(&haiku); err != nil {
		log.Println("Failed insert")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Insert OK")
}

func getHaikus(w http.ResponseWriter, r *http.Request) {

}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func withDB(db *mgo.Session) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dbsession := db.Copy()
			defer dbsession.Close() // clean up
			context.Set(r, "database", dbsession)
			h.ServeHTTP(w, r)
		})
	}
}
