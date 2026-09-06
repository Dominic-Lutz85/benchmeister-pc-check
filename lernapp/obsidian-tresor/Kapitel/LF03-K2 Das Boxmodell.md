---
lernfeld: 3
kapitel: 2
dauer: 30
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-03
---

# Das Boxmodell

**Lernfeld 03, Kapitel 2** · 30 Minuten · [[LF03 CSS Grundlagen|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du berechnest im Kopf, wie breit ein Element am Ende wirklich ist.

## Das musst du wissen

- Jedes Element ist eine Schachtel aus vier Ringen: Inhalt, **padding** (Innenabstand), **border** (Rahmen), **margin** (Aussenabstand).
- `box-sizing:border-box` gehoert in jedes Projekt ganz oben hin. Damit meint `width:300px` tatsaechlich 300 Pixel inklusive Rahmen und Innenabstand, statt 300 plus Zugaben.
- Aussenabstaende benachbarter Elemente **fallen zusammen** (margin collapsing): 20px unten und 30px oben ergeben 30px, nicht 50px. Das ueberrascht jeden Anfaenger genau einmal.
- Deshalb ist es meist ruhiger, Abstaende zwischen Geschwistern per `gap` im Flex- oder Grid-Behaelter zu setzen statt per margin an jedem Kind.

## Aufgabe

> [!todo] Selber machen
> Setze eine Box auf 300px Breite mit 20px padding und 5px border, einmal mit und einmal ohne border-box, und miss beide in den Entwicklerwerkzeugen.

## Selbstpruefung

> [!question]
> Wie breit ist eine Box mit width:200px, padding:16px, border:2px bei content-box?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: css box model erklaert border-box deutsch](https://www.youtube.com/results?search_query=css+box+model+erklaert+border-box+deutsch) · empfohlen: Kevin Powell / Programmieren lernen
- [MDN: Das Boxmodell](https://developer.mozilla.org/de/docs/Learn_web_development/Core/Styling_basics/Box_model)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF03-K1 Selektoren, Kaskade, Spezifitaet]] · [[LF03 CSS Grundlagen|Lernfeld]] · [[00 Start]] · [[LF03-K3 Einheiten und Farbe]] →
