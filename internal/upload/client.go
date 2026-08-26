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

// Zugangsdaten der öffentlichen BenchMeister-Datenbank. Das ist KEIN
// Geheimnis: Dieser Schlüssel steckt genauso im Browser-Code der Website
// und ist bei Supabase ausdrücklich dafür gedacht, öffentlich zu sein. Der
// Schutz kommt nicht aus Geheimhaltung, sondern aus den Regeln in der
// Datenbank selbst: Mit diesem Schlüssel lässt sich ein Ergebnis anlegen
// und genau ein Ergebnis über seinen Token abrufen, sonst nichts. Ein
// Auslesen der Tabelle ist damit nicht möglich.
const (
	supabaseURL     = "https://hdmyjymgpwwcphtmgkjp.supabase.co"
	supabaseAnonKey = "sb_publishable_ReU8DKrHUdgIq8gzDP2vSg_226RUak5"
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
	rumpf, err := json.Marshal(baueAnfrage(s, marktforschung))
	if err != nil {
		return "", err
	}

	anfrageObjekt, err := http.NewRequest(
		http.MethodPost,
		supabaseURL+"/rest/v1/rpc/submit_scan_result",
		bytes.NewReader(rumpf),
	)
	if err != nil {
		return "", err
	}
	anfrageObjekt.Header.Set("apikey", supabaseAnonKey)
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

	if antwort.StatusCode < 200 || antwort.StatusCode >= 300 {
		return "", fmt.Errorf(
			"BenchMeister hat die Daten abgelehnt (Status %d): %s",
			antwort.StatusCode, strings.TrimSpace(string(inhalt)),
		)
	}

	// Die Funktion gibt den Token als JSON-Zeichenkette zurück,
	// einschließlich Anführungszeichen.
	var token string
	if err := json.Unmarshal(inhalt, &token); err != nil {
		return "", fmt.Errorf("unerwartete Antwort von BenchMeister: %s", string(inhalt))
	}
	if token == "" {
		return "", fmt.Errorf("BenchMeister hat keinen Ergebnis-Schlüssel zurückgegeben")
	}

	return ergebnisBasisURL + token, nil
}
