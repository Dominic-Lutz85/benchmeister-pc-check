package last

import (
	"context"
	"math"
	"runtime"
	"sort"
	"time"
)

/*
Der Lasttest selbst.

WAS ER MISST, UND WAS AUSDRUECKLICH NICHT
=========================================
Er misst, ob der Prozessor seinen Takt unter Dauerlast haelt. Er misst
NICHT, wie schnell der Rechner ist. Das ist ein Unterschied, den fast
jedes Werkzeug verwischt: Eine Punktzahl sagt vor allem, was gekauft
wurde, ein Taktverlauf sagt, wie es dem Rechner damit geht.

Faellt der Takt nach einer Minute deutlich ab, liegt es an Kuehlung,
Power-Limit, verstaubtem Kuehler oder alter Waermeleitpaste. Welches
davon, sagt der Verlauf nicht, aber dass etwas ist, sagt er sicher.

DIE SICHERHEITSREGELN, DIE HIER IM CODE STEHEN
==============================================
Ein Lasttest auf fremden Rechnern ist etwas anderes als einer auf dem
eigenen. Deshalb sind die Grenzen nicht Einstellungen, sondern fest
verdrahtet:

 1. HOECHSTDAUER. Laenger als drei Minuten laeuft nichts, auch wenn es
    jemand anders aufruft. Zwei Minuten reichen fuer die Aussage,
    laenger bringt nur Waerme.

 2. KEINE EXTREMLAST. Gerechnet wird gewoehnliche Gleitkomma- und
    Ganzzahlarbeit. Kein AVX-512, keine handgeschriebenen
    Vektorschleifen, nichts, was auf ein Power-Limit zielt. Leitsatz:
    Der Test darf nicht haerter sein als ein anspruchsvolles Spiel.
    Dann geht niemand ein Risiko ein, das er nicht ohnehin taeglich
    eingeht.

 3. EIN KERN BLEIBT FREI, sofern mehr als zwei da sind. Damit bleibt
    der Rechner bedienbar und der Abbruchknopf erreichbar. Das kostet
    ein paar Prozent Aussagekraft und ist es wert.

 4. JEDERZEIT ABBRECHBAR ueber den Context. Kein Aufruf ohne
    Abbruchmoeglichkeit, deshalb ist ctx der erste Parameter und nicht
    optional.

 5. NICHTS WIRD EINGESTELLT. Kein Uebertakten, keine Spannung, keine
    Luefterkurve. Dieses Paket schreibt nichts ins System. Das ist die
    Grenze, hinter der aus einer Messung ein Eingriff wird, und die
    wird nicht ueberschritten.
*/

// Hoechstdauer, siehe Regel 1. Bewusst als Konstante und nicht als
// Parameter mit Vorgabewert: Was aenderbar ist, wird geaendert.
const Hoechstdauer = 3 * time.Minute

// Abstand zwischen zwei Taktmessungen. Unter 500 ms liefert der
// PDH-Zaehler keinen neuen Wert, siehe takt.go.
const messabstand = 500 * time.Millisecond

// Messpunkt ist ein Taktwert mit dem Zeitpunkt seit Beginn.
type Messpunkt struct {
	Nach time.Duration
	Mhz  int
}

// Ergebnis eines Lastlaufs.
type Ergebnis struct {
	Verlauf []Messpunkt

	// SpitzeMhz ist der hoechste gemessene Takt, also das, was der
	// Prozessor kann, wenn er noch kalt ist.
	SpitzeMhz int

	// DauerMhz ist der Median der letzten Drittelmessungen, also der
	// Takt, den er unter Dauerlast wirklich haelt. Median statt
	// Mittelwert, weil ein einzelner Einbruch durch einen
	// Hintergrunddienst sonst das Ergebnis faerbt.
	DauerMhz int

	// GehaltenProzent ist DauerMhz im Verhaeltnis zu SpitzeMhz. Der
	// eigentliche Befund.
	GehaltenProzent float64

	// Abgebrochen ist wahr, wenn der Lauf vorzeitig endete.
	Abgebrochen bool
}

// Beurteilung ist die Einordnung des gehaltenen Takts.
type Beurteilung string

const (
	// BeurteilungStabil: haelt praktisch den vollen Takt.
	BeurteilungStabil Beurteilung = "stabil"
	// BeurteilungLeichterAbfall: normal bei vielen Serienkuehlern.
	BeurteilungLeichterAbfall Beurteilung = "leichter Abfall"
	// BeurteilungDrosselt: da geht Leistung verloren.
	BeurteilungDrosselt Beurteilung = "drosselt deutlich"
)

/*
Die Schwellen.

Ein gewisser Rueckgang ist NORMAL und kein Mangel: Prozessoren boosten
im kalten Zustand hoeher, als sie dauerhaft halten koennen, das ist
Absicht des Herstellers. Erst ein deutlicher Abfall zeigt ein Problem.

Deshalb liegt die untere Schwelle bei 85 Prozent und nicht bei 95. Ein
Fehlalarm kostet mehr als ein uebersehener Fall, siehe die Regeln in
CLAUDE.md.
*/
const (
	schwelleStabil = 0.95
	schwelleLeicht = 0.85
)

// Einordnen uebersetzt den gehaltenen Anteil in einen Befund.
func Einordnen(gehalten float64) Beurteilung {
	switch {
	case gehalten >= schwelleStabil:
		return BeurteilungStabil
	case gehalten >= schwelleLeicht:
		return BeurteilungLeichterAbfall
	default:
		return BeurteilungDrosselt
	}
}

// arbeiten erzeugt Last, bis der Context endet.
//
// Bewusst gewoehnliche Rechenarbeit mit einer Abhaengigkeitskette:
// Jeder Schritt braucht das Ergebnis des vorherigen, damit der
// Uebersetzer die Schleife nicht wegoptimiert. Kein Vektorcode, siehe
// Regel 2.
func arbeiten(ctx context.Context) {
	x := 1.000001
	i := 0
	for {
		// Nicht bei jedem Durchlauf nachsehen, das waere teurer als die
		// Rechnung selbst.
		if i%4096 == 0 {
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
		x = math.Sqrt(x*1.0000003 + 1.0000001)
		if x > 1e9 {
			x = 1.000001
		}
		i++
	}
}

// arbeiter sagt, wie viele Kerne belastet werden. Siehe Regel 3.
func arbeiter() int {
	n := runtime.NumCPU()
	if n > 2 {
		return n - 1
	}
	return n
}

/*
Ausfuehren laesst den Test laufen.

nenntaktMhz kommt aus dem Scan (Win32_Processor.MaxClockSpeed). Ohne
brauchbaren Nenntakt gibt es keine Messung, denn der Zaehler liefert nur
Prozentwerte, und ein Prozentwert ohne Bezug ist keine Zahl.

Der Aufrufer ist dafuer zustaendig, VORHER zu fragen. Dieses Paket
startet nichts von selbst, aber es kann auch nicht pruefen, ob jemand
zugestimmt hat.
*/
func Ausfuehren(ctx context.Context, nenntaktMhz int, dauer time.Duration) (*Ergebnis, error) {
	if dauer > Hoechstdauer {
		dauer = Hoechstdauer
	}

	messer, err := NeuerTaktMesser(nenntaktMhz)
	if err != nil {
		return nil, err
	}
	defer messer.Schliessen()

	laufCtx, stopp := context.WithTimeout(ctx, dauer)
	defer stopp()

	for i := 0; i < arbeiter(); i++ {
		go arbeiten(laufCtx)
	}

	ergebnis := &Ergebnis{}
	beginn := time.Now()
	uhr := time.NewTicker(messabstand)
	defer uhr.Stop()

	for {
		select {
		case <-laufCtx.Done():
			// Vom Aufrufer abgebrochen, nicht durch das Zeitlimit?
			ergebnis.Abgebrochen = ctx.Err() != nil
			auswerten(ergebnis)
			return ergebnis, nil
		case <-uhr.C:
			mhz, _, err := messer.Messen()
			if err != nil {
				continue
			}
			ergebnis.Verlauf = append(ergebnis.Verlauf, Messpunkt{
				Nach: time.Since(beginn).Round(time.Millisecond),
				Mhz:  mhz,
			})
		}
	}
}

// auswerten fuellt Spitze, Dauertakt und den gehaltenen Anteil.
func auswerten(e *Ergebnis) {
	if len(e.Verlauf) == 0 {
		return
	}

	for _, p := range e.Verlauf {
		if p.Mhz > e.SpitzeMhz {
			e.SpitzeMhz = p.Mhz
		}
	}

	// Letztes Drittel als Dauerlast-Fenster, mindestens ein Wert.
	ab := len(e.Verlauf) * 2 / 3
	if ab >= len(e.Verlauf) {
		ab = len(e.Verlauf) - 1
	}
	letzte := make([]int, 0, len(e.Verlauf)-ab)
	for _, p := range e.Verlauf[ab:] {
		letzte = append(letzte, p.Mhz)
	}
	sort.Ints(letzte)
	e.DauerMhz = letzte[len(letzte)/2]

	if e.SpitzeMhz > 0 {
		e.GehaltenProzent = float64(e.DauerMhz) / float64(e.SpitzeMhz)
	}
}
