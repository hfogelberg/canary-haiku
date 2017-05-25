package main

import (
	"log"
	"net/http"
)

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("*** Auth Middleware ***")

	next(w, r)
}
