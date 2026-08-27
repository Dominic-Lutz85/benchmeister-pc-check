package pruefung

import "testing"

// Die Teilenummern hier sind echt und stammen vom Entwicklungsrechner
// (ausgelesen am 27.08.2026), nicht ausgedacht. Genau an ihnen hat sich
// gezeigt, dass Windows die Sollgeschwindigkeit nicht meldet und sie nur
// in der Teilenummer steht.
func TestSollTaktAusTeilenummer(t *testing.T) {
	faelle := []struct {
		teilenummer string
		erwartet    int
		warum       string
	}{
		{"CMK32GX4M2E3200C16", 3200, "Corsair DDR4-3200, echte Teilenummer"},
		{"F4-3200C16-8GVKB", 3200, "G.Skill DDR4-3200, echte Teilenummer"},
		{"CMK32GX5M2B6000C36", 6000, "Corsair DDR5-6000"},
		{"KF556C36-16", 0, "kein vierstelliger Takt erkennbar"},
		{"", 0, "leere Teilenummer"},
		{"BELIEBIG1234XYZ", 0, "1234 ist keine bekannte Taktstufe"},
	}
	for _, f := range faelle {
		if got := sollTaktAusTeilenummer(f.teilenummer); got != f.erwartet {
			t.Errorf("%s: sollTaktAusTeilenummer(%q) = %d, erwartet %d",
				f.warum, f.teilenummer, got, f.erwartet)
		}
	}
}

// Waechter gegen zu gierige Erkennung: Enthaelt eine Teilenummer mehrere
// bekannte Taktstufen, ist keine davon sicher. Dann lieber gar kein
// Befund als ein falscher.
func TestMehrdeutigeTeilenummerGibtNichts(t *testing.T) {
	if got := sollTaktAusTeilenummer("XX3200YY3600ZZ"); got != 0 {
		t.Errorf("mehrdeutige Teilenummer sollte 0 ergeben, ergab %d", got)
	}
}

func TestSpeicherErkenntAbgeschaltetesProfil(t *testing.T) {
	// Genau die Bestueckung des Entwicklungsrechners: vier Riegel,
	// alle fuer 3200 gebaut, alle mit 2133 laufend.
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "P0 CHANNEL A"},
		{KapazitaetBytes: 8589934592, TaktMhz: 2133, Teilenummer: "F4-3200C16-8GVKB", Kanal: "P0 CHANNEL A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "P0 CHANNEL B"},
		{KapazitaetBytes: 8589934592, TaktMhz: 2133, Teilenummer: "F4-3200C16-8GVKB", Kanal: "P0 CHANNEL B"},
	}
	befunde := Speicher(riegel)

	var profil, gemischt bool
	for _, b := range befunde {
		if b.Titel == "Arbeitsspeicher läuft unter seiner Sollgeschwindigkeit" {
			profil = true
			if b.Schwere != Hinweis {
				t.Error("abgeschaltetes Speicherprofil muss ein Hinweis sein")
			}
		}
		if b.Titel == "Verschiedene Speicherriegel gemischt" {
			gemischt = true
		}
	}
	if !profil {
		t.Error("das abgeschaltete Speicherprofil wurde nicht erkannt")
	}
	if !gemischt {
		t.Error("die gemischten Riegel wurden nicht erkannt")
	}
}

func TestSpeicherSchweigtWennAllesPasst(t *testing.T) {
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "B"},
	}
	if befunde := Speicher(riegel); len(befunde) != 0 {
		t.Errorf("bei einwandfreier Bestueckung darf nichts gemeldet werden, kam: %+v", befunde)
	}
}

func TestSpeicherErkenntEinzelnenRiegel(t *testing.T) {
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "A"},
	}
	befunde := Speicher(riegel)
	if len(befunde) != 1 || befunde[0].Titel != "Nur ein Speicherriegel verbaut" {
		t.Errorf("Einkanalbetrieb wurde nicht erkannt, kam: %+v", befunde)
	}
}

func TestSpeicherOhneTeilenummerMeldetKeinProfil(t *testing.T) {
	// Manche Riegel melden gar keine Teilenummer. Dann ist der Soll-Wert
	// unbekannt, und es darf nichts behauptet werden.
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, Teilenummer: "", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, Teilenummer: "", Kanal: "B"},
	}
	if befunde := Speicher(riegel); len(befunde) != 0 {
		t.Errorf("ohne Teilenummer darf nichts gemeldet werden, kam: %+v", befunde)
	}
}

func TestLaufwerke(t *testing.T) {
	// Bestueckung des Entwicklungsrechners plus zwei erfundene Faelle.
	laufwerke := []Laufwerk{
		{Name: "VendorCo ProductCode", MedienArt: 0, BusArt: busUSB, Bytes: 62914560000},
		{Name: "KINGSTON SKC3000S1024G", MedienArt: medienSSD, BusArt: busNVMe},
		{Name: "Samsung SSD 970 EVO Plus 2TB", MedienArt: medienSSD, BusArt: busNVMe},
		{Name: "Crucial BX500", MedienArt: medienSSD, BusArt: busSATA},
		{Name: "WDC WD10EZEX", MedienArt: medienHDD, BusArt: busSATA},
	}
	befunde := Laufwerke(laufwerke)

	if len(befunde) != 2 {
		t.Fatalf("erwartet: genau zwei Befunde (SATA-SSD und Festplatte), kam: %+v", befunde)
	}
	// USB-Geraete und NVMe-SSDs duerfen nichts ausloesen.
	for _, b := range befunde {
		if b.Titel != "SSD am SATA-Anschluss" && b.Titel != "Magnetfestplatte verbaut" {
			t.Errorf("unerwarteter Befund: %s", b.Titel)
		}
	}
}
