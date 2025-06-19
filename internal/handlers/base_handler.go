package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// renderTemplate lädt base.html und ein Seiten-Template aus dem richtigen Verzeichnis und rendert es.
func renderTemplate(w http.ResponseWriter, page string, data interface{}) error {
	// Aktuelles Arbeitsverzeichnis ermitteln (in Docker: /app)
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Arbeitsverzeichnis konnte nicht ermittelt werden: %v", err)
		http.Error(w, "Serverfehler", http.StatusInternalServerError)
		return err
	}

	// Mögliche Template-Verzeichnisse relativ zu wd
	dirs := []string{
		filepath.Join(wd, "templates"),
		filepath.Join(wd, "internal", "templates"),
	}

	var tplDir string
	for _, d := range dirs {
		if info, statErr := os.Stat(d); statErr == nil && info.IsDir() {
			tplDir = d
			break
		}
	}
	if tplDir == "" {
		log.Printf("Templates nicht gefunden in %v", dirs)
		http.Error(w, "Templates nicht gefunden", http.StatusInternalServerError)
		return os.ErrNotExist
	}

	// Pfade zu base.html und Seiten-Template
	basePath := filepath.Join(tplDir, "base.html")
	pagePath := filepath.Join(tplDir, page)

	// Parse base.html und Seiten-Template
	tmpl, err := template.ParseFiles(basePath, pagePath)
	if err != nil {
		log.Printf("Fehler beim Parsen der Templates: %v", err)
		http.Error(w, "Fehler beim Laden der Seite", http.StatusInternalServerError)
		return err
	}

	// Render base.html mit den Block-Definitionen aus dem Seiten-Template
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Fehler beim Rendern von %s: %v", page, err)
		http.Error(w, "Fehler beim Rendern der Seite", http.StatusInternalServerError)
	}
	return err
}
