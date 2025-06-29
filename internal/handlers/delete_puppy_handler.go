package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func DeletePuppyHandler(w http.ResponseWriter, r *http.Request, driver neo4j.DriverWithContext) {
	if r.Method != http.MethodPost {
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID fehlt", http.StatusBadRequest)
		return
	}

	// Neo4j-Session
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	// Cypher-Abfrage zum Löschen des Welpen
	_, err := session.Run(context.Background(),
		"MATCH (p:Puppy {id: $id}) DETACH DELETE p",
		map[string]interface{}{"id": id},
	)
	if err != nil {
		log.Printf("Fehler beim Löschen des Welpen: %v", err)
		http.Error(w, "Fehler beim Löschen des Welpen", http.StatusInternalServerError)
		return
	}

	// Optional: Erfolgsmeldung setzen
	// Zurück zur Übersicht
	http.Redirect(w, r, "/admin/puppies", http.StatusSeeOther)
}
