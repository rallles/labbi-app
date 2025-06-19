package handlers

import (
	"context"
	"log"
	"net/http"
	"net/smtp"

	"labbi-app/internal/config"
	"labbi-app/internal/database"
)

// ContactHandler zeigt das Kontaktformular (GET) und verarbeitet es (POST)
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Formular anzeigen
		renderTemplate(w, "contact.html", nil)
		return
	}

	// POST: Formular verarbeiten
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ungültige Eingaben", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	message := r.FormValue("message")

	// 1. Neo4j speichern
	cfg := config.LoadConfig()
	driver, err := database.NewNeo4jDriver(cfg)
	if err != nil {
		http.Error(w, "Datenbankverbindung fehlgeschlagen", http.StatusInternalServerError)
		log.Printf("Neo4j init error: %v", err)
		return
	}
	defer driver.Close(context.Background())

	session := driver.NewSession(context.Background(), database.DefaultSessionConfig())
	defer session.Close(context.Background())

	_, err = session.Run(context.Background(),
		"CREATE (c:Contact {name: $name, email: $email, phone: $phone, message: $message, ts: datetime()})",
		map[string]interface{}{
			"name":    name,
			"email":   email,
			"phone":   phone,
			"message": message,
		},
	)
	if err != nil {
		http.Error(w, "Fehler beim Speichern", http.StatusInternalServerError)
		log.Printf("Neo4j save error: %v", err)
		return
	}

	// 2. E-Mail senden (optional)
	errorMail := sendContactMail(name, email, phone, message)
	if errorMail != nil {
		log.Printf("E-Mail-Versand fehlgeschlagen: %v", errorMail)
	}

	// 3. Ergebnis anzeigen
	data := struct {
		Success bool
		Name    string
	}{
		Success: errorMail == nil,
		Name:    name,
	}
	renderTemplate(w, "contact_result.html", data)
}

// sendContactMail versendet eine Benachrichtigungs-E-Mail
func sendContactMail(name, email, phone, msg string) error {
	cfg := config.LoadConfig()
	from := cfg.Neo4jUser + "@example.com" // oder eigenes Absender-Email
	to := cfg.Neo4jUser + "@example.com"   // Ziel-Adresse
	auth := smtp.PlainAuth("", cfg.Neo4jUser, cfg.Neo4jPassword, "localhost")

	subject := "Neue Kontaktanfrage von " + name
	body := "Name: " + name + "\n" +
		"E-Mail: " + email + "\n" +
		"Telefon: " + phone + "\n\n" +
		"Nachricht:\n" + msg

	msgData := []byte("Subject: " + subject + "\r\n" +
		"\r\n" + body)

	return smtp.SendMail("localhost:25", auth, from, []string{to}, msgData)
}
