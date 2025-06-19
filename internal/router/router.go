package router

import (
	"labbi-app/internal/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func SetupRoutes(mux *http.ServeMux, driver neo4j.DriverWithContext) {
	// 1) Spezifische Seiten zuerst
	mux.HandleFunc("/about", handlers.AboutHandler)
	mux.HandleFunc("/dogs", handlers.DogsHandler)
	mux.HandleFunc("/puppies", handlers.PuppiesHandler)
	//mux.HandleFunc("/list-puppies", handlers.ListPuppiesHandler)
	//mux.HandleFunc("/admin",    handlers.AdminDashboard)
	mux.HandleFunc("/contact", handlers.ContactHandler)
	mux.HandleFunc("/impressum", handlers.ImpressumHandler)

	// 2) Statische Dateien
	// 1) Arbeitsverzeichnis ermitteln
	// wir starten in cmd/, also eine Ebene hoch
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Arbeitsverzeichnis: %v", err)
	}
	projectRoot := filepath.Dir(wd)
	staticDir := filepath.Join(projectRoot, "static")
	// 2) Statische Dateien
	mux.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))),
	)

	// 3) Nur exakt "/" → HomeHandler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			handlers.HomeHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}
