package pruefung

import (
	"strings"
	"testing"
)

// Die Teilenummern hier sind echt und stammen vom Entwicklungsrechner
// (ausgelesen am 27.08.2026), nicht ausgedacht. Genau an ihnen hat sich
// gezeigt, dass Windows die Sollgeschwindigkeit nicht meldet und sie nur
// in der Teilenummer steht.
func TestSollTaktAusTeilenummer(t *testing.T) {
	faelle := []struct {
		teilenummer string
		erwartet    int
		warum       string
	}{
		{"CMK32GX4M2E3200C16", 3200, "Corsair DDR4-3200, echte Teilenummer"},
		{"F4-3200C16-8GVKB", 3200, "G.Skill DDR4-3200, echte Teilenummer"},
		{"CMK32GX5M2B6000C36", 6000, "Corsair DDR5-6000"},
		{"KF556C36-16", 0, "kein vierstelliger Takt erkennbar"},
		{"", 0, "leere Teilenummer"},
		{"BELIEBIG1234XYZ", 0, "1234 ist keine bekannte Taktstufe"},
	}
	for _, f := range faelle {
		if got := sollTaktAusTeilenummer(f.teilenummer); got != f.erwartet {
			t.Errorf("%s: sollTaktAusTeilenummer(%q) = %d, erwartet %d",
				f.warum, f.teilenummer, got, f.erwartet)
		}
	}
}

// Waechter gegen zu gierige Erkennung: Enthaelt eine Teilenummer mehrere
// bekannte Taktstufen, ist keine davon sicher. Dann lieber gar kein
// Befund als ein falscher.
func TestMehrdeutigeTeilenummerGibtNichts(t *testing.T) {
	if got := sollTaktAusTeilenummer("XX3200YY3600ZZ"); got != 0 {
		t.Errorf("mehrdeutige Teilenummer sollte 0 ergeben, ergab %d", got)
	}
}

func TestSpeicherErkenntAbgeschaltetesProfil(t *testing.T) {
	// Genau die Bestueckung des Entwicklungsrechners: vier Riegel,
	// alle fuer 3200 gebaut, alle mit 2133 laufend.
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "P0 CHANNEL A"},
		{KapazitaetBytes: 8589934592, TaktMhz: 2133, Teilenummer: "F4-3200C16-8GVKB", Kanal: "P0 CHANNEL A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "P0 CHANNEL B"},
		{KapazitaetBytes: 8589934592, TaktMhz: 2133, Teilenummer: "F4-3200C16-8GVKB", Kanal: "P0 CHANNEL B"},
	}
	befunde := Speicher(riegel)

	var profil, gemischt bool
	for _, b := range befunde {
		if b.Titel == "Arbeitsspeicher läuft womöglich unter seiner Sollgeschwindigkeit" {
			profil = true
			if b.Schwere != Hinweis {
				t.Error("abgeschaltetes Speicherprofil muss ein Hinweis sein")
			}
		}
		if b.Titel == "Verschiedene Speicherriegel gemischt" {
			gemischt = true
		}
	}
	if !profil {
		t.Error("das abgeschaltete Speicherprofil wurde nicht erkannt")
	}
	if !gemischt {
		t.Error("die gemischten Riegel wurden nicht erkannt")
	}
}

func TestSpeicherSchweigtWennAllesPasst(t *testing.T) {
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "B"},
	}
	if befunde := Speicher(riegel); len(befunde) != 0 {
		t.Errorf("bei einwandfreier Bestueckung darf nichts gemeldet werden, kam: %+v", befunde)
	}
}

func TestSpeicherErkenntEinzelnenRiegel(t *testing.T) {
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "A"},
	}
	befunde := Speicher(riegel)
	if len(befunde) != 1 || befunde[0].Titel != "Nur ein Speicherriegel verbaut" {
		t.Errorf("Einkanalbetrieb wurde nicht erkannt, kam: %+v", befunde)
	}
}

func TestSpeicherOhneTeilenummerMeldetKeinProfil(t *testing.T) {
	// Manche Riegel melden gar keine Teilenummer. Dann ist der Soll-Wert
	// unbekannt, und es darf nichts behauptet werden.
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, Teilenummer: "", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, Teilenummer: "", Kanal: "B"},
	}
	if befunde := Speicher(riegel); len(befunde) != 0 {
		t.Errorf("ohne Teilenummer darf nichts gemeldet werden, kam: %+v", befunde)
	}
}

func TestLaufwerke(t *testing.T) {
	// Bestueckung des Entwicklungsrechners plus zwei erfundene Faelle.
	laufwerke := []Laufwerk{
		{Name: "VendorCo ProductCode", MedienArt: 0, BusArt: busUSB, Bytes: 62914560000},
		{Name: "KINGSTON SKC3000S1024G", MedienArt: medienSSD, BusArt: busNVMe},
		{Name: "Samsung SSD 970 EVO Plus 2TB", MedienArt: medienSSD, BusArt: busNVMe},
		{Name: "Crucial BX500", MedienArt: medienSSD, BusArt: busSATA},
		{Name: "WDC WD10EZEX", MedienArt: medienHDD, BusArt: busSATA},
	}
	befunde := Laufwerke(laufwerke)

	if len(befunde) != 2 {
		t.Fatalf("erwartet: genau zwei Befunde (SATA-SSD und Festplatte), kam: %+v", befunde)
	}
	// USB-Geraete und NVMe-SSDs duerfen nichts ausloesen.
	for _, b := range befunde {
		if b.Titel != "SSD hängt am SATA-Anschluss" && b.Titel != "Magnetfestplatte verbaut" {
			t.Errorf("unerwarteter Befund: %s", b.Titel)
		}
	}
}

// Waechter fuer die Ueberschrift der Ergebnisseite. Dort steht ueber dem
// einen Block "ohne dass du etwas kaufen musst". Jeder Laufwerks-Befund
// laeuft aber auf einen Neukauf hinaus, keiner davon darf also in diesem
// Block landen. Genau das war der Fehler, den ein PCGH-Moderator am
// 28.08.2026 angestrichen hat.
func TestLaufwerksbefundeSindZukauf(t *testing.T) {
	laufwerke := []Laufwerk{
		{Name: "Crucial BX500", MedienArt: medienSSD, BusArt: busSATA},
		{Name: "WDC WD10EZEX", MedienArt: medienHDD, BusArt: busSATA},
	}
	for _, b := range Laufwerke(laufwerke) {
		if b.Schwere != Zukauf {
			t.Errorf("%q muss als Zukauf gelten, ist aber %q", b.Titel, b.Schwere)
		}
	}
}

// DER FALL, DER DAS GANZE AUSGELOEST HAT, und die Korrektur dazu.
//
// Am 28.08.2026 meldete ein Moderator im PCGH-Forum einen Fehlalarm:
// Sein Speicher lief mit 5600 MT/s, das Programm behauptete 4800.
//
// Die erste Reparatur war, bei unterschiedlichen Taktfeldern gar nichts
// mehr zu sagen. Das war fachlich falsch. Ein Nutzer (NullPointerEx) hat
// im selben Faden erklaert, dass die beiden Felder gar nicht dasselbe
// meinen: Speed ist der JEDEC-Grundtakt OHNE Profil,
// ConfiguredClockSpeed der real eingestellte. Ein Unterschied ist also
// kein Widerspruch, sondern der Normalfall bei aktivem XMP.
//
// Seit 29.08.2026 gilt deshalb: ConfiguredClockSpeed zaehlt, Speed ist
// nur der Rueckfall.
func TestAktivesProfilErzeugtKeinenFehlalarm(t *testing.T) {
	// Genau der Fall aus dem Forum: JEDEC 4800, tatsaechlich 5600, und
	// die Teilenummer nennt 5200. Der Speicher laeuft ueber seinem
	// Profil, es darf also KEIN Hinweis kommen.
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 4800, TaktMhzZweiter: 5600,
			Teilenummer: "TESTKIT5200C40", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 4800, TaktMhzZweiter: 5600,
			Teilenummer: "TESTKIT5200C40", Kanal: "B"},
	}
	for _, b := range Speicher(riegel) {
		if b.Schwere == Hinweis {
			t.Errorf("bei aktivem Profil darf kein Hinweis kommen, kam: %q", b.Titel)
		}
	}
}

// Die Gegenprobe: Profil AUS, beide Felder nennen den Grundtakt, die
// Teilenummer nennt mehr. Genau dafuer gibt es dieses Programm.
func TestAbgeschaltetesProfilWirdWeiterhinErkannt(t *testing.T) {
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 4800, TaktMhzZweiter: 4800,
			Teilenummer: "TESTKIT5200C40", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 4800, TaktMhzZweiter: 4800,
			Teilenummer: "TESTKIT5200C40", Kanal: "B"},
	}
	var gefunden bool
	for _, b := range Speicher(riegel) {
		if b.Schwere == Hinweis {
			gefunden = true
			if !strings.Contains(b.Feststellung, "4800") {
				t.Error("der Befund muss den gemeldeten Ist-Takt nennen")
			}
		}
	}
	if !gefunden {
		t.Error("abgeschaltetes Profil muss weiterhin einen Hinweis ergeben")
	}
}

// Fehlt ConfiguredClockSpeed ganz (manche Boards fuellen das Feld
// nicht), muss Speed als Rueckfall einspringen. Ohne diesen Rueckfall
// waere der wichtigste Befund auf solchen Rechnern tot.
func TestOhneZweitesFeldGiltSpeed(t *testing.T) {
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, TaktMhzZweiter: 0,
			Teilenummer: "CMK32GX4M2E3200C16", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, TaktMhzZweiter: 0,
			Teilenummer: "CMK32GX4M2E3200C16", Kanal: "B"},
	}
	var gefunden bool
	for _, b := range Speicher(riegel) {
		if b.Schwere == Hinweis && strings.Contains(b.Feststellung, "2133") {
			gefunden = true
		}
	}
	if !gefunden {
		t.Error("ohne ConfiguredClockSpeed muss Speed als Rueckfall greifen")
	}
}

// Die uebrigen Speicherpruefungen haengen nicht am Takt und muessen auch
// dann laufen, wenn zum Takt nichts zu sagen ist.
func TestTaktpruefungStopptDieUebrigenPruefungenNicht(t *testing.T) {
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 4800, TaktMhzZweiter: 5600,
			Teilenummer: "TESTKIT5200C40", Kanal: "A"},
	}
	var einzeln bool
	for _, b := range Speicher(riegel) {
		if b.Titel == "Nur ein Speicherriegel verbaut" {
			einzeln = true
		}
	}
	if !einzeln {
		t.Error("der Einkanal-Befund muss auch bei unklarem Takt kommen")
	}
}

// Uebereinstimmende Felder heissen NICHT, dass der Wert stimmt: Auch
// beide koennen denselben falschen Grundtakt melden. Deshalb muss der
// Befund den Weg zum Gegenpruefen nennen, und er darf nicht auf den
// Task-Manager verweisen, der dieselbe Tabelle liest und den Fehler nur
// bestaetigen wuerde.
func TestTaktbefundNenntDieGegenpruefung(t *testing.T) {
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, TaktMhzZweiter: 2133,
			Teilenummer: "CMK32GX4M2E3200C16", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, TaktMhzZweiter: 2133,
			Teilenummer: "CMK32GX4M2E3200C16", Kanal: "B"},
	}
	befunde := Speicher(riegel)
	if len(befunde) != 1 {
		t.Fatalf("erwartet: genau der Takt-Befund, kam: %+v", befunde)
	}
	if !strings.Contains(befunde[0].Empfehlung, "CPU-Z") {
		t.Error("der Befund muss sagen, womit man ihn nachpruefen kann")
	}
	if !strings.Contains(befunde[0].Empfehlung, "Task-Manager") {
		t.Error("der Befund muss davor warnen, dass der Task-Manager hier nichts beweist")
	}
	if strings.Contains(befunde[0].Feststellung, "laufen aber mit") {
		t.Error("der Befund darf den Ist-Takt nicht als Tatsache behaupten, wir kennen ihn nicht sicher")
	}
}

/*
 * Der Halbe-Rate-Fall, gefunden bei der Durchsicht fremder Foren am
 * 29.08.2026: Windows meldet den Speichertakt auf manchen Rechnern
 * halbiert (1600 statt 3200), waehrend CPU-Z den vollen Wert zeigt.
 * Ohne Sonderbehandlung haette das Programm daraus einen sehr lauten,
 * sehr falschen Befund gemacht.
 */
func TestSpeicherHalbeRate(t *testing.T) {
	// CMK32GX4M2E3200C16 heisst Soll 3200.
	const teil = "CMK32GX4M2E3200C16"

	faelle := []struct {
		name          string
		ist           uint32
		halbeErwartet bool
		zuLangsamErw  bool
	}{
		{"genau die Haelfte", 1600, true, false},
		{"knapp neben der Haelfte", 1700, true, false},
		{"deutlich drunter, aber nicht die Haelfte", 1200, false, true},
		{"echter XMP-Fall", 2133, false, true},
		{"alles in Ordnung", 3200, false, false},
	}

	for _, f := range faelle {
		t.Run(f.name, func(t *testing.T) {
			befunde := Speicher([]Riegel{
				{TaktMhzZweiter: f.ist, Teilenummer: teil},
				{TaktMhzZweiter: f.ist, Teilenummer: teil},
			})

			var halbe, zuLangsam bool
			for _, b := range befunde {
				switch b.Titel {
				case "Windows meldet genau den halben Speichertakt":
					halbe = true
				case "Arbeitsspeicher läuft womöglich unter seiner Sollgeschwindigkeit":
					zuLangsam = true
				}
			}
			if halbe != f.halbeErwartet {
				t.Errorf("Halbe-Rate-Befund: %v erwartet, %v bekommen", f.halbeErwartet, halbe)
			}
			if zuLangsam != f.zuLangsamErw {
				t.Errorf("Zu-langsam-Befund: %v erwartet, %v bekommen", f.zuLangsamErw, zuLangsam)
			}
			// Nie beide gleichzeitig: Zwei Meldungen zur selben Zahl, die
			// Verschiedenes behaupten, sind schlimmer als eine.
			if halbe && zuLangsam {
				t.Error("beide Takt-Befunde gleichzeitig, das widerspricht sich")
			}
		})
	}
}

// Der Sonderfall darf die uebrigen Pruefungen nicht verschlucken. Das
// waere beim ersten Entwurf beinahe passiert, dort stand ein return.
func TestHalbeRateVerschlucktDenRestNicht(t *testing.T) {
	befunde := Speicher([]Riegel{
		{TaktMhzZweiter: 1600, Teilenummer: "CMK32GX4M2E3200C16"},
	})

	var einRiegel bool
	for _, b := range befunde {
		if b.Titel == "Nur ein Speicherriegel verbaut" {
			einRiegel = true
		}
	}
	if !einRiegel {
		t.Error("Der Einkanal-Befund fehlt, der Sonderfall hat ihn verschluckt")
	}
}
