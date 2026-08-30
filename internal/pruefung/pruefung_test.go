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

// Bis zum 29.08.2026 hiess dieser Test "SchweigtWennAllesPasst" und
// verlangte GAR keinen Befund. Die Erwartung hat sich bewusst gedreht,
// nach einem Vorschlag von "Misanthrop68" im PCGH-Forum: Schweigen laesst
// sich von Uebersehen nicht unterscheiden. Was gleich bleibt: Es darf
// kein HINWEIS kommen, also nichts, was nach Handlungsbedarf aussieht.
func TestSpeicherBestaetigtWennAllesPasst(t *testing.T) {
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "B"},
	}
	befunde := Speicher(riegel)
	if len(befunde) != 1 || befunde[0].Schwere != Bestaetigung {
		t.Fatalf("erwartet genau eine Bestaetigung, kam: %+v", befunde)
	}
	for _, b := range befunde {
		if b.Schwere == Hinweis {
			t.Errorf("bei einwandfreier Bestueckung darf kein Hinweis kommen: %q", b.Titel)
		}
	}
}

// Laeuft der Speicher ueber seinem JEDEC-Grundtakt, hat das Board mehr
// eingestellt als die Werksvorgabe. Nur DANN darf von einem Profil die
// Rede sein. Genau der Fall des Moderators: Speed 4800, konfiguriert 5600.
func TestSpeicherNenntDasProfilNurWennErkennbar(t *testing.T) {
	mitProfil := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 4800, TaktMhzZweiter: 5600, Teilenummer: "CMK32GX5M2B5600C36", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 4800, TaktMhzZweiter: 5600, Teilenummer: "CMK32GX5M2B5600C36", Kanal: "B"},
	}
	befunde := Speicher(mitProfil)
	if len(befunde) != 1 || befunde[0].Schwere != Bestaetigung {
		t.Fatalf("erwartet genau eine Bestaetigung, kam: %+v", befunde)
	}
	if !strings.Contains(befunde[0].Hintergrund, "4800") {
		t.Errorf("der Grundtakt ohne Profil muss genannt werden: %q", befunde[0].Hintergrund)
	}
	if !strings.Contains(befunde[0].Hintergrund, "EXPO") {
		t.Errorf("bei erkennbarem Profil muss es benannt werden: %q", befunde[0].Hintergrund)
	}

	// Gegenprobe: Ohne Unterschied zwischen den Feldern ist nicht
	// feststellbar, ob ein Profil laeuft. Dann darf keines behauptet
	// werden, auch nicht vorsichtig.
	ohneProfil := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "B"},
	}
	befunde = Speicher(ohneProfil)
	if len(befunde) != 1 {
		t.Fatalf("erwartet genau eine Bestaetigung, kam: %+v", befunde)
	}
	for _, wort := range []string{"XMP", "EXPO", "D.O.C.P"} {
		if strings.Contains(befunde[0].Hintergrund, wort) {
			t.Errorf("ohne Beleg darf %q nicht behauptet werden: %q", wort, befunde[0].Hintergrund)
		}
	}
}

func TestSpeicherErkenntEinzelnenRiegel(t *testing.T) {
	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 6000, Teilenummer: "CMK32GX5M2B6000C36", Kanal: "A"},
	}
	var einkanal bool
	for _, b := range Speicher(riegel) {
		if b.Titel == "Nur ein Speicherriegel verbaut" {
			einkanal = true
		}
	}
	if !einkanal {
		t.Errorf("Einkanalbetrieb wurde nicht erkannt, kam: %+v", Speicher(riegel))
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
	// Seit dem 30.08.2026 steht die Gegenpruefung im Hintergrund, also
	// im zugeklappten Teil. Sie muss trotzdem DA sein: Sie ist der
	// Grund, warum dieses Programm nach einem Fehlalarm bei einem
	// PCGH-Moderator nicht mehr behauptet, sondern zum Nachpruefen
	// auffordert. Verschieben ja, weglassen nie.
	if !strings.Contains(befunde[0].Hintergrund, "CPU-Z") {
		t.Error("der Befund muss sagen, womit man ihn nachpruefen kann")
	}
	if !strings.Contains(befunde[0].Hintergrund, "Task-Manager") {
		t.Error("der Befund muss davor warnen, dass der Task-Manager hier nichts beweist")
	}
	// Und die Empfehlung muss ohne den Hintergrund handlungsfaehig
	// machen. Wer nicht aufklappt, muss trotzdem wissen was zu tun ist.
	if !strings.Contains(befunde[0].Empfehlung, "XMP") {
		t.Error("die Empfehlung muss auch ohne Aufklappen sagen, was zu tun ist")
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

/*
 * Die Laengenbremse, angelegt am 30.08.2026.
 *
 * ANLASS: "Hellhammer" im PCGH-Forum, sinngemaess: zu viel begruendet,
 * zu ausformuliert, das sei laecherlich weil unnoetig. Nachgemessen gab
 * ihm recht. Der laengste Befund hatte 590 Zeichen, also fuenf Saetze
 * fuer einen einzigen Fund. Bei drei Funden las man eine halbe Seite.
 *
 * WARUM ALS TEST UND NICHT ALS VORSATZ: Die Texte sind nicht an einem
 * Tag lang geworden. Jeder einzelne Satz kam als Reaktion auf eine
 * berechtigte Ruecksmeldung dazu, jeder war fuer sich richtig, und
 * niemand hat je die Summe angesehen. Genau so wachsen solche Texte
 * wieder, wenn nichts sie aufhaelt.
 *
 * Die Grenze gilt nur fuer den SICHTBAREN Teil. Der Hintergrund darf
 * lang sein, er ist zugeklappt und wird nur gelesen, wenn jemand es
 * wissen will.
 */
func TestBefundeBleibenKurz(t *testing.T) {
	// 220 Zeichen sind etwa zwei Saetze. Wer mehr braucht, hat etwas
	// im sichtbaren Teil stehen, das in den Hintergrund gehoert.
	const grenze = 220

	riegel := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 2133, TaktMhzZweiter: 2133,
			Teilenummer: "CMK32GX4M2E3200C16", Kanal: "A"},
	}
	laufwerke := []Laufwerk{
		{Name: "ST2000DM008", MedienArt: 3, BusArt: 11, Bytes: 2000398934016},
		{Name: "Samsung 870 EVO", MedienArt: 4, BusArt: 11, Bytes: 1000204886016},
	}
	alle := append(Speicher(riegel), Laufwerke(laufwerke)...)
	if len(alle) < 3 {
		t.Fatalf("zu wenige Befunde zum Pruefen, kam: %d", len(alle))
	}

	for _, b := range alle {
		sichtbar := len([]rune(b.Feststellung)) + len([]rune(b.Empfehlung))
		if sichtbar > grenze {
			t.Errorf(
				"%q: sichtbarer Teil ist %d Zeichen, erlaubt sind %d. "+
					"Was nicht zum Handeln noetig ist, gehoert in Hintergrund.",
				b.Titel, sichtbar, grenze)
		}
		// Umgekehrt: Ein Befund ohne Empfehlung ist keiner.
		if b.Empfehlung == "" {
			t.Errorf("%q hat keine Empfehlung", b.Titel)
		}
	}
}

/*
 * Kanal-Erkennung, angelegt am 30.08.2026 auf Vorschlag von "cryon1c"
 * im PCGH-Forum.
 *
 * SEIN PUNKT: Zwei Riegel koennen beide im selben Kanal stecken. Dann
 * laeuft der Speicher im Einkanalbetrieb, obwohl zwei drinstecken, und
 * "dies passiert auch oft und hat einen noch groesseren Einfluss als
 * ein fehlendes XMP-Profil". Bis dahin hat das Programm nur GEZAEHLT
 * und diesen Fall durchgelassen.
 *
 * DER SCHWIERIGE TEIL ist nicht das Erkennen, sondern das Schweigen.
 * Die Beschriftung kommt vom Board-Hersteller und ist oft nichtssagend.
 * "BANK 0" und "BANK 1" sehen nach zwei Kanaelen aus, sind auf vielen
 * Boards aber zwei Slots DESSELBEN Kanals. Wer das als Kanal liest,
 * meldet Einkanalbetrieb, wo keiner ist, und produziert genau den
 * Fehlalarm, der uns beim Speichertakt schon zweimal blamiert hat.
 * Deshalb hat kanalAus() zwei Rueckgaben und nicht eine.
 */
func TestKanalAus(t *testing.T) {
	faelle := []struct {
		eingabe string
		kanal   string
		erkannt bool
		warum   string
	}{
		{"P0 CHANNEL A", "A", true, "die haeufigste Schreibweise"},
		{"P0 CHANNEL B", "B", true, "zweiter Kanal"},
		{"Channel C", "C", true, "Kleinschreibung und ohne Praefix"},
		{"  P1 CHANNEL D  ", "D", true, "mit Leerzeichen drumherum"},
		{"BANK 0", "", false, "nennt keinen Kanal, sondern einen Slot"},
		{"BANK 1", "", false, "dito, und sieht truegerisch nach Kanal 2 aus"},
		{"", "", false, "manche Boards liefern gar nichts"},
		{"DIMM_A1", "", false, "Slot-Bezeichnung, kein Kanal"},
		{"CHANNEL", "", false, "das Wort allein sagt nichts"},
		{"CHANNEL 1", "", false, "Ziffer statt Buchstabe, nicht auswertbar"},
	}
	for _, f := range faelle {
		k, ok := kanalAus(f.eingabe)
		if ok != f.erkannt || k != f.kanal {
			t.Errorf("%q (%s): erwartet %q/%v, kam %q/%v",
				f.eingabe, f.warum, f.kanal, f.erkannt, k, ok)
		}
	}
}

func TestAlleRiegelImSelbenKanal(t *testing.T) {
	imSelben := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 3200, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "P0 CHANNEL A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 3200, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "P0 CHANNEL A"},
	}
	var gefunden bool
	for _, b := range Speicher(imSelben) {
		if b.Titel == "Alle Riegel stecken im selben Kanal" {
			gefunden = true
			if !strings.Contains(b.Feststellung, "Kanal A") {
				t.Errorf("der Befund muss den Kanal nennen: %q", b.Feststellung)
			}
			if !strings.Contains(b.Empfehlung, "umstecken") {
				t.Errorf("die Empfehlung muss das Umstecken nennen: %q", b.Empfehlung)
			}
		}
	}
	if !gefunden {
		t.Error("zwei Riegel im selben Kanal wurden nicht erkannt")
	}
}

func TestRichtigVerteilteRiegelBleibenStill(t *testing.T) {
	verteilt := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 3200, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "P0 CHANNEL A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 3200, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "P0 CHANNEL B"},
	}
	for _, b := range Speicher(verteilt) {
		if b.Titel == "Alle Riegel stecken im selben Kanal" {
			t.Error("bei sauber verteilten Riegeln darf nichts gemeldet werden")
		}
	}
}

/*
 * Die wichtigste Gegenprobe: Wo der Kanal nicht ablesbar ist, wird
 * geschwiegen. Lieber ein uebersehener Fall als ein Fehlalarm, denn
 * ein Fehlalarm kostet Vertrauen in ALLE anderen Befunde mit.
 */
func TestOhneKanalangabeWirdNichtsBehauptet(t *testing.T) {
	unklar := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 3200, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "BANK 0"},
		{KapazitaetBytes: 17179869184, TaktMhz: 3200, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "BANK 1"},
	}
	for _, b := range Speicher(unklar) {
		if b.Titel == "Alle Riegel stecken im selben Kanal" {
			t.Error("ohne lesbare Kanalangabe darf kein Einkanalbetrieb behauptet werden")
		}
	}

	// Auch der gemischte Fall zaehlt: Ein Riegel mit lesbarer Angabe,
	// einer ohne. Dann ist die Lage unklar, also wird geschwiegen.
	halb := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 3200, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "P0 CHANNEL A"},
		{KapazitaetBytes: 17179869184, TaktMhz: 3200, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "BANK 1"},
	}
	for _, b := range Speicher(halb) {
		if b.Titel == "Alle Riegel stecken im selben Kanal" {
			t.Error("bei teilweise unlesbaren Angaben darf nichts behauptet werden")
		}
	}
}

// Ein einzelner Riegel ist der andere Befund und darf nicht doppelt
// gemeldet werden. Zwei Meldungen zur selben Sache sind schlimmer als
// eine, das war schon beim halben Speichertakt so.
func TestEinzelnerRiegelMeldetNurEinmal(t *testing.T) {
	einer := []Riegel{
		{KapazitaetBytes: 17179869184, TaktMhz: 3200, Teilenummer: "CMK32GX4M2E3200C16", Kanal: "P0 CHANNEL A"},
	}
	var kanal, einzeln int
	for _, b := range Speicher(einer) {
		switch b.Titel {
		case "Alle Riegel stecken im selben Kanal":
			kanal++
		case "Nur ein Speicherriegel verbaut":
			einzeln++
		}
	}
	if kanal != 0 {
		t.Error("bei einem einzelnen Riegel ist der Kanal-Befund unsinnig")
	}
	if einzeln != 1 {
		t.Errorf("der Einzelriegel-Befund muss genau einmal kommen, kam %dx", einzeln)
	}
}
