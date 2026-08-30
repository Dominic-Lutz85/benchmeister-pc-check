// Prüft, wie das Programm mit den Antworten der Website umgeht.
//
// Das ist die einzige Stelle, die überhaupt ins Netz spricht, und sie
// wurde in 1.0.3 umgestellt: vorher direkt an die Datenbank, jetzt an
// benchmeister.de/api/scan. Damit hat sich auch das Antwortformat
// geändert, und genau da entstehen die Fehler, die niemand bemerkt, weil
// sie nur im Fehlerfall auftreten.
package upload

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/scan"
)

func beispielScan() *scan.ScanResult {
	vram := 12
	takt := 6000
	return &scan.ScanResult{
		CPUName:           "AMD Ryzen 7 7800X3D 8-Core Processor",
		CPUCores:          8,
		CPUThreads:        16,
		GPUName:           "NVIDIA GeForce RTX 4070",
		GPUVramGb:         &vram,
		RamTotalGb:        32,
		RamSpeedMhz:       &takt,
		StorageType:       "NVMe-SSD",
		StorageCapacityGb: 2000,
		ResolutionWidth:   2560,
		ResolutionHeight:  1440,
	}
}

func TestSendenGibtErgebnisadresseZurueck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"token":"e52e7455-57d1-4035-bfa4-adedba4482f3"}`)
	}))
	defer server.Close()

	adresse, err := senden(server.URL, beispielScan(), false)
	if err != nil {
		t.Fatalf("unerwarteter Fehler: %v", err)
	}
	erwartet := ergebnisBasisURL + "e52e7455-57d1-4035-bfa4-adedba4482f3"
	if adresse != erwartet {
		t.Errorf("Adresse = %q, erwartet %q", adresse, erwartet)
	}
}

func TestSendenReichtDenKlartextFehlerDurch(t *testing.T) {
	// Der wichtigste Fall. Wer sein Ergebnis absenden will und es klappt
	// nicht, muss lesen können, warum. Eine nackte Statusnummer hilft
	// niemandem, und "in ein paar Minuten noch einmal" ist ein Rat, den
	// man befolgen kann.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"fehler":"Zu viele Übertragungen in kurzer Zeit. Bitte in ein paar Minuten noch einmal versuchen."}`)
	}))
	defer server.Close()

	_, err := senden(server.URL, beispielScan(), false)
	if err == nil {
		t.Fatal("erwartet: Fehler")
	}
	if !strings.Contains(err.Error(), "in ein paar Minuten") {
		t.Errorf("Fehlertext der Website fehlt: %v", err)
	}
	// Die Statusnummer soll NICHT zusätzlich davorstehen, sie macht die
	// Meldung nur unverständlicher.
	if strings.Contains(err.Error(), "429") {
		t.Errorf("Statusnummer sollte nicht in der Meldung stehen: %v", err)
	}
}

func TestSendenBeiFehlerOhneTextNenntDenStatus(t *testing.T) {
	// Wenn die Website gar nichts Lesbares liefert (etwa eine
	// Fehlerseite des Betreibers), ist die Statusnummer besser als nichts.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = io.WriteString(w, "<html>Fehler</html>")
	}))
	defer server.Close()

	_, err := senden(server.URL, beispielScan(), false)
	if err == nil || !strings.Contains(err.Error(), "502") {
		t.Errorf("erwartet: Meldung mit Status 502, bekommen: %v", err)
	}
}

func TestSendenMeldetFehlendenSchluessel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	if _, err := senden(server.URL, beispielScan(), false); err == nil {
		t.Error("erwartet: Fehler, wenn kein Ergebnis-Schlüssel zurückkommt")
	}
}

func TestSendenSchicktKeinenDatenbankSchluesselMehr(t *testing.T) {
	// Seit 1.0.3 steckt kein Zugangsschlüssel mehr im Programm. Bliebe
	// versehentlich einer im Kopf der Anfrage, wäre er wieder öffentlich
	// auslesbar, und der ganze Umbau wäre umsonst.
	var kopfzeilen http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		kopfzeilen = r.Header.Clone()
		_, _ = io.WriteString(w, `{"token":"x"}`)
	}))
	defer server.Close()

	_, _ = senden(server.URL, beispielScan(), false)
	if kopfzeilen.Get("apikey") != "" || kopfzeilen.Get("Authorization") != "" {
		t.Errorf("Anfrage traegt noch Zugangsdaten: %v", kopfzeilen)
	}
}

func TestGesendetWirdNurWasDieVorschauZeigt(t *testing.T) {
	// Das Werkzeug zeigt vor dem Absenden einen Bereich "Rohdaten: genau
	// das wuerde uebertragen". Wenn der Rumpf davon abweicht, ist dieses
	// Versprechen gebrochen, und das ist der Kern des ganzen Programms.
	var rumpf []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rumpf, _ = io.ReadAll(r.Body)
		_, _ = io.WriteString(w, `{"token":"x"}`)
	}))
	defer server.Close()

	s := beispielScan()
	s.MainboardName = "ASRock B550M Pro4"

	vorschau, err := AlsJSON(s, true)
	if err != nil {
		t.Fatalf("Vorschau fehlgeschlagen: %v", err)
	}
	if _, err := senden(server.URL, s, true); err != nil {
		t.Fatalf("Senden fehlgeschlagen: %v", err)
	}

	var ausVorschau, gesendet map[string]any
	if err := json.Unmarshal([]byte(vorschau), &ausVorschau); err != nil {
		t.Fatalf("Vorschau ist kein JSON: %v", err)
	}
	if err := json.Unmarshal(rumpf, &gesendet); err != nil {
		t.Fatalf("Gesendetes ist kein JSON: %v", err)
	}
	if len(ausVorschau) != len(gesendet) {
		t.Fatalf("Vorschau hat %d Felder, gesendet wurden %d", len(ausVorschau), len(gesendet))
	}
	for feld, wert := range ausVorschau {
		if gesendet[feld] != wert {
			t.Errorf("Feld %s: Vorschau %v, gesendet %v", feld, wert, gesendet[feld])
		}
	}
	// Das Mainboard wird angezeigt, aber ausdruecklich nicht uebertragen.
	if strings.Contains(string(rumpf), "ASRock") {
		t.Error("Mainboard wurde uebertragen, darf es aber nicht")
	}
}

/*
 * Waechter ueber den Umfang der Uebertragung, angelegt am 30.08.2026.
 *
 * ANLASS: Beim Ergaenzen der Steckplatzzahl fiel auf, dass ein neues
 * Feld im ScanResult theoretisch mit hochgeladen werden koennte. Es
 * ging gut, weil baueAnfrage() seine Felder einzeln aufzaehlt und
 * nicht die ganze Struktur weiterreicht. Aber gutgegangen ist keine
 * Absicherung.
 *
 * WARUM DAS DER WICHTIGSTE TEST DES PROGRAMMS IST: Die zentrale Zusage
 * lautet, dass vorher woertlich auf dem Bildschirm steht, was
 * uebertragen wird. Ein Feld, das still dazukommt, bricht diese Zusage,
 * ohne dass es jemandem auffiele. Kein Absturz, keine Fehlermeldung,
 * nur eine Zeile mehr im Datenpaket.
 *
 * Der Test friert die Liste ein. Wer ein Feld ergaenzt, muss ihn
 * anfassen, und dabei faellt die Frage an: Steht das auch in der
 * Vorschau?
 */
func TestNurZugesagteFelderWerdenUebertragen(t *testing.T) {
	// Genau die Felder, die in der Vorschau woertlich genannt werden.
	// Diese Liste ist eine Zusage an die Nutzer, kein Implementierungs-
	// detail.
	// Das p_-Praefix stammt vom Datenbankaufruf submit_scan_result, die
	// Namen dahinter sind die aus der Vorschau.
	erlaubt := map[string]bool{
		"p_cpu_name_raw":            true,
		"p_cpu_cores":               true,
		"p_cpu_threads":             true,
		"p_gpu_name_raw":            true,
		"p_gpu_vram_gb":             true,
		"p_ram_total_gb":            true,
		"p_ram_speed_mhz":           true,
		"p_storage_type":            true,
		"p_storage_capacity_gb":     true,
		"p_resolution_width":        true,
		"p_resolution_height":       true,
		"p_consent_show_result":     true,
		"p_consent_market_research": true,
	}

	roh, err := json.Marshal(baueAnfrage(&scan.ScanResult{
		CPUName:         "AMD Ryzen 7 7800X3D",
		RamTotalGb:      32,
		RamSteckplaetze: 4, // rein oertlich, darf NICHT auftauchen
		Riegel:          []scan.RiegelInfo{{KapazitaetBytes: 1}},
		MainboardName:   "ASUS TUF B650-PLUS", // ebenfalls oertlich
	}, false))
	if err != nil {
		t.Fatalf("Anfrage laesst sich nicht serialisieren: %v", err)
	}

	var felder map[string]any
	if err := json.Unmarshal(roh, &felder); err != nil {
		t.Fatalf("Anfrage laesst sich nicht lesen: %v", err)
	}

	for name := range felder {
		if !erlaubt[name] {
			t.Errorf(
				"%q wird uebertragen, steht aber nicht in der Liste der "+
					"zugesagten Felder. Entweder gehoert es in die Vorschau "+
					"(preview.html.tmpl) und dann auch hier hinein, oder es "+
					"gehoert gar nicht in die Uebertragung.", name)
		}
	}

	// Die Gegenrichtung: Ein zugesagtes Feld darf nicht stillschweigend
	// verschwinden. Sonst stimmt die Vorschau auch nicht mehr.
	for name := range erlaubt {
		if _, da := felder[name]; !da {
			t.Errorf("%q ist zugesagt, fehlt aber in der Uebertragung", name)
		}
	}
}
