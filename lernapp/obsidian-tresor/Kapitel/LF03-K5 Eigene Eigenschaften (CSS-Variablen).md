---
lernfeld: 3
kapitel: 5
dauer: 25
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-03
---

# Eigene Eigenschaften (CSS-Variablen)

**Lernfeld 03, Kapitel 5** · 25 Minuten · [[LF03 CSS Grundlagen|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du legst ein Farb- und Abstandssystem an, das an einer Stelle aenderbar ist.

## Das musst du wissen

- Deklariert wird mit zwei Bindestrichen im `:root`, benutzt mit `var()`: `--akzent:#0C6B5E`, dann `color:var(--akzent)`.
- Der grosse Unterschied zu Variablen in Praeprozessoren: diese leben im Browser. Man kann sie zur Laufzeit umschreiben, und genau so baut man einen hellen und einen dunklen Modus mit einem einzigen Regelblock.
- `var(--x, blau)` nimmt blau, falls `--x` nicht gesetzt ist. Nuetzlich fuer Komponenten, die auch ohne Systemdatei funktionieren sollen.
- Vergib Namen nach **Rolle**, nicht nach Aussehen: `--flaeche` und `--akzent` statt `--hellgrau` und `--gruen`. Sonst heisst deine Variable im dunklen Modus hellgrau und ist dunkel.

## Aufgabe

> [!todo] Selber machen
> Zieh alle Farben deiner Uebungsseite in `:root` und stell danach in einer Media Query fuer den dunklen Modus nur die Variablen um, keine einzige Komponentenregel.

## Selbstpruefung

> [!question]
> Warum ist `--farbe-blau` ein schlechterer Name als `--farbe-akzent`?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: css custom properties variablen dark mode tutorial](https://www.youtube.com/results?search_query=css+custom+properties+variablen+dark+mode+tutorial) · empfohlen: Kevin Powell
- [MDN: Eigene CSS-Eigenschaften](https://developer.mozilla.org/de/docs/Web/CSS/--*)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF03-K4 Typografie in CSS]] · [[LF03 CSS Grundlagen|Lernfeld]] · [[00 Start]] · [[LF04-K1 Normalfluss, display und Positionierung]] →
