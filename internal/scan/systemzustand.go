package scan

import (
	"os"
	"strings"

	"github.com/yusufpapurcu/wmi"
)

/*
Drei Angaben zum Zustand des laufenden Systems, angelegt am 31.08.2026
nach Vorschlaegen von "cryon1c" im PCGH-Forum: volle Systemplatte,
mehrere Schutzprogramme, Netzwerk unter Gigabit.

Alle drei sind keine Hardware-Fragen, sondern Zustandsfragen. Und alle
drei bleiben auf dem Rechner, siehe SystemzustandInfo in types.go.

WAS HIER BEWUSST NICHT PASSIERT: Keine Laufwerksbuchstaben, keine
Datentraegernamen, keine MAC-Adressen, keine Adapternamen, keine IPs.
Gebraucht werden Zahlen, nicht Kennungen.
*/

type win32LogicalDisk struct {
	Size      uint64
	FreeSpace uint64
}

type securityCenterAntiVirus struct {
	DisplayName string
}

type win32NetworkAdapter struct {
	Speed uint64
}

// Systemzustand liest die drei Werte. Jeder Teil scheitert fuer sich
// allein: Faellt einer aus, bleiben die anderen brauchbar, und der
// jeweilige "Erkannt"-Schalter bleibt false.
func Systemzustand() SystemzustandInfo {
	var z SystemzustandInfo
	systemplatte(&z)
	virenschutz(&z)
	netzwerk(&z)
	return z
}

/*
Freier Platz auf dem Laufwerk, auf dem Windows liegt.

NUR DIESES EINE LAUFWERK, und zwar aus zwei Gruenden. Erstens ist es
das einzige, dessen Fuellstand die Geschwindigkeit des ganzen Rechners
betrifft: Windows legt dort Auslagerungsdatei und Zwischenspeicher an.
Zweitens geht niemanden etwas an, wie voll die Datengraeber sind.

Der Buchstabe kommt aus der Umgebungsvariablen SystemDrive statt fest
aus "C:". Auf einem Rechner, der von D: startet, waere die feste
Annahme schlicht falsch.
*/
func systemplatte(z *SystemzustandInfo) {
	laufwerk := strings.TrimSpace(os.Getenv("SystemDrive"))
	if laufwerk == "" {
		return
	}

	var liste []win32LogicalDisk
	// DriveType=3 ist eine feste Platte. Ohne die Einschraenkung
	// koennte bei einem ungewoehnlichen SystemDrive ein Netzlaufwerk
	// antworten.
	q := "select Size, FreeSpace from Win32_LogicalDisk where DeviceID='" +
		laufwerk + "' and DriveType=3"
	if err := wmi.Query(q, &liste); err != nil || len(liste) == 0 {
		return
	}

	d := liste[0]
	// Groesse 0 kommt bei nicht bereiten Datentraegern vor. Daraus
	// "0 Prozent frei" zu rechnen waere ein Fehlalarm.
	if d.Size == 0 {
		return
	}

	const gb = 1024 * 1024 * 1024
	z.SystemplatteGesamtGb = int(d.Size / gb)
	z.SystemplatteFreiGb = int(d.FreeSpace / gb)
	z.SystemplatteErkannt = true
}

/*
Namen der registrierten Echtzeit-Schutzprogramme.

WARUM NUR DIE NAMEN UND NICHT DER ZUSTAND: Der Namensraum
root\SecurityCenter2 ist von Microsoft nicht dokumentiert, und das Feld
productState ist ein Bitfeld, dessen Bedeutung nur aus Beobachtung
rekonstruiert ist. Verschiedene Quellen legen einzelne Bits
unterschiedlich aus.

Damit ist es fuer Regel 1 dieses Programms ungeeignet: Wer daraus
"laeuft" oder "laeuft nicht" ableitet, baut einen Fehlalarm auf einer
Vermutung. Gezaehlt wird deshalb nur, was registriert IST, und der
Befund sagt das auch so.

Windows Defender erscheint hier mit. Das ist richtig so, denn zusammen
mit einem zweiten Programm ist gerade er der haeufige Fall.
*/
func virenschutz(z *SystemzustandInfo) {
	var liste []securityCenterAntiVirus
	q := "select displayName from AntiVirusProduct"
	// Eigener Namensraum. Auf Windows-Server-Ausgaben fehlt er
	// vollstaendig, dann bleibt VirenschutzErkannt false und es gibt
	// keinen Befund.
	if err := wmi.QueryNamespace(q, &liste, `root\SecurityCenter2`); err != nil {
		return
	}

	z.VirenschutzErkannt = true
	for _, p := range liste {
		name := strings.TrimSpace(p.DisplayName)
		if name != "" {
			z.Virenschutz = append(z.Virenschutz, name)
		}
	}
}

/*
Aushandlungsgeschwindigkeit der kabelgebundenen Verbindung.

NUR KABEL, kein WLAN. AdapterTypeId 0 ist Ethernet 802.3, 9 ist
Wireless, und das ist dokumentiert. Bei WLAN waere eine niedrigere
Aushandlung voellig normal und kein Zeichen fuer ein Problem, ein
Befund darauf waere Unsinn.

NetConnectionStatus=2 heisst verbunden. Ein Adapter, an dem kein Kabel
steckt, meldet Fantasiewerte.
*/
func netzwerk(z *SystemzustandInfo) {
	var liste []win32NetworkAdapter
	q := "select Speed from Win32_NetworkAdapter " +
		"where PhysicalAdapter=TRUE and NetConnectionStatus=2 and AdapterTypeId=0"
	if err := wmi.Query(q, &liste); err != nil || len(liste) == 0 {
		return
	}

	// Bei mehreren verbundenen Kabeln zaehlt die schnellste Verbindung.
	// Ein zweiter, langsamer Anschluss fuer ein Messgeraet oder eine
	// Kamera ist kein Grund fuer eine Warnung.
	var schnellste uint64
	for _, a := range liste {
		if a.Speed > schnellste {
			schnellste = a.Speed
		}
	}
	if schnellste == 0 {
		return
	}

	z.NetzwerkMbit = int(schnellste / 1000000)
	z.NetzwerkErkannt = true
}
