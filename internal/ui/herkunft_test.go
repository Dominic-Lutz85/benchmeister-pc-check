package ui

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/pruefung"
)

// Diese Tests halten die Abwehr gegen DNS-Rebinding fest, siehe die
// ausfuehrliche Begruendung in herkunft.go. Sie sind bewusst so
// geschrieben, dass der entscheidende Fall (fremder Host-Name) sofort
// erkennbar ist.
func TestHerkunftErlaubt(t *testing.T) {
	const port = "54321"

	faelle := []struct {
		name     string
		host     string
		herkunft string
		erwartet bool
	}{
		{
			name:     "eigene Seite, wie das Programm sie oeffnet",
			host:     "127.0.0.1:54321",
			herkunft: "http://127.0.0.1:54321",
			erwartet: true,
		},
		{
			name:     "eigene Seite ohne Origin-Kopf",
			host:     "127.0.0.1:54321",
			herkunft: "",
			erwartet: true,
		},
		{
			name:     "von Hand als localhost eingetippt",
			host:     "localhost:54321",
			herkunft: "http://localhost:54321",
			erwartet: true,
		},
		{
			// DER Fall, um den es geht: Bei DNS-Rebinding zeigt der
			// Domainname des Angreifers auf 127.0.0.1, im Host-Kopf steht
			// aber weiterhin sein Name.
			name:     "DNS-Rebinding: fremder Name zeigt auf 127.0.0.1",
			host:     "boese.example:54321",
			herkunft: "http://boese.example:54321",
			erwartet: false,
		},
		{
			name:     "fremde Seite ohne Origin-Kopf",
			host:     "boese.example:54321",
			herkunft: "",
			erwartet: false,
		},
		{
			name:     "richtiger Host, aber fremder Origin",
			host:     "127.0.0.1:54321",
			herkunft: "https://boese.example",
			erwartet: false,
		},
		{
			name:     "anderer Port auf demselben Rechner",
			host:     "127.0.0.1:9999",
			herkunft: "http://127.0.0.1:9999",
			erwartet: false,
		},
		{
			// Ein Name, der mit unserer Adresse anfaengt, aber eine andere
			// Domain ist. Faellt auf, wenn jemand die Pruefung spaeter zu
			// einem simplen Praefix-Vergleich vereinfacht.
			name:     "Name beginnt nur mit unserer Adresse",
			host:     "127.0.0.1.boese.example:54321",
			herkunft: "",
			erwartet: false,
		},
		{
			name:     "Host ohne Port",
			host:     "127.0.0.1",
			herkunft: "",
			erwartet: false,
		},
	}

	for _, f := range faelle {
		t.Run(f.name, func(t *testing.T) {
			anfrage := httptest.NewRequest(http.MethodPost, "/consent", nil)
			anfrage.Host = f.host
			if f.herkunft != "" {
				anfrage.Header.Set("Origin", f.herkunft)
			}

			if got := herkunftErlaubt(anfrage, port); got != f.erwartet {
				t.Errorf("herkunftErlaubt(Host=%q, Origin=%q) = %v, erwartet %v",
					f.host, f.herkunft, got, f.erwartet)
			}
		})
	}
}

// Der Schutz muss um ALLE Pfade liegen, nicht nur um /consent: Auf der
// Startseite steht das vollstaendige Scan-Ergebnis.
func TestNurEigeneHerkunftSchuetztAuchDieStartseite(t *testing.T) {
	const port = "54321"

	durchgelassen := false
	huelle := nurEigeneHerkunft(port, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) { durchgelassen = true }))

	for _, pfad := range []string{"/", "/consent"} {
		durchgelassen = false
		anfrage := httptest.NewRequest(http.MethodGet, pfad, nil)
		anfrage.Host = "boese.example:54321"
		schreiber := httptest.NewRecorder()

		huelle.ServeHTTP(schreiber, anfrage)

		if durchgelassen {
			t.Errorf("%s: fremde Herkunft wurde durchgelassen", pfad)
		}
		if schreiber.Code != http.StatusForbidden {
			t.Errorf("%s: Code %d, erwartet %d", pfad, schreiber.Code, http.StatusForbidden)
		}
	}
}

// Waechter fuer die Trennung der beiden Anzeigebloecke. Ueber dem einen
// steht "ohne dass du etwas kaufen musst", darunter darf also nichts
// landen, was auf einen Neukauf hinauslaeuft. Ein PCGH-Moderator hat
// genau diesen Widerspruch am 28.08.2026 gefunden.
func TestNachKostenTrenntZukaufAb(t *testing.T) {
	alle := []pruefung.Befund{
		{Schwere: pruefung.Hinweis, Titel: "Speicherprofil aus"},
		{Schwere: pruefung.Zukauf, Titel: "Magnetfestplatte verbaut"},
		{Schwere: pruefung.Anmerkung, Titel: "Riegel gemischt"},
		{Schwere: pruefung.Zukauf, Titel: "SSD haengt am SATA-Anschluss"},
	}
	kostenlos, kostenpflichtig := nachKosten(alle)

	if len(kostenlos) != 2 || len(kostenpflichtig) != 2 {
		t.Fatalf("erwartet 2 und 2, kam %d und %d", len(kostenlos), len(kostenpflichtig))
	}
	for _, b := range kostenlos {
		if b.Schwere == pruefung.Zukauf {
			t.Errorf("%q kostet Geld und darf nicht unter der Gratis-Ueberschrift stehen", b.Titel)
		}
	}
	for _, b := range kostenpflichtig {
		if b.Schwere != pruefung.Zukauf {
			t.Errorf("%q gehoert nicht in den Zukauf-Block", b.Titel)
		}
	}
}

// Ohne Befunde duerfen beide Bloecke leer bleiben, damit die Seite nicht
// zwei leere Karten zeigt.
func TestNachKostenOhneBefunde(t *testing.T) {
	kostenlos, kostenpflichtig := nachKosten(nil)
	if len(kostenlos) != 0 || len(kostenpflichtig) != 0 {
		t.Error("ohne Befunde muessen beide Listen leer sein")
	}
}
