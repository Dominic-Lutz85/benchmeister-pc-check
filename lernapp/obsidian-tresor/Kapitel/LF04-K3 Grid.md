---
lernfeld: 4
kapitel: 3
dauer: 45
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-04
---

# Grid

**Lernfeld 04, Kapitel 3** · 45 Minuten · [[LF04 Layout|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du baust zweidimensionale Layouts, die sich ohne Media Query anpassen.

## Das musst du wissen

- Grid denkt in Zeilen **und** Spalten gleichzeitig. `grid-template-columns` legt die Spalten fest, `fr` ist der uebrige Platz als Anteil.
- Die wichtigste Zeile ueberhaupt: `grid-template-columns:repeat(auto-fill,minmax(240px,1fr))`. Damit passt sich ein Kartenraster von allein an jede Breite an, ohne eine einzige Media Query.
- Bereiche benennen macht komplexe Layouts lesbar: `grid-template-areas` zeichnet das Layout als Text im Code. Wer das einmal gemacht hat, will es nicht mehr anders.
- Faustregel: Grid fuer das Seitengeruest und Raster, Flexbox fuer alles darin. Sie konkurrieren nicht, sie ergaenzen sich.

## Aufgabe

> [!todo] Selber machen
> Baue eine Seite mit Kopf, Seitenleiste, Inhalt und Fuss ueber `grid-template-areas`, und darunter ein Kartenraster mit der auto-fill-Zeile. Zieh das Fenster schmal.

## Selbstpruefung

> [!question]
> Was ist der Unterschied zwischen `auto-fill` und `auto-fit`?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: css grid tutorial deutsch grid-template-areas](https://www.youtube.com/results?search_query=css+grid+tutorial+deutsch+grid-template-areas) · empfohlen: Kevin Powell / Programmieren lernen
- [MDN: Grid](https://developer.mozilla.org/de/docs/Learn_web_development/Core/CSS_layout/Grids)
- [Grid Garden (Spiel)](https://cssgridgarden.com/#de)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF04-K2 Flexbox]] · [[LF04 Layout|Lernfeld]] · [[00 Start]] · [[LF04-K4 Responsive Media und Container Queries]] →
