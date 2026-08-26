// Package ui zeigt das Scan-Ergebnis lokal an und nimmt die Zustimmung
// entgegen.
//
// Warum eine HTML-Seite im normalen Browser statt eines eigenen
// Programmfensters: Ein Fenster-Baukasten für Go wäre ein zweiter
// Oberflächen-Baukasten neben dem der Website, dauerhaft zu pflegen, und
// würde das Programm deutlich größer machen. Vor allem aber ist eine
// HTML-Seite für jede:n einsehbar (Rechtsklick, Seitenquelltext anzeigen),
// ein eigenes Fenster wäre für Nicht-Fachleute eine Blackbox. Bei einem
// Programm, dessen ganzer Sinn Vertrauen ist, wiegt das schwerer als die
// paar Zeilen Server-Code hier.
//
// Der Server hört ausschließlich auf 127.0.0.1, ist also von außen nicht
// erreichbar, und beendet sich nach der einen Antwort selbst.
package ui

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/TwerkiTwerk/benchmeister-pc-check/internal/scan"
	"github.com/TwerkiTwerk/benchmeister-pc-check/internal/upload"
)

//go:embed assets/preview.html.tmpl
var vorlagen embed.FS

// Falls jemand das Fenster einfach offen liegen lässt, soll das Programm
// nicht ewig weiterlaufen. Nach dieser Zeit ohne Entscheidung beendet es
// sich von selbst, ohne etwas gesendet zu haben.
const wartezeit = 15 * time.Minute

type vorlagenDaten struct {
	Scan         *scan.ScanResult
	RohdatenJSON string
}

type zustimmung struct {
	Anzeigen       bool `json:"anzeigen"`
	Marktforschung bool `json:"marktforschung"`
}

type antwort struct {
	URL    string `json:"url,omitempty"`
	Fehler string `json:"fehler,omitempty"`
}

// Anzeigen startet den lokalen Server, öffnet die Seite im Standardbrowser
// und wartet auf die Entscheidung. Gibt die Adresse der Ergebnisseite
// zurück, oder einen leeren Text, wenn nichts übertragen wurde.
func Anzeigen(ergebnis *scan.ScanResult) (string, error) {
	rohdaten, err := upload.AlsJSON(ergebnis, false)
	if err != nil {
		return "", err
	}

	vorlage, err := template.ParseFS(vorlagen, "assets/preview.html.tmpl")
	if err != nil {
		return "", err
	}

	// Port 0 heißt: Betriebssystem sucht einen freien Port aus. So kollidiert
	// das Programm nie mit etwas, das schon läuft.
	lauscher, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("lokaler Anzeigedienst konnte nicht starten: %w", err)
	}
	adresse := fmt.Sprintf("http://%s", lauscher.Addr().String())

	var (
		einmal      sync.Once
		fertig      = make(chan struct{})
		ergebnisURL string
		sendeFehler error
	)
	beenden := func() { einmal.Do(func() { close(fertig) }) }

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = vorlage.Execute(w, vorlagenDaten{Scan: ergebnis, RohdatenJSON: rohdaten})
	})

	mux.HandleFunc("/consent", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(antwort{Fehler: "falsche Anfrageart"})
			return
		}

		var z zustimmung
		if err := json.NewDecoder(r.Body).Decode(&z); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(antwort{Fehler: "Anfrage nicht lesbar"})
			return
		}

		// Zweite Sperre neben der Oberfläche: Ohne die erste Zustimmung
		// wird hier nichts gesendet, egal was der Browser schickt. Die
		// dritte Sperre sitzt in der Datenbank selbst.
		if !z.Anzeigen {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(antwort{Fehler: "ohne Zustimmung wird nichts übertragen"})
			return
		}

		url, err := upload.Senden(ergebnis, z.Marktforschung)
		if err != nil {
			sendeFehler = err
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(antwort{Fehler: err.Error()})
			// Nicht beenden: Vielleicht war nur kurz das Netz weg, dann
			// soll ein zweiter Versuch möglich sein.
			return
		}

		ergebnisURL = url
		_ = json.NewEncoder(w).Encode(antwort{URL: url})

		// Kurz warten, damit der Browser die Antwort noch bekommt und
		// weiterleiten kann, bevor der Server zumacht.
		go func() {
			time.Sleep(1500 * time.Millisecond)
			beenden()
		}()
	})

	server := &http.Server{Handler: mux}
	go func() { _ = server.Serve(lauscher) }()

	imBrowserOeffnen(adresse)
	fmt.Println("Ergebnis im Browser geöffnet:", adresse)
	fmt.Println("Falls sich nichts öffnet, ruf die Adresse von Hand auf.")

	select {
	case <-fertig:
	case <-time.After(wartezeit):
		fmt.Println("Keine Entscheidung getroffen, es wurde nichts übertragen.")
	}

	abschluss, abbrechen := context.WithTimeout(context.Background(), 3*time.Second)
	defer abbrechen()
	_ = server.Shutdown(abschluss)

	if ergebnisURL == "" && sendeFehler != nil {
		return "", sendeFehler
	}
	return ergebnisURL, nil
}

// Öffnet eine Adresse im eingestellten Standardbrowser. "start" ist ein
// eingebauter Befehl der Windows-Eingabeaufforderung, deshalb der Umweg
// über cmd. Das leere Anführungszeichen-Paar ist nötig, weil start sonst
// die Adresse als Fenstertitel deuten würde.
func imBrowserOeffnen(adresse string) {
	_ = exec.Command("cmd", "/c", "start", "", adresse).Start()
}
