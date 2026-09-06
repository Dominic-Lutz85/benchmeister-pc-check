package ui

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/last"
	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/scan"
	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/upload"
)

/*
Die Endpunkte für Lasttest und Ranglisten-Eintrag, angelegt am 06.09.2026.

WARUM ABFRAGEN STATT EINER LANGEN ANTWORT
=========================================
Der naheliegende Weg wäre ein einziger Aufruf, der zwei Minuten offen
bleibt und am Ende das Ergebnis liefert. Er wäre kürzer und wäre falsch:
Eine Seite, auf der zwei Minuten lang nichts passiert, sieht kaputt aus,
und abbrechen ließe sich auch nichts.

Deshalb drei kleine Endpunkte. Der Browser startet, fragt im Sekundentakt
den Stand ab und kann jederzeit abbrechen. Das kostet ein paar Zeilen
mehr und ist der Unterschied zwischen einem Vorgang, den man sieht, und
einer Sanduhr.

WER DEN LASTTEST STARTEN DARF
=============================
Nur wer vorher zugestimmt hat, und die Zustimmung wird hier geprüft und
nicht nur in der Oberfläche. Dieselbe Regel wie bei /consent: Was
Wirkung auf dem Rechner oder im Netz hat, wird an jeder Stelle geprüft,
an der es vorbeikommt.
*/

// lastZustand hält den laufenden oder letzten Lasttest.
//
// Ein Programmlauf, ein Test: Wer zweimal startet, bricht den ersten ab.
// Zwei gleichzeitige Lasttests würden sich gegenseitig die Rechenzeit
// wegnehmen und beide Messungen wertlos machen.
type lastZustand struct {
	mu      sync.Mutex
	laeuft  bool
	beginn  time.Time
	dauer   time.Duration
	abbruch context.CancelFunc
	// Der zuletzt gemessene Takt, waehrend der Lauf noch laeuft. Ohne
	// ihn zeigt die Oberflaeche zwei Minuten lang Nullen, denn das
	// Ergebnis entsteht erst am Ende.
	letzterMhz int
	ergebnis   *last.Ergebnis
	fehler     string
}

type standAntwort struct {
	Laeuft      bool    `json:"laeuft"`
	Sekunden    int     `json:"sekunden"`
	Gesamt      int     `json:"gesamt"`
	AktuellMhz  int     `json:"aktuell_mhz"`
	SpitzeMhz   int     `json:"spitze_mhz"`
	DauerMhz    int     `json:"dauer_mhz"`
	Gehalten    float64 `json:"gehalten"`
	Urteil      string  `json:"urteil"`
	Abgebrochen bool    `json:"abgebrochen"`
	Fertig      bool    `json:"fertig"`
	Fehler      string  `json:"fehler"`
}

/** Zwei Minuten. Lang genug, dass eine schwache Kühlung sichtbar wird,
 *  kurz genug, dass niemand daneben sitzt und wartet. */
const lastdauer = 120 * time.Second

func (z *lastZustand) starten(nenntaktMhz int) {
	z.mu.Lock()
	// Ein zweiter Start bricht den ersten ab, statt danebenzulaufen.
	if z.abbruch != nil {
		z.abbruch()
	}
	ctx, abbrechen := context.WithCancel(context.Background())
	z.abbruch = abbrechen
	z.laeuft = true
	z.beginn = time.Now()
	z.dauer = lastdauer
	z.ergebnis = nil
	z.letzterMhz = 0
	z.fehler = ""
	z.mu.Unlock()

	go func() {
		ergebnis, err := last.AusfuehrenMit(ctx, nenntaktMhz, lastdauer,
			func(p last.Messpunkt) {
				z.mu.Lock()
				z.letzterMhz = p.Mhz
				z.mu.Unlock()
			})
		z.mu.Lock()
		defer z.mu.Unlock()
		z.laeuft = false
		if err != nil {
			z.fehler = err.Error()
			return
		}
		z.ergebnis = ergebnis
	}()
}

func (z *lastZustand) stand() standAntwort {
	z.mu.Lock()
	defer z.mu.Unlock()

	a := standAntwort{
		Laeuft:     z.laeuft,
		Gesamt:     int(z.dauer.Seconds()),
		Fehler:     z.fehler,
		Sekunden:   int(time.Since(z.beginn).Seconds()),
		AktuellMhz: z.letzterMhz,
	}
	if z.beginn.IsZero() {
		a.Sekunden = 0
	}

	if z.ergebnis != nil {
		e := z.ergebnis
		a.Fertig = true
		a.SpitzeMhz = e.SpitzeMhz
		a.DauerMhz = e.DauerMhz
		a.Gehalten = e.GehaltenProzent
		a.Urteil = string(last.Einordnen(e.GehaltenProzent))
		a.Abgebrochen = e.Abgebrochen

	}
	return a
}

func (z *lastZustand) abbrechen() {
	z.mu.Lock()
	defer z.mu.Unlock()
	if z.abbruch != nil {
		z.abbruch()
	}
}

/*
registriereLasttest hängt die vier Endpunkte an den Server.

ergebnis liefert Prozessorname, Grafikkarte und Nenntakt, die alle drei
schon ausgelesen sind. Der Lasttest fragt nichts Neues am System ab, er
rechnet nur.
*/
func registriereLasttest(mux *http.ServeMux, ergebnis *scan.ScanResult) {
	zustand := &lastZustand{}
	betriebssystem := scan.Betriebssystem()

	schreibe := func(w http.ResponseWriter, wert any, code int) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if code != http.StatusOK {
			w.WriteHeader(code)
		}
		_ = json.NewEncoder(w).Encode(wert)
	}

	mux.HandleFunc("/lasttest/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			schreibe(w, antwort{Fehler: "falsche Anfrageart"}, http.StatusMethodNotAllowed)
			return
		}

		var wunsch struct {
			Zugestimmt bool `json:"zugestimmt"`
		}
		if err := json.NewDecoder(r.Body).Decode(&wunsch); err != nil {
			schreibe(w, antwort{Fehler: "Anfrage nicht lesbar"}, http.StatusBadRequest)
			return
		}
		// Zweite Sperre neben der Oberfläche. Ein Lasttest ist das
		// Einzige, was dieses Programm auf dem Rechner überhaupt
		// anrichtet, also wird die Zustimmung auch hier geprüft.
		if !wunsch.Zugestimmt {
			schreibe(w, antwort{Fehler: "ohne Zustimmung wird nicht gemessen"}, http.StatusBadRequest)
			return
		}
		if ergebnis.CPUNenntaktMhz <= 0 {
			schreibe(w, antwort{
				Fehler: "Windows meldet für diesen Prozessor keinen Nenntakt. Ohne ihn lässt sich der gemessene Takt nicht einordnen.",
			}, http.StatusUnprocessableEntity)
			return
		}

		zustand.starten(ergebnis.CPUNenntaktMhz)
		schreibe(w, zustand.stand(), http.StatusOK)
	})

	mux.HandleFunc("/lasttest/stand", func(w http.ResponseWriter, _ *http.Request) {
		schreibe(w, zustand.stand(), http.StatusOK)
	})

	mux.HandleFunc("/lasttest/abbrechen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			schreibe(w, antwort{Fehler: "falsche Anfrageart"}, http.StatusMethodNotAllowed)
			return
		}
		zustand.abbrechen()
		schreibe(w, zustand.stand(), http.StatusOK)
	})

	mux.HandleFunc("/rangliste/eintragen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			schreibe(w, antwort{Fehler: "falsche Anfrageart"}, http.StatusMethodNotAllowed)
			return
		}

		var wunsch struct {
			Alias      string `json:"alias"`
			Zugestimmt bool   `json:"zugestimmt"`
		}
		if err := json.NewDecoder(r.Body).Decode(&wunsch); err != nil {
			schreibe(w, antwort{Fehler: "Anfrage nicht lesbar"}, http.StatusBadRequest)
			return
		}
		if !wunsch.Zugestimmt {
			schreibe(w, antwort{Fehler: "ohne Zustimmung wird nichts eingetragen"}, http.StatusBadRequest)
			return
		}

		zustand.mu.Lock()
		e := zustand.ergebnis
		zustand.mu.Unlock()

		// Ohne Messung gibt es nichts einzutragen. Das kann nur passieren,
		// wenn jemand die Oberfläche umgeht, aber dann soll es einen
		// klaren Satz geben statt einer leeren Zeile in der Liste.
		if e == nil || e.SpitzeMhz == 0 {
			schreibe(w, antwort{Fehler: "Es liegt keine fertige Messung vor."}, http.StatusBadRequest)
			return
		}
		if e.Abgebrochen {
			schreibe(w, antwort{
				Fehler: "Die Messung wurde abgebrochen. Eine abgebrochene Messung sagt nichts über Dauerlast aus.",
			}, http.StatusBadRequest)
			return
		}

		schluessel, err := upload.EintragSenden(upload.Messung{
			Alias:             wunsch.Alias,
			CPUName:           ergebnis.CPUName,
			GPUName:           ergebnis.GPUName,
			Betriebssystem:    betriebssystem,
			NenntaktMhz:       ergebnis.CPUNenntaktMhz,
			SpitzeMhz:         e.SpitzeMhz,
			DauerMhz:          e.DauerMhz,
			MessdauerSekunden: int(lastdauer.Seconds()),
			Einwilligung:      true,
		})
		if err != nil {
			// Der Satz der Website geht unverändert weiter, siehe
			// upload/rangliste.go. "Der Name ist schon vergeben" ist
			// keine Panne, sondern eine Entscheidung für den Menschen.
			schreibe(w, antwort{Fehler: err.Error()}, http.StatusBadGateway)
			return
		}

		schreibe(w, map[string]string{"loeschschluessel": schluessel}, http.StatusOK)
	})
}
