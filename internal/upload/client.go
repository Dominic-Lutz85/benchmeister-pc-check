// Package upload überträgt ein Scan-Ergebnis an BenchMeister, aber nur
// nach ausdrücklicher Zustimmung. Das ist die EINZIGE Stelle im ganzen
// Programm, die überhaupt eine Netzwerkverbindung aufbaut.
//
// Wer prüfen will, ob dieses Programm heimlich Daten sendet, muss genau
// diese Datei lesen und danach im restlichen Quelltext nach "http" suchen.
// Außer dem lokalen Server auf 127.0.0.1 (siehe internal/ui) gibt es
// nichts weiter.
package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/scan"
)

// Die einzige Adresse, mit der dieses Programm spricht.
//
// SEIT 1.0.3 GEHT DAS NICHT MEHR AN DIE DATENBANK DIREKT. Vorher stand
// hier der öffentliche Datenbankschlüssel, und das Programm schrieb
// selbst in die Tabelle. Kein Geheimnisverrat, der Schlüssel ist bei
// Supabase für die Öffentlichkeit gedacht. Aber: Er steht im Quelltext
// dieses Programms, das öffentlich einsehbar ist. Wer ihn herauskopiert,
// konnte damit beliebig viele erfundene Rechner in die Statistik
// schreiben.
//
// Das wäre ausgerechnet bei diesen Daten der teuerste Schaden. Die
// Hardware-Statistik soll später Herstellern etwas belegen, und ein
// Bestand mit erfundenen Zeilen belegt nichts. Man sieht ihm hinterher
// auch nicht an, welche Zeilen echt waren.
//
// Jetzt nimmt die Website das Ergebnis entgegen, prüft es und begrenzt,
// wie viel in kurzer Zeit von derselben Stelle ankommt. Der Schlüssel
// zum Schreiben liegt nur noch dort auf dem Server. Für die Person am
// Rechner ändert sich nichts, außer dass Fehlermeldungen jetzt
// verständlich sind statt roher Datenbankantworten.
const (
	uploadURL        = "https://www.benchmeister.de/api/scan"
	ergebnisBasisURL = "https://benchmeister.de/pc-check-tool/ergebnis/"
)

// Anfrage bildet exakt die Parameter der Datenbank-Funktion
// submit_scan_result ab. Diese Struktur IST der vollständige Umfang
// dessen, was den Rechner verlässt, es gibt keinen zweiten Aufruf und
// keine zusätzlichen Felder.
//
// Beachte: Das Mainboard aus dem Scan-Ergebnis taucht hier bewusst NICHT
// auf. Es wird im Programm angezeigt, aber nicht übertragen.
type anfrage struct {
	CPUNameRaw string `json:"p_cpu_name_raw"`
	CPUCores   int    `json:"p_cpu_cores"`
	CPUThreads int    `json:"p_cpu_threads"`

	GPUNameRaw string `json:"p_gpu_name_raw"`
	GPUVramGb  *int   `json:"p_gpu_vram_gb"`

	RamTotalGb  int  `json:"p_ram_total_gb"`
	RamSpeedMhz *int `json:"p_ram_speed_mhz"`

	StorageType       string `json:"p_storage_type"`
	StorageCapacityGb int    `json:"p_storage_capacity_gb"`

	ResolutionWidth  int `json:"p_resolution_width"`
	ResolutionHeight int `json:"p_resolution_height"`

	ConsentShowResult     bool `json:"p_consent_show_result"`
	ConsentMarketResearch bool `json:"p_consent_market_research"`
}

// AlsJSON gibt die Nutzdaten so zurück, wie sie gesendet würden. Genau
// damit füttert die Oberfläche ihren aufklappbaren "Rohdaten"-Bereich:
// Was dort steht, ist nicht nachgebaut oder beschönigt, sondern derselbe
// Text, der später wirklich übertragen wird.
func AlsJSON(s *scan.ScanResult, marktforschung bool) (string, error) {
	daten, err := json.MarshalIndent(baueAnfrage(s, marktforschung), "", "  ")
	if err != nil {
		return "", err
	}
	return string(daten), nil
}

func baueAnfrage(s *scan.ScanResult, marktforschung bool) anfrage {
	return anfrage{
		CPUNameRaw:            s.CPUName,
		CPUCores:              s.CPUCores,
		CPUThreads:            s.CPUThreads,
		GPUNameRaw:            s.GPUName,
		GPUVramGb:             s.GPUVramGb,
		RamTotalGb:            s.RamTotalGb,
		RamSpeedMhz:           s.RamSpeedMhz,
		StorageType:           s.StorageType,
		StorageCapacityGb:     s.StorageCapacityGb,
		ResolutionWidth:       s.ResolutionWidth,
		ResolutionHeight:      s.ResolutionHeight,
		ConsentShowResult:     true, // ohne diese Zustimmung wird gar nicht erst gesendet
		ConsentMarketResearch: marktforschung,
	}
}

// Senden überträgt das Ergebnis und gibt die Adresse der Ergebnisseite
// zurück.
//
// marktforschung entspricht dem zweiten, freiwilligen Häkchen. Es ist ein
// eigener Parameter und kein Teil des Scan-Ergebnisses, damit im Code
// sichtbar bleibt, dass das eine bewusste Entscheidung der Nutzerin oder
// des Nutzers ist und kein Nebeneffekt.
func Senden(s *scan.ScanResult, marktforschung bool) (string, error) {
	return senden(uploadURL, s, marktforschung)
}

// senden nimmt die Adresse als Parameter, damit der Test sie auf einen
// eigenen Server umlenken kann. Sie bleibt trotzdem eine Konstante mit
// genau einem Aufrufer, siehe Senden oben: Wer wissen will, wohin dieses
// Programm sendet, findet weiterhin genau eine Adresse im Quelltext.
func senden(adresse string, s *scan.ScanResult, marktforschung bool) (string, error) {
	rumpf, err := json.Marshal(baueAnfrage(s, marktforschung))
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

	// Die Website antwortet immer mit einem JSON-Objekt: bei Erfolg mit
	// "token", sonst mit "fehler" in verständlichem Deutsch.
	var antwortDaten struct {
		Token  string `json:"token"`
		Fehler string `json:"fehler"`
	}
	_ = json.Unmarshal(inhalt, &antwortDaten)

	if antwort.StatusCode < 200 || antwort.StatusCode >= 300 {
		// Den Text der Website weiterreichen, wenn es einen gibt. Er sagt
		// der Person, was zu tun ist ("bitte in ein paar Minuten noch
		// einmal"), während eine nackte Statusnummer niemandem hilft.
		if antwortDaten.Fehler != "" {
			return "", fmt.Errorf("%s", antwortDaten.Fehler)
		}
		return "", fmt.Errorf(
			"BenchMeister hat die Daten abgelehnt (Status %d): %s",
			antwort.StatusCode, strings.TrimSpace(string(inhalt)),
		)
	}

	if antwortDaten.Token == "" {
		return "", fmt.Errorf("BenchMeister hat keinen Ergebnis-Schlüssel zurückgegeben")
	}

	return ergebnisBasisURL + antwortDaten.Token, nil
}
