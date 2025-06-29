package handlers

import (
	"context"
	"fmt"
	"labbi-app/internal/models"
	"labbi-app/internal/utils"

	"net/http"
	"strconv"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func EditPuppyFormHandler(w http.ResponseWriter, r *http.Request, driver neo4j.DriverWithContext) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID fehlt", http.StatusBadRequest)
		return
	}

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	result, err := session.Run(context.Background(),
		`MATCH (p:Puppy {id: $id}) RETURN p`, map[string]interface{}{"id": id})
	if err != nil || !result.Next(context.Background()) {
		http.Error(w, "Welpe nicht gefunden", http.StatusNotFound)
		return
	}
	rec := result.Record()
	p := rec.Values[0].(neo4j.Node)
	props := p.Props

	// Baue ein Welpe-Struct für das Formular
	welpe := models.Welpe{
		ID:           props["id"].(string),
		Name:         props["name"].(string),
		Geburtsdatum: utils.DateToString(props["geburtsdatum"]),
		Geschlecht:   props["geschlecht"].(string),
		Farbe:        models.Fellfarbe(props["farbe"].(string)),
		Gewicht:      props["gewicht"].(float64),
		Charakter:    props["charakter"].(string),
		Geimpft:      props["geimpft"].(bool),
		Gechippt:     props["gechippt"].(bool),
		Entwurmt:     props["entwurmt"].(bool),
		Eltern:       toStringSlice(props["eltern"]),
		Notizen:      props["notizen"].(string),
		Bilder:       toStringSlice(props["bilder"]),
	}
	// (Hilfsfunktion `toStringSlice` siehe weiter unten.)

	// Edit-Formular anzeigen
	if err := renderAdminTemplate(w, "admin/admin_puppies_edit.html", welpe); err != nil {
		http.Error(w, "Fehler beim Anzeigen des Edit-Formulars", http.StatusInternalServerError)
	}
}

func toStringSlice(val interface{}) []string {
	var res []string
	if arr, ok := val.([]interface{}); ok {
		for _, v := range arr {
			res = append(res, fmt.Sprintf("%v", v))
		}
	}
	return res
}

func EditPuppySaveHandler(w http.ResponseWriter, r *http.Request, driver neo4j.DriverWithContext) {
	if r.Method != http.MethodPost {
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
		return
	}
	// Werte aus dem Formular holen
	id := r.FormValue("id")
	name := r.FormValue("name")
	geburtsdatum := r.FormValue("geburtsdatum")
	geschlecht := r.FormValue("geschlecht")
	farbe := r.FormValue("farbe")
	gewicht, _ := strconv.ParseFloat(r.FormValue("gewicht"), 64)
	charakter := r.FormValue("charakter")
	geimpft := r.FormValue("geimpft") == "on"
	gechippt := r.FormValue("gechippt") == "on"
	entwurmt := r.FormValue("entwurmt") == "on"
	eltern := strings.Split(r.FormValue("eltern"), ",")
	for i := range eltern {
		eltern[i] = strings.TrimSpace(eltern[i])
	}
	notizen := r.FormValue("notizen")

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err := session.Run(context.Background(),
		`MATCH (p:Puppy {id: $id})
        SET p.name = $name,
            p.geburtsdatum = $geburtsdatum,
            p.geschlecht = $geschlecht,
            p.farbe = $farbe,
            p.gewicht = $gewicht,
            p.charakter = $charakter,
            p.geimpft = $geimpft,
            p.gechippt = $gechippt,
            p.entwurmt = $entwurmt,
            p.eltern = $eltern,
            p.notizen = $notizen`,
		map[string]interface{}{
			"id":           id,
			"name":         name,
			"geburtsdatum": geburtsdatum,
			"geschlecht":   geschlecht,
			"farbe":        farbe,
			"gewicht":      gewicht,
			"charakter":    charakter,
			"geimpft":      geimpft,
			"gechippt":     gechippt,
			"entwurmt":     entwurmt,
			"eltern":       eltern,
			"notizen":      notizen,
		},
	)
	if err != nil {
		http.Error(w, "Fehler beim Speichern", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/puppies", http.StatusSeeOther)
}
