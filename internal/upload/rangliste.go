package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

/*
Sendeweg für die Rangliste, angelegt am 06.09.2026.

WARUM EIN ZWEITER WEG UND NICHT DER BESTEHENDE
==============================================
Der Scan-Weg oben überträgt, was verbaut ist, und bekommt einen Link auf
eine Ergebnisseite zurück. Hier geht eine Messung hinaus und ein
Löschschlüssel kommt zurück, und beides gehört zu einer anderen
Entscheidung des Menschen davor: Beim Scan sagt er "zeig mir mein
Ergebnis", hier sagt er "trag mich in eine öffentliche Liste ein".

Zwei Entscheidungen, zwei Wege. Sie in einen Aufruf zu legen hieße, dass
eine Zustimmung für beides gilt, und genau das soll sie nicht.

WAS HIER NICHT MITGESCHICKT WIRD: der gehaltene Anteil in Prozent. Das
ist die Zahl, nach der die Liste sortiert, also die einzige, die zu
fälschen sich lohnt. Sie wird auf dem Server aus Spitzen- und Dauertakt
gerechnet. Wir schicken Messwerte, keine Platzierung.
*/

const ranglisteURL = "https://www.benchmeister.de/api/rangliste"

// Messung ist genau das, was den Rechner für die Rangliste verlässt.
// Diese Struktur IST der vollständige Umfang, es gibt kein zweites Feld
// und keinen zweiten Aufruf.
type Messung struct {
	Alias          string `json:"p_alias"`
	CPUName        string `json:"p_cpu_name"`
	GPUName        string `json:"p_gpu_name"`
	Betriebssystem string `json:"p_betriebssystem"`

	NenntaktMhz       int `json:"p_nenntakt_mhz"`
	SpitzeMhz         int `json:"p_spitze_mhz"`
	DauerMhz          int `json:"p_dauer_mhz"`
	MessdauerSekunden int `json:"p_messdauer_sekunden"`

	// Ohne diese Zustimmung lehnt schon die Datenbank ab. Sie steht hier
	// als Feld und nicht als stiller Vorgabewert, damit sie beim Aufruf
	// sichtbar gesetzt werden muss.
	Einwilligung bool `json:"p_einwilligung"`
}

// EintragSenden trägt eine Messung in die Rangliste ein und gibt den
// Löschschlüssel zurück.
//
// Der Schlüssel ist der einzige Weg, den Eintrag später wieder
// loszuwerden. Wer ihn verliert, kommt an seinen Eintrag nicht mehr
// heran, deshalb muss der Aufrufer ihn anzeigen und nicht verschlucken.
func EintragSenden(m Messung) (loeschschluessel string, err error) {
	return eintragSenden(ranglisteURL, m)
}

// eintragSenden nimmt die Adresse als Parameter, damit der Test sie auf
// einen eigenen Server umlenken kann. Dieselbe Aufteilung wie bei
// Senden/senden oben, aus demselben Grund: Im Quelltext steht genau eine
// echte Adresse.
func eintragSenden(adresse string, m Messung) (string, error) {
	rumpf, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	anfrageObjekt, err := http.NewRequest(
		http.MethodPost,
		adresse,
		bytes.NewReader(rumpf),
	)
	if err != nil {
		return "", err
	}
	anfrageObjekt.Header.Set("Content-Type", "application/json")

	klient := &http.Client{Timeout: 20 * time.Second}
	antwort, err := klient.Do(anfrageObjekt)
	if err != nil {
		return "", fmt.Errorf("keine Verbindung zu BenchMeister: %w", err)
	}
	defer antwort.Body.Close()

	inhalt, err := io.ReadAll(io.LimitReader(antwort.Body, 4096))
	if err != nil {
		return "", err
	}

	var antwortDaten struct {
		Loeschschluessel string `json:"loeschschluessel"`
		Fehler           string `json:"fehler"`
	}
	_ = json.Unmarshal(inhalt, &antwortDaten)

	if antwort.StatusCode < 200 || antwort.StatusCode >= 300 {
		/*
		 * Der Text der Website geht unverändert weiter. Er ist auf diesen
		 * Fall zugeschnitten: "Der Name X ist schon vergeben" sagt genau,
		 * was zu tun ist, während "Status 409" niemandem hilft. Gerade
		 * beim Namen ist das der Unterschied zwischen "such dir einen
		 * anderen aus" und "das Programm ist kaputt".
		 */
		if antwortDaten.Fehler != "" {
			return "", fmt.Errorf("%s", antwortDaten.Fehler)
		}
		return "", fmt.Errorf(
			"BenchMeister hat den Eintrag abgelehnt (Status %d): %s",
			antwort.StatusCode, strings.TrimSpace(string(inhalt)),
		)
	}

	if antwortDaten.Loeschschluessel == "" {
		return "", fmt.Errorf("BenchMeister hat keinen Löschschlüssel zurückgegeben")
	}
	return antwortDaten.Loeschschluessel, nil
}
