package ui

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"

	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/pruefung"
	"github.com/Dominic-Lutz85/benchmeister-pc-check/internal/scan"
)

/*
 * Haelt die Barrierefreiheit der Ergebnisseite fest, angelegt am
 * 07.09.2026.
 *
 * WARUM DAS EINEN TEST WERT IST: Die fuenf Punkte hier waren alle einmal
 * verletzt, und keiner davon faellt beim Ansehen auf. Eine zu blasse
 * Schrift sieht nach Absicht aus, eine fehlende Ueberschrift sieht nach
 * gar nichts aus, und ein Meldebereich ohne Rolle sieht voellig normal
 * aus. Genau solche Sachen kommen beim naechsten Umbau der Vorlage
 * zurueck, wenn sie niemand festhaelt.
 *
 * Die Zielgruppe sind Leute, die solchen Programmen misstrauen. Wer die
 * Seite nicht bedienen kann, kann auch nicht nachpruefen, was sie
 * behauptet, und dann ist das ganze Transparenzversprechen nichts wert.
 */

// beispielBefunde liefert je einen Befund pro Schweregrad. Alle vier
// zusammen, weil die Vorlage sie unterschiedlich behandelt.
func beispielBefunde() (kostenlos, kostenpflichtig, laeuft, alle []pruefung.Befund) {
	kostenlos = []pruefung.Befund{
		{
			Schwere:      pruefung.Hinweis,
			Titel:        "Der Arbeitsspeicher läuft unter seiner Geschwindigkeit",
			Feststellung: "Gebaut für 6000 MT/s, läuft mit 4800 MT/s.",
			Empfehlung:   "Im BIOS das Profil einschalten.",
		},
		{
			Schwere:      pruefung.Anmerkung,
			Titel:        "Verschiedene Speicherriegel gemischt",
			Feststellung: "Zwei verschiedene Riegeltypen verbaut.",
			Empfehlung:   "Nichts tun, solange alles stabil läuft.",
		},
	}
	kostenpflichtig = []pruefung.Befund{
		{
			Schwere:      pruefung.Zukauf,
			Titel:        "Nur ein Speicherriegel verbaut",
			Feststellung: "Ein Riegel in einem Board mit vier Steckplätzen.",
			Empfehlung:   "Einen zweiten, baugleichen Riegel ergänzen.",
		},
	}
	laeuft = []pruefung.Befund{
		{
			Schwere:      pruefung.Bestaetigung,
			Titel:        "Die Systemplatte ist eine NVMe-SSD",
			Feststellung: "Windows liegt auf einer NVMe-SSD.",
			Empfehlung:   "Hier ist nichts zu tun.",
		},
	}
	alle = append(append(append([]pruefung.Befund{}, kostenlos...), kostenpflichtig...), laeuft...)
	return
}

// volleSeite rendert die Vorlage mit allem, was mehrzeilig werden kann.
func volleSeite(t *testing.T) string {
	t.Helper()
	kostenlos, kostenpflichtig, laeuft, alle := beispielBefunde()
	return rendere(t, vorlagenDaten{
		Scan: &scan.ScanResult{
			CPUName: "AMD Ryzen 7 7800X3D", CPUCores: 8, CPUThreads: 16,
			GPUName: "NVIDIA GeForce RTX 5060 Ti", RamTotalGb: 16,
			MainboardName:   "B650 Beispiel",
			ResolutionWidth: 2560, ResolutionHeight: 1440,
		},
		Befunde: alle, Kostenlos: kostenlos,
		Kostenpflichtig: kostenpflichtig, Laeuft: laeuft,
		Laufwerke: []LaufwerkZeile{
			{Beschriftung: "NVMe-SSD, 2 TB (System)"},
			{Beschriftung: "SATA-SSD, 1 TB"},
			{Beschriftung: "Festplatte, 3 TB"},
		},
		Pcie: []PcieZeile{
			{Beschriftung: "Grafikkarte: PCIe 4.0 x8"},
			{Beschriftung: "NVMe-SSD: PCIe 4.0 x4"},
		},
	})
}

/*
 * WCAG 1.4.3 verlangt fuer Text unter 18,5px ein Kontrastverhaeltnis von
 * mindestens 4,5:1. --faint stand bis zum 07.09.2026 auf #7d7488 und kam
 * damit auf 4,15:1 gegen den Kartengrund. Sichtbar ist das nicht, es
 * liest sich einfach nur "dezent".
 *
 * Der Test rechnet gegen den Kartengrund --panel, weil dort die meisten
 * Stellen mit dieser Farbe stehen (Tabellenkoepfe, Fusszeilen). Gegen
 * den dunkleren Seitenhintergrund ist der Wert immer besser.
 */
func TestFarbenHabenGenugKontrast(t *testing.T) {
	css := vorlageLesen(t)

	faelle := []struct {
		vordergrund, hintergrund string
		mindestens               float64
		wozu                     string
	}{
		{"--faint", "--panel", 4.5, "Tabellenköpfe und Fußzeilen"},
		{"--muted", "--panel", 4.5, "Feststellungen und Empfehlungen"},
		{"--fg", "--panel", 4.5, "Fließtext"},
		{"--accent", "--panel", 4.5, "Links und Aufklapper"},
		{"--ok", "--panel", 4.5, "Marke einer Bestätigung"},
		{"--warn", "--panel", 4.5, "Marke eines Hinweises"},
	}

	for _, f := range faelle {
		vg := farbeAus(t, css, f.vordergrund)
		hg := farbeAus(t, css, f.hintergrund)
		v := kontrast(vg, hg)
		if v < f.mindestens {
			t.Errorf("%s auf %s (%s) ergibt %.2f:1, verlangt sind %.1f:1",
				f.vordergrund, f.hintergrund, f.wozu, v, f.mindestens)
		}
	}
}

/*
 * WCAG 2.4.7: Wer die Seite mit der Tabulatortaste bedient, muss sehen,
 * wo er steht. Auf diesem dunklen Grund ist der Standardrahmen des
 * Browsers praktisch unsichtbar, es braucht also eine eigene Regel.
 */
func TestFokusIstSichtbar(t *testing.T) {
	css := vorlageLesen(t)
	if !strings.Contains(css, ":focus-visible") {
		t.Error("keine :focus-visible-Regel in der Vorlage")
	}
	regel := regexp.MustCompile(`:focus-visible\s*\{[^}]*outline\s*:`)
	if !regel.MatchString(css) {
		t.Error(":focus-visible ist da, setzt aber kein outline")
	}
}

/*
 * WCAG 1.3.1: Vorlesewerkzeuge bauen aus den Ueberschriften ein
 * Sprungverzeichnis. Standen die Befundtitel als div da, kam man nur
 * durch Durchlesen an sie heran.
 */
func TestBefundtitelSindUeberschriften(t *testing.T) {
	html := volleSeite(t)
	_, _, _, alle := beispielBefunde()

	for _, b := range alle {
		gesucht := "<h3 class=\"titel\">" + b.Titel + "</h3>"
		if !strings.Contains(html, gesucht) {
			t.Errorf("Titel %q steht nicht als h3 auf der Seite", b.Titel)
		}
	}
	if strings.Contains(html, "<div class=\"titel\">") {
		t.Error("es gibt noch Befundtitel als div")
	}
}

/*
 * WCAG 1.4.1: In der Karte "kostet dich gerade Leistung" stehen zwei
 * Schweregrade nebeneinander, Hinweise und Anmerkungen. Bis zum
 * 07.09.2026 unterschied sie allein die Farbe des Streifens links. Wer
 * Farben nicht unterscheidet, sah zwei gleiche Kaesten.
 */
func TestSchweregradStehtNichtNurInDerFarbe(t *testing.T) {
	html := volleSeite(t)

	for _, wort := range []string{"Kostet Leistung", "Nur zur Einordnung"} {
		if !strings.Contains(html, "<p class=\"marke\">"+wort+"</p>") {
			t.Errorf("die Marke %q fehlt auf der Seite", wort)
		}
	}

	// Und sie muessen zum Streifen passen, sonst widersprechen sich Wort
	// und Farbe.
	hinweis := regexp.MustCompile(
		`class="befund hinweis">\s*(?:\n\s*)?<p class="marke">Kostet Leistung</p>`)
	if !hinweis.MatchString(html) {
		t.Error("der gelbe Streifen und die Marke \"Kostet Leistung\" gehoeren nicht zusammen")
	}
}

/*
 * Ein <th> ohne Inhalt meldet ein Vorlesewerkzeug als leer und liest die
 * Zelle daneben ohne jeden Bezug vor. Bei mehreren Laufwerken und
 * mehreren PCIe-Geraeten war genau das der Fall.
 */
func TestKeineLeerenTabellenkoepfe(t *testing.T) {
	html := volleSeite(t)

	leer := regexp.MustCompile(`<th[^>]*>\s*</th>`)
	if treffer := leer.FindAllString(html, -1); len(treffer) > 0 {
		t.Errorf("%d leere Kopfzelle(n) in der Tabelle: %q", len(treffer), treffer)
	}

	// Drei Laufwerke, ein Kopf darueber.
	if !strings.Contains(html, `rowspan="3">Laufwerke</th>`) {
		t.Error("der Kopf ueber den Laufwerken spannt nicht ueber alle drei Zeilen")
	}
	if !strings.Contains(html, `rowspan="2">PCIe-Anbindung</th>`) {
		t.Error("der Kopf ueber den PCIe-Zeilen spannt nicht ueber beide Zeilen")
	}
}

/*
 * WCAG 4.1.3: Die Meldungen erscheinen erst nach einem Klick. Ohne Rolle
 * bemerkt ein Vorlesewerkzeug die Aenderung nicht, es liest nur, was
 * beim Laden da war. Der Nutzer wartet dann auf eine Antwort, die schon
 * auf dem Bildschirm steht.
 */
func TestMeldungenWerdenAngesagt(t *testing.T) {
	html := volleSeite(t)

	for _, id := range []string{"meldung", "eintrag-meldung", "mess-urteil"} {
		// Die Rolle muss im selben Tag stehen wie die id.
		tag := regexp.MustCompile(`<[^>]*id="` + regexp.QuoteMeta(id) + `"[^>]*>`)
		gefunden := tag.FindString(html)
		if gefunden == "" {
			t.Errorf("Element mit id %q gibt es nicht mehr", id)
			continue
		}
		if !strings.Contains(gefunden, `role="status"`) {
			t.Errorf("%s ist kein Meldebereich: %s", id, gefunden)
		}
	}

	// Fehler sollen lauter angesagt werden als Erfolge.
	if !strings.Contains(html, `setAttribute("role", "alert")`) {
		t.Error("Fehlermeldungen wechseln nicht auf role=alert")
	}
}

/*
 * Der Sprung auf die Ergebnisseite stand bis zum 07.09.2026 unmittelbar
 * hinter der Erfolgsmeldung. Damit war die Meldung zwar da, aber niemand
 * las sie, der Browser war schon weg. Ein Vorlesewerkzeug meldet einen
 * geaenderten Bereich erst nach einer kurzen Pause und kam gar nicht
 * dazu.
 */
func TestSprungLaesstZeitZumLesen(t *testing.T) {
	html := volleSeite(t)

	sofort := regexp.MustCompile(`meldung\.appendChild\(verweis\);\s*\n\s*window\.location\.href`)
	if sofort.MatchString(html) {
		t.Error("der Sprung passiert unmittelbar nach der Meldung")
	}
	verzoegert := regexp.MustCompile(`setTimeout\([\s\S]*?window\.location\.href[\s\S]*?,\s*(\d+)\)`)
	m := verzoegert.FindStringSubmatch(html)
	if m == nil {
		t.Fatal("kein verzoegerter Sprung auf die Ergebnisseite gefunden")
	}
	var ms int
	fmt.Sscanf(m[1], "%d", &ms)
	if ms < 2000 {
		t.Errorf("der Sprung kommt nach %d ms, das reicht zum Lesen nicht", ms)
	}
}

// ---------------------------------------------------------------------
// Hilfsmittel
// ---------------------------------------------------------------------

func vorlageLesen(t *testing.T) string {
	t.Helper()
	roh, err := vorlagen.ReadFile("assets/preview.html.tmpl")
	if err != nil {
		t.Fatalf("Vorlage laesst sich nicht lesen: %v", err)
	}
	return string(roh)
}

// farbeAus holt den Wert einer CSS-Variablen aus der Vorlage.
func farbeAus(t *testing.T, css, name string) [3]float64 {
	t.Helper()
	re := regexp.MustCompile(regexp.QuoteMeta(name) + `:\s*#([0-9a-fA-F]{6})`)
	m := re.FindStringSubmatch(css)
	if m == nil {
		t.Fatalf("Farbe %s steht nicht in der Vorlage", name)
	}
	var r, g, b int
	fmt.Sscanf(m[1], "%02x%02x%02x", &r, &g, &b)
	return [3]float64{float64(r), float64(g), float64(b)}
}

// kontrast rechnet das Verhaeltnis nach WCAG 2.x aus. Bewusst von Hand
// und nicht ueber eine Bibliothek: Es sind sechs Zeilen, und eine
// Abhaengigkeit mehr in einem Programm, das die Leute selbst
// nachbauen sollen, kostet mehr als sie einbringt.
func kontrast(a, b [3]float64) float64 {
	la, lb := helligkeit(a), helligkeit(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

func helligkeit(f [3]float64) float64 {
	k := func(c float64) float64 {
		c /= 255
		if c <= 0.04045 {
			return c / 12.92
		}
		return math.Pow((c+0.055)/1.055, 2.4)
	}
	return 0.2126*k(f[0]) + 0.7152*k(f[1]) + 0.0722*k(f[2])
}
