package handlers

import (
	"log"
	"net/http"
)

// AdminDashboardHandler zeigt Form-Links und optional eine Status-Meldung.
func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Nur GET zulassen
	if r.Method != http.MethodGet {
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Query-Param auslesen
	//success := r.URL.Query().Get("success") == "true"

	data2 := struct {
		SuccessMessage string
	}{}

	// Daten an das Template reichen
	//data := struct {
	//	Success bool
	//}{
	//	Success: success,
	//}
	// Erfolgsmeldung aus Query-Parameter übernehmen:
	if r.URL.Query().Get("success") == "true" {
		data2.SuccessMessage = "Der Welpe wurde erfolgreich hinzugefügt!"
	}

	// Admin-Dashboard im Admin-Layout rendern
	if err := renderAdminTemplate(w, "admin/admin_dashboard.html", data2); err != nil {
		log.Printf("Fehler beim Rendern des Admin-Dashboards: %v", err)
		http.Error(w, "Fehler beim Rendern des Admin-Dashboards", http.StatusInternalServerError)
	}
}
