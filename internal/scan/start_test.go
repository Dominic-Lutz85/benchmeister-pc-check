package scan

import "testing"

/*
Ruft jede Abfrage genau so auf, wie main.go es beim Start tut.

WARUM ES DIESEN TEST GIBT: Fassung 1.0.13 stuerzte auf jedem Rechner
beim Start ab. Ursache war ein Typ in einer WMI-Struktur:
ChassisTypes ist int32, deklariert war []uint16. Die Bibliothek wandelt
ueber Reflection um und bricht bei einem unerwarteten Typ mit einem
Panic ab, nicht mit einem Fehler.

Alle anderen Tests waren gruen. Sie pruefen die Auswertung, und die
bekommt ihre Werte als einfache Go-Strukturen uebergeben, damit sie
ohne Windows laufen. Genau diese Trennung, die die Pruefungen testbar
macht, hat die kaputte Abfrage verdeckt.

Gemeldet von "Misanthrop68" im PCGH-Forum, keine halbe Stunde nach der
Veroeffentlichung: "Bei mir bricht das Programm ohne weitere Angaben
beim Starten ab."

Der Test prueft bewusst KEINE Werte. Was hier herauskommt, haengt vom
Rechner ab. Geprueft wird nur, dass es ueberhaupt herauskommt.
*/
func TestAbfragenLaufenOhneAbsturz(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Absturz beim Auslesen, das Programm wuerde beim Start "+
				"abbrechen: %v", r)
		}
	}()

	bauform := BauformErmitteln()
	zustand := Systemzustand()
	ergebnis, err := Auslesen()

	if err != nil {
		t.Logf("Auslesen meldet einen Fehler, das ist erlaubt: %v", err)
	} else {
		t.Logf("CPU=%q GPU=%q", ergebnis.CPUName, ergebnis.GPUName)
	}
	t.Logf("Bauform=%d Platte=%v Netz=%v Schutz=%v",
		bauform, zustand.SystemplatteErkannt, zustand.NetzwerkErkannt,
		zustand.VirenschutzErkannt)
}
