// Package pruefung vergleicht Soll und Ist der verbauten Hardware und
// benennt ungenutztes Potenzial.
//
// WARUM DAS UND NICHT EIN BENCHMARK
// =================================
// Ein Benchmark misst, wie schnell ein Rechner ist, und braucht dafür
// Volllast, eine Vergleichsdatenbank mit tausenden Ergebnissen und für
// Grafikkarten eine Render-Schleife, die in Go auf Windows nicht
// vernünftig zu bauen ist. Außerdem widerspräche er der ausdrücklichen
// Zusage dieses Programms: keine Belastungstests, keine Volllast, keine
// Administratorrechte.
//
// Diese Prüfung tut etwas anderes und für die meisten Menschen
// Nützlicheres: Sie beantwortet nicht "welche Punktzahl habe ich",
// sondern "hole ich aus dem, was ich schon besitze, überhaupt alles
// heraus". Sie misst nichts, sie vergleicht nur, was Windows ohnehin
// meldet. Kein zusätzlicher Zugriff, keine Last, keine neuen Rechte.
//
// Der Anlass ist ein echter Fund: Auf dem Entwicklungsrechner selbst
// liefen am 27.08.2026 vier Riegel mit 2133 MT/s, obwohl alle vier für
// 3200 gebaut sind. Rund ein Drittel Speicherbandbreite lag brach, ohne
// dass es jemandem aufgefallen wäre.
//
// WAS HIER NICHT PASSIERT: Es werden keine zusätzlichen Felder erhoben,
// die eine Person oder einen Rechner wiedererkennbar machen. Die
// Teilenummern der Speicherriegel werden gelesen, aber NICHT übertragen,
// sie dienen ausschließlich der lokalen Anzeige. Siehe die Liste in
// scan/types.go.
package pruefung

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Schwere sagt, wie dringend ein Befund ist.
type Schwere string

const (
	// Hinweis: kostet spürbar Leistung, lässt sich meist ohne Zukauf beheben.
	Hinweis Schwere = "hinweis"
	// Anmerkung: einordnend, kein Fehler.
	Anmerkung Schwere = "anmerkung"
)

// Befund ist ein einzelnes Ergebnis der Prüfung.
type Befund struct {
	Schwere Schwere
	Titel   string
	// Was festgestellt wurde, in einem Satz.
	Feststellung string
	// Was man dagegen tun kann. Bewusst konkret, nicht "optimieren Sie".
	Empfehlung string
}

// Riegel ist ein einzelner Speicherbaustein, so wie Windows ihn meldet.
type Riegel struct {
	KapazitaetBytes uint64
	// Tatsächlich eingestellter Takt in MT/s.
	TaktMhz uint32
	// Teilenummer, z.B. "CMK32GX4M2E3200C16". Aus ihr lässt sich die
	// Soll-Geschwindigkeit lesen, siehe sollTaktAusTeilenummer.
	Teilenummer string
	// Kanal-Beschriftung des Mainboards, z.B. "P0 CHANNEL A".
	Kanal string
}

// Laufwerk ist ein Datenträger, so wie Windows ihn meldet.
type Laufwerk struct {
	Name string
	// MSFT_PhysicalDisk.MediaType: 3 = HDD, 4 = SSD.
	MedienArt uint16
	// MSFT_PhysicalDisk.BusType: 11 = SATA, 17 = NVMe, 7 = USB.
	BusArt uint16
	Bytes  uint64
}

const (
	medienHDD = 3
	medienSSD = 4

	busSATA = 11
	busUSB  = 7
	busNVMe = 17
)

/*
 * Erlaubte Soll-Takte. Bewusst eine feste Liste statt einer Spanne:
 *
 * Teilenummern sind voller Zahlen, die nichts mit dem Takt zu tun haben.
 * In "CMK32GX4M2E3200C16" stecken 32 (Kapazität), 4 (DDR4), 3200 (Takt)
 * und 16 (Latenz). Eine Suche nach "irgendeine vierstellige Zahl" würde
 * regelmäßig danebengreifen. Nur echte, im Handel übliche Taktstufen
 * gelten deshalb als Treffer.
 *
 * Lieber keinen Befund als einen falschen: Wer hier eine Stufe vermisst,
 * ergänzt sie, statt die Prüfung aufzuweichen.
 */
var bekannteTakte = map[int]bool{
	// DDR4
	2133: true, 2400: true, 2666: true, 2933: true, 3000: true,
	3200: true, 3466: true, 3600: true, 3733: true, 4000: true,
	4266: true, 4400: true, 4800: true,
	// DDR5 (4800 und 5200 kommen in beiden vor)
	5200: true, 5600: true, 6000: true, 6400: true, 6600: true,
	7200: true, 7600: true, 8000: true, 8400: true,
}

var zahlenMuster = regexp.MustCompile(`\d{4}`)

// sollTaktAusTeilenummer liest die Soll-Geschwindigkeit aus einer
// Teilenummer. Gibt 0 zurück, wenn sich nichts Eindeutiges finden lässt.
//
// Warum überhaupt dieser Umweg: Windows meldet die Soll-Geschwindigkeit
// NICHT. Auf dem Testrechner gaben sowohl Speed als auch
// ConfiguredClockSpeed 2133 an, also beide den Ist-Wert. Die 3200 stehen
// ausschließlich in der Teilenummer. Nachgemessen am 27.08.2026.
//
// Kommen mehrere bekannte Taktstufen in einer Teilenummer vor, gilt keine
// als sicher und die Funktion gibt 0 zurück. Ein Befund, der auf einer
// geratenen Zahl beruht, wäre schlimmer als gar keiner.
func sollTaktAusTeilenummer(teilenummer string) int {
	treffer := map[int]bool{}
	for _, fund := range zahlenMuster.FindAllString(teilenummer, -1) {
		wert, err := strconv.Atoi(fund)
		if err == nil && bekannteTakte[wert] {
			treffer[wert] = true
		}
	}
	if len(treffer) != 1 {
		return 0
	}
	for wert := range treffer {
		return wert
	}
	return 0
}

// Speicher prüft die Riegel auf ungenutztes Potenzial.
func Speicher(riegel []Riegel) []Befund {
	if len(riegel) == 0 {
		return nil
	}
	var befunde []Befund

	// ---- 1. Läuft der Speicher unter seiner Sollgeschwindigkeit? ----
	//
	// Das ist der wichtigste Befund des ganzen Programms. XMP (Intel) und
	// EXPO (AMD) sind ab Werk im BIOS AUS, der Speicher läuft dann mit der
	// vom Prozessor vorgegebenen Grundgeschwindigkeit. Sehr viele Rechner
	// laufen jahrelang so, ohne dass es jemandem auffällt.
	var istTakt uint32
	for _, r := range riegel {
		if r.TaktMhz > 0 && (istTakt == 0 || r.TaktMhz < istTakt) {
			istTakt = r.TaktMhz
		}
	}
	sollTakt := 0
	for _, r := range riegel {
		// Der niedrigste Soll-Wert zählt: Bei gemischten Riegeln läuft
		// alles mit dem langsamsten, mehr ist also gar nicht erreichbar.
		if s := sollTaktAusTeilenummer(r.Teilenummer); s > 0 {
			if sollTakt == 0 || s < sollTakt {
				sollTakt = s
			}
		}
	}
	// 5 Prozent Abstand, damit kleine Abweichungen (2133 gegen 2132)
	// keinen Befund auslösen.
	if sollTakt > 0 && istTakt > 0 && float64(istTakt) < float64(sollTakt)*0.95 {
		/*
		 * Die Zahl muss als GEWINN formuliert sein, nicht als Verlust.
		 *
		 * Bei 2133 statt 3200 ist sollTakt/istTakt = 1,5. Daraus "50
		 * Prozent Bandbreite bleiben ungenutzt" zu machen waere falsch:
		 * Das klaenge, als fehle die Haelfte, tatsaechlich fehlt ein
		 * Drittel des Moeglichen. Richtig ist "50 Prozent mehr waeren
		 * drin", denn das ist der Zuwachs bezogen auf den Ist-Zustand.
		 *
		 * Genau die Sorte Zahlendreher, gegen die dieses Projekt
		 * antritt. Aufgefallen beim ersten echten Durchlauf.
		 */
		gewinn := int(float64(sollTakt)/float64(istTakt)*100) - 100
		befunde = append(befunde, Befund{
			Schwere: Hinweis,
			Titel:   "Arbeitsspeicher läuft unter seiner Sollgeschwindigkeit",
			Feststellung: fmt.Sprintf(
				"Deine Riegel sind für %d MT/s gebaut, laufen aber mit %d MT/s. "+
					"Mit der Sollgeschwindigkeit wären rund %d Prozent mehr Bandbreite drin.",
				sollTakt, istTakt, gewinn),
			Empfehlung: "Im BIOS das Speicherprofil einschalten: bei Intel heißt es XMP, " +
				"bei AMD EXPO oder D.O.C.P. Kostet nichts und ist in zwei Minuten erledigt. " +
				"Läuft der Rechner danach unruhig, das Profil wieder ausschalten.",
		})
	}

	// ---- 2. Nur ein Riegel, also Einkanalbetrieb ----
	if len(riegel) == 1 {
		befunde = append(befunde, Befund{
			Schwere: Hinweis,
			Titel:   "Nur ein Speicherriegel verbaut",
			Feststellung: "Mit einem einzelnen Riegel läuft der Speicher im Einkanalbetrieb. " +
				"Ein zweiter gleicher Riegel bringt spürbar mehr, vor allem bei Grafikeinheiten im Prozessor.",
			Empfehlung: "Einen zweiten, baugleichen Riegel ergänzen. Zwei mal 8 GB sind " +
				"schneller als einmal 16 GB.",
		})
	}

	// ---- 3. Gemischte Riegel ----
	//
	// Kein Fehler, aber der häufigste Grund, warum das Speicherprofil
	// nicht stabil läuft. Gehört deshalb dazu, wenn Befund 1 auftritt.
	verschiedene := map[string]bool{}
	for _, r := range riegel {
		if t := strings.TrimSpace(r.Teilenummer); t != "" {
			verschiedene[t] = true
		}
	}
	if len(verschiedene) > 1 {
		befunde = append(befunde, Befund{
			Schwere: Anmerkung,
			Titel:   "Verschiedene Speicherriegel gemischt",
			Feststellung: fmt.Sprintf(
				"Es sind %d verschiedene Riegeltypen verbaut. Das läuft, ist aber der "+
					"häufigste Grund dafür, dass sich das Speicherprofil nicht stabil aktivieren lässt.",
				len(verschiedene)),
			Empfehlung: "Wenn das Profil aus Punkt eins nicht stabil läuft, liegt es " +
				"vermutlich hieran. Ein einheitlicher Satz Riegel löst das.",
		})
	}

	return befunde
}

// Laufwerke prüft die Datenträger.
func Laufwerke(laufwerke []Laufwerk) []Befund {
	var befunde []Befund

	for _, l := range laufwerke {
		// USB-Geräte gehören nicht in diese Betrachtung, das sind Sticks
		// und externe Platten.
		if l.BusArt == busUSB {
			continue
		}
		if l.MedienArt == medienHDD {
			befunde = append(befunde, Befund{
				Schwere: Hinweis,
				Titel:   "Magnetfestplatte verbaut",
				Feststellung: fmt.Sprintf(
					"%s ist eine klassische Festplatte mit Scheiben. Beim Starten und "+
						"Laden ist sie um ein Vielfaches langsamer als jede SSD.",
					strings.TrimSpace(l.Name)),
				Empfehlung: "Falls Windows darauf liegt: Der Umzug auf eine SSD ist die " +
					"spürbarste Aufrüstung überhaupt, deutlich mehr als ein schnellerer Prozessor.",
			})
			continue
		}
		if l.MedienArt == medienSSD && l.BusArt == busSATA {
			befunde = append(befunde, Befund{
				Schwere: Anmerkung,
				Titel:   "SSD am SATA-Anschluss",
				Feststellung: fmt.Sprintf(
					"%s ist eine SSD, hängt aber am SATA-Anschluss. Der begrenzt auf "+
						"etwa 550 MB/s, eine NVMe-SSD im M.2-Steckplatz schafft ein Vielfaches.",
					strings.TrimSpace(l.Name)),
				Empfehlung: "Im Alltag fällt der Unterschied weniger auf als der Sprung von " +
					"Festplatte auf SSD. Lohnt sich vor allem beim Umgang mit großen Dateien.",
			})
		}
	}

	return befunde
}

// Alle führt sämtliche Prüfungen zusammen.
func Alle(riegel []Riegel, laufwerke []Laufwerk) []Befund {
	return append(Speicher(riegel), Laufwerke(laufwerke)...)
}
