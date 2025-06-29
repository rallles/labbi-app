package handlers

import (
	"context"
	"labbi-app/internal/config"
	"labbi-app/internal/database"
	"log"
	"net/http"

	"labbi-app/internal/models"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func ListPuppiesHandler(w http.ResponseWriter, r *http.Request) {
	// Neo4j-Treiber initialisieren
	cfg := config.LoadConfig()
	driver, err := database.NewNeo4jDriver(cfg)
	if err != nil {
		log.Printf("Neo4j-Driver konnte nicht initialisiert werden: %v", err)
		http.Error(w, "Serverfehler", http.StatusInternalServerError)
		return
	}
	defer driver.Close(context.Background())

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	result, err := session.Run(context.Background(),
		"MATCH (p:Puppy) RETURN p.name, p.birthdate, p.gender, p.description, p.images",
		nil)
	if err != nil {
		http.Error(w, "Fehler beim Laden der Daten", http.StatusInternalServerError)
		return
	}

	var puppies []models.Puppy

	for result.Next(context.Background()) {
		record := result.Record()
		imagesInterface, _ := record.Get("p.images")
		images := []string{}
		if imgList, ok := imagesInterface.([]interface{}); ok {
			for _, img := range imgList {
				if str, ok := img.(string); ok {
					images = append(images, str)
				}
			}
		}

		puppies = append(puppies, models.Puppy{
			Name:        record.Values[0].(string),
			Birthdate:   record.Values[1].(string),
			Gender:      record.Values[2].(string),
			Description: record.Values[3].(string),
			Images:      images,
		})
	}

	renderTemplate(w, "list-puppies.html", puppies)
}
