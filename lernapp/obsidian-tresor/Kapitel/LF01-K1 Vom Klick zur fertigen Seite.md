---
lernfeld: 1
kapitel: 1
dauer: 25
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-01
---

# Vom Klick zur fertigen Seite

**Lernfeld 01, Kapitel 1** · 25 Minuten · [[LF01 Wie das Web funktioniert|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du kannst in eigenen Worten erklaeren, was zwischen der Eingabe einer Adresse und dem sichtbaren Bild passiert.

## Das musst du wissen

- Eine Adresse ist ein Wegweiser, kein Ort. **DNS** uebersetzt `beispiel.de` in eine IP-Adresse, das ist die eigentliche Hausnummer im Netz.
- Der Browser stellt eine **Anfrage** (Request), der Server schickt eine **Antwort** (Response). In der Antwort steht ein Statuscode: 200 heisst gefunden, 404 nicht gefunden, 301 umgezogen, 500 Server kaputt.
- Die erste Antwort ist fast immer nur das HTML. Darin stehen Verweise auf CSS, Bilder und Skripte, die der Browser danach einzeln nachfordert. Jede dieser Dateien ist eine eigene Anfrage, und genau das ist spaeter der Hebel fuer Ladezeit.
- **HTTPS** ist HTTP mit Verschluesselung. Ohne ist heute keine Seite mehr vertretbar, Browser markieren sie als unsicher.

## Aufgabe

> [!todo] Selber machen
> Oeffne eine Seite, die du magst, druecke F12, geh auf den Reiter Netzwerk und lade neu. Zaehle, wie viele Anfragen laufen und welche am laengsten dauert. Schreib beides in die Notiz.

## Selbstpruefung

> [!question]
> Warum sieht man bei einer Seite oft erst Text und danach erst die Bilder?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: wie funktioniert das internet einfach erklaert HTTP DNS](https://www.youtube.com/results?search_query=wie+funktioniert+das+internet+einfach+erklaert+HTTP+DNS) · empfohlen: Simply Explained / Doktor Whatson
- [MDN: Wie das Web funktioniert](https://developer.mozilla.org/de/docs/Learn_web_development/Getting_started/Web_standards/How_the_web_works)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
[[LF01 Wie das Web funktioniert|Lernfeld]] · [[00 Start]] · [[LF01-K2 Dateien, Ordner, Pfade]] →
