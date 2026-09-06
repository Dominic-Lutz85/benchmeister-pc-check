---
lernfeld: 3
kapitel: 3
dauer: 30
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-03
---

# Einheiten und Farbe

**Lernfeld 03, Kapitel 3** · 30 Minuten · [[LF03 CSS Grundlagen|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du waehlst Einheiten begruendet statt aus Gewohnheit.

## Das musst du wissen

- `rem` haengt an der Grundschriftgroesse des Browsers und respektiert damit die Einstellung des Nutzers. `px` tut das nicht. Fuer Schriftgroessen deshalb rem, fuer Rahmen und Haarlinien ruhig px.
- `%` bezieht sich auf das Elternelement, `vw/vh` auf das Fenster, `ch` auf die Breite einer Null in der aktuellen Schrift. `ch` ist ideal fuer Zeilenlaengen: `max-width:65ch` ist besser begruendet als 700px.
- Farbe: `hsl()` ist von Hand steuerbar (Farbton, Saettigung, Helligkeit), `oklch()` ist neuer und rechnet in wahrgenommener Helligkeit. Wer in oklch nur den Farbton dreht, behaelt gleich helle Farben, in hsl nicht.
- Farbe hat immer ein Gegenstueck: den Kontrast zum Hintergrund. Diese Pruefung kommt in Lernfeld 7 und entscheidet mit, ob eine Farbe ueberhaupt in Frage kommt.

## Aufgabe

> [!todo] Selber machen
> Baue eine Farbpalette aus fuenf Werten in hsl, dann dieselbe in oklch. Aendere in beiden nur den Farbton und vergleiche, ob die Helligkeit gleich bleibt.

## Selbstpruefung

> [!question]
> Warum kann eine Seite mit px-Schriftgroessen fuer sehbehinderte Nutzer ein Problem sein?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: css einheiten rem em vw ch erklaert](https://www.youtube.com/results?search_query=css+einheiten+rem+em+vw+ch+erklaert) · empfohlen: Kevin Powell
- [MDN: Werte und Einheiten](https://developer.mozilla.org/de/docs/Learn_web_development/Core/Styling_basics/Values_and_units)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF03-K2 Das Boxmodell]] · [[LF03 CSS Grundlagen|Lernfeld]] · [[00 Start]] · [[LF03-K4 Typografie in CSS]] →
