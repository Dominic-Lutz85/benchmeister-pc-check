package ui

import (
	"strings"
	"testing"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/scan"
)

/*
 * Die PCIe-Anbindung ist der einzige Abschnitt, der Zahlen zeigt, ohne
 * sie zu bewerten. Das ist kein Versehen, sondern die einzige ehrliche
 * Form: Im Leerlauf sind die Werte niedriger als unter Last, und der
 * von Windows gemeldete Hoechstwert stimmt nachweislich nicht (siehe
 * scan.PcieGeraet).
 *
 * Genau deshalb braucht dieser Abschnitt einen Test. Zahlen ohne
 * Einordnung sind hier schlimmer als gar keine Zahlen: "x8" liest sich
 * fuer die meisten wie ein Mangel, auch wenn es keiner ist. Wer die
 * Erklaerung spaeter herausnimmt oder die Zeilen woandershin kopiert,
 * soll einen roten Test bekommen und keinen zufriedenen Leser weniger.
 */

func beispielGeraete() []scan.PcieGeraet {
	// Echte Werte, am 01.09.2026 auf dem Entwicklungsrechner gemessen.
	// Die Karte ist der interessante Fall: x8 ist bei ihr richtig, weil
	// sie von Haus aus nur acht Leitungen hat.
	return []scan.PcieGeraet{
		{Name: "NVIDIA GeForce RTX 5060 Ti", Breite: 8, Generation: 4},
		{Name: "Standardmäßiger NVM Express-Controller", Breite: 4, Generation: 3},
	}
}

func TestPcieZeilenNennenBreiteUndGeneration(t *testing.T) {
	zeilen := pcieZeilen(&scan.ScanResult{Pcie: beispielGeraete()})

	if len(zeilen) != 2 {
		t.Fatalf("zwei Zeilen erwartet, %d bekommen", len(zeilen))
	}
	will := "NVIDIA GeForce RTX 5060 Ti: x8, PCIe 4.0"
	if zeilen[0].Beschriftung != will {
		t.Errorf("Beschriftung falsch.\nwill: %s\nist:  %s", will, zeilen[0].Beschriftung)
	}
}

// Ein Geraet ohne Namen ergaebe eine Zeile, die mit ": x4" beginnt.
func TestPcieZeilenLassenNamenloseWeg(t *testing.T) {
	zeilen := pcieZeilen(&scan.ScanResult{Pcie: []scan.PcieGeraet{
		{Name: "   ", Breite: 4, Generation: 4},
	}})
	if len(zeilen) != 0 {
		t.Fatalf("keine Zeile erwartet, %d bekommen: %q", len(zeilen), zeilen)
	}
}

/*
 * Der eigentliche Waechter: Die Zahlen duerfen nie ohne die Erklaerung
 * erscheinen, warum daraus keine Warnung wird.
 */
func TestPcieZahlenNieOhneErklaerung(t *testing.T) {
	html := rendere(t, vorlagenDaten{
		Scan: &scan.ScanResult{},
		Pcie: pcieZeilen(&scan.ScanResult{Pcie: beispielGeraete()}),
	})

	if !strings.Contains(html, "x8, PCIe 4.0") {
		t.Fatal("die Anbindung der Grafikkarte fehlt in der Uebersicht")
	}
	if !strings.Contains(html, "Warum hier keine Bewertung steht") {
		t.Error("die Zahlen stehen da, die Einordnung dazu fehlt")
	}
	// Die beiden Gruende muessen beide dastehen. Einer allein erklaert
	// nur die Haelfte, und die weggelassene Haelfte ist jeweils die, auf
	// die jemand hereinfaellt.
	if !strings.Contains(html, "Stromsparen senkt die Werte") {
		t.Error("Grund 1 fehlt: niedrigere Werte im Leerlauf")
	}
	if !strings.Contains(html, "Der gemeldete Höchstwert stimmt nicht") {
		t.Error("Grund 2 fehlt: unbrauchbarer Hoechstwert von Windows")
	}
}

// Ohne auslesbare Werte faellt der Abschnitt weg, statt leer dazustehen.
func TestOhnePcieKeinAbschnitt(t *testing.T) {
	html := rendere(t, vorlagenDaten{Scan: &scan.ScanResult{}})
	if strings.Contains(html, "PCIe-Anbindung") {
		t.Error("ohne Werte darf die Zeile nicht erscheinen")
	}
}
