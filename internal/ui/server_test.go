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

/*
 * Haelt die beiden Stellen zusammen, die den Ist-Takt berechnen.
 *
 * Die Regel steht zweimal im Programm: einmal in
 * scan.RiegelInfo.IstTaktMhz (fuer Anzeige und Uebertragung) und einmal
 * in pruefung.Speicher (fuer den Befund). Das ist Absicht, die beiden
 * Pakete sollen nichts voneinander wissen. Genau diese Doppelung war
 * aber am 29.08.2026 der Fehler: Die Pruefung rechnete mit
 * ConfiguredClockSpeed, Anzeige und Upload weiterhin nur mit Speed. Auf
 * derselben Seite stand dann im Befund 5600 und in der Tabelle darueber
 * 4800.
 *
 * Der Test prueft nicht die Formel, sondern das, was zaehlt: Der Befund
 * kommt genau dann, wenn der von scan berechnete Ist-Takt unter der
 * Sollgeschwindigkeit liegt. Laufen die beiden Rechnungen auseinander,
 * faellt er.
 */
func TestScanUndPruefungRechnenGleich(t *testing.T) {
	// CMK32GX5M2B5200C40 heisst: Soll 5200. Die 5 Prozent Toleranz der
	// Pruefung liegen bei 4940, deshalb ist 5000 knapp in Ordnung.
	const teil = "CMK32GX5M2B5200C40"

	faelle := []struct {
		name       string
		speed      uint32
		konfig     uint32
		befundNoet bool
	}{
		{"Profil an, laeuft ueber Soll", 4800, 5600, false},
		{"Profil an, laeuft auf Soll", 4800, 5200, false},
		{"Profil aus, Board fuellt beide Felder", 4800, 4800, true},
		{"Board fuellt das zweite Feld nicht", 4800, 0, true},
		{"Board fuellt nur das zweite Feld", 0, 5200, false},
		{"knapp innerhalb der Toleranz", 4800, 5000, false},
		{"knapp ausserhalb der Toleranz", 4800, 4900, true},
	}

	for _, f := range faelle {
		t.Run(f.name, func(t *testing.T) {
			roh := scan.RiegelInfo{
				KapazitaetBytes: 16 << 30,
				TaktMhz:         f.speed,
				TaktMhzZweiter:  f.konfig,
				Teilenummer:     teil,
			}
			e := &scan.ScanResult{Riegel: []scan.RiegelInfo{roh, roh}}

			var befund bool
			for _, b := range pruefung.Speicher(riegelFuerPruefung(e)) {
				if b.Titel == "Arbeitsspeicher läuft womöglich unter seiner Sollgeschwindigkeit" {
					befund = true
				}
			}

			// Erste Probe: Stimmt das Verhalten mit der Tabelle ueberein?
			// Faengt es ab, wenn BEIDE Rechnungen gemeinsam falsch werden.
			if befund != f.befundNoet {
				t.Errorf("Speed %d, ConfiguredClockSpeed %d: Befund %v erwartet, %v bekommen",
					f.speed, f.konfig, f.befundNoet, befund)
			}

			// Zweite Probe, das ist die eigentliche Klammer: Aus dem Wert,
			// den scan berechnet, muss sich der Befund vorhersagen lassen.
			// Weicht eine der beiden Rechnungen ab, passt die Vorhersage
			// nicht mehr und dieser Test faellt.
			const sollTakt = 5200
			const toleranz = 0.95
			ist := roh.IstTaktMhz()
			vorhergesagt := ist > 0 && float64(ist) < sollTakt*toleranz
			if befund != vorhergesagt {
				t.Errorf("scan.IstTaktMhz sagt %d, daraus folgt Befund %v, "+
					"die Pruefung sagt aber %v. Die beiden Rechnungen sind "+
					"auseinandergelaufen (Speed %d, ConfiguredClockSpeed %d)",
					ist, vorhergesagt, befund, f.speed, f.konfig)
			}
		})
	}
}
