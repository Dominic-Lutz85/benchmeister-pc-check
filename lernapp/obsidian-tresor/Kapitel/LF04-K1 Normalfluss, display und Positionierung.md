---
lernfeld: 4
kapitel: 1
dauer: 30
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-04
---

# Normalfluss, display und Positionierung

**Lernfeld 04, Kapitel 1** · 30 Minuten · [[LF04 Layout|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du verstehst, was der Browser von allein tut, bevor du eingreifst.

## Das musst du wissen

- Ohne CSS stapelt der Browser Blockelemente untereinander und setzt Inline-Elemente nebeneinander in den Textfluss. Das ist der **Normalfluss**, und er ist oft schon fast richtig.
- `position:relative` verschiebt ein Element, laesst seinen Platz aber frei. `absolute` nimmt es aus dem Fluss und orientiert es am naechsten positionierten Vorfahren. `fixed` haelt es am Fenster, `sticky` laesst es mitlaufen und dann kleben.
- `sticky` funktioniert nur, wenn ein Versatz gesetzt ist (`top:0`) und kein Vorfahre `overflow:hidden` hat. Das ist der Grund, warum es scheinbar grundlos nicht klebt.
- Absolute Positionierung ist fuer Ausnahmen da: ein Abzeichen auf einer Karte, ein Schliessen-Kreuz. Ein ganzes Layout so zu bauen, faellt spaetestens auf dem Handy auseinander.

## Aufgabe

> [!todo] Selber machen
> Setz ein Abzeichen in die obere rechte Ecke einer Karte mit relative und absolute. Danach mach eine Kopfzeile sticky und finde heraus, was sie kaputt macht.

## Selbstpruefung

> [!question]
> Woran erkennt ein absolut positioniertes Element, worauf es sich bezieht?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: css position relative absolute sticky erklaert deutsch](https://www.youtube.com/results?search_query=css+position+relative+absolute+sticky+erklaert+deutsch) · empfohlen: Kevin Powell
- [MDN: Positionierung](https://developer.mozilla.org/de/docs/Learn_web_development/Core/CSS_layout/Positioning)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF03-K5 Eigene Eigenschaften (CSS-Variablen)]] · [[LF04 Layout|Lernfeld]] · [[00 Start]] · [[LF04-K2 Flexbox]] →
