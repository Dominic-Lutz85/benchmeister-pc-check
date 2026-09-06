---
lernfeld: 4
kapitel: 5
dauer: 25
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-04
---

# Fliessende Groessen: clamp, min, max

**Lernfeld 04, Kapitel 5** · 25 Minuten · [[LF04 Layout|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du laesst Schrift und Abstaende stufenlos mitwachsen statt in Spruengen.

## Das musst du wissen

- `clamp(klein, wunsch, gross)` ist der Arbeitspferd-Aufruf: `font-size:clamp(1.75rem, 4vw, 3rem)` waechst mit dem Fenster, bleibt aber in Grenzen.
- Der mittlere Wert sollte immer einen `rem`-Anteil enthalten (etwa `1rem + 2vw`), sonst ignoriert die Groesse die Zoomeinstellung des Nutzers. Das ist ein echter Barrierefreiheitsfehler, den man leicht uebersieht.
- Dasselbe funktioniert fuer Abstaende und Breiten: `width:min(100%, 65ch)` heisst 'so breit wie noetig, aber nie ueber 65 Zeichen'.
- Weniger Umbruchpunkte bedeuten weniger Stellen, an denen etwas kaputtgehen kann. Fliessende Werte sind deshalb nicht nur eleganter, sondern robuster.

## Aufgabe

> [!todo] Selber machen
> Ersetze in deiner Uebungsseite alle Schriftgroessen der Ueberschriften durch clamp-Werte und pruefe bei 320px, 768px und 1600px Fensterbreite.

## Selbstpruefung

> [!question]
> Warum darf im mittleren clamp-Wert nicht nur vw stehen?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: css clamp fluid typography tutorial](https://www.youtube.com/results?search_query=css+clamp+fluid+typography+tutorial) · empfohlen: Kevin Powell / Utopia
- [MDN: clamp()](https://developer.mozilla.org/de/docs/Web/CSS/clamp)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF04-K4 Responsive Media und Container Queries]] · [[LF04 Layout|Lernfeld]] · [[00 Start]] · [[LF04-K6 Layoutmuster, die immer wiederkommen]] →
