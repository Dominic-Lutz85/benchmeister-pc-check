package ui

import (
	"testing"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/pruefung"
	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/scan"
)

/*
 * Dieser Test haelt einen Fehler fest, der zweimal ausgeliefert wurde.
 *
 * Version 1.0.5 sollte ConfiguredClockSpeed zum Ist-Takt machen. Das
 * Auslesen tat es, die Pruefung tat es, aber riegelFuerPruefung liess
 * das Feld beim Umkopieren weg. Ergebnis: Der Wert kam nie an, der
 * Rueckfall auf Speed griff jedes Mal, und das Programm meldete
 * unveraendert den JEDEC-Grundtakt.
 *
 * Der Fall stammt nicht aus der Fantasie, sondern aus einem Screenshot
 * von Misanthrop68 im PCGH-Forum am 29.08.2026: zwei Riegel
 * CMK32GX5M2B5200C40, Speed 4800, ConfiguredClockSpeed 5600. Sein
 * Speicher lief also ueber dem Profil, und ausgerechnet ihm riet das
 * Programm, endlich XMP einzuschalten.
 *
 * Der Test prueft deshalb nicht nur, dass das eine Feld ankommt,
 * sondern gleich das Verhalten, das daran haengt: Bei diesen Werten
 * darf KEIN Takt-Befund entstehen.
 */
func TestRiegelFuerPruefungBehaeltDenZweitenTakt(t *testing.T) {
	e := &scan.ScanResult{
		Riegel: []scan.RiegelInfo{
			{KapazitaetBytes: 16 << 30, TaktMhz: 4800, TaktMhzZweiter: 5600, Teilenummer: "CMK32GX5M2B5200C40", Kanal: "P0 CHANNEL A"},
			{KapazitaetBytes: 16 << 30, TaktMhz: 4800, TaktMhzZweiter: 5600, Teilenummer: "CMK32GX5M2B5200C40", Kanal: "P0 CHANNEL B"},
		},
	}

	riegel := riegelFuerPruefung(e)
	if len(riegel) != 2 {
		t.Fatalf("zwei Riegel erwartet, %d bekommen", len(riegel))
	}
	for i, r := range riegel {
		if r.TaktMhzZweiter != 5600 {
			t.Errorf("Riegel %d: ConfiguredClockSpeed ging beim Umwandeln verloren, "+
				"TaktMhzZweiter ist %d statt 5600", i, r.TaktMhzZweiter)
		}
	}

	// Soll 5200 aus der Teilenummer, Ist 5600 aus ConfiguredClockSpeed:
	// Der Speicher laeuft ueber dem Profil. Ein Befund waere ein
	// Falschalarm, und genau der stand im Forum.
	for _, b := range pruefung.Speicher(riegel) {
		if b.Titel == "Arbeitsspeicher läuft womöglich unter seiner Sollgeschwindigkeit" {
			t.Errorf("Falschalarm: %s", b.Feststellung)
		}
	}
}

// Gegenprobe, damit der Test oben nicht dadurch gruen wird, dass die
// Takt-Pruefung ueberhaupt nichts mehr meldet: Steht das Profil aus,
// muss der Befund weiterhin kommen.
func TestRiegelOhneProfilMeldetWeiterhin(t *testing.T) {
	e := &scan.ScanResult{
		Riegel: []scan.RiegelInfo{
			{KapazitaetBytes: 16 << 30, TaktMhz: 4800, TaktMhzZweiter: 4800, Teilenummer: "CMK32GX5M2B5200C40"},
			{KapazitaetBytes: 16 << 30, TaktMhz: 4800, TaktMhzZweiter: 4800, Teilenummer: "CMK32GX5M2B5200C40"},
		},
	}

	var gefunden bool
	for _, b := range pruefung.Speicher(riegelFuerPruefung(e)) {
		if b.Titel == "Arbeitsspeicher läuft womöglich unter seiner Sollgeschwindigkeit" {
			gefunden = true
		}
	}
	if !gefunden {
		t.Error("Bei 4800 statt 5200 muss der Takt-Befund kommen, er kam nicht")
	}
}
