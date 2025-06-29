package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// AddPuppyHandler verarbeitet das Admin-Formular (POST) und speichert den neuen Welpen in Neo4j.
func AddPuppyHandler(w http.ResponseWriter, r *http.Request, driver neo4j.DriverWithContext) {
	// Nur POST zulassen
	if r.Method != http.MethodPost {
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Multipart-Form parsen (max. 20 MB im RAM)
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		log.Printf("Fehler beim Lesen des Formulars: %v", err)
		http.Error(w, "Fehler beim Lesen des Formulars", http.StatusBadRequest)
		return
	}

	// Formularwerte auslesen
	name := r.FormValue("name")
	birthdate := r.FormValue("geburtsdatum")
	gender := r.FormValue("geschlecht")
	color := r.FormValue("farbe")
	weight := r.FormValue("gewicht")
	character := r.FormValue("charakter")
	vaccinated := r.FormValue("geimpft") == "true"
	chipped := r.FormValue("gechippt") == "true"
	dewormed := r.FormValue("entwurmung") == "true"
	notes := r.FormValue("notizen")

	// Eltern (Checkboxen können mehrfach übergeben werden)
	parents := r.Form["eltern"]

	// Bilder verarbeiten
	files := r.MultipartForm.File["images"]
	imagePaths, err := saveUploadedImages(files)
	if err != nil {
		log.Printf("Fehler beim Speichern der Bilder: %v", err)
		http.Error(w, "Fehler beim Speichern der Bilder", http.StatusInternalServerError)
		return
	}

	// Neo4j-Verbindung und Insert
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	puppyID := uuid.NewString()
	// Welpenknoten erstellen
	_, err = session.Run(context.Background(),
		`CREATE (p:Puppy {
		    id: $id,
		    name: $name,
		    geburtsdatum: date($birthdate),
		    geschlecht: $gender,
		    farbe: $color,
		    gewicht: toFloat($weight),
		    charakter: $character,
		    geimpft: $vaccinated,
		    gechippt: $chipped,
		    entwurmt: $dewormed,
		    eltern: $parents,
		    notizen: $notes,
		    bilder: $images
		})`, map[string]interface{}{
			"id":         puppyID,
			"name":       name,
			"birthdate":  birthdate,
			"gender":     gender,
			"color":      color,
			"weight":     weight,
			"character":  character,
			"vaccinated": vaccinated,
			"chipped":    chipped,
			"dewormed":   dewormed,
			"parents":    parents,
			"notes":      notes,
			"images":     imagePaths,
		})

	// Elternbeziehungen anlegen, falls vorhanden
	if len(parents) > 0 {
		for _, parent := range parents {
			_, err = session.Run(context.Background(),
				`MATCH (p:Puppy {id: $puppyID}), (parent:Puppy {id: $parentID})
				CREATE (p)-[:HAS_PARENT]->(parent)`,
				map[string]interface{}{
					"puppyID":  puppyID,
					"parentID": parent,
				})
		}
	}
	// Fehler beim Speichern in Neo4j
	if err != nil {
		log.Printf("Fehler beim Speichern in Neo4j: %v", err)
		http.Error(w, "Fehler beim Speichern", http.StatusInternalServerError)
		return
	}

	// Bei Erfolg zurück zum Dashboard
	http.Redirect(w, r, "/admin?success=true", http.StatusSeeOther)
}

// saveUploadedImages speichert alle hochgeladenen Dateien und liefert ihre Pfade zurück
func saveUploadedImages(files []*multipart.FileHeader) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// Gehe vom Arbeitsverzeichnis (z.B. .../labbi-app/cmd) eine Ebene nach oben zum Projekt-Stammverzeichnis:
	projectRoot := filepath.Dir(wd)
	imageDir := filepath.Join(projectRoot, "static", "images")
	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		return nil, err
	}

	var paths []string
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		ext := filepath.Ext(fh.Filename)
		name := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().UnixNano(), ext)
		target := filepath.Join(imageDir, name)

		out, err := os.Create(target)
		if err != nil {
			return nil, err
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			return nil, err
		}
		paths = append(paths, name)
	}
	return paths, nil
}
