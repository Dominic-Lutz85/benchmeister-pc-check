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
// erreichbar, und beendet sich nach der einen Antwort selbst. Zusätzlich
// nimmt er nur Anfragen an, die nachweislich von seiner eigenen Seite
// kommen, siehe herkunft.go.
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
	"strings"
	"sync"
	"time"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/pruefung"
	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/scan"
	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/upload"
)

//go:embed assets/preview.html.tmpl
var vorlagen embed.FS

// Falls jemand das Fenster einfach offen liegen lässt, soll das Programm
// nicht ewig weiterlaufen. Nach dieser Zeit ohne Entscheidung beendet es
// sich von selbst, ohne etwas gesendet zu haben.
const wartezeit = 15 * time.Minute

// Wie lange nach dem Schliessen der Seite gewartet wird, bevor sich das
// Programm beendet. Deckt den Fall ab, dass die Seite nur neu geladen
// wurde: Dann ist sie in dieser Zeit laengst wieder da. Grosszuegig
// bemessen, weil ein zu frueher Abbruch schlimmer waere als drei
// Sekunden Nachlauf.
const nachfrist = 3 * time.Second

type vorlagenDaten struct {
	Scan         *scan.ScanResult
	RohdatenJSON string
	// Befunde der Plausibilitaetspruefung. Rein lokal: Sie entstehen aus
	// Angaben, die gar nicht uebertragen werden. Dient in der Vorlage nur
	// noch der Frage "gab es ueberhaupt etwas".
	Befunde []pruefung.Befund
	// Dieselben Befunde, getrennt nach der Frage, ob die Behebung Geld
	// kostet. Sie stehen in der Anzeige unter verschiedenen Ueberschriften,
	// siehe Begruendung in preview.html.tmpl.
	Kostenlos       []pruefung.Befund
	Kostenpflichtig []pruefung.Befund
	// Alle Datentraeger als fertige Zeilen fuer die Hardware-Uebersicht.
	Laufwerke []LaufwerkZeile
}

// LaufwerkZeile ist ein Datentraeger, fertig beschriftet fuer die Anzeige.
type LaufwerkZeile struct {
	Beschriftung string
}

/*
 * Baut die Zeilen fuer die Laufwerks-Uebersicht.
 *
 * Ergaenzt am 29.08.2026 nach einer Rueckmeldung im PCGH-Forum: Dort
 * hatte jemand fuenf Datentraeger im Rechner und sah in der Uebersicht
 * nur einen, ausgerechnet die 3-TB-Festplatte. Grund war, dass die
 * Uebersicht bis dahin nur StorageType und StorageCapacityGb zeigte, und
 * die beschreiben allein den GROESSTEN Datentraeger (siehe
 * scan/scan.go). Die Befunde darunter sprachen sehr wohl von seinen
 * SSDs, das Programm schien sich also selbst zu widersprechen.
 *
 * Die Angaben lagen die ganze Zeit vor, sie wurden nur nicht gezeigt.
 */
func laufwerksZeilen(e *scan.ScanResult) []LaufwerkZeile {
	zeilen := make([]LaufwerkZeile, 0, len(e.Laufwerke))
	for _, l := range e.Laufwerke {
		name := strings.TrimSpace(l.Name)
		if name == "" {
			name = "Unbenanntes Laufwerk"
		}

		// Medienart() statt l.MedienArt: repariert NVMe-SSDs, deren Art
		// Windows als "unspecified" meldet. Siehe scan/types.go.
		var art string
		switch l.Medienart() {
		case 3:
			art = "Festplatte"
		case 4:
			art = "SSD"
		default:
			art = "Art unbekannt"
		}

		// Der Anschluss steht nur dann dabei, wenn er etwas aussagt. "SSD
		// am SATA" ist eine Information, "SSD am Bus 0" waere Rauschen.
		switch l.BusArt {
		case 11:
			art += " am SATA-Anschluss"
		case 17:
			art += " im M.2-Steckplatz (NVMe)"
		case 7:
			art = "USB-Laufwerk"
		}

		gb := l.Bytes / (1000 * 1000 * 1000)
		zeilen = append(zeilen, LaufwerkZeile{
			Beschriftung: fmt.Sprintf("%s, %d GB, %s", name, gb, art),
		})
	}
	return zeilen
}

// nachKosten teilt die Befunde in die beiden Anzeigebloecke.
func nachKosten(alle []pruefung.Befund) (kostenlos, kostenpflichtig []pruefung.Befund) {
	for _, b := range alle {
		if b.Schwere == pruefung.Zukauf {
			kostenpflichtig = append(kostenpflichtig, b)
			continue
		}
		kostenlos = append(kostenlos, b)
	}
	return kostenlos, kostenpflichtig
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

	// Die Pruefung laeuft rein rechnerisch auf bereits ausgelesenen
	// Angaben: kein zusaetzlicher Zugriff aufs System, keine Last, keine
	// neuen Rechte. Siehe internal/pruefung.
	befunde := pruefung.Alle(riegelFuerPruefung(ergebnis), laufwerkeFuerPruefung(ergebnis))

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
	// Fuer die Herkunftspruefung, siehe herkunft.go.
	_, port, err := net.SplitHostPort(lauscher.Addr().String())
	if err != nil {
		return "", fmt.Errorf("Port des lokalen Dienstes nicht lesbar: %w", err)
	}

	// ergebnisURL und sendeFehler werden im Handler geschrieben (eigene
	// Goroutine) und unten im Hauptablauf gelesen. Ohne Sperre ist das ein
	// Wettlauf: Beim Erfolgsweg sorgt zwar das Schliessen des Kanals fuer
	// eine saubere Reihenfolge, beim Fehlerweg passiert das aber gerade
	// NICHT, dort laeuft das Programm in die Wartezeit und liest den Wert
	// ohne jede Absprache. In der Praxis faellt das kaum auf, das Go-
	// Speichermodell garantiert es aber nicht, und "go test -race" wuerde
	// es anzeigen. Sicherheitsdurchsicht 27.08.2026.
	var (
		einmal      sync.Once
		fertig      = make(chan struct{})
		sperre      sync.Mutex
		ergebnisURL string
		sendeFehler error
	)
	beenden := func() { einmal.Do(func() { close(fertig) }) }

	/*
	 * Abschied mit Nachfrist, siehe /verlassen weiter unten.
	 *
	 * gehtVielleicht() startet den Countdown, bleibt() bricht ihn ab. Die
	 * Zaehlernummer verhindert einen Wettlauf: Ein Countdown beendet das
	 * Programm nur, wenn seither weder ein neuer Countdown gestartet noch
	 * die Seite neu geladen wurde. Ohne den Zaehler koennte ein alter,
	 * bereits abgebrochener Countdown das Programm noch abschiessen.
	 */
	var (
		abschiedSperre sync.Mutex
		abschiedNummer int
	)
	// Wie viele geoeffnete Seiten sich gerade melden, siehe /anwesend.
	var offeneSeiten int
	bleibt := func() {
		abschiedSperre.Lock()
		abschiedNummer++
		abschiedSperre.Unlock()
	}
	gehtVielleicht := func() {
		abschiedSperre.Lock()
		abschiedNummer++
		meine := abschiedNummer
		abschiedSperre.Unlock()
		go func() {
			time.Sleep(nachfrist)
			abschiedSperre.Lock()
			nochAktuell := abschiedNummer == meine
			abschiedSperre.Unlock()
			if nochAktuell {
				beenden()
			}
		}()
	}
	merken := func(url string, fehler error) {
		sperre.Lock()
		defer sperre.Unlock()
		ergebnisURL, sendeFehler = url, fehler
	}
	lesen := func() (string, error) {
		sperre.Lock()
		defer sperre.Unlock()
		return ergebnisURL, sendeFehler
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// Ein erneuter Aufruf der Seite hebt einen laufenden Abschied auf.
		// Das ist der Neuladen-Fall, siehe /verlassen.
		bleibt()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		kostenlos, kostenpflichtig := nachKosten(befunde)
		_ = vorlage.Execute(w, vorlagenDaten{
			Scan:            ergebnis,
			RohdatenJSON:    rohdaten,
			Befunde:         befunde,
			Kostenlos:       kostenlos,
			Kostenpflichtig: kostenpflichtig,
			Laufwerke:       laufwerksZeilen(ergebnis),
		})
	})

	/*
	 * Die Seite meldet sich ab, wenn sie geschlossen wird. Ohne das lief
	 * das Programm nach dem Schliessen des Browserfensters noch bis zu 15
	 * Minuten im Hintergrund weiter, mit offenem lokalem Server.
	 *
	 * Gefunden hat das ein Moderator im PCGH-Forum am 28.08.2026: "Warum
	 * bleibt der Task nach den schliessen vom Programm weiter aktiv?" Fuer
	 * ein Programm, das ausdruecklich damit wirbt, keinen Dienst zu
	 * hinterlassen, ist das der denkbar schlechteste Eindruck, und der
	 * Einwand ist vollkommen berechtigt: Aus Sicht der Benutzerin IST die
	 * Browserseite das Programm.
	 *
	 * Die Nachfrist ist noetig, weil pagehide auch beim Neuladen feuert.
	 * Kommt in dieser Zeit wieder ein Aufruf der Seite herein, war es ein
	 * Neuladen und der Abschied wird verworfen.
	 */
	mux.HandleFunc("/verlassen", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		gehtVielleicht()
	})

	// Gegenstueck zu /verlassen fuer den Zurueck-Knopf: Holt der Browser
	// die Seite aus seinem Vor-und-Zurueck-Speicher, gibt es keinen neuen
	// Aufruf von "/", der den Abschied abblasen koennte. Diese Meldung tut
	// es stattdessen.
	mux.HandleFunc("/bleibt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		bleibt()
	})

	/*
	 * Die eigentliche Anwesenheitsmeldung: eine Anfrage, die nie
	 * beantwortet wird. Solange die Seite offen ist, steht diese
	 * Verbindung. Wird der Reiter geschlossen, bricht sie ab, und der
	 * Server erfaehrt es im selben Moment ueber den Anfrage-Kontext.
	 *
	 * WARUM NICHT der naheliegende Weg. Zuerst stand hier nur ein
	 * sendBeacon im pagehide-Ereignis. Im Test am 28.08.2026 zeigte sich:
	 * Beim harten Schliessen eines Reiters kam die Meldung nicht an, das
	 * Programm lief weiter. Ein Herzschlag per Zeitgeber waere auch nichts
	 * geworden, den drosseln Browser in unsichtbaren Reitern auf einen
	 * Schlag pro Minute und frieren ihn spaeter ganz ein. Eine stehende
	 * Verbindung ist von beidem nicht betroffen.
	 *
	 * Der Zaehler ist noetig, weil die Seite auch mehrfach offen sein
	 * kann. Erst wenn die LETZTE Verbindung faellt, geht das Programm.
	 */
	mux.HandleFunc("/anwesend", func(w http.ResponseWriter, r *http.Request) {
		// Kopfzeilen sofort abschicken, damit der Browser die Verbindung
		// als stehend ansieht und nicht auf einen Rumpf wartet.
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if abspueler, ok := w.(http.Flusher); ok {
			abspueler.Flush()
		}

		abschiedSperre.Lock()
		offeneSeiten++
		abschiedNummer++ // zaehlt wie bleibt(): ein laufender Abschied ist hinfaellig
		abschiedSperre.Unlock()

		select {
		case <-r.Context().Done():
		case <-fertig:
			return
		}

		abschiedSperre.Lock()
		offeneSeiten--
		keineMehrDa := offeneSeiten <= 0
		abschiedSperre.Unlock()
		if keineMehrDa {
			gehtVielleicht()
		}
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
			merken("", err)
			w.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(w).Encode(antwort{Fehler: err.Error()})
			// Nicht beenden: Vielleicht war nur kurz das Netz weg, dann
			// soll ein zweiter Versuch möglich sein.
			return
		}

		merken(url, nil)
		_ = json.NewEncoder(w).Encode(antwort{URL: url})

		// Kurz warten, damit der Browser die Antwort noch bekommt und
		// weiterleiten kann, bevor der Server zumacht.
		go func() {
			time.Sleep(1500 * time.Millisecond)
			beenden()
		}()
	})

	// Ergaenzt am 27.08.2026 nach dem Sicherheitsaudit: Ohne diese Huelle
	// koennte eine fremde Webseite per DNS-Rebinding mit dem lokalen
	// Dienst sprechen, solange das Programm laeuft. Begruendung im Detail
	// in herkunft.go.
	server := &http.Server{Handler: nurEigeneHerkunft(port, mux)}
	go func() { _ = server.Serve(lauscher) }()

	imBrowserOeffnen(adresse)
	fmt.Println("Ergebnis im Browser geöffnet:", adresse)
	fmt.Println("Falls sich nichts öffnet, ruf die Adresse von Hand auf.")
	fmt.Println("Zum Beenden reicht es, die Seite im Browser zu schließen.")

	select {
	case <-fertig:
	case <-time.After(wartezeit):
		fmt.Println("Keine Entscheidung getroffen, es wurde nichts übertragen.")
	}

	abschluss, abbrechen := context.WithTimeout(context.Background(), 3*time.Second)
	defer abbrechen()
	_ = server.Shutdown(abschluss)

	url, fehler := lesen()
	if url == "" && fehler != nil {
		return "", fehler
	}
	return url, nil
}

// Öffnet eine Adresse im eingestellten Standardbrowser. "start" ist ein
// eingebauter Befehl der Windows-Eingabeaufforderung, deshalb der Umweg
// über cmd. Das leere Anführungszeichen-Paar ist nötig, weil start sonst
// die Adresse als Fenstertitel deuten würde.
func imBrowserOeffnen(adresse string) {
	_ = exec.Command("cmd", "/c", "start", "", adresse).Start()
}

// Die beiden Umwandler halten scan und pruefung voneinander unabhaengig:
// Das Auslese-Paket muss nichts ueber die Pruefung wissen und umgekehrt.
//
// ACHTUNG BEIM ERGAENZEN VON FELDERN. Genau hier ist der Fix von 1.0.5
// verendet: Das Auslesen holte ConfiguredClockSpeed korrekt, die
// Pruefung wertete es korrekt aus, aber dieser Umwandler dazwischen
// kopierte es nicht mit. TaktMhzZweiter blieb damit immer 0, der
// Rueckfall auf Speed griff jedes Mal, und das Programm meldete
// weiterhin den JEDEC-Grundtakt als Ist-Wert.
//
// Aufgefallen am 29.08.2026 an einem Screenshot von Misanthrop68 im
// PCGH-Forum: CMK32GX5M2B5200C40, Speed 4800, ConfiguredClockSpeed
// 5600. Sein Speicher lief also sogar ueber dem Profil, das Programm
// behauptete trotzdem 4800 und riet ihm, XMP einzuschalten.
//
// Die Lehre: Ein stiller Umwandler ist die gefaehrlichste Sorte Code.
// Er bricht nichts, er laesst nur etwas weg, und der Fehler sieht
// danach genau so aus wie vor der Reparatur. Der Test in
// server_test.go haelt diese Stelle jetzt fest.
func riegelFuerPruefung(e *scan.ScanResult) []pruefung.Riegel {
	aus := make([]pruefung.Riegel, 0, len(e.Riegel))
	for _, r := range e.Riegel {
		aus = append(aus, pruefung.Riegel{
			KapazitaetBytes: r.KapazitaetBytes,
			TaktMhz:         r.TaktMhz,
			TaktMhzZweiter:  r.TaktMhzZweiter,
			Teilenummer:     r.Teilenummer,
			Kanal:           r.Kanal,
		})
	}
	return aus
}

func laufwerkeFuerPruefung(e *scan.ScanResult) []pruefung.Laufwerk {
	aus := make([]pruefung.Laufwerk, 0, len(e.Laufwerke))
	for _, l := range e.Laufwerke {
		aus = append(aus, pruefung.Laufwerk{
			Name: l.Name,
			// Medienart() statt l.MedienArt, damit die Pruefung dieselbe
			// Art sieht wie die Anzeige darueber. Sonst haette eine als
			// "unspecified" gemeldete SATA-SSD in der Uebersicht anders
			// dagestanden als im Befund.
			MedienArt: l.Medienart(),
			BusArt:    l.BusArt,
			Bytes:     l.Bytes,
		})
	}
	return aus
}
