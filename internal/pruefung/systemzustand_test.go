package pruefung

import (
	"strings"
	"testing"
)

// Die 220-Zeichen-Grenze gilt auch fuer die neuen Befunde, siehe
// TestBefundeBleibenKurz. Hier wird bewusst der schlimmste Fall
// geprueft: mehrere Schutzprogramme mit sehr langen Namen.
func TestZustandBefundeBleibenKurz(t *testing.T) {
	const grenze = 220

	faelle := []struct {
		name string
		z    Zustand
	}{
		{"volle Platte", Zustand{
			SystemplatteFreiGb: 8, SystemplatteGesamtGb: 931, SystemplatteErkannt: true}},
		{"winzige Platte", Zustand{
			SystemplatteFreiGb: 1, SystemplatteGesamtGb: 119, SystemplatteErkannt: true}},
		{"zwei Schutzprogramme", Zustand{
			VirenschutzErkannt: true,
			Virenschutz:        []string{"Windows-Sicherheit", "Avast Antivirus"}}},
		{"vier Schutzprogramme mit langen Namen", Zustand{
			VirenschutzErkannt: true,
			Virenschutz: []string{
				"Kaspersky Endpoint Security for Windows",
				"Bitdefender Antivirus Plus Total Security",
				"Malwarebytes Premium Real-Time Protection",
				"Windows-Sicherheit",
			}}},
		{"100 Mbit", Zustand{NetzwerkMbit: 100, NetzwerkErkannt: true}},
		{"10 Mbit", Zustand{NetzwerkMbit: 10, NetzwerkErkannt: true}},
	}

	for _, f := range faelle {
		for _, b := range AlleZustand(f.z) {
			laenge := len([]rune(b.Feststellung + " " + b.Empfehlung))
			if laenge > grenze {
				t.Errorf("%s: Befund %q hat %d Zeichen, erlaubt sind %d.\n%s %s",
					f.name, b.Titel, laenge, grenze, b.Feststellung, b.Empfehlung)
			}
		}
	}
}

// Regel 1 aus CLAUDE.md: Wo nichts sicher ablesbar ist, wird
// geschwiegen. Ein nicht gesetzter "Erkannt"-Schalter darf NIE zu
// einem Befund fuehren, sonst meldet das Programm auf jedem Rechner,
// auf dem eine Abfrage fehlschlaegt, einen Fehlalarm.
func TestOhneMessungKeinBefund(t *testing.T) {
	leer := Zustand{}
	if b := AlleZustand(leer); len(b) != 0 {
		t.Errorf("Ohne jede Messung darf es keinen Befund geben, waren aber %d: %+v", len(b), b)
	}

	// Der gefaehrlichste Einzelfall: Die Abfrage lief nicht, alle Werte
	// stehen auf 0. "0 von 0 GB frei" und "0 Mbit" duerfen nicht als
	// Alarm durchgehen.
	nullwerte := Zustand{
		SystemplatteFreiGb: 0, SystemplatteGesamtGb: 0,
		NetzwerkMbit: 0,
	}
	if b := AlleZustand(nullwerte); len(b) != 0 {
		t.Errorf("Nullwerte ohne Erkennung duerfen keinen Befund geben, waren aber %d", len(b))
	}
}

func TestSystemplatteSchwellen(t *testing.T) {
	faelle := []struct {
		name     string
		frei     int
		gesamt   int
		erwartet bool
	}{
		{"halb voll, alles gut", 250, 500, false},
		{"knapp ueber zehn Prozent", 60, 500, false},
		{"unter zehn Prozent", 40, 500, true},
		// Grosse Platte: zehn Prozent waeren 400 GB, das waere albern.
		// Hier greift die Regel nicht, weil absolut genug frei ist.
		{"grosse Platte, 300 GB frei", 300, 4000, false},
		// Kleine Platte: elf Prozent, aber absolut zu wenig fuer ein
		// Windows-Update.
		{"kleine Platte, nur 13 GB frei", 13, 119, true},
	}

	for _, f := range faelle {
		z := Zustand{
			SystemplatteFreiGb: f.frei, SystemplatteGesamtGb: f.gesamt,
			SystemplatteErkannt: true}
		got := len(Systemplatte(z)) > 0
		if got != f.erwartet {
			t.Errorf("%s (%d von %d GB): Befund=%v, erwartet=%v",
				f.name, f.frei, f.gesamt, got, f.erwartet)
		}
	}
}

// Ein einzelnes Schutzprogramm ist der Normalfall und kein Befund.
// Windows Defender laeuft auf praktisch jedem Rechner.
func TestEinSchutzprogrammIstKeinBefund(t *testing.T) {
	z := Zustand{VirenschutzErkannt: true, Virenschutz: []string{"Windows-Sicherheit"}}
	if b := Virenschutz(z); len(b) != 0 {
		t.Errorf("Ein einzelnes Schutzprogramm darf keinen Befund geben, war aber: %+v", b)
	}
}

// Der Befund darf NICHT behaupten, die Programme liefen gleichzeitig.
// Der Zustand ist aus dem undokumentierten Bitfeld nicht sicher
// ablesbar, siehe Kommentar in systemzustand.go.
func TestVirenschutzBehauptetKeinenZustand(t *testing.T) {
	z := Zustand{VirenschutzErkannt: true,
		Virenschutz: []string{"Windows-Sicherheit", "Avast Antivirus"}}
	b := Virenschutz(z)
	if len(b) != 1 {
		t.Fatalf("erwartet genau ein Befund, waren %d", len(b))
	}
	sichtbar := b[0].Feststellung + " " + b[0].Empfehlung
	for _, wort := range []string{"laufen gleichzeitig", "sind aktiv", "beide aktiv"} {
		if strings.Contains(strings.ToLower(sichtbar), wort) {
			t.Errorf("Der sichtbare Teil behauptet %q, das ist nicht messbar: %s", wort, sichtbar)
		}
	}
	if !strings.Contains(sichtbar, "kennt") {
		t.Errorf("Der Befund sollte von registrierten Programmen sprechen: %s", sichtbar)
	}
}

// Einzahl und Mehrzahl bei der Zusammenfassung. Der Probelauf am
// 31.08.2026 lieferte "und 1 weitere", was sich wie ein Tippfehler
// liest.
func TestNamenKurzGrammatik(t *testing.T) {
	drei := []string{"A", "B", "C"}
	if got := namenKurz(drei); !strings.Contains(got, "und ein weiteres") {
		t.Errorf("bei genau einem Rest erwartet \"und ein weiteres\", war: %s", got)
	}
	vier := []string{"A", "B", "C", "D"}
	if got := namenKurz(vier); !strings.Contains(got, "und 2 weitere") {
		t.Errorf("bei zwei Rest erwartet \"und 2 weitere\", war: %s", got)
	}
	zwei := []string{"A", "B"}
	if got := namenKurz(zwei); strings.Contains(got, "weiter") {
		t.Errorf("bei zwei Namen darf nichts zusammengefasst werden, war: %s", got)
	}
}

// Gigabit ist der Normalfall und kein Befund.
func TestGigabitIstKeinBefund(t *testing.T) {
	for _, mbit := range []int{1000, 2500, 10000} {
		z := Zustand{NetzwerkMbit: mbit, NetzwerkErkannt: true}
		if b := Netzwerk(z); len(b) != 0 {
			t.Errorf("%d Mbit/s darf keinen Befund geben, war aber: %+v", mbit, b)
		}
	}
}
