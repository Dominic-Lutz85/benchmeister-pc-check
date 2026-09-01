package scan

import (
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

/*
Mit wie vielen Leitungen und in welcher PCIe-Generation Grafikkarte und
NVMe-SSDs gerade angebunden sind.

ANLASS (01.09.2026): "Logos_Atum" im PCGH-Forum. Er hatte beim
Zusammenbau zwei Karten in einen x8- und einen x16-Steckplatz gesteckt
und wollte wissen, wie viele Leitungen noch frei sind und wie schnell
eine nachgeruestete M.2 laufen wuerde.

WAS DAVON GEHT UND WAS NICHT, denn das ist nicht dasselbe:

  Geht:      In welchem Modus laeuft, was JETZT drinsteckt.
  Geht nicht: Wie viele Leitungen frei sind, und wie schnell etwas
              laufen wuerde, das noch gar nicht drinsteckt.

Das Zweite haengt daran, was sich auf dem konkreten Board welche
Leitungen teilt. Ob der zweite Steckplatz auf x8 faellt, sobald der
dritte M.2-Platz belegt wird, steht im Handbuch dieses einen Boards und
in keiner Schnittstelle. Fuer beliebige Rechner koennten wir das nicht
mitliefern, und raten kommt nach Regel 1 nicht in Frage.

WARUM NICHT UEBER WMI wie alles andere hier: WMI kennt die
Aushandlungsbreite nicht. Windows legt sie als Geraeteeigenschaft ab,
erreichbar nur ueber SetupAPI.

WARUM EIN EIGENER AUFRUF und nicht windows.SetupDiGetDeviceProperty aus
x/sys: Der dortige Bequemlichkeits-Wrapper wandelt ausschliesslich
Zeichenketten um und gibt fuer alles andere "unimplemented property
type" zurueck. Die Werte hier sind Zahlen.
*/

// Eigenschaftsgruppe der PCI-Geraete,
// {3ab22e31-8264-4b4e-9af5-a8d2d8e33e62}. Darin:
//
//	 9  CurrentLinkSpeed   ausgehandelte Generation
//	10  CurrentLinkWidth   ausgehandelte Zahl der Leitungen
//	11  MaxLinkSpeed       gemeldeter Hoechstwert
//	12  MaxLinkWidth       gemeldeter Hoechstwert
//
// Die beiden Hoechstwerte werden bewusst NICHT ausgelesen, Begruendung
// bei PcieGeraet in types.go.
var pcieGruppe = windows.DEVPROPGUID{
	Data1: 0x3ab22e31, Data2: 0x8264, Data3: 0x4b4e,
	Data4: [8]byte{0x9a, 0xf5, 0xa8, 0xd2, 0xd8, 0xe3, 0x3e, 0x62},
}

const (
	pidAktuelleGeneration = 9
	pidAktuelleBreite     = 10
)

// DEVPROP_TYPE_UINT32. Nicht aus der Dokumentation uebernommen, sondern
// am 01.09.2026 an einem echten Rechner nachgemessen: Alle vier
// Eigenschaften kamen mit diesem Typ und genau vier Byte zurueck.
const typUint32 = 0x07

var (
	setupapi                 = windows.NewLazySystemDLL("setupapi.dll")
	procGetDeviceProperty    = setupapi.NewProc("SetupDiGetDevicePropertyW")
	interessanteGeraeteArten = map[string]bool{
		// Geraeteklassen, wie Windows sie meldet. Am 01.09.2026 an einem
		// Rechner mit zwoelf PCIe-Geraeten gegengeprueft: Diese beiden
		// sind die, nach denen jemand fragt. Alles andere war Innenleben
		// (Speicherbruecken, USB-Controller, der Sicherheitsprozessor)
		// und meldet Werte, die niemandem etwas sagen.
		//
		// Nuetzlicher Nebeneffekt der Filterung: Die Tonausgabe einer
		// Grafikkarte haengt am selben Link wie die Karte und meldet
		// deshalb dieselben Zahlen noch einmal. Sie faellt hier als
		// Klasse "System" von selbst heraus, sonst stuende jede Karte
		// doppelt in der Uebersicht.
		"Display":     true,
		"SCSIAdapter": true,
	}
)

/*
Liest eine Zahl-Eigenschaft eines Geraets.

Gibt IMMER auch ein ok zurueck, wie kanalAus() und die Zustandswerte.
Fehlende Werte sind hier der Normalfall und kein Fehler: Bruecken und
Host-Controller melden diese Eigenschaften gar nicht, Windows antwortet
dann mit "Element nicht gefunden". Wer daraus eine 0 macht, schreibt
"x0" in die Uebersicht.
*/
func pcieZahl(satz windows.DevInfo, geraet *windows.DevInfoData, pid uint32) (uint32, bool) {
	schluessel := windows.DEVPROPKEY{FmtID: pcieGruppe, PID: windows.DEVPROPID(pid)}
	var typ uint32
	var puffer [16]byte
	var noetig uint32

	erfolg, _, _ := syscall.SyscallN(procGetDeviceProperty.Addr(),
		uintptr(satz), uintptr(unsafe.Pointer(geraet)),
		uintptr(unsafe.Pointer(&schluessel)), uintptr(unsafe.Pointer(&typ)),
		uintptr(unsafe.Pointer(&puffer[0])), uintptr(len(puffer)),
		uintptr(unsafe.Pointer(&noetig)), 0)

	// Typ und Laenge werden geprueft und nicht vorausgesetzt. Wenn eine
	// kuenftige Windows-Ausgabe hier etwas anderes liefert, soll der Wert
	// fehlen und nicht falsch sein.
	if erfolg == 0 || typ != typUint32 || noetig != 4 {
		return 0, false
	}
	return uint32(puffer[0]) | uint32(puffer[1])<<8 |
		uint32(puffer[2])<<16 | uint32(puffer[3])<<24, true
}

func alsText(wert interface{}) string {
	if s, ok := wert.(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

// pcieAnbindung fuellt e.Pcie. Ohne Rueckgabewert wie arbeitsspeicher()
// und laufwerk(): Wenn nichts auslesbar ist, bleibt die Liste leer und
// die Uebersicht zeigt den Abschnitt gar nicht erst.
func pcieAnbindung(e *ScanResult) {
	// "PCI" schraenkt auf den PCI-Bus ein, DIGCF_PRESENT auf das, was
	// gerade wirklich steckt. Ohne das Zweite stuenden auch Geraete in
	// der Liste, die frueher einmal verbaut waren.
	satz, err := windows.SetupDiGetClassDevsEx(nil, "PCI", 0,
		windows.DIGCF_PRESENT|windows.DIGCF_ALLCLASSES, 0, "")
	if err != nil {
		return
	}
	defer windows.SetupDiDestroyDeviceInfoList(satz)

	for i := 0; ; i++ {
		geraet, err := windows.SetupDiEnumDeviceInfo(satz, i)
		if err != nil {
			break
		}

		klasse := alsText(regEigenschaft(satz, geraet, windows.SPDRP_CLASS))
		if !interessanteGeraeteArten[klasse] {
			continue
		}

		breite, okBreite := pcieZahl(satz, geraet, pidAktuelleBreite)
		generation, okGeneration := pcieZahl(satz, geraet, pidAktuelleGeneration)
		// Beide oder keiner. Eine Breite ohne Generation waere eine halbe
		// Angabe, und halbe Angaben laden zum Weiterdenken ein.
		if !okBreite || !okGeneration || breite == 0 || generation == 0 {
			continue
		}

		name := alsText(regEigenschaft(satz, geraet, windows.SPDRP_DEVICEDESC))
		if name == "" {
			continue
		}

		e.Pcie = append(e.Pcie, PcieGeraet{
			Name:       name,
			Breite:     breite,
			Generation: generation,
		})
	}
}

// Kleine Huelle, damit der Fehlerfall oben nicht dreimal auftaucht.
func regEigenschaft(satz windows.DevInfo, geraet *windows.DevInfoData, was windows.SPDRP) interface{} {
	wert, err := windows.SetupDiGetDeviceRegistryProperty(satz, geraet, was)
	if err != nil {
		return nil
	}
	return wert
}
