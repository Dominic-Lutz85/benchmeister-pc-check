package scan

import (
	"strings"

	"github.com/yusufpapurcu/wmi"
)

/*
Welches Windows hier laeuft, angelegt am 06.09.2026 fuer die Rangliste.

WOZU DAS GUT IST
================
Es geht NICHT in die Wertung ein. Es steht in der Liste daneben, und
zwar aus einem Grund, der erst mit vielen Eintraegen traegt: Sobald
genug Messungen zusammenkommen, laesst sich ablesen, ob dieselbe
Hardware unter verschiedenen Windows-Staenden anders laeuft. Das ist
eine Frage, die staendig gestellt und praktisch nie mit Zahlen
beantwortet wird.

Waere die Fassung Teil der Wertung, waere genau diese Auswertung ein
Zirkelschluss.

WAS HIER NICHT AUSGELESEN WIRD: kein Rechnername, kein Benutzername,
keine Installationskennung, kein Produktschluessel. Win32_OperatingSystem
haelt all das bereit, und nichts davon hat mit der Frage zu tun, wie ein
Prozessor seinen Takt haelt.
*/

type win32OperatingSystem struct {
	Caption     string
	BuildNumber string
}

// Betriebssystem liefert eine Zeile wie "Windows 11 Pro 26100".
//
// Leer, wenn die Abfrage nicht durchgeht. Ein leerer Wert ist besser als
// ein geratener: In der Rangliste steht dann "unbekannt", und das ist
// ehrlich.
func Betriebssystem() string {
	var liste []win32OperatingSystem
	q := "select Caption, BuildNumber from Win32_OperatingSystem"
	if err := wmi.Query(q, &liste); err != nil || len(liste) == 0 {
		return ""
	}

	// Caption steht als "Microsoft Windows 11 Pro" in der Tabelle. Das
	// "Microsoft" davor traegt keine Information und kostet in der
	// Tabellenspalte nur Platz.
	name := strings.TrimSpace(liste[0].Caption)
	name = strings.TrimPrefix(name, "Microsoft ")

	bau := strings.TrimSpace(liste[0].BuildNumber)
	if name == "" {
		return ""
	}
	if bau == "" {
		return name
	}
	// Die Baunummer entscheidet, nicht der Name: "Windows 11" gibt es
	// als 22000 und als 26100, und dazwischen liegen zwei Jahre
	// Aenderungen am Scheduler.
	return name + " " + bau
}
