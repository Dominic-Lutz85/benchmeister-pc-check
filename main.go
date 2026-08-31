// BenchMeister PC-Check: liest einmalig aus, welche Hardware in diesem
// Rechner steckt, zeigt das Ergebnis lokal an und überträgt es nur, wenn
// man ausdrücklich zustimmt.
//
// Was dieses Programm bewusst NICHT tut:
//   - nichts installieren, nichts in die Registry schreiben,
//     keinen Autostart einrichten, keinen Dienst hinterlassen
//   - keine Administratorrechte anfordern
//   - keine Belastungstests fahren (kein Hochtakten, keine Volllast)
//   - keine Kennungen erfassen, über die sich ein Rechner oder eine Person
//     wiedererkennen ließe (siehe internal/scan/types.go)
//   - vor der Zustimmung keine einzige Netzwerkverbindung aufbauen
//
// Der letzte Punkt lässt sich mit einer Firewall nachprüfen, und genau
// dazu laden wir ausdrücklich ein.
package main

import (
	"fmt"
	"os"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/scan"
	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/ui"
)

// Version dieses Programms. Ergaenzt am 27.08.2026: Vorher stand die
// Versionsnummer nur im Git-Etikett und auf der Website, das Programm
// selbst kannte sie nicht. Bei einer Rueckmeldung ("bei mir kommt Fehler
// X") liess sich damit nicht feststellen, welcher Stand ueberhaupt
// laeuft. Beim Veroeffentlichen einer neuen Fassung hier mitziehen.
const version = "1.0.12"

func main() {
	fmt.Println("BenchMeister PC-Check", version)
	fmt.Println("Lese aus, was in diesem Rechner steckt ...")

	ergebnis, err := scan.Auslesen()
	if err != nil {
		abbrechen(err)
	}

	fmt.Printf("Gefunden: %s, %s\n", ergebnis.CPUName, ergebnis.GPUName)

	// Zustand des laufenden Systems: volle Systemplatte, mehrere
	// Schutzprogramme, Netzwerk unter Gigabit. Steht bewusst getrennt
	// vom Scan-Ergebnis, weil nichts davon uebertragen wird, siehe
	// SystemzustandInfo in internal/scan/types.go.
	//
	// Faellt eine der drei Abfragen aus, bleibt der jeweilige Wert
	// unerkannt und es gibt keinen Befund. Ein Fehler ist das nicht,
	// deshalb gibt Systemzustand() auch keinen zurueck.
	zustand := scan.Systemzustand()

	url, err := ui.Anzeigen(ergebnis, zustand)
	if err != nil {
		abbrechen(err)
	}

	if url == "" {
		fmt.Println("Es wurde nichts übertragen. Bis zum nächsten Mal.")
		return
	}
	fmt.Println("Dein Ergebnis:", url)
}

// Bei einem Fehler nicht einfach zumachen: Wird das Programm per
// Doppelklick gestartet, verschwindet das Fenster sofort und niemand
// erfährt, was schiefgelaufen ist.
func abbrechen(err error) {
	fmt.Fprintln(os.Stderr, "\nDas hat nicht geklappt:", err)
	fmt.Fprintln(os.Stderr, "\nZum Schließen die Eingabetaste drücken.")
	fmt.Scanln()
	os.Exit(1)
}
