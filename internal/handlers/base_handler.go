package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
		filepath.Join(wd, "templates"),                           // direkt unter dem Arbeitsverzeichnis
		filepath.Join(wd, "internal", "templates"),               // direkt unter cwd/internal/templates
		filepath.Join(filepath.Dir(wd), "templates"),             // eine Ebene höher: Projekt-Root/templates
		filepath.Join(filepath.Dir(wd), "internal", "templates"), // eine Ebene höher: Projekt-Root/internal/templates
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

// renderAdminTemplate lädt admin_base.html und das spezifizierte Admin-Template
// und sucht den Templates-Ordner in gängigen Projektstrukturen.
func renderAdminTemplate(w http.ResponseWriter, page string, data interface{}) error {
	// Arbeitsverzeichnis ermitteln
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Arbeitsverzeichnis konnte nicht ermittelt werden: %v", err)
		http.Error(w, "Serverfehler", http.StatusInternalServerError)
		return err
	}

	// Mögliche Template-Pfade relativ zu wd
	dirs := []string{
		filepath.Join(wd, "templates"),                           // ./templates
		filepath.Join(wd, "internal", "templates"),               // ./internal/templates
		filepath.Join(filepath.Dir(wd), "templates"),             // ../templates
		filepath.Join(filepath.Dir(wd), "internal", "templates"), // ../internal/templates
	}

	var tplDir string
	for _, d := range dirs {
		if info, statErr := os.Stat(d); statErr == nil && info.IsDir() {
			tplDir = d
			break
		}
	}
	if tplDir == "" {
		log.Printf("Admin-Templates nicht gefunden in: %v", dirs)
		http.Error(w, "Templates nicht gefunden", http.StatusInternalServerError)
		return os.ErrNotExist
	}

	// Pfade zu admin_base.html und Seiten-Template
	basePath := filepath.Join(tplDir, "admin_base.html")
	pagePath := filepath.Join(tplDir, page) // z.B. "admin/add_puppy.html"

	// --- HIER: Template mit FuncMap vorbereiten ---
	funcMap := template.FuncMap{
		"join": strings.Join,
		// weitere Funktionen nach Bedarf
	}

	tmpl, err := template.New("admin_base.html").Funcs(funcMap).ParseFiles(basePath, pagePath)
	if err != nil {
		log.Printf("Fehler beim Parsen der Admin-Templates (%s & %s): %v", basePath, pagePath, err)
		http.Error(w, "Fehler beim Laden der Seite", http.StatusInternalServerError)
		return err
	}

	// Render admin_base.html mit dem content-Block aus page
	err = tmpl.ExecuteTemplate(w, "admin_base.html", data)
	if err != nil {
		log.Printf("Fehler beim Rendern der Admin-Seite %s: %v", page, err)
		http.Error(w, "Fehler beim Rendern der Seite", http.StatusInternalServerError)
	}
	return err
}
