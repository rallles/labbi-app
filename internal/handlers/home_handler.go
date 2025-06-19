package handlers

import (
	"log"
	"net/http"
)

// HomeHandler rendert die Startseite.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Startseite aufgerufen")
	renderTemplate(w, "index.html", nil)
}
