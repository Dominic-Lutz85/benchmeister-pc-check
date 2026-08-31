package pruefung

import (
	"strings"
	"testing"
)

// Ein einzelner Riegel und ein Riegel im falschen Kanal sind die zwei
// Faelle, deren Rat bei einem Notebook nicht befolgbar ist.
func nurEinRiegel() []Riegel {
	return []Riegel{{KapazitaetBytes: 8589934592, TaktMhz: 2400,
		TaktMhzZweiter: 2400, Kanal: "A"}}
}

func TestMobilTauschtDenRatAus(t *testing.T) {
	stationaer := FuerBauform(Speicher(nurEinRiegel()), false)
	mobil := FuerBauform(Speicher(nurEinRiegel()), true)

	var textStationaer, textMobil string
	for _, b := range stationaer {
		if strings.Contains(b.Titel, "ein Speicherriegel") {
			textStationaer = b.Empfehlung
		}
	}
	for _, b := range mobil {
		if strings.Contains(b.Titel, "ein Speicherriegel") {
			textMobil = b.Empfehlung
		}
	}

	if textStationaer == "" || textMobil == "" {
		t.Fatalf("Befund fehlt. stationaer=%q mobil=%q", textStationaer, textMobil)
	}
	if textStationaer == textMobil {
		t.Errorf("Bei einem Notebook muss der Rat ein anderer sein, war beide Male: %s", textMobil)
	}
	// Der springende Punkt: Nicht zum Kauf raten, ohne dass klar ist,
	// ob ueberhaupt ein Steckplatz frei ist.
	if !strings.Contains(strings.ToLower(textMobil), "nachsehen") &&
		!strings.Contains(strings.ToLower(textMobil), "prüfen") {
		t.Errorf("Der mobile Rat sollte zum Nachsehen auffordern, war: %s", textMobil)
	}
}

// Regel 1: Bei unbekannter Bauform darf sich nichts aendern. Die
// Gehaeuseart kommt aus einer Herstellertabelle und fehlt oft.
func TestOhneMobilBleibtAllesGleich(t *testing.T) {
	vorher := Speicher(nurEinRiegel())
	nachher := FuerBauform(vorher, false)
	if len(vorher) != len(nachher) {
		t.Fatalf("Zahl der Befunde hat sich geaendert: %d zu %d", len(vorher), len(nachher))
	}
	for i := range vorher {
		if vorher[i].Empfehlung != nachher[i].Empfehlung {
			t.Errorf("Befund %q wurde veraendert, obwohl nicht mobil", vorher[i].Titel)
		}
	}
}

/*
Waechter: Wer zum Nachruesten oder Umstecken raet, braucht eine mobile
Fassung.

ANLASS (31.08.2026): "PCGH_Jacky" im PCGH-Forum hat darauf hingewiesen,
dass Aufruest-Empfehlungen bei Notebooks anders behandelt werden
muessen. Ohne diesen Test faellt ein kuenftiger Befund mit demselben
Fehler nicht auf, denn er sieht voellig normal aus.
*/
func TestRatZumNachruestenHatMobileFassung(t *testing.T) {
	// Woerter, die eine koerperliche Handlung am Geraet verlangen.
	verdaechtig := []string{"ergänzen", "umstecken", "nachrüsten", "dazustecken"}

	riegelFaelle := [][]Riegel{
		nurEinRiegel(),
		// Zwei Riegel, beide im selben Kanal.
		{
			{KapazitaetBytes: 8589934592, TaktMhz: 3200, TaktMhzZweiter: 3200, Kanal: "A"},
			{KapazitaetBytes: 8589934592, TaktMhz: 3200, TaktMhzZweiter: 3200, Kanal: "A"},
		},
	}

	for _, riegel := range riegelFaelle {
		for _, b := range Speicher(riegel) {
			for _, wort := range verdaechtig {
				if !strings.Contains(strings.ToLower(b.Empfehlung), wort) {
					continue
				}
				if b.EmpfehlungMobil == "" {
					t.Errorf("Befund %q raet zu %q, hat aber keine mobile Fassung. "+
						"Bei einem Notebook ist der Speicher oft verlötet oder der "+
						"Steckplatz nicht erreichbar.", b.Titel, wort)
				}
			}
		}
	}
}

// Die 220-Zeichen-Grenze gilt auch fuer die mobile Fassung. Sie ist
// laenger formuliert und reisst die Grenze schneller.
func TestMobileFassungBleibtKurz(t *testing.T) {
	const grenze = 220

	riegelFaelle := [][]Riegel{
		nurEinRiegel(),
		{
			{KapazitaetBytes: 8589934592, TaktMhz: 3200, TaktMhzZweiter: 3200, Kanal: "A"},
			{KapazitaetBytes: 8589934592, TaktMhz: 3200, TaktMhzZweiter: 3200, Kanal: "A"},
		},
	}

	for _, riegel := range riegelFaelle {
		for _, b := range FuerBauform(Speicher(riegel), true) {
			laenge := len([]rune(b.Feststellung + " " + b.Empfehlung))
			if laenge > grenze {
				t.Errorf("Mobile Fassung von %q hat %d Zeichen, erlaubt sind %d:\n%s %s",
					b.Titel, laenge, grenze, b.Feststellung, b.Empfehlung)
			}
		}
	}
}
