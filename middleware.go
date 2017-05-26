package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

func (connection *Connection) AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("*** Auth Middleware ***")

	store := sessions.NewCookieStore([]byte(SessionsSecret))
	session, _ := store.Get(r, SessionName)

	val := session.Values["haikuToken"]

	log.Println("Cookie fetched from store")

	if val != nil {
		log.Println("There is a token")
		token := val.(string)
		log.Print(token)
		log.Println("")

		username := GetUsernameFromToken(token)
		user, err := connection.UsernameIsInDb(username)
		if err != nil {
			log.Println("Error checking username in Db")
			log.Print(err)
		}

		log.Println("User is known in Db")
		log.Print(user)

		log.Println("Time to check if token matches user")
		log.Print(user.Tokens)
		log.Println("")
		for _, t := range user.Tokens {
			log.Println("Checking token:")
			log.Println(t)
			if strings.TrimSpace(token) == strings.TrimSpace(t) {
				log.Println("Token matches. Ok to continue")
				next(w, r)
			}
		}

		log.Println("Token is not Ok. Back to Login")
		// connection.GetAdmin(w, r)
		next(w, r)

	}

	// No cookie in store. Redirect to login
	// connection.GetLogin(w, r)

	http.Redirect(w, r, "/login", http.StatusFound)
	next(w, r)
}
