package handlers

import (
	"log"
	"net/http"
)

// HomeHandler rendert die Startseite.
func DogsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Startseite aufgerufen")
	renderTemplate(w, "dogs.html", nil)
}
