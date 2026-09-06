package upload

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func beispielMessung() Messung {
	return Messung{
		Alias:             "Tueftler",
		CPUName:           "AMD Ryzen 7 5800X3D",
		GPUName:           "NVIDIA GeForce RTX 5060 Ti",
		Betriebssystem:    "Windows 10 Pro 19045",
		NenntaktMhz:       3801,
		SpitzeMhz:         4694,
		DauerMhz:          4673,
		MessdauerSekunden: 120,
		Einwilligung:      true,
	}
}

/*
Der wichtigste Test dieser Datei.

Was den Rechner verlässt, muss genau das sein, was auf dem Bildschirm
stand, und nichts weiter. Ein Feld, das sich hier einschleicht, fällt
sonst niemandem auf, denn die Übertragung ist unsichtbar. Der Test
zählt deshalb die Schlüssel im JSON ab, statt nur einzelne zu prüfen.
*/
func TestNurDieseFelderVerlassenDenRechner(t *testing.T) {
	rohdaten, err := json.Marshal(beispielMessung())
	if err != nil {
		t.Fatalf("laesst sich nicht verpacken: %v", err)
	}

	var felder map[string]any
	if err := json.Unmarshal(rohdaten, &felder); err != nil {
		t.Fatalf("laesst sich nicht auspacken: %v", err)
	}

	erwartet := map[string]bool{
		"p_alias": true, "p_cpu_name": true, "p_gpu_name": true,
		"p_betriebssystem": true, "p_nenntakt_mhz": true, "p_spitze_mhz": true,
		"p_dauer_mhz": true, "p_messdauer_sekunden": true, "p_einwilligung": true,
	}

	for name := range felder {
		if !erwartet[name] {
			t.Errorf("unerwartetes Feld geht hinaus: %q", name)
		}
	}
	for name := range erwartet {
		if _, da := felder[name]; !da {
			t.Errorf("Feld fehlt: %q", name)
		}
	}

	// Die Zahl, nach der die Liste sortiert, wird auf dem Server
	// gerechnet. Ginge sie hier hinaus, koennte man sich seinen Platz
	// aussuchen.
	if _, da := felder["p_gehalten_prozent"]; da {
		t.Error("der gehaltene Anteil darf nicht mitgeschickt werden")
	}
}

func TestLoeschschluesselKommtZurueck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("falscher Content-Type: %q", r.Header.Get("Content-Type"))
		}
		rumpf, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(rumpf), "Tueftler") {
			t.Errorf("der Name fehlt im Rumpf: %s", rumpf)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"loeschschluessel":"11111111-2222-3333-4444-555555555555"}`))
	}))
	defer server.Close()

	schluessel, err := eintragSenden(server.URL, beispielMessung())
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	if schluessel != "11111111-2222-3333-4444-555555555555" {
		t.Errorf("falscher Schluessel: %q", schluessel)
	}
}

/*
Der vergebene Name ist kein Fehler des Programms, sondern eine
Entscheidung, die der Mensch treffen muss. Der Satz der Website muss
deshalb unveraendert bei ihm ankommen. Ein "Status 409" waere die
schlechteste aller Antworten: Er sieht aus wie ein Defekt.
*/
func TestVergebenerNameKommtImKlartextAn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"fehler":"Der Name „Tueftler“ ist schon vergeben. Such dir bitte einen anderen aus."}`))
	}))
	defer server.Close()

	_, err := eintragSenden(server.URL, beispielMessung())
	if err == nil {
		t.Fatal("ein vergebener Name muss einen Fehler ergeben")
	}
	if !strings.Contains(err.Error(), "schon vergeben") {
		t.Errorf("der Satz der Website kam nicht an: %v", err)
	}
	if strings.Contains(err.Error(), "409") {
		t.Errorf("die Statusnummer steht in der Meldung: %v", err)
	}
}

func TestAntwortOhneSchluesselIstEinFehler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	if _, err := eintragSenden(server.URL, beispielMessung()); err == nil {
		t.Error("ohne Loeschschluessel darf das nicht als Erfolg gelten")
	}
}
