package scan

import (
	"strings"
	"testing"
)

/*
Der Test kann den WMI-Aufruf nicht nachstellen, also prueft er das, was
sich pruefen laesst: dass die Ausgabe entweder leer ist oder eine
brauchbare Zeile ergibt, und dass sie nichts Persoenliches enthaelt.

Der zweite Teil ist der wichtigere. Win32_OperatingSystem liefert neben
dem Namen auch Rechnername, Benutzername und Installationskennung. Wer
die Abfrage spaeter erweitert, um "noch eben" ein Feld mitzunehmen,
soll hier anecken, bevor es jemandem auffaellt, der seinen Rechnernamen
in einer oeffentlichen Liste findet.
*/
func TestBetriebssystemGibtNurNameUndBaunummer(t *testing.T) {
	got := Betriebssystem()

	if got == "" {
		t.Skip("kein Windows oder WMI nicht erreichbar")
	}

	if strings.HasPrefix(got, "Microsoft ") {
		t.Errorf("das Praefix Microsoft steht noch drin: %q", got)
	}
	if got != strings.TrimSpace(got) {
		t.Errorf("Leerzeichen am Rand: %q", got)
	}

	// Ein Benutzer- oder Rechnername enthaelt fast immer einen
	// Backslash (DOMAENE\name) oder ein at-Zeichen.
	for _, zeichen := range []string{"\\", "@"} {
		if strings.Contains(got, zeichen) {
			t.Errorf("verdaechtiges Zeichen %q in %q, kommt da ein Name mit?", zeichen, got)
		}
	}

	// Eine sinnvolle Zeile ist kurz. Wer versehentlich mehrere Felder
	// aneinanderhaengt, landet schnell darueber.
	if len(got) > 80 {
		t.Errorf("die Zeile ist mit %d Zeichen zu lang: %q", len(got), got)
	}
}
