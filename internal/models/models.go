package models

type Dog struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Born   string `json:"born"`
	Gender string `json:"gender"`
}

type Buyer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type Purchase struct {
	Date  string `json:"date"`
	Price int    `json:"price"`
}

type Puppy struct {
	Name        string
	Birthdate   string
	Gender      string
	Description string
	Images      []string
}

type Welpe struct {
	ID           string    `json:"id"`           // Eindeutige ID (z. B. UUID)
	Name         string    `json:"name"`         // Name des Welpen
	Farbe        Fellfarbe `json:"farbe"`        // Fellfarbe
	Geburtsdatum string    `json:"geburtsdatum"` // Geburtsdatum NORMALNO time.Time nicht string
	Geschlecht   string    `json:"geschlecht"`   // "männlich" oder "weiblich"
	Gewicht      float64   `json:"gewicht"`      // Gewicht in kg
	Charakter    string    `json:"charakter"`    // z. B. "verspielt", "ruhig"
	Geimpft      bool      `json:"geimpft"`      // Impfstatus
	Gechippt     bool      `json:"gechippt"`     // Chip vorhanden
	Entwurmt     bool      `json:"entwurmt"`     // Entwurmung erfolgt
	BildURL      string    `json:"bildUrl"`      // Link zu einem Bild
	Eltern       []string  `json:"eltern"`       // IDs der Elterntiere (z. B. Gandalf, Anna)
	Notizen      string    `json:"notizen"`      // Freitext für Besonderheiten
	Bilder       []string  `json:"bilder"`       // Liste von Bild-URLs
}

// Fellfarbe für Labrador Retriever
// Diese Typdefinition ermöglicht es, nur vordefinierte Farben zu verwenden
type Fellfarbe string

const (
	FellfarbeUnbekannt Fellfarbe = ""
	FarbeSchwarz       Fellfarbe = "schwarz"
	FarbeGelb          Fellfarbe = "gelb"
	FarbeBraun         Fellfarbe = "braun"
	FarbeFoxRed        Fellfarbe = "fox red" // Variante von Gelb
	FarbeSilber        Fellfarbe = "silber"  // nicht FCI-anerkannt
	FarbeChampagner    Fellfarbe = "champagner"
	FarbeCharcoal      Fellfarbe = "charcoal" // verdünntes Schwarz
)

// IstGueltigeFarbe prüft, ob die gegebene Fellfarbe gültig ist
func IstGueltigeFarbe(f Fellfarbe) bool {
	switch f {
	case FarbeSchwarz, FarbeGelb, FarbeBraun, FarbeFoxRed, FarbeSilber, FarbeChampagner, FarbeCharcoal:
		return true
	default:
		return false
	}
}
