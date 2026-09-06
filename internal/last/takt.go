// Package last misst, was ein Rechner unter Beanspruchung tatsaechlich
// leistet, statt nur abzulesen, was verbaut ist.
//
// WARUM ES DIESES PAKET GIBT
// ==========================
// Bis hierher liest das Programm nur ab, was Windows meldet. Das
// beantwortet "was steckt drin", aber nicht die Frage, die jeder
// eigentlich hat: Laeuft die Kiste so, wie sie sollte, oder habe ich ein
// Montagsmodell? Zwei Rechner mit derselben Teileliste koennen sich um
// zwanzig Prozent unterscheiden, und man sieht es keiner Teileliste an.
//
// Die Messgroesse ist bewusst NICHT die Temperatur, sondern der
// GEHALTENE TAKT. Grund: Temperatur ist ohne Kernel-Treiber auf den
// meisten Rechnern gar nicht lesbar, der Takt dagegen schon. Und der
// Takt ist ohnehin die ehrlichere Zahl, denn er ist das Ergebnis von
// allem zusammen: Kuehlung, Power-Limit, Gehaeusebelueftung,
// Waermeleitpaste, Staub. Faellt er unter Dauerlast ab, stimmt etwas
// nicht, egal woran es liegt.
//
// DER STOLPERSTEIN, DER DAS HIER NOETIG MACHT
// ===========================================
// Windows-Leistungsindikatoren heissen in jeder Sprache anders. Auf
// einem deutschen System liefert
//
//	\Processor Information(_Total)\% Processor Performance
//
// den Fehler "Das angegebene Objekt wurde nicht gefunden", weil der
// Zaehler dort \Prozessorinformationen(_Total)\% Prozessorleistung
// heisst. Am 06.09.2026 genau so nachgestellt.
//
// Die Loesung ist PdhAddEnglishCounter: Diese Funktion nimmt IMMER die
// englische Schreibweise, unabhaengig von der Systemsprache, und
// uebersetzt intern. Damit steht im Quelltext ein Name, der auf jedem
// Windows funktioniert, ohne Zaehler-Indizes nachschlagen zu muessen.
//
// WARUM NICHT Win32_Processor.CurrentClockSpeed
// =============================================
// Weil das Feld gerundet und traege ist. Am 06.09.2026 gegengeprueft:
// Win32_Processor meldete 3801 MHz (den Nennwert), der Zaehler
// dagegen 121,7 Prozent Leistung, also tatsaechlich rund 4627 MHz. Wer
// das WMI-Feld nimmt, sieht einen Boost nie und einen Einbruch auch
// nicht.
package last

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	pdh = windows.NewLazySystemDLL("pdh.dll")

	procOpenQuery         = pdh.NewProc("PdhOpenQueryW")
	procAddEnglishCounter = pdh.NewProc("PdhAddEnglishCounterW")
	procCollectQueryData  = pdh.NewProc("PdhCollectQueryData")
	procGetFormattedValue = pdh.NewProc("PdhGetFormattedCounterValue")
	procCloseQuery        = pdh.NewProc("PdhCloseQuery")
)

// Der Zaehler, der die tatsaechliche Leistung in Prozent des Nenntakts
// angibt. Ueber 100 bedeutet Boost, unter 100 Drosselung oder Leerlauf.
const zaehlerLeistung = `\Processor Information(_Total)\% Processor Performance`

// PDH_FMT_DOUBLE, siehe pdh.h.
const fmtDouble = 0x00000200

// pdhWert entspricht PDH_FMT_COUNTERVALUE. Der Wert ist eine Union;
// bei fmtDouble liegt an dieser Stelle ein float64.
type pdhWert struct {
	Status uint32
	_      uint32 // Ausrichtung, damit der Wert auf 8 Byte faellt
	Wert   float64
}

// TaktMesser liest wiederholt den tatsaechlichen Prozessortakt.
//
// Bewusst als Struktur statt als Funktion: Ein PDH-Ratenzaehler braucht
// zwei Abfragen, um einen Wert zu liefern, die erste dient nur als
// Bezugspunkt. Wer das als Einzelfunktion baut, oeffnet fuer jeden
// Messpunkt eine neue Abfrage und bekommt jedes Mal eine Null zurueck.
type TaktMesser struct {
	abfrage uintptr
	zaehler uintptr
	// NenntaktMhz kommt aus Win32_Processor.MaxClockSpeed und ist der
	// Bezugswert, auf den sich die Prozentangabe bezieht.
	NenntaktMhz int
}

// NeuerTaktMesser oeffnet die Abfrage und nimmt gleich den ersten,
// verworfenen Messpunkt.
func NeuerTaktMesser(nenntaktMhz int) (*TaktMesser, error) {
	if nenntaktMhz <= 0 {
		return nil, fmt.Errorf("kein brauchbarer Nenntakt: %d MHz", nenntaktMhz)
	}

	var abfrage uintptr
	rc, _, _ := procOpenQuery.Call(0, 0, uintptr(unsafe.Pointer(&abfrage)))
	if rc != 0 {
		return nil, fmt.Errorf("PdhOpenQuery fehlgeschlagen: 0x%X", rc)
	}

	pfad, err := syscall.UTF16PtrFromString(zaehlerLeistung)
	if err != nil {
		procCloseQuery.Call(abfrage)
		return nil, err
	}

	var zaehler uintptr
	rc, _, _ = procAddEnglishCounter.Call(
		abfrage,
		uintptr(unsafe.Pointer(pfad)),
		0,
		uintptr(unsafe.Pointer(&zaehler)),
	)
	if rc != 0 {
		procCloseQuery.Call(abfrage)
		return nil, fmt.Errorf("Leistungsindikator nicht verfuegbar: 0x%X", rc)
	}

	m := &TaktMesser{abfrage: abfrage, zaehler: zaehler, NenntaktMhz: nenntaktMhz}

	// Erster Aufruf: setzt nur den Bezugspunkt, liefert noch nichts.
	procCollectQueryData.Call(abfrage)
	return m, nil
}

// Schliessen gibt die Abfrage frei. Fehlt der Aufruf, bleibt ein
// Handle offen, solange das Programm laeuft.
func (m *TaktMesser) Schliessen() {
	if m.abfrage != 0 {
		procCloseQuery.Call(m.abfrage)
		m.abfrage = 0
	}
}

// Messen liefert den aktuellen Takt in MHz und den Prozentwert.
//
// Zwischen zwei Aufrufen muessen mindestens einige hundert Millisekunden
// liegen, sonst hat der Zaehler keine neue Grundlage und meldet den
// alten Wert.
func (m *TaktMesser) Messen() (mhz int, prozent float64, err error) {
	rc, _, _ := procCollectQueryData.Call(m.abfrage)
	if rc != 0 {
		return 0, 0, fmt.Errorf("PdhCollectQueryData fehlgeschlagen: 0x%X", rc)
	}

	var wert pdhWert
	rc, _, _ = procGetFormattedValue.Call(
		m.zaehler,
		uintptr(fmtDouble),
		0,
		uintptr(unsafe.Pointer(&wert)),
	)
	if rc != 0 {
		return 0, 0, fmt.Errorf("PdhGetFormattedCounterValue fehlgeschlagen: 0x%X", rc)
	}

	prozent = wert.Wert
	mhz = int(float64(m.NenntaktMhz) * prozent / 100)
	return mhz, prozent, nil
}
