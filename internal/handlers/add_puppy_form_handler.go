package handlers

import (
	"log"
	"net/http"
)

// AddPuppyFormHandler rendert das Formular zum Anlegen eines neuen Welpen.
func AddPuppyFormHandler(w http.ResponseWriter, r *http.Request) {
	// Nur GET-Anfragen erlauben
	if r.Method != http.MethodGet {
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Formular-Template rendern
	// Erwarteter Pfad: templates/admin/add_puppy.html
	err := renderAdminTemplate(w, "admin/add_puppy.html", nil)
	if err != nil {
		log.Printf("Fehler beim Anzeigen des Add-Puppy-Formulars: %v", err)
	}
}
