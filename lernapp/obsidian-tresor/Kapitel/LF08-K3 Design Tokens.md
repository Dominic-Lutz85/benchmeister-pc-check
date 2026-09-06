---
lernfeld: 8
kapitel: 3
dauer: 30
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-08
---

# Design Tokens

**Lernfeld 08, Kapitel 3** · 30 Minuten · [[LF08 Werkzeuge und Design-Systeme|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Gestaltung und Umsetzung teilen dieselben Werte unter denselben Namen.

## Das musst du wissen

- Ein **Token** ist ein benannter Wert: `farbe.akzent`, `abstand.3`, `radius.klein`. In Figma sind das Variablen, im Code CSS-Eigenschaften. Gleicher Name auf beiden Seiten.
- Zwei Ebenen bewaehren sich: **Grundtokens** (die Farbe petrol-600) und **Rollentokens** (die Akzentfarbe der Oberflaeche). Komponenten benutzen nur Rollentokens, nie Grundtokens.
- Genau dadurch wird ein dunkler Modus zur Kleinigkeit: nur die Zuordnung der Rollen aendert sich, kein einziger Baustein.
- Tokens sind erst dann etwas wert, wenn beide Seiten dieselben Namen benutzen. Eine Liste, die nur im Entwurf lebt, ist Dekoration.

## Aufgabe

> [!todo] Selber machen
> Lege in Figma Variablen fuer Farben und Abstaende an, benenne sie nach Rollen und uebertrage exakt dieselben Namen in deinen `:root`-Block.

## Selbstpruefung

> [!question]
> Warum sollte eine Komponente nie direkt auf einen Grundtoken zugreifen?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: design tokens erklaert figma variables css](https://www.youtube.com/results?search_query=design+tokens+erklaert+figma+variables+css) · empfohlen: Figma / Design Systems
- [Figma: Variablen](https://help.figma.com/hc/de/articles/15339657135383)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF08-K2 Auto Layout, Komponenten, Varianten]] · [[LF08 Werkzeuge und Design-Systeme|Lernfeld]] · [[00 Start]] · [[LF08-K4 Ein kleines Design-System bauen]] →
