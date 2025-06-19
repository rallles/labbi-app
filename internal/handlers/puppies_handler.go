package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// PuppiesHandler rendert die Welpen-Seite.
func PuppiesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Welpen-Seite aufgerufen")
	renderTemplate(w, "puppies.html", nil)
}

// MakePuppiesHandler erstellt einen HTTP-Handler für die Welpen-Seite.
func MakePuppiesHandler(driver neo4j.DriverWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Neue Session anlegen
		session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
		defer session.Close(context.Background())

		// 2. Cypher-Abfrage
		result, err := session.Run(context.Background(),
			"MATCH (p:Puppy) RETURN p.name AS name, p.images AS images", nil)
		if err != nil {
			http.Error(w, "DB-Fehler", http.StatusInternalServerError)
			return
		}

		// 3. Aufbereiten und Rendern
		var puppies []struct {
			Name   string
			Images []string
		}
		for result.Next(context.Background()) {
			rec := result.Record()
			name, _ := rec.Get("name")
			images, _ := rec.Get("images")
			var imgs []string
			if n, ok := name.(string); ok {
				if ims, ok := images.([]string); ok {
					imgs = ims
				}
				puppies = append(puppies, struct {
					Name   string
					Images []string
				}{Name: n, Images: imgs})
			}
		}
		renderTemplate(w, "puppies.html", puppies)
	}
}
