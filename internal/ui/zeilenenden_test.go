package ui

import (
	"bytes"
	"testing"
)

/*
 * Haelt fest, dass die eingebettete Ergebnisseite mit LF endet.
 *
 * WARUM DAS EIN TEST IST UND NICHT NUR EIN EINTRAG IN .gitattributes:
 * Die Vorgabe dort greift beim Auschecken. Speichert danach ein Editor
 * die Datei mit CRLF, faellt das nirgends auf. git normalisiert beim
 * Lesen wieder auf LF, "git status" bleibt sauber, und die Datei liegt
 * trotzdem falsch auf der Platte. Genau dieser Zustand lag am
 * 07.09.2026 vor, und zwar bei fuenf Dateien.
 *
 * Das ist hier kein Schoenheitsfehler. Diese Vorlage wandert per
 * go:embed unveraendert in die exe. Nachgemessen: derselbe Quelltext,
 * einmal mit LF und einmal mit CRLF gebaut, ergab zwei exe-Dateien mit
 * 7.339.490 abweichenden Bytes und 1024 Byte Groessenunterschied.
 *
 * Die exe ist nicht signiert, Windows warnt beim Start. Das einzige
 * belastbare Gegenargument lautet "bau es selbst und vergleiche die
 * Pruefsumme mit meiner". Wer das mit einer CRLF-Arbeitskopie tut,
 * bekommt eine Abweichung und muss annehmen, dass an der
 * veroeffentlichten Datei etwas nicht stimmt.
 *
 * Der Test liest ueber das eingebettete Dateisystem, also genau die
 * Bytes, die spaeter im Programm stehen, nicht die von der Platte.
 */
func TestVorlageHatKeineWagenruecklaeufe(t *testing.T) {
	roh, err := vorlagen.ReadFile("assets/preview.html.tmpl")
	if err != nil {
		t.Fatalf("Vorlage laesst sich nicht lesen: %v", err)
	}
	if n := bytes.Count(roh, []byte("\r\n")); n > 0 {
		t.Errorf("die eingebettete Vorlage hat %d CRLF-Zeilenenden. "+
			"Damit weicht die gebaute exe von der aus dem Baulauf ab. "+
			"Zum Beheben: Datei loeschen und mit \"git checkout -- "+
			"internal/ui/assets/preview.html.tmpl\" neu holen.", n)
	}
}
