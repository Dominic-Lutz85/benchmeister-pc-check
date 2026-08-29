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
