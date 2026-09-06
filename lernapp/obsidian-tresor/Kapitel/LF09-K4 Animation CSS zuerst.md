---
lernfeld: 9
kapitel: 4
dauer: 30
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-09
---

# Animation: CSS zuerst

**Lernfeld 09, Kapitel 4** · 30 Minuten · [[LF09 JavaScript fuer Gestalter|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du animierst fluessig und weisst, welche Eigenschaften teuer sind.

## Das musst du wissen

- `transform` und `opacity` sind billig, weil der Browser sie ohne Neuberechnung des Layouts zeichnen kann. `width`, `height`, `top` und `margin` zu animieren, ruckelt.
- Dauer: 150 bis 300 Millisekunden fuer kleine Uebergaenge. Alles ueber 400 fuehlt sich zaeh an, egal wie schoen es aussieht.
- Animation soll etwas **erklaeren**: woher kam das Element, was hat sich geaendert. Bewegung ohne Aussage ist Laerm und veraltet schnell.
- `prefers-reduced-motion` respektieren, immer. Und: nicht alles gleichzeitig animieren. Ein orchestrierter Moment wirkt staerker als zehn Effekte.

## Aufgabe

> [!todo] Selber machen
> Baue einen Uebergang fuer ein Aufklapp-Menue, einmal mit height (ruckelt) und einmal mit transform oder grid-template-rows. Vergleiche im Leistungsreiter.

## Selbstpruefung

> [!question]
> Welche zwei CSS-Eigenschaften kann man am guenstigsten animieren?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: css animation performance transform opacity tutorial](https://www.youtube.com/results?search_query=css+animation+performance+transform+opacity+tutorial) · empfohlen: Kevin Powell / Chrome for Developers
- [MDN: Uebergaenge](https://developer.mozilla.org/de/docs/Web/CSS/CSS_transitions/Using_CSS_transitions)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF09-K3 Formulare, Pruefung, Zustaende]] · [[LF09 JavaScript fuer Gestalter|Lernfeld]] · [[00 Start]] · [[LF09-K5 Wann ein Framework, wann nicht]] →
