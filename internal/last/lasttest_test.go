package last

import (
	"context"
	"runtime"
	"testing"
	"time"
)

// Die Schwellen entscheiden, ob jemand seinen Rechner aufschraubt.
// Verschiebt sie jemand, soll das auffallen und nicht durchrutschen.
func TestEinordnenTrifftDieSchwellen(t *testing.T) {
	faelle := []struct {
		name     string
		gehalten float64
		will     Beurteilung
	}{
		{"voller Takt", 1.0, BeurteilungStabil},
		{"knapp ueber der Stabil-Schwelle", 0.951, BeurteilungStabil},
		{"genau auf der Stabil-Schwelle", 0.95, BeurteilungStabil},
		{"knapp darunter", 0.949, BeurteilungLeichterAbfall},
		{"genau auf der unteren Schwelle", 0.85, BeurteilungLeichterAbfall},
		{"knapp darunter", 0.849, BeurteilungDrosselt},
		{"deutlicher Einbruch", 0.6, BeurteilungDrosselt},
	}

	for _, f := range faelle {
		t.Run(f.name, func(t *testing.T) {
			if got := Einordnen(f.gehalten); got != f.will {
				t.Errorf("%.3f ergab %q, erwartet %q", f.gehalten, got, f.will)
			}
		})
	}
}

/*
Der Dauertakt ist der MEDIAN des letzten Drittels, nicht der letzte
Messwert und nicht der Mittelwert.

Warum das einen eigenen Test verdient: Mit dem letzten Messwert haengt
der ganze Befund an einer einzigen Zahl, und ein Hintergrunddienst, der
in der letzten halben Sekunde anspringt, wuerde eine gesunde Kuehlung
als Drosselung melden. Genau die Sorte Fehlalarm, die laut CLAUDE.md
teurer ist als ein uebersehener Fall.
*/
func TestDauertaktIgnoriertEinzelneAusreisser(t *testing.T) {
	e := &Ergebnis{Verlauf: []Messpunkt{
		{Nach: 500 * time.Millisecond, Mhz: 4700},
		{Nach: 1 * time.Second, Mhz: 4690},
		{Nach: 1500 * time.Millisecond, Mhz: 4680},
		{Nach: 2 * time.Second, Mhz: 4670},
		{Nach: 2500 * time.Millisecond, Mhz: 4660},
		// Letztes Drittel: zwei gute Werte, ein Ausreisser.
		{Nach: 3 * time.Second, Mhz: 4650},
		{Nach: 3500 * time.Millisecond, Mhz: 1200},
		{Nach: 4 * time.Second, Mhz: 4655},
	}}

	auswerten(e)

	if e.SpitzeMhz != 4700 {
		t.Errorf("Spitze %d MHz, erwartet 4700", e.SpitzeMhz)
	}
	if e.DauerMhz != 4650 {
		t.Errorf("Dauertakt %d MHz, erwartet 4650 (Median), der Ausreisser hat durchgeschlagen", e.DauerMhz)
	}
	if Einordnen(e.GehaltenProzent) != BeurteilungStabil {
		t.Errorf("ein einzelner Einbruch hat das Urteil gekippt: %.3f", e.GehaltenProzent)
	}
}

func TestEchteDrosselungWirdErkannt(t *testing.T) {
	// Takt faellt ueber den Lauf deutlich ab, so sieht ein zugesetzter
	// Kuehler aus.
	e := &Ergebnis{Verlauf: []Messpunkt{
		{Mhz: 4700}, {Mhz: 4600}, {Mhz: 4300},
		{Mhz: 4000}, {Mhz: 3700}, {Mhz: 3500},
		{Mhz: 3400}, {Mhz: 3390}, {Mhz: 3380},
	}}

	auswerten(e)

	if Einordnen(e.GehaltenProzent) != BeurteilungDrosselt {
		t.Errorf("Abfall von 4700 auf 3380 MHz wurde als %q eingeordnet (%.3f)",
			Einordnen(e.GehaltenProzent), e.GehaltenProzent)
	}
}

func TestLeerenVerlaufUeberlebtDieAuswertung(t *testing.T) {
	e := &Ergebnis{}
	auswerten(e)
	if e.SpitzeMhz != 0 || e.DauerMhz != 0 || e.GehaltenProzent != 0 {
		t.Error("ohne Messwerte darf nichts behauptet werden")
	}
}

// Regel 3: Ein Kern bleibt frei, damit der Rechner bedienbar bleibt.
func TestEinKernBleibtFrei(t *testing.T) {
	n := runtime.NumCPU()
	got := arbeiter()

	if n > 2 && got != n-1 {
		t.Errorf("%d Kerne, belastet werden %d, erwartet %d", n, got, n-1)
	}
	if got < 1 {
		t.Error("es muss mindestens ein Arbeiter laufen")
	}
	if got > n {
		t.Errorf("mehr Arbeiter (%d) als Kerne (%d)", got, n)
	}
}

// Regel 1: Die Hoechstdauer ist eine harte Grenze, kein Vorschlag.
func TestHoechstdauerIstNichtUeberschreitbar(t *testing.T) {
	if Hoechstdauer > 3*time.Minute {
		t.Errorf("Hoechstdauer steht auf %v, das ist mehr als die zugesagten drei Minuten", Hoechstdauer)
	}
}

/*
Regel 4: Der Lauf muss sich jederzeit abbrechen lassen.

Der Test bricht nach 1,2 Sekunden ab und prueft, dass der Aufruf
tatsaechlich zurueckkommt, statt die volle Dauer zu laufen. Ohne diese
Zusicherung waere ein Abbruchknopf in der Oberflaeche eine Attrappe.
*/
func TestLaufBrichtAufZurufAb(t *testing.T) {
	if testing.Short() {
		t.Skip("kurzer Durchlauf, erzeugt bewusst keine Last")
	}

	ctx, stopp := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1200 * time.Millisecond)
		stopp()
	}()

	beginn := time.Now()
	e, err := Ausfuehren(ctx, 3800, 30*time.Second)
	gebraucht := time.Since(beginn)

	if err != nil {
		t.Skipf("Leistungsindikator auf diesem System nicht verfuegbar: %v", err)
	}
	if gebraucht > 5*time.Second {
		t.Errorf("Abbruch dauerte %v, der Lauf haengt", gebraucht)
	}
	if !e.Abgebrochen {
		t.Error("der Abbruch wurde nicht als solcher vermerkt")
	}
}
