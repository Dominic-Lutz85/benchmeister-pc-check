---
lernfeld: 4
kapitel: 6
dauer: 30
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-04
---

# Layoutmuster, die immer wiederkommen

**Lernfeld 04, Kapitel 6** · 30 Minuten · [[LF04 Layout|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du hast eine kleine Sammlung erprobter Geruest-Loesungen im Kopf.

## Das musst du wissen

- **Der Stapel:** alles untereinander mit gleichem Abstand. Ein Flex-Behaelter mit `flex-direction:column` und `gap`, mehr nicht. Deckt gefuehlt die Haelfte aller Faelle ab.
- **Mitte mit Rand:** `width:min(100% - 2rem, 68rem); margin-inline:auto`. Zentriert und haelt gleichzeitig Luft zum Bildschirmrand.
- **Klebender Fuss:** Seite als Grid mit `grid-template-rows:auto 1fr auto` und `min-height:100dvh`. Der Fuss sitzt unten, auch wenn die Seite kurz ist. `dvh` statt `vh`, weil die Handy-Adressleiste sonst dazwischenfunkt.
- **Seitenleiste, die umbricht:** Flexbox mit `flex-wrap` und einer Grundbreite an der Leiste. Sie rutscht von allein unter den Inhalt, wenn es eng wird.

## Aufgabe

> [!todo] Selber machen
> Baue diese vier Muster in eine einzige Uebungsdatei und speicher sie als deine persoenliche Vorlage. Du wirst sie oft brauchen.

## Selbstpruefung

> [!question]
> Warum ist `100dvh` auf dem Handy besser als `100vh`?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: css layout patterns every layout tutorial](https://www.youtube.com/results?search_query=css+layout+patterns+every+layout+tutorial) · empfohlen: Kevin Powell / Every Layout
- [Web.dev: Zehn moderne Layouts](https://web.dev/patterns/layout/)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF04-K5 Fliessende Groessen clamp, min, max]] · [[LF04 Layout|Lernfeld]] · [[00 Start]] · [[LF05-K1 Typografie als Handwerk]] →
