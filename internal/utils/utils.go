package utils

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

// dateToString konvertiert verschiedene Datentypen in einen String im Format "YYYY-MM-DD".
func DateToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case dbtype.Date:
		return v.Time().Format("2006-01-02")
	case time.Time:
		return v.Format("2006-01-02")
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

var AdminTemplates *template.Template

func InitAdminTemplates() {
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	// Hole das aktuelle Arbeitsverzeichnis (wird meist cmd/)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fehler beim Bestimmen des Arbeitsverzeichnisses: %v", err)
	}

	// Gehe eins hoch ins Projektverzeichnis
	projectRoot := filepath.Dir(wd)

	// Baue absolute Pfade zu den Templates
	base := filepath.Join(projectRoot, "internal", "templates", "admin_base.html")
	edit := filepath.Join(projectRoot, "internal", "templates", "admin", "admin_puppies_edit.html")
	// Weitere Templates nach Bedarf

	AdminTemplates, err = template.New("admin_base.html").Funcs(funcMap).ParseFiles(
		base, edit,
	)
	if err != nil {
		log.Fatalf("Fehler beim Parsen der Admin-Templates: %v", err)
	}
}
