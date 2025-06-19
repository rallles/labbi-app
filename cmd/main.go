package main

import (
	"context"
	"labbi-app/internal/config"
	"labbi-app/internal/database"
	"labbi-app/internal/router"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// .env laden – ignoriert Fehler, wenn keine .env vorhanden ist
	_ = godotenv.Load()

	// 1. Konfiguration laden
	cfg := config.LoadConfig()
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ":8080" // Standardport
	}

	// 2. Neo4j-Driver initialisieren
	driver, err := database.NewNeo4jDriver(cfg)
	if err != nil {
		log.Fatalf("Neo4j-Driver konnte nicht initialisiert werden: %v", err)
	}
	defer driver.Close(context.Background())

	// 3. ServeMux und Routing aufsetzen
	mux := http.NewServeMux()
	router.SetupRoutes(mux, driver)

	// 4. Server starten
	log.Printf("Labbi-App läuft auf %s", cfg.ServerAddress)
	err = http.ListenAndServe(cfg.ServerAddress, mux)
	if err != nil {
		log.Fatalf("Server konnte nicht gestartet werden: %v", err)
	}
}
