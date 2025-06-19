package handlers

import (
	"log"
	"net/http"
)

// AboutHandler rendert die "Über uns"-Seite.
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("About-Seite aufgerufen")
	renderTemplate(w, "about.html", nil)
}
