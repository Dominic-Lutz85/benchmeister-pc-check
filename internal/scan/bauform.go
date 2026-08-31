package scan

import "github.com/yusufpapurcu/wmi"

/*
Bauform des Geraets: mobil oder stationaer. Angelegt am 31.08.2026 nach
einem Hinweis von "PCGH_Jacky" im PCGH-Forum.

WARUM DAS NOETIG WURDE: Zwei Empfehlungen des Programms sind bei einem
Notebook nicht befolgbar. "Einen zweiten Riegel ergaenzen" geht nicht,
wenn der Speicher verloetet ist oder kein Steckplatz frei liegt.
"Einen Riegel in den anderen Kanal umstecken" geht nicht, wenn der
zweite Steckplatz unter dem verklebten Boden sitzt.

Ein Rat, dem jemand nicht folgen KANN, ist schlimmer als kein Rat: Er
kostet das Vertrauen in alle anderen Befunde mit. Genau die Regel, die
in CLAUDE.md an erster Stelle steht.

DIE ZAHL IST NICHT VERHANDELBAR, aber auch nicht sicher: ChassisTypes
kommt aus der SMBIOS-Tabelle, die der Hersteller fuellt. Microsoft
schreibt dazu in der eigenen Dokumentation zur Geraeteerkennung, dass
Hersteller keine einheitliche Methode benutzen und nicht jeder Rechner
diese Angabe liefert. Deshalb drei Zustaende statt zwei: mobil,
stationaer, unbekannt.

WAS BEI EINEM IRRTUM PASSIERT, in beide Richtungen durchdacht:

  Notebook nicht erkannt  -> alter Text, also der Stand von gestern.
                             Kein neuer Schaden.
  Desktop faelschlich als -> Rat wird vorsichtiger ("erst nachsehen,
  Notebook erkannt           ob ein Steckplatz frei ist"). Das ist bei
                             einem Desktop ueberfluessig, aber nicht
                             falsch.

Beide Fehlerrichtungen sind harmlos. Deshalb ist diese Erkennung
vertretbar, obwohl die Quelle unzuverlaessig ist.

NICHT VERSUCHT: zu erkennen, ob der Speicher verloetet ist. Dafuer gibt
es keinen Weg. Win32_PhysicalMemory.FormFactor kennt SODIMM (12), aber
keinen Wert fuer "fest verbaut", und die Steckplatzzahl aus SMBIOS
zaehlt verloetete Bausteine auf vielen Geraeten als Steckplatz mit.
Deshalb sagt der Text bei einem Notebook nicht "geht nicht", sondern
"sieh nach, ob es geht".
*/

type win32SystemEnclosure struct {
	ChassisTypes []uint16
}

// Bauform eines Rechners.
type Bauform int

const (
	BauformUnbekannt Bauform = iota
	BauformStationaer
	BauformMobil
)

/*
Werte aus der Microsoft-Dokumentation zu Win32_SystemEnclosure,
Stand 31.08.2026 nachgeschlagen statt aus dem Gedaechtnis:

	 8 Portable      11 Hand Held    30 Tablet
	 9 Laptop        14 Sub Notebook 31 Convertible
	10 Notebook                      32 Detachable

Alles davon ist ein Geraet, bei dem Aufruesten nicht selbstverstaendlich
ist. Docking Station (12) gehoert bewusst NICHT dazu: Dort steckt
gegebenenfalls ein Notebook drin, die Station selbst hat keinen
Speicher.
*/
var mobileGehaeuse = map[uint16]bool{
	8: true, 9: true, 10: true, 11: true, 14: true,
	30: true, 31: true, 32: true,
}

// Stationaere Bauformen. Alles andere (Server-Gehaeuse, Sonderformen)
// bleibt unbekannt, statt geraten zu werden.
var stationaereGehaeuse = map[uint16]bool{
	3: true, 4: true, 5: true, 6: true, 7: true, 13: true, 15: true, 16: true,
}

/*
BauformErmitteln liest die Gehaeuseart aus.

ChassisTypes ist ein Feld, kein Einzelwert, und manche Geraete melden
mehrere Eintraege. Ein mobiler Eintrag genuegt: Wer "Notebook" und
"Main System Chassis" gleichzeitig meldet, ist ein Notebook.
*/
func BauformErmitteln() Bauform {
	var liste []win32SystemEnclosure
	if err := wmi.Query("select ChassisTypes from Win32_SystemEnclosure", &liste); err != nil {
		return BauformUnbekannt
	}

	stationaerGesehen := false
	for _, e := range liste {
		for _, t := range e.ChassisTypes {
			if mobileGehaeuse[t] {
				return BauformMobil
			}
			if stationaereGehaeuse[t] {
				stationaerGesehen = true
			}
		}
	}
	if stationaerGesehen {
		return BauformStationaer
	}
	return BauformUnbekannt
}
