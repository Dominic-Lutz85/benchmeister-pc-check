package ui

import (
	"net"
	"net/http"
	"strings"
)

// herkunftErlaubt prüft, ob eine Anfrage wirklich von der eigenen Seite
// kommt, und nicht von einer fremden Webseite im selben Browser.
//
// WOGEGEN DAS SCHÜTZT: DNS-Rebinding. Der Server hört zwar nur auf
// 127.0.0.1 und ist damit aus dem Netz nicht erreichbar. Eine präparierte
// Webseite kann den Browser aber dazu bringen, ihren eigenen Domainnamen
// auf 127.0.0.1 aufzulösen. Für den Browser sieht das dann wie dieselbe
// Herkunft aus, er darf also nicht nur senden, sondern auch die Antwort
// lesen. Ohne die Prüfung hier könnte eine solche Seite, solange das
// Programm läuft:
//
//   - den Upload auslösen und dabei die Marktforschungs-Zustimmung
//     selbst setzen, obwohl die Person sie nicht gegeben hat. Das wäre
//     nicht nur ein Sicherheits-, sondern ein Einwilligungsproblem: Die
//     Datenschutzerklärung sagt ausdrücklich zu, dass diese Zustimmung
//     getrennt und freiwillig ist.
//   - die Startseite auslesen, auf der das vollständige Scan-Ergebnis
//     steht, und den Abruf-Token aus der Antwort mitnehmen.
//
// WIE DIE PRÜFUNG WIRKT: Entscheidend ist der Host-Kopf. Bei einem
// Rebinding-Angriff trägt er den Domainnamen des Angreifers, weil der
// Browser dort einträgt, was in der Adresszeile stand, nicht wohin es
// aufgelöst wurde. Nur unsere eigene Seite schickt "127.0.0.1:<port>".
// Damit ist der Angriff abgewehrt, unabhängig davon, ob ein
// Origin-Kopf mitkommt.
//
// Der Origin-Kopf wird zusätzlich geprüft, aber nur wenn er da ist:
// Browser senden ihn bei POST inzwischen auch bei gleicher Herkunft,
// verlassen wollen wir uns darauf aber nicht.
func herkunftErlaubt(r *http.Request, port string) bool {
	if !hostErlaubt(r.Host, port) {
		return false
	}

	herkunft := r.Header.Get("Origin")
	if herkunft == "" {
		return true
	}
	rest, gefunden := strings.CutPrefix(herkunft, "http://")
	return gefunden && hostErlaubt(rest, port)
}

// hostErlaubt lässt genau den eigenen Zuhörer durch, unter beiden
// Schreibweisen. "localhost" ist mit drin, weil jemand die Adresse von
// Hand eintippen könnte (der Hinweis dazu steht in der Konsolenausgabe).
// Das schwächt nichts ab: Ein Rebinding-Angriff braucht einen eigenen
// Domainnamen, und "localhost" lässt sich nicht auf einen fremden Rechner
// umbiegen.
func hostErlaubt(host, port string) bool {
	name, hostPort, err := net.SplitHostPort(host)
	if err != nil || hostPort != port {
		return false
	}
	return name == "127.0.0.1" || strings.EqualFold(name, "localhost")
}

// nurEigeneHerkunft legt sich um den gesamten Server. Bewusst um ALLE
// Pfade, nicht nur um /consent: Auf der Startseite steht das vollständige
// Scan-Ergebnis, das ist genauso schützenswert wie der Upload selbst.
func nurEigeneHerkunft(port string, weiter http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !herkunftErlaubt(r, port) {
			http.Error(w, "Diese Anfrage kommt nicht von der Ergebnisseite dieses Programms.",
				http.StatusForbidden)
			return
		}
		weiter.ServeHTTP(w, r)
	})
}
