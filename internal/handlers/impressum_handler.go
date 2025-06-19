package handlers

import (
	"log"
	"net/http"
)

// ImpressumHandler rendert die Impressum-Seite.
func ImpressumHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Impressum-Seite aufgerufen")
	renderTemplate(w, "impressum.html", nil)
}
