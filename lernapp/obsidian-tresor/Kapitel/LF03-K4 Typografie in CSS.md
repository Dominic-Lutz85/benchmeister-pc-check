---
lernfeld: 3
kapitel: 4
dauer: 35
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-03
---

# Typografie in CSS

**Lernfeld 03, Kapitel 4** · 35 Minuten · [[LF03 CSS Grundlagen|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du stellst Text so ein, dass er ueber laengere Strecken angenehm lesbar bleibt.

## Das musst du wissen

- Schrift laedt man ueber `@font-face` oder einen Dienst. Immer eine echte Rueckfallkette angeben (`"Meine Schrift", Georgia, serif`), sonst zeigt der Browser bei einem Fehler die Standardschrift und die Seite kippt.
- `font-display:swap` sorgt dafuer, dass der Text sofort in der Rueckfallschrift erscheint, statt unsichtbar zu bleiben, bis die Schrift geladen ist.
- Zeilenhoehe: etwa 1,5 bei Fliesstext, deutlich enger (1,1 bis 1,25) bei grossen Ueberschriften. Je groesser die Schrift, desto weniger Durchschuss braucht sie.
- `text-wrap:balance` fuer Ueberschriften verteilt die Woerter gleichmaessig auf die Zeilen und verhindert eine einsame letzte Silbe. Eine Zeile Code, sichtbarer Unterschied.

## Aufgabe

> [!todo] Selber machen
> Setze denselben Absatz dreimal: mit 45, 65 und 110 Zeichen Zeilenlaenge. Lies alle drei laut und entscheide selbst.

## Selbstpruefung

> [!question]
> Warum braucht eine 48px-Ueberschrift eine kleinere Zeilenhoehe als ein 17px-Absatz?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: web typografie css line-height schriftgroesse tutorial](https://www.youtube.com/results?search_query=web+typografie+css+line-height+schriftgroesse+tutorial) · empfohlen: Kevin Powell / Design Course
- [MDN: Text gestalten](https://developer.mozilla.org/de/docs/Learn_web_development/Core/Text_styling/Fundamentals)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF03-K3 Einheiten und Farbe]] · [[LF03 CSS Grundlagen|Lernfeld]] · [[00 Start]] · [[LF03-K5 Eigene Eigenschaften (CSS-Variablen)]] →
