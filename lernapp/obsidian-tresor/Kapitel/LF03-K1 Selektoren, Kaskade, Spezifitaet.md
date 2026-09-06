---
lernfeld: 3
kapitel: 1
dauer: 35
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-03
---

# Selektoren, Kaskade, Spezifitaet

**Lernfeld 03, Kapitel 1** · 35 Minuten · [[LF03 CSS Grundlagen|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du kannst erklaeren, warum eine Regel gewinnt und eine andere nicht.

## Das musst du wissen

- CSS entscheidet Konflikte in dieser Reihenfolge: **Herkunft**, dann **Spezifitaet**, dann **Reihenfolge im Code**. Bei gleicher Spezifitaet gewinnt die spaetere Regel.
- Spezifitaet als Zahlentripel lesen: ID = (1,0,0), Klasse oder Attribut oder Pseudoklasse = (0,1,0), Element = (0,0,1). Eine ID schlaegt damit hundert Klassen, und das ist der Grund, warum man IDs in CSS meidet.
- `!important` ist kein Werkzeug, sondern eine Kapitulation. Wer es einmal setzt, braucht es beim naechsten Mal wieder. Loesung ist fast immer ein einfacherer Selektor, nicht ein staerkerer.
- **Vererbung** ist etwas anderes als Kaskade: Schrift, Farbe und Zeilenhoehe geben sich an Kindelemente weiter, Abstaende und Rahmen nicht. Deshalb setzt man Schrift einmal am `body` und ist fertig.

## Aufgabe

> [!todo] Selber machen
> Bau eine Seite mit drei Regeln, die absichtlich um dasselbe Element streiten. Sag vorher voraus, welche gewinnt, und pruefe es in den Entwicklerwerkzeugen.

## Selbstpruefung

> [!question]
> Warum gilt die Empfehlung, in CSS keine IDs als Selektoren zu benutzen?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: css spezifitaet kaskade erklaert deutsch](https://www.youtube.com/results?search_query=css+spezifitaet+kaskade+erklaert+deutsch) · empfohlen: Kevin Powell
- [MDN: Kaskade und Vererbung](https://developer.mozilla.org/de/docs/Learn_web_development/Core/Styling_basics/Handling_conflicts)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF02-K6 Semantik die richtigen Bausteine]] · [[LF03 CSS Grundlagen|Lernfeld]] · [[00 Start]] · [[LF03-K2 Das Boxmodell]] →
