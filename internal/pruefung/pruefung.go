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
	// Hinweis: kostet spürbar Leistung, lässt sich ohne Zukauf beheben.
	Hinweis Schwere = "hinweis"
	// Anmerkung: einordnend, kein Fehler.
	Anmerkung Schwere = "anmerkung"
	// Zukauf: stimmt zwar, kostet aber Geld. Muss vom Rest getrennt
	// stehen, sonst ist die Überschrift der Ergebnisseite gelogen.
	//
	// Eingeführt am 28.08.2026 nach einer Rückmeldung im PCGH-Forum: Über
	// allen Befunden stand "Das hier kostet dich gerade Leistung, ohne
	// dass du etwas kaufen musst", darunter aber der Rat, eine SATA-SSD
	// gegen eine NVMe zu tauschen. Das ist ein Neukauf, und umstecken geht
	// nicht einmal, die Bauformen sind verschieden. Wer so etwas liest,
	// glaubt dem Rest der Seite auch nicht mehr.
	Zukauf Schwere = "zukauf"
	// Bestaetigung: hier ist nichts zu tun, und genau das ist die
	// Nachricht.
	//
	// Eingefuehrt am 29.08.2026 auf Vorschlag von "Misanthrop68" im
	// PCGH-Forum. Seine Frage: Waere es nicht besser, einem Neuling auch
	// zu sagen, DASS der Speicher mit aktivem Profil laeuft? "So hat er
	// keine Rueckmeldung und koennte verunsichert sein."
	//
	// Der Punkt sitzt. Bisher schwieg das Programm zum Speichertakt,
	// sobald er in Ordnung war. Wer nicht weiss, wie das Programm
	// arbeitet, kann Schweigen nicht von Uebersehen unterscheiden. Es gab
	// zwar den Fall "gar nichts gefunden", der etwas Positives sagte, aber
	// der greift nur, wenn ueberhaupt kein Befund vorliegt. Wer Spiele auf
	// einer Festplatte hat, sah zum Speicher weiterhin nichts.
	//
	// Bestaetigungen stehen in einer EIGENEN Karte, nicht zwischen den
	// Funden. Sonst waere die Ueberschrift "Das hier kostet dich gerade
	// Leistung" wieder gelogen, derselbe Fehler wie beim Zukauf-Fall.
	Bestaetigung Schwere = "bestaetigung"
)

// Befund ist ein einzelnes Ergebnis der Prüfung.
type Befund struct {
	Schwere Schwere
	Titel   string
	// Was festgestellt wurde, in einem Satz.
	Feststellung string
	// Was man dagegen tun kann. Bewusst konkret, nicht "optimieren Sie".
	Empfehlung string
	// Abweichende Empfehlung fuer mobile Geraete, leer wenn dieselbe
	// gilt. Angelegt am 31.08.2026 nach einem Hinweis von "PCGH_Jacky"
	// im PCGH-Forum.
	//
	// WARUM AM BEFUND UND NICHT IN EINER NACHBEARBEITUNG: Der Befund
	// weiss selbst am besten, ob sein Rat bei einem Notebook noch
	// stimmt. Eine Funktion, die spaeter nach Titeln sucht und Texte
	// austauscht, waere beim ersten umformulierten Titel still kaputt.
	//
	// Ausgetauscht wird in FuerBauform(), und ein Test haelt fest, dass
	// jeder Rat zum Nachruesten eine mobile Fassung mitbringt.
	EmpfehlungMobil string
	// Warum das so ist, wie man gegenprueft, welche Ausnahmen es gibt.
	//
	// Steht in der Anzeige in einem ZUGEKLAPPTEN Block. Eingefuehrt am
	// 29.08.2026 nach einer Ruecksmeldung von "Hellhammer" im
	// PCGH-Forum: zu viel begruendet, zu ausformuliert.
	//
	// Er hatte recht. Der laengste Befund hatte 590 Zeichen, fuenf
	// Saetze fuer einen Fund. Falsch war aber nicht die Ehrlichkeit,
	// sondern dass alles auf einer Ebene stand. Jeder dieser Saetze
	// hatte mal einen Fehler verhindert, keiner davon gehoert
	// geloescht. Sie gehoeren nur nicht alle nach vorn.
	//
	// FAUSTREGEL BEIM SCHREIBEN: Feststellung und Empfehlung muessen
	// allein ausreichen, um richtig zu handeln. Alles, was nur der
	// braucht, der es genau wissen will, kommt hierher.
	Hintergrund string
}

// Riegel ist ein einzelner Speicherbaustein, so wie Windows ihn meldet.
type Riegel struct {
	KapazitaetBytes uint64
	// Win32_PhysicalMemory.Speed in MT/s.
	TaktMhz uint32
	// Win32_PhysicalMemory.ConfiguredClockSpeed in MT/s, 0 wenn das Board
	// das Feld nicht füllt. Dient allein dazu, dem ersten Wert zu
	// widersprechen, siehe Speicher().
	TaktMhzZweiter uint32
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

/*
 * Der Weg zum Gegenprüfen, wörtlich gleich in beiden Takt-Befunden.
 *
 * Der Windows-Task-Manager ist hier ausdrücklich NICHT der richtige Ort,
 * auch wenn er die Geschwindigkeit anzeigt: Er liest dieselbe
 * SMBIOS-Tabelle wie dieses Programm und zeigt deshalb denselben Wert,
 * richtig oder falsch. Ihn zu empfehlen wäre eine Scheinbestätigung.
 *
 * CPU-Z liest den Speichercontroller direkt aus und ist damit die
 * unabhängige zweite Meinung. Es zeigt die halbe Rate (DDR heißt Double
 * Data Rate), deshalb steht die Verdoppelung ausdrücklich dabei.
 */
// Die Gegenpruef-Erklaerung. Stand bis zum 29.08.2026 mitten in den
// Empfehlungen und machte sie allein schon 486 Zeichen lang. Jetzt ist
// sie Hintergrund: Wer nur wissen will was zu tun ist, liest den kurzen
// Satz oben, wer misstrauisch ist, klappt auf.
const gegenpruefung = "Manche Mainboards tragen hier nur den Grundtakt ein, obwohl der " +
	"Speicher längst schneller läuft. Der Windows-Task-Manager hilft dabei nicht, der liest " +
	"dieselbe Tabelle wie ich. Verlässlich ist das kostenlose CPU-Z: Reiter Memory, den Wert " +
	"bei DRAM Frequency mal zwei nehmen, das ergibt die MT/s. Kommt dort etwas anderes heraus " +
	"als hier, liegt der Fehler bei mir, und ich würde mich sehr über eine Nachricht über " +
	"benchmeister.de/kontakt freuen."

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
/*
 * Liest aus der Kanal-Beschriftung den Prozessor und den Kanal heraus.
 *
 * Windows liefert in BankLabel Zeichenketten wie "P0 CHANNEL A" oder
 * "BANK 0". Die erste Form nennt den Kanal, die zweite nicht. Welche
 * ein Board schreibt, entscheidet sein Hersteller.
 *
 * Deshalb gibt es drei Rueckgaben und nicht eine: das Praefix vor dem
 * Wort CHANNEL, den Kanal selbst, und die Angabe ob ueberhaupt etwas
 * erkennbar war. Wo nichts zu erkennen ist, wird nichts behauptet.
 * Genau diese Unterscheidung hat beim Speichertakt gefehlt und zu
 * einem Fehlalarm bei einem Moderator gefuehrt.
 *
 * WARUM DAS PRAEFIX EIGENS ZURUECKKOMMT, ergaenzt am 30.08.2026:
 * Es gibt Rechner mit zwei Prozessoren, und jeder hat eigene Kanaele.
 * Die Beschriftung lautet dann P0 CHANNEL A und P1 CHANNEL A. Das sind
 * ZWEI verschiedene Kanaele an zwei Prozessoren. Die erste Fassung las
 * nur den Buchstaben und haette bei so einer Maschine Einkanalbetrieb
 * gemeldet, obwohl alles richtig steckt.
 *
 * Gesucht wird bewusst nur nach dem Wort CHANNEL mit einem Buchstaben
 * dahinter. "BANK 0" und "BANK 1" sehen zwar nach zwei Kanaelen aus,
 * sind aber auf vielen Boards zwei Slots DESSELBEN Kanals. Wer das
 * verwechselt, meldet Einkanalbetrieb, wo keiner ist.
 */
func kanalAus(beschriftung string) (praefix, kanal string, erkannt bool) {
	gross := strings.ToUpper(strings.TrimSpace(beschriftung))
	i := strings.Index(gross, "CHANNEL")
	if i < 0 {
		return "", "", false
	}
	rest := strings.TrimSpace(gross[i+len("CHANNEL"):])
	if rest == "" {
		return "", "", false
	}
	// Der Kanal ist genau ein Buchstabe: A, B, C, D.
	buchstabe := rest[0]
	if buchstabe < 'A' || buchstabe > 'D' {
		return "", "", false
	}
	// Alles vor CHANNEL benennt den Prozessor, falls das Board ihn
	// nennt. Bei Einzelsockel-Boards steht dort oft "P0" oder nichts.
	return strings.TrimSpace(gross[:i]), string(buchstabe), true
}

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
	//
	// WELCHES FELD DEN IST-TAKT ENTHAELT, korrigiert am 29.08.2026.
	//
	// Am 28.08. meldete ein Moderator im PCGH-Forum einen Fehlalarm: Sein
	// Speicher lief mit 5600 MT/s, dieses Programm behauptete 4800. Die
	// erste Reparatur war, bei widerspruechlichen Feldern gar nichts mehr
	// zu sagen. Das war fachlich falsch, und ein Nutzer (NullPointerEx)
	// hat im selben Faden erklaert warum:
	//
	//   Speed                = hoechster nativer JEDEC-Takt, OHNE Profil
	//   ConfiguredClockSpeed = der reale Takt, wie das BIOS ihn eingestellt
	//                          hat (XMP, EXPO, manuell)
	//
	// Ein Unterschied zwischen beiden ist also kein Widerspruch, sondern
	// der NORMALFALL bei aktivem Profil: Speed 4800 und
	// ConfiguredClockSpeed 6400 heisst schlicht "XMP laeuft". Die alte
	// Fassung antwortete ausgerechnet dann "ich weiss es nicht", wenn die
	// Daten eindeutig waren, und war damit fuer Einsteiger wertlos.
	//
	// Richtig ist: ConfiguredClockSpeed nehmen, wenn vorhanden, sonst
	// Speed als Rueckfall. Der Rest bleibt wie gehabt der Abgleich gegen
	// die Teilenummer.
	//
	// Was bleibt: Verlassen kann man sich auf keines der Felder blind.
	// Was in SMBIOS Typ 17 landet, entscheidet das BIOS des Boards, und
	// manche ueberschreiben Speed mit dem konfigurierten Wert. Deshalb
	// nennt der Befund weiterhin den Weg zum Gegenpruefen.
	// ConfiguredClockSpeed hat Vorrang, Speed ist nur der Rückfall.
	// Begründung siehe der grosse Kommentarblock oben.
	//
	// ACHTUNG: Dieselbe Regel steht ein zweites Mal in
	// scan.RiegelInfo.IstTaktMhz, weil dieses Paket bewusst nichts ueber
	// das Auslese-Paket weiss. Genau so eine Doppelung hat am 29.08.2026
	// dazu gefuehrt, dass der Befund 5600 sagte und die Tabelle darueber
	// 4800. Wer hier etwas aendert, muss es dort auch aendern. Der Test
	// TestScanUndPruefungRechnenGleich in internal/ui haelt beide
	// zusammen und faellt, wenn sie auseinanderlaufen.
	var istTakt uint32
	// Der Grundtakt OHNE Profil, also Win32_PhysicalMemory.Speed allein.
	// Wird nur fuer die Bestaetigung weiter unten gebraucht: Liegt der
	// Ist-Takt darueber, hat das Board nachweislich mehr eingestellt als
	// die Werksvorgabe, und nur dann darf von einem Profil die Rede sein.
	// Fuer die Fehlersuche taugt das Feld nicht, siehe Kommentar oben.
	var jedecTakt uint32
	for _, r := range riegel {
		takt := r.TaktMhzZweiter
		if takt == 0 {
			takt = r.TaktMhz
		}
		// Bei gemischten Riegeln zählt der langsamste, mit dem läuft alles.
		if takt > 0 && (istTakt == 0 || takt < istTakt) {
			istTakt = takt
		}
		if r.TaktMhz > 0 && (jedecTakt == 0 || r.TaktMhz < jedecTakt) {
			jedecTakt = r.TaktMhz
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
	/*
	 * SONDERFALL: Windows meldet fast genau die HAELFTE.
	 *
	 * Ergaenzt am 29.08.2026 nach Durchsicht fremder Foren-Faeden. In
	 * mehreren dokumentierten Faellen zeigte Windows den Speichertakt
	 * halbiert an (1600 statt 3200), waehrend CPU-Z den richtigen Wert
	 * meldete. Der Grund ist die alte Verwechslung von MHz und MT/s: DDR
	 * uebertraegt zweimal je Takt, 3200 MT/s sind 1600 MHz. Microsoft hat
	 * die Beschriftung im Task-Manager deshalb 2024 von "MHz" auf "MT/s"
	 * geaendert, das Grundproblem gibt es aber weiter.
	 *
	 * Ohne diesen Zweig haette das Programm daraus einen sehr lauten
	 * Befund gemacht: "Deine Riegel sind fuer 3200 gebaut, laufen aber mit
	 * 1600, rund 100 Prozent mehr waeren drin." Bei jemandem, dessen
	 * Speicher voellig richtig laeuft. Also genau die Sorte Falschalarm,
	 * die dieses Programm schon zweimal blamiert hat.
	 *
	 * BEWUSST NICHT UNTERDRUECKT, sondern anders formuliert: Ein Riegel
	 * KANN wirklich mit der halben Rate laufen, DDR4-1600 gibt es. Beide
	 * Faelle sehen von hier aus identisch aus. Deshalb nennt der Befund
	 * beide Moeglichkeiten und die eine Messung, die sie unterscheidet,
	 * statt sich fuer eine zu entscheiden.
	 */
	// Kein frueher Ausstieg: Die Pruefungen auf Einkanalbetrieb und
	// gemischte Riegel weiter unten muessen trotzdem laufen.
	halbeRate := false
	if sollTakt > 0 && istTakt > 0 {
		verhaeltnis := float64(istTakt) / float64(sollTakt)
		if verhaeltnis >= 0.45 && verhaeltnis <= 0.55 {
			halbeRate = true
			befunde = append(befunde, Befund{
				Schwere: Hinweis,
				Titel:   "Windows meldet genau den halben Speichertakt",
				Feststellung: fmt.Sprintf(
					"Deine Riegel sind für %d MT/s gebaut, Windows meldet %d. Genau die Hälfte.",
					sollTakt, istTakt),
				Empfehlung: "Mit CPU-Z gegenprüfen: Reiter Memory, den Wert bei DRAM " +
					"Frequency mal zwei. Stimmt er, ist alles gut. Wenn nicht, im BIOS " +
					"das Speicherprofil einschalten (XMP, EXPO oder D.O.C.P.).",
				Hintergrund: "Genau die Hälfte hat erfahrungsgemäß zwei mögliche Gründe. " +
					"Entweder zählt Windows die reine Taktrate statt der Übertragungen pro " +
					"Sekunde. Speicher überträgt zweimal je Takt, aus 3200 werden dann 1600, " +
					"und dann ist alles in Ordnung. Oder dein Speicher läuft wirklich halb so " +
					"schnell. Beides sieht von hier aus gleich aus, deshalb die Gegenprobe. " +
					gegenpruefung,
			})
		}
	}

	// 5 Prozent Abstand, damit kleine Abweichungen (2133 gegen 2132)
	// keinen Befund auslösen. Nicht bei halber Rate: Dazu steht oben
	// bereits ein Befund, und zwei Meldungen zur selben Zahl, die
	// verschiedene Dinge behaupten, sind schlimmer als eine.
	if !halbeRate && sollTakt > 0 && istTakt > 0 && float64(istTakt) < float64(sollTakt)*0.95 {
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
			Titel:   "Arbeitsspeicher läuft womöglich unter seiner Sollgeschwindigkeit",
			Feststellung: fmt.Sprintf(
				"Deine Riegel sind für %d MT/s gebaut, Windows meldet %d. Rund %d Prozent "+
					"mehr Bandbreite wären drin.",
				sollTakt, istTakt, gewinn),
			Empfehlung: "Im BIOS das Speicherprofil einschalten: XMP bei Intel, " +
				"EXPO oder D.O.C.P. bei AMD. Dauert zwei Minuten.",
			Hintergrund: "Vorher gegenprüfen lohnt sich. " + gegenpruefung +
				" Läuft der Rechner nach dem Einschalten unruhig, das Profil wieder " +
				"ausschalten. Dann vertragen sich die Riegel nicht mit der Einstellung, " +
				"und daran ist nichts kaputt.",
		})
	}

	/*
	 * Die Gegenrichtung: Der Takt stimmt, und das wird auch gesagt.
	 *
	 * Bedingung ist bewusst "nicht nennenswert unter Soll" und nicht
	 * "ueber der JEDEC-Basis". Beides waere richtig, aber nur das erste
	 * beantwortet die Frage, die sich der Nutzer stellt: Hole ich hier
	 * alles heraus?
	 *
	 * WAS HIER NICHT BEHAUPTET WIRD: dass XMP oder EXPO eingeschaltet
	 * ist. Sicher wissen wir das nicht, das BIOS koennte auch von Hand
	 * eingestellt sein, und bei Riegeln ohne Profil ist der Grundtakt
	 * bereits der Sollwert. Deshalb steht da die Messung ("laeuft mit X,
	 * gebaut fuer Y") und nur dort ein Wort zum Profil, wo der Ist-Takt
	 * tatsaechlich ueber dem JEDEC-Grundtakt liegt. Das ist der einzige
	 * Fall, in dem sich das aus den Daten ableiten laesst.
	 */
	if !halbeRate && sollTakt > 0 && istTakt > 0 && float64(istTakt) >= float64(sollTakt)*0.95 {
		feststellung := fmt.Sprintf(
			"Deine Riegel sind für %d MT/s gebaut und laufen laut Windows mit %d MT/s.",
			sollTakt, istTakt)
		empfehlung := "Hier ist nichts zu tun."
		hintergrund := "Mehr geben diese Riegel nicht her, ohne sie über ihre Vorgabe " +
			"hinaus zu betreiben."
		if jedecTakt > 0 && istTakt > jedecTakt {
			hintergrund = fmt.Sprintf(
				"Ohne Speicherprofil würden sie mit %d MT/s laufen. Dein Mainboard hat "+
					"also mehr eingestellt: XMP, EXPO oder D.O.C.P. ist an, oder jemand "+
					"hat von Hand nachgeholfen.",
				jedecTakt)
		}
		befunde = append(befunde, Befund{
			Schwere:      Bestaetigung,
			Titel:        "Der Arbeitsspeicher läuft auf Sollgeschwindigkeit",
			Feststellung: feststellung,
			Empfehlung:   empfehlung,
			Hintergrund:  hintergrund,
		})
	}

	/*
	 * ---- 2a. Mehrere Riegel, aber alle im selben Kanal ----
	 *
	 * Der Fall, den das Programm bis zum 30.08.2026 uebersehen hat, weil
	 * es nur die Riegel GEZAEHLT hat. Zwei Riegel nebeneinander im
	 * selben Kanal laufen im Einkanalbetrieb, genau wie ein einzelner.
	 * Von aussen sieht der Rechner voll bestueckt aus.
	 *
	 * Anders als beim fehlenden Speicherprofil hilft hier kein
	 * BIOS-Schalter: Die Riegel muessen umgesteckt werden. Das ist ein
	 * Handgriff und kostet nichts, deshalb steht der Befund unter den
	 * kostenlosen.
	 */
	kanaele := map[string]bool{}
	prozessoren := map[string]bool{}
	kanalUnklar := false
	for _, r := range riegel {
		if p, k, ok := kanalAus(r.Kanal); ok {
			kanaele[k] = true
			prozessoren[p] = true
		} else {
			kanalUnklar = true
		}
	}
	/*
	 * Bei mehr als einem Prozessor wird geschwiegen, ergaenzt am
	 * 30.08.2026 nach Dominics Hinweis auf Dual-Sockel-Systeme.
	 *
	 * Dort ist "laeuft Dual-Channel" eine Frage PRO Prozessor, und ob
	 * die Riegel sinnvoll auf beide verteilt sind, ist noch eine
	 * dritte. Aus der Beschriftung allein laesst sich das nicht
	 * beantworten. Solche Maschinen sind selten und ihre Besitzer
	 * wissen in aller Regel, was sie tun.
	 */
	if len(riegel) > 1 && !kanalUnklar && len(prozessoren) == 1 && len(kanaele) == 1 {
		befunde = append(befunde, Befund{
			Schwere: Hinweis,
			Titel:   "Alle Riegel stecken im selben Kanal",
			Feststellung: fmt.Sprintf(
				"%d Riegel, aber alle in Kanal %s. Damit läuft der Speicher im Einkanalbetrieb.",
				len(riegel), einzigerKanal(kanaele)),
			Empfehlung: "Einen Riegel in einen Steckplatz des anderen Kanals umstecken. " +
				"Welcher das ist, steht im Handbuch des Mainboards.",
			EmpfehlungMobil: "Umstecken hilft, wenn beide Steckplätze erreichbar sind. " +
				"Bei Notebooks liegt oft nur einer unter einer Klappe.",
			Hintergrund: "Zwei Kanäle arbeiten nebeneinander und verdoppeln die " +
				"Bandbreite. Stecken alle Riegel im selben, bleibt die Hälfte " +
				"ungenutzt, obwohl der Rechner voll bestückt aussieht. Welche " +
				"Steckplätze die richtigen sind, steht im Handbuch deines " +
				"Mainboards und ist meist auch auf dem Board selbst aufgedruckt. " +
				"Feste Regeln gibt es nicht: Auf manchen Boards sind es der " +
				"zweite und vierte Steckplatz, auf anderen der erste und dritte. " +
				"Anders als beim Speicherprofil hilft hier kein Schalter im BIOS, " +
				"die Riegel müssen wirklich umgesteckt werden.",
		})
	}

	// ---- 2. Nur ein Riegel, also Einkanalbetrieb ----
	if len(riegel) == 1 {
		befunde = append(befunde, Befund{
			Schwere:      Hinweis,
			Titel:        "Nur ein Speicherriegel verbaut",
			Feststellung: "Mit einem einzelnen Riegel läuft der Speicher im Einkanalbetrieb.",
			Empfehlung:   "Einen zweiten, baugleichen Riegel ergänzen.",
			EmpfehlungMobil: "Vor dem Kauf nachsehen, ob ein Steckplatz frei ist: " +
				"Bei vielen Notebooks ist der Speicher fest verlötet.",
			Hintergrund: "Zwei Riegel arbeiten nebeneinander und verdoppeln die Bandbreite. " +
				"Zwei mal 8 GB sind deshalb schneller als einmal 16 GB. Am deutlichsten " +
				"merkt man das bei Prozessoren mit eingebauter Grafik, weil die sich den " +
				"Arbeitsspeicher mit dem Rest teilen müssen.",
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
				"Es sind %d verschiedene Riegeltypen verbaut. Das läuft, ist aber kein Idealzustand.",
				len(verschiedene)),
			Empfehlung: "Nichts tun, solange alles stabil läuft.",
			Hintergrund: "Gemischte Riegel sind der häufigste Grund dafür, dass sich das " +
				"Speicherprofil nicht stabil aktivieren lässt. Wenn du oben den Hinweis zum " +
				"Profil bekommen hast und es nach dem Einschalten hakt, liegt es vermutlich " +
				"hieran. Ein einheitlicher Satz Riegel löst das, ist aber ein Neukauf und " +
				"lohnt sich nur, wenn dich wirklich etwas stört.",
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
				Schwere: Zukauf,
				Titel:   "Magnetfestplatte verbaut",
				Feststellung: fmt.Sprintf(
					"%s ist eine Festplatte mit Scheiben, um ein Vielfaches langsamer "+
						"als jede SSD.",
					strings.TrimSpace(l.Name)),
				// Der zweite Satz ist am 28.08.2026 dazugekommen. Vorher
				// stand hier nur der Aufrüst-Rat, und ein Moderator im
				// PCGH-Forum hat zu Recht angemerkt, dass ein Ratschlag
				// ohne Wissen um den Verwendungszweck keiner ist. Für
				// Sicherungen und Archive ist eine Festplatte genau
				// richtig, und das muss dastehen.
				// Formulierung am 29.08.2026 nach einem Vorschlag von
				// "Misanthrop68" im PCGH-Forum geschaerft: Erst sagen, wofuer
				// die Platte NICHT taugt, dann wofuer sie weiterhin gut ist.
				// Vorher stand die Entwarnung hinten und wurde ueberlesen.
				Empfehlung: "Liegen Windows oder Spiele darauf, lohnt der Umzug auf eine " +
					"SSD. Ist es nur der Datenkeller, ist alles richtig.",
				Hintergrund: "Beim Starten und beim Laden macht sich das am deutlichsten " +
					"bemerkbar. Für Windows und Programme ist eine Festplatte heute keine gute " +
					"Grundlage mehr, dort kostet sie bei jedem Start und jedem Ladebildschirm " +
					"Zeit. Für Sicherungen, Archive, Fotos und Musik ist sie dagegen völlig in " +
					"Ordnung und pro Terabyte konkurrenzlos günstig. Es kommt also darauf an, " +
					"was darauf liegt, und das sehe ich von hier aus nicht.",
			})
			continue
		}
		if l.MedienArt == medienSSD && l.BusArt == busSATA {
			befunde = append(befunde, Befund{
				Schwere: Zukauf,
				Titel:   "SSD hängt am SATA-Anschluss",
				Feststellung: fmt.Sprintf(
					"%s hängt am SATA-Anschluss. Der begrenzt auf etwa 550 MB/s, eine "+
						"NVMe im M.2-Steckplatz schafft ein Vielfaches.",
					strings.TrimSpace(l.Name)),
				// Der erste Satz stellt klar, dass das kein Handgriff ist.
				// Vorher las sich der Befund, als könne man die Platte
				// einfach umstecken, was bei verschiedenen Bauformen
				// natürlich nicht geht. PCGH-Rückmeldung vom 28.08.2026.
				// Ebenfalls nach dem Vorschlag von "Misanthrop68" (29.08.2026):
				// Der Hinweis soll klar sagen, dass es KEINEN Grund zum Handeln
				// gibt, und den einen Zeitpunkt nennen, an dem es sich lohnt.
				Empfehlung: "Nichts tun. Erst wenn diese SSD ohnehin ersetzt wird, " +
					"lieber gleich eine M.2 nehmen.",
				Hintergrund: "Umstecken geht nicht, SATA und M.2 sind zwei verschiedene " +
					"Bauformen, es liefe also auf einen Neukauf hinaus. Einen Grund dafür " +
					"gibt es gerade nicht: Im Alltag merkt man den Unterschied kaum, das ist " +
					"etwas ganz anderes als der Sprung von Festplatte auf SSD. Eine M.2 ist " +
					"außerdem schneller und braucht keinen Platz im Gehäuse.",
			})
		}
	}

	return befunde
}

// Alle führt sämtliche Prüfungen zusammen.
func Alle(riegel []Riegel, laufwerke []Laufwerk) []Befund {
	return append(Speicher(riegel), Laufwerke(laufwerke)...)
}

/*
FuerBauform tauscht Empfehlungen aus, die bei einem mobilen Geraet
nicht befolgbar waeren.

Angewendet wird das erst nach der Pruefung, damit die Pruefungen selbst
nichts von Gehaeusen wissen muessen und ohne Windows testbar bleiben.

BEI UNBEKANNTER BAUFORM WIRD NICHTS GEAENDERT. Die Gehaeuseart kommt aus
einer Tabelle, die der Hersteller fuellt, und liegt nicht immer vor.
Dann bleibt der bisherige Text stehen, der fuer die Mehrheit stimmt.
*/
func FuerBauform(befunde []Befund, mobil bool) []Befund {
	if !mobil {
		return befunde
	}
	angepasst := make([]Befund, len(befunde))
	for i, b := range befunde {
		if b.EmpfehlungMobil != "" {
			b.Empfehlung = b.EmpfehlungMobil
		}
		angepasst[i] = b
	}
	return angepasst
}

// Der einzige Kanal aus der Menge. Wird nur aufgerufen, wenn genau
// einer drin ist.
func einzigerKanal(m map[string]bool) string {
	for k := range m {
		return k
	}
	return "?"
}
