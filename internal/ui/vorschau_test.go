package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/pruefung"
	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/scan"
)

/*
Schreibt die fertige Oberflaeche als HTML-Datei heraus, damit man sie
beim Gestalten ansehen kann, ohne jedes Mal das ganze Programm laufen
zu lassen.

WOZU: Das Programm scannt beim Start echte Hardware und oeffnet danach
den Browser. Fuer eine Aenderung an einer Ueberschrift ist das ein
langer Weg, und die Fundliste sieht auf dem Entwicklungsrechner immer
gleich aus. Dieser Test rendert dieselbe Vorlage mit einem Beispielfall,
der alle Zustaende gleichzeitig zeigt: kostenlose Funde,
kostenpflichtige, Bestaetigungen, mehrere Laufwerke und PCIe-Zeilen.

AUFRUF (schreibt nach vorschau.html im Projektordner):

	VORSCHAU=1 go test ./internal/ui/ -run TestVorschauSchreiben

Ohne die Umgebungsvariable tut der Test nichts. So laeuft er in der
normalen Testrunde mit, ohne bei jedem Durchgang eine Datei anzulegen.

DIE DATEN HIER SIND ERFUNDEN, und das ist genau der Unterschied zum
Programm selbst: Dort steht die eiserne Regel, dass nichts angezeigt
wird, was nicht gemessen wurde. Hier geht es nicht um eine Aussage
ueber einen Rechner, sondern um das Aussehen der Seite. Deshalb heisst
die Datei vorschau.html und liegt in .gitignore.
*/
func TestVorschauSchreiben(t *testing.T) {
	if os.Getenv("VORSCHAU") == "" {
		t.Skip("nur mit VORSCHAU=1, siehe Kommentar")
	}

	kostenlos := []pruefung.Befund{
		{
			Schwere:      pruefung.Hinweis,
			Titel:        "Der Arbeitsspeicher läuft unter seiner Geschwindigkeit",
			Feststellung: "Deine Riegel sind für 6000 MT/s gebaut, laufen laut Windows aber mit 4800 MT/s.",
			Empfehlung:   "Im BIOS das Profil EXPO oder XMP einschalten.",
			Hintergrund:  "Ohne Profil läuft DDR5 auf dem JEDEC-Grundtakt von 4800 MT/s. Das ist kein Defekt, sondern die Werkseinstellung.",
		},
		{
			Schwere:      pruefung.Hinweis,
			Titel:        "Der Energieplan steht auf Ausbalanciert",
			Feststellung: "Windows nutzt den Plan „Ausbalanciert“.",
			Empfehlung:   "Für Spiele auf „Höchstleistung“ umstellen.",
		},
	}

	kostenpflichtig := []pruefung.Befund{
		{
			Schwere:      pruefung.Zukauf,
			Titel:        "Nur ein Speicherriegel verbaut",
			Feststellung: "Ein Riegel mit 16 GB in einem Board mit vier Steckplätzen.",
			Empfehlung:   "Einen zweiten, baugleichen Riegel ergänzen.",
			Hintergrund:  "Mit einem Riegel läuft der Speicher im Einkanalbetrieb. In Spielen kostet das je nach Titel zwischen 5 und 20 Prozent Bildrate.",
		},
	}

	laeuft := []pruefung.Befund{
		{
			Schwere:      pruefung.Bestaetigung,
			Titel:        "Die Systemplatte ist eine NVMe-SSD",
			Feststellung: "Windows liegt auf einer NVMe-SSD mit 2 TB.",
			Empfehlung:   "Hier ist nichts zu tun.",
		},
	}

	alle := append(append(append([]pruefung.Befund{}, kostenlos...), kostenpflichtig...), laeuft...)

	daten := vorlagenDaten{
		Scan: &scan.ScanResult{
			CPUName:           "AMD Ryzen 7 7800X3D",
			CPUCores:          8,
			CPUThreads:        16,
			GPUName:           "NVIDIA GeForce RTX 5060 Ti",
			RamTotalGb:        16,
			StorageType:       "SSD",
			StorageCapacityGb: 2000,
			ResolutionWidth:   2560,
			ResolutionHeight:  1440,
		},
		RohdatenJSON:    "{\n  \"cpu_name\": \"AMD Ryzen 7 7800X3D\",\n  \"gpu_name\": \"NVIDIA GeForce RTX 5060 Ti\"\n}",
		Befunde:         alle,
		Kostenlos:       kostenlos,
		Kostenpflichtig: kostenpflichtig,
		Laeuft:          laeuft,
		Laufwerke: []LaufwerkZeile{
			{Beschriftung: "NVMe-SSD, 2 TB (System)"},
			{Beschriftung: "SATA-SSD, 1 TB"},
		},
	}

	html := rendere(t, daten)

	ziel := filepath.Join("..", "..", "vorschau.html")
	if err := os.WriteFile(ziel, []byte(html), 0o644); err != nil {
		t.Fatalf("Vorschau laesst sich nicht schreiben: %v", err)
	}
	abs, _ := filepath.Abs(ziel)
	t.Logf("Vorschau geschrieben: %s", abs)
}
