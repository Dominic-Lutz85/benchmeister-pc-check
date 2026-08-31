package pruefung

import (
	"fmt"
	"strings"
)

/*
Drei Befunde zum Zustand des Systems statt zur Hardware, angelegt am
31.08.2026 nach Vorschlaegen von "cryon1c" im PCGH-Forum.

Sein Punkt war treffend: Das sind "generelle Probleme bei
Konfigurationen, die aber keine Fehlermeldungen verursachen. Kosten
aber Leistung." Genau dafuer ist dieses Programm da.
*/

// Zustand sind die drei ausgelesenen Werte. Bewusst als eigener Typ
// hier statt eines Imports aus scan: Das Pruefpaket kennt keine
// WMI-Abfragen und laesst sich damit ohne Windows testen.
type Zustand struct {
	SystemplatteFreiGb   int
	SystemplatteGesamtGb int
	SystemplatteErkannt  bool

	Virenschutz        []string
	VirenschutzErkannt bool

	NetzwerkMbit    int
	NetzwerkErkannt bool
}

/*
SCHWELLEN, und warum ausgerechnet diese.

Es gibt zwei verschiedene Gruende, warum wenig Platz schadet, und sie
brauchen zwei verschiedene Massstaebe:

DER ANTEIL zaehlt fuer die SSD selbst. Sie verteilt Schreibzugriffe auf
freie Bloecke, und je weniger davon uebrig sind, desto mehr muss sie
vorher umraeumen. Das macht sie spuerbar langsamer.

DIE ABSOLUTE ZAHL zaehlt fuer Windows. Ein groesseres Update braucht
ungefaehr 20 GB am Stueck, unabhaengig davon, wie gross die Platte ist.

WARUM DIE PROZENTREGEL EINE OBERGRENZE BRAUCHT, gefunden beim Schreiben
der Tests am 31.08.2026: Auf einer 4-TB-Platte sind 300 GB frei knapp
acht Prozent. Die reine Prozentregel haette hier gewarnt, obwohl 300 GB
freie Bloecke fuer jedes Wear-Leveling reichlich sind und fuer jedes
Update ohnehin. Ein Fehlalarm auf einem voellig gesunden Rechner.

Deshalb greift die Prozentregel nur, solange auch absolut wenig frei
ist. Die 20-GB-Regel gilt daneben immer.
*/
const (
	freiProzentGrenze = 10
	// Oberhalb dieser Menge ist der Anteil egal, dann ist genug da.
	freiReichlichGb = 100
	// Darunter wird es unabhaengig von der Plattengroesse eng.
	freiGbGrenze = 20
)

// Systemplatte warnt vor einer vollen Windows-Platte.
func Systemplatte(z Zustand) []Befund {
	if !z.SystemplatteErkannt || z.SystemplatteGesamtGb == 0 {
		return nil
	}

	anteil := z.SystemplatteFreiGb * 100 / z.SystemplatteGesamtGb
	anteilKnapp := anteil < freiProzentGrenze && z.SystemplatteFreiGb < freiReichlichGb
	absolutKnapp := z.SystemplatteFreiGb < freiGbGrenze
	if !anteilKnapp && !absolutKnapp {
		return nil
	}

	return []Befund{{
		Schwere: Hinweis,
		Titel:   "Windows-Laufwerk ist fast voll",
		Feststellung: fmt.Sprintf(
			"Auf dem Laufwerk mit Windows sind noch %d von %d GB frei, also %d Prozent.",
			z.SystemplatteFreiGb, z.SystemplatteGesamtGb, anteil),
		Empfehlung: "Platz schaffen, bis wieder rund ein Zehntel frei ist.",
		Hintergrund: "Eine fast volle SSD wird langsamer, und zwar aus einem bauartbedingten " +
			"Grund: Sie verteilt Schreibzugriffe auf freie Bloecke, und je weniger davon " +
			"uebrig sind, desto mehr muss sie vorher umraeumen. Dazu kommt Windows selbst. " +
			"Es legt auf diesem Laufwerk die Auslagerungsdatei und Zwischenspeicher an und " +
			"braucht fuer ein groesseres Update ungefaehr 20 GB am Stueck. Fehlt der Platz, " +
			"bricht das Update ab oder laeuft sehr lange. Am meisten bringt meist der " +
			"Downloads-Ordner, danach alte Spiele. Die Datentraegerbereinigung von Windows " +
			"raeumt zusaetzlich alte Update-Reste weg.",
	}}
}

/*
Virenschutz warnt, wenn mehr als ein Schutzprogramm registriert ist.

WAS HIER BEWUSST NICHT BEHAUPTET WIRD: dass beide gleichzeitig LAUFEN.
Der Namensraum root\SecurityCenter2 ist von Microsoft nicht
dokumentiert, und das Feld mit dem Zustand ist ein Bitfeld, dessen
Bedeutung nur aus Beobachtung stammt. Wer daraus "aktiv" ableitet,
baut einen Fehlalarm auf einer Vermutung, siehe Regel 1 in CLAUDE.md.

Deshalb steht im Befund "registriert", und der Hintergrund sagt, dass
auch Reste einer alten Installation dazugehoeren koennen. Das ist
ehrlicher und fuehrt zur selben Handlung: nachsehen.
*/
func Virenschutz(z Zustand) []Befund {
	if !z.VirenschutzErkannt || len(z.Virenschutz) < 2 {
		return nil
	}

	return []Befund{{
		Schwere: Hinweis,
		Titel:   "Mehrere Schutzprogramme registriert",
		Feststellung: fmt.Sprintf(
			"Windows kennt %d Schutzprogramme: %s.",
			len(z.Virenschutz), namenKurz(z.Virenschutz)),
		Empfehlung: "Nachsehen, ob wirklich zwei gleichzeitig laufen, und eines abschalten.",
		Hintergrund: "Zwei Schutzprogramme, die gleichzeitig jede Datei pruefen, machen " +
			"dieselbe Arbeit doppelt und geraten sich dabei ins Gehege. Das kostet vor allem " +
			"beim Starten von Programmen und beim Kopieren spuerbar Zeit, manchmal blockieren " +
			"sie sich auch gegenseitig. Wichtig: Diese Liste zeigt, was bei Windows ANGEMELDET " +
			"ist, nicht was gerade laeuft. Ob ein Programm aktiv ist, meldet Windows nur in " +
			"einem Feld, das Microsoft nicht dokumentiert, und darauf verlasse ich mich nicht. " +
			"Es kann also auch der Rest einer alten Installation sein, die nie sauber entfernt " +
			"wurde. Nachsehen laesst sich das unter Einstellungen, Datenschutz und Sicherheit, " +
			"Windows-Sicherheit.",
	}}
}

/*
Netzwerk warnt, wenn eine Kabelverbindung unter Gigabit aushandelt.

Der haeufigste Fall dahinter ist ein Kabel oder ein Anschluss, der nur
100 Mbit schafft, oder ein Stromsparmodus im Router, der nachts
herunterschaltet und morgens nicht zurueckkommt. Beides faellt im
Alltag nicht auf, bis eine grosse Datei kopiert wird.

Nur Kabelverbindungen, siehe scan/systemzustand.go: Bei WLAN waere
eine niedrigere Aushandlung normal.
*/
func Netzwerk(z Zustand) []Befund {
	if !z.NetzwerkErkannt || z.NetzwerkMbit >= 1000 {
		return nil
	}

	return []Befund{{
		Schwere: Hinweis,
		Titel:   "Netzwerk läuft unter Gigabit",
		Feststellung: fmt.Sprintf(
			"Die Kabelverbindung ist mit %d Mbit/s ausgehandelt, üblich sind 1000.",
			z.NetzwerkMbit),
		Empfehlung: "Anderes Netzwerkkabel probieren und den Rechner einmal neu verbinden.",
		Hintergrund: "Ausgehandelt wird zwischen Netzwerkkarte und Gegenstelle, und beide " +
			"einigen sich auf das langsamere Ende. Die haeufigste Ursache ist ein altes oder " +
			"beschaedigtes Kabel: Fuer Gigabit muessen alle acht Adern durchgaengig sein, bei " +
			"vier reicht es nur fuer 100 Mbit, und ein geknicktes Kabel verhaelt sich genauso. " +
			"Die zweite haeufige Ursache sind Stromsparmodi in Routern, die die Geschwindigkeit " +
			"herunterschalten und nicht von selbst zuruecksetzen. Am Internetanschluss merkt " +
			"man davon meist nichts, weil der ohnehin langsamer ist. Auffaellig wird es beim " +
			"Kopieren im Heimnetz und beim Streamen im eigenen Netz.",
	}}
}

/*
namenKurz haelt die Aufzaehlung so kurz, dass Feststellung und
Empfehlung zusammen unter der 220-Zeichen-Grenze bleiben.

Ohne das reisst der Befund bei mehreren Programmen mit langen Namen die
Grenze. "Kaspersky Endpoint Security for Windows" allein sind 38
Zeichen, drei solche Namen sprengen jedes Mass. Aufgefallen beim
Schreiben der Tests, nicht erst beim Nutzer.

Ab dem dritten Namen wird zusammengefasst. Die genaue Zahl steht
ohnehin am Satzanfang, und wer nachsehen will, findet die Liste in den
Windows-Einstellungen vollstaendig.
*/
func namenKurz(namen []string) string {
	const hoechstens = 2
	if len(namen) <= hoechstens {
		return strings.Join(namen, ", ")
	}
	rest := len(namen) - hoechstens
	// Einzahl und Mehrzahl. Im Probelauf stand "und 1 weitere", und das
	// liest sich wie ein Tippfehler. Bei einem Programm, dessen ganzer
	// Anspruch Sorgfalt ist, faellt so etwas negativ auf.
	if rest == 1 {
		return strings.Join(namen[:hoechstens], ", ") + " und ein weiteres"
	}
	return strings.Join(namen[:hoechstens], ", ") +
		fmt.Sprintf(" und %d weitere", rest)
}

// AlleZustand fuehrt die drei Zustandspruefungen zusammen.
func AlleZustand(z Zustand) []Befund {
	var befunde []Befund
	befunde = append(befunde, Systemplatte(z)...)
	befunde = append(befunde, Virenschutz(z)...)
	befunde = append(befunde, Netzwerk(z)...)
	return befunde
}
