// internal/handlers/admin_puppies_table.go
package handlers

import (
	"context"
	"fmt"
	"labbi-app/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Hilfsfunktionen zur Typumwandlung
func toString(val interface{}) string {
	if val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", val)
}

func toFloat(val interface{}) float64 {
	if val == nil {
		return 0.0
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int64:
		return float64(v)
	case int:
		return float64(v)
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	}
	return 0.0
}

func toBool(val interface{}) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	if s, ok := val.(string); ok {
		return s == "true"
	}
	return false
}

func toDateString(val interface{}) string {
	// Für Template: "2006-01-02"
	if val == nil {
		return ""
	}
	switch t := val.(type) {
	case time.Time:
		return t.Format("2006-01-02")
	case string:
		return t
	}
	return fmt.Sprintf("%v", val)
}

func toFellfarbe(val interface{}) models.Fellfarbe {
	if val == nil {
		return models.FellfarbeUnbekannt
	}
	if s, ok := val.(string); ok {
		return models.Fellfarbe(s)
	}
	return models.Fellfarbe(fmt.Sprintf("%v", val))
}

// Haupt-Handler
func ListPuppiesAdminHandler(w http.ResponseWriter, r *http.Request, driver neo4j.DriverWithContext) {
	if r.Method != http.MethodGet {
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
		return
	}

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	cypher := `
        MATCH (p:Puppy)
        RETURN p.id AS id,
               p.name AS name,
               p.geburtsdatum AS geburtsdatum,
               p.geschlecht AS geschlecht,
               p.farbe AS farbe,
               p.gewicht AS gewicht,
               p.charakter AS charakter,
               p.geimpft AS geimpft,
               p.gechippt AS gechippt,
               p.entwurmt AS entwurmt,
               p.eltern AS eltern,
               p.notizen AS notizen,
               p.bilder AS bilder
        ORDER BY p.geburtsdatum DESC
    `
	result, err := session.Run(context.Background(), cypher, nil)
	if err != nil {
		log.Printf("Fehler beim Abfragen der Welpen: %v", err)
		http.Error(w, "Fehler beim Laden der Daten", http.StatusInternalServerError)
		return
	}

	var puppies []models.Welpe
	for result.Next(context.Background()) {
		rec := result.Record()

		// Eltern als []string
		parents := []string{}
		if vals, ok := rec.Get("eltern"); ok && vals != nil {
			if slice, ok2 := vals.([]interface{}); ok2 {
				for _, v := range slice {
					parents = append(parents, toString(v))
				}
			}
		}

		// Bilder als []string
		images := []string{}
		if vals, ok := rec.Get("bilder"); ok && vals != nil {
			if slice, ok2 := vals.([]interface{}); ok2 {
				for _, v := range slice {
					images = append(images, toString(v))
				}
			}
		}

		idVal, _ := rec.Get("id")
		nameVal, _ := rec.Get("name")
		geburtsdatumVal, _ := rec.Get("geburtsdatum")
		geschlechtVal, _ := rec.Get("geschlecht")
		farbeVal, _ := rec.Get("farbe")
		gewichtVal, _ := rec.Get("gewicht")
		charakterVal, _ := rec.Get("charakter")
		geimpftVal, _ := rec.Get("geimpft")
		gechipptVal, _ := rec.Get("gechippt")
		entwurmtVal, _ := rec.Get("entwurmt")
		notizenVal, _ := rec.Get("notizen")

		puppy := models.Welpe{
			ID:           toString(idVal),
			Name:         toString(nameVal),
			Geburtsdatum: toDateString(geburtsdatumVal),
			Geschlecht:   toString(geschlechtVal),
			Farbe:        toFellfarbe(farbeVal),
			Gewicht:      toFloat(gewichtVal),
			Charakter:    toString(charakterVal),
			Geimpft:      toBool(geimpftVal),
			Gechippt:     toBool(gechipptVal),
			Entwurmt:     toBool(entwurmtVal),
			Eltern:       parents,
			Notizen:      toString(notizenVal),
			Bilder:       images,
		}
		puppies = append(puppies, puppy)
	}
	if err = result.Err(); err != nil {
		log.Printf("Fehler beim Auslesen der Welpen: %v", err)
		http.Error(w, "Fehler beim Auslesen der Daten", http.StatusInternalServerError)
		return
	}

	if err := renderAdminTemplate(w, "admin/admin_puppies_table.html", puppies); err != nil {
		log.Printf("Fehler beim Rendern der Welpen-Liste: %v", err)
		http.Error(w, "Interner Fehler beim Rendern", http.StatusInternalServerError)
	}
}
