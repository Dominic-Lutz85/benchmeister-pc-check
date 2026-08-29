package ui

import (
	"html/template"
	"strings"
	"testing"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/pruefung"
	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/scan"
)

/*
 * Sichert die Aufteilung der Ergebniskarten ab, angelegt am 29.08.2026
 * zusammen mit dem Bestaetigungs-Befund.
 *
 * WARUM DAS EINEN TEST WERT IST: Die Vorlage ist die einzige Stelle, an
 * der die Schweregrade wieder zusammenlaufen. Wer dort eine Bedingung
 * verdreht, bekommt keinen Uebersetzungsfehler, sondern eine Seite, auf
 * der ueber einem einwandfrei laufenden Bauteil steht "Das hier kostet
 * dich gerade Leistung". Genau dieser Fehler ist am 28.08.2026 schon
 * einmal passiert (mit den Zukauf-Empfehlungen) und wurde von einem
 * Moderator im PCGH-Forum angestrichen. Ein zweites Mal soll ihn ein
 * Test finden und kein Leser.
 */

func rendere(t *testing.T, daten vorlagenDaten) string {
	t.Helper()
	v, err := template.ParseFS(vorlagen, "assets/preview.html.tmpl")
	if err != nil {
		t.Fatalf("Vorlage laesst sich nicht lesen: %v", err)
	}
	var b strings.Builder
	if err := v.Execute(&b, daten); err != nil {
		t.Fatalf("Vorlage laesst sich nicht rendern: %v", err)
	}
	return b.String()
}

func TestBestaetigungStehtNichtUnterDerFundUeberschrift(t *testing.T) {
	bestaetigung := pruefung.Befund{
		Schwere:      pruefung.Bestaetigung,
		Titel:        "Der Arbeitsspeicher läuft auf Sollgeschwindigkeit",
		Feststellung: "Deine Riegel sind für 5600 MT/s gebaut und laufen laut Windows mit 5600 MT/s.",
		Empfehlung:   "Hier ist nichts zu tun.",
	}
	fund := pruefung.Befund{
		Schwere:      pruefung.Hinweis,
		Titel:        "Nur ein Speicherriegel verbaut",
		Feststellung: "Einkanalbetrieb.",
		Empfehlung:   "Zweiten Riegel ergänzen.",
	}
	alle := []pruefung.Befund{bestaetigung, fund}
	kostenlos, kostenpflichtig, laeuft := nachKosten(alle)
	html := rendere(t, vorlagenDaten{
		Scan:            &scan.ScanResult{},
		Befunde:         alle,
		Kostenlos:       kostenlos,
		Kostenpflichtig: kostenpflichtig,
		Laeuft:          laeuft,
	})

	if !strings.Contains(html, bestaetigung.Titel) {
		t.Error("die Bestaetigung fehlt auf der Seite")
	}
	if !strings.Contains(html, "Das läuft schon so, wie es soll.") {
		t.Error("die Ueberschrift der Bestaetigungs-Karte fehlt")
	}

	// Der Kern: Die Bestaetigung muss NACH der Fund-Ueberschrift stehen
	// und darf nicht zwischen den Funden liegen.
	fundUeberschrift := strings.Index(html, "Das hier kostet dich gerade Leistung")
	gutUeberschrift := strings.Index(html, "Das läuft schon so, wie es soll.")
	bestaetigungPos := strings.Index(html, bestaetigung.Titel)
	if fundUeberschrift < 0 || gutUeberschrift < 0 {
		t.Fatal("eine der beiden Ueberschriften fehlt")
	}
	if bestaetigungPos < gutUeberschrift {
		t.Error("die Bestaetigung steht vor ihrer eigenen Ueberschrift, also in der falschen Karte")
	}
}

// Solange irgendein Befund vorliegt, auch eine blosse Bestaetigung, darf
// die Karte "Nichts gefunden" nicht erscheinen. Sonst stuenden zwei
// Aussagen zum selben Sachverhalt nebeneinander.
func TestKeineLeerKarteWennBestaetigungVorliegt(t *testing.T) {
	alle := []pruefung.Befund{{
		Schwere: pruefung.Bestaetigung,
		Titel:   "Der Arbeitsspeicher läuft auf Sollgeschwindigkeit",
	}}
	kostenlos, kostenpflichtig, laeuft := nachKosten(alle)
	html := rendere(t, vorlagenDaten{
		Scan:            &scan.ScanResult{},
		Befunde:         alle,
		Kostenlos:       kostenlos,
		Kostenpflichtig: kostenpflichtig,
		Laeuft:          laeuft,
	})
	if strings.Contains(html, "Nichts gefunden.") {
		t.Error("bei vorliegender Bestaetigung darf die Leer-Karte nicht erscheinen")
	}
	// Und die Fund-Karte auch nicht, sie waere leer.
	if strings.Contains(html, "Das hier kostet dich gerade Leistung") {
		t.Error("ohne echte Funde darf die Fund-Ueberschrift nicht erscheinen")
	}
}

/*
 * Der Aufklapper, angelegt am 30.08.2026.
 *
 * Der Hintergrund eines Befundes muss in einem ZUGEKLAPPTEN Block
 * landen, nicht offen im Text. Sonst waere die Umstellung wirkungslos:
 * Die Begruendungen stuenden wieder da, nur an anderer Stelle, und der
 * Befund waere so lang wie vorher.
 *
 * Geprueft wird auch, dass <details> KEIN open traegt. Ein aufgeklapptes
 * details sieht im Quelltext fast gleich aus und macht den ganzen
 * Umbau zunichte.
 */
func TestHintergrundStehtImAufklapper(t *testing.T) {
	befund := pruefung.Befund{
		Schwere:      pruefung.Hinweis,
		Titel:        "Nur ein Speicherriegel verbaut",
		Feststellung: "Mit einem einzelnen Riegel läuft der Speicher im Einkanalbetrieb.",
		Empfehlung:   "Einen zweiten, baugleichen Riegel ergänzen.",
		Hintergrund:  "Zwei Riegel arbeiten nebeneinander und verdoppeln die Bandbreite.",
	}
	alle := []pruefung.Befund{befund}
	kostenlos, kostenpflichtig, laeuft := nachKosten(alle)
	html := rendere(t, vorlagenDaten{
		Scan:            &scan.ScanResult{},
		Befunde:         alle,
		Kostenlos:       kostenlos,
		Kostenpflichtig: kostenpflichtig,
		Laeuft:          laeuft,
	})

	if !strings.Contains(html, befund.Hintergrund) {
		t.Error("der Hintergrund fehlt auf der Seite")
	}
	if !strings.Contains(html, "<summary>Warum?</summary>") {
		t.Error("der Aufklapper fehlt")
	}
	if strings.Contains(html, `<details class="mehr" open`) {
		t.Error("der Aufklapper darf nicht offen starten, sonst ist er wirkungslos")
	}

	// Der Hintergrund muss INNERHALB des details stehen, nicht davor.
	// Auf die eigene Klasse eingrenzen: Die Vorlage hat noch einen
	// zweiten Aufklapper fuer die Rohdaten, und der steht weiter oben.
	auf := strings.Index(html, `<details class="mehr"`)
	// Ab dem oeffnenden Tag suchen. Sonst findet man das schliessende
	// des Rohdaten-Aufklappers, der weiter oben steht.
	zu := auf + strings.Index(html[auf:], "</details>")
	pos := strings.Index(html, befund.Hintergrund)
	if auf < 0 || zu < 0 {
		t.Fatal("kein details-Block im Ergebnis")
	}
	if pos < auf || pos > zu {
		t.Error("der Hintergrund steht ausserhalb des Aufklappers")
	}
}

// Ein Befund ohne Hintergrund darf keinen leeren Aufklapper erzeugen.
// Ein "Warum?" das nichts erklaert, ist schlimmer als keines.
func TestOhneHintergrundKeinAufklapper(t *testing.T) {
	alle := []pruefung.Befund{{
		Schwere:      pruefung.Hinweis,
		Titel:        "Irgendein Fund",
		Feststellung: "Etwas ist auffällig.",
		Empfehlung:   "Etwas tun.",
	}}
	kostenlos, kostenpflichtig, laeuft := nachKosten(alle)
	html := rendere(t, vorlagenDaten{
		Scan:            &scan.ScanResult{},
		Befunde:         alle,
		Kostenlos:       kostenlos,
		Kostenpflichtig: kostenpflichtig,
		Laeuft:          laeuft,
	})
	if strings.Contains(html, `<details class="mehr"`) {
		t.Error("ohne Hintergrund darf kein leerer Aufklapper erscheinen")
	}
}
