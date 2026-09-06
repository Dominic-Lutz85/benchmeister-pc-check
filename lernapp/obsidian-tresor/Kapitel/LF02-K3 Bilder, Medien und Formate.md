---
lernfeld: 2
kapitel: 3
dauer: 30
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-02
---

# Bilder, Medien und Formate

**Lernfeld 02, Kapitel 3** · 30 Minuten · [[LF02 HTML Struktur und Bedeutung|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du baust Bilder so ein, dass sie schnell laden und auch ohne Sicht verstaendlich sind.

## Das musst du wissen

- Das `alt`-Attribut beschreibt, was auf dem Bild **passiert**, nicht wie die Datei heisst. Rein dekorative Bilder bekommen `alt=""`, damit Vorleseprogramme sie ueberspringen.
- Formate: **SVG** fuer Logos und Symbole (beliebig skalierbar, winzig), **WebP oder AVIF** fuer Fotos, JPEG als Rueckfall. PNG nur, wenn Transparenz mit harten Kanten noetig ist.
- Immer `width` und `height` angeben, auch wenn CSS die Groesse bestimmt. Sonst springt das Layout beim Nachladen, und genau das misst Google als Ruckeln.
- `loading="lazy"` bei allem, was beim Laden noch nicht sichtbar ist. Beim ersten grossen Bild oben aber nicht, das soll ja sofort kommen.

## Aufgabe

> [!todo] Selber machen
> Bau eine Seite mit drei Fotos. Schreib fuer jedes ein alt-Attribut, das jemandem am Telefon das Bild erklaert. Lass die Datei danach durch ein Kompressionswerkzeug laufen und vergleiche die Groesse.

## Selbstpruefung

> [!question]
> Wann ist ein leeres alt-Attribut die richtige Wahl, und wann ein Fehler?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: bilder fuer webseiten optimieren webp alt text](https://www.youtube.com/results?search_query=bilder+fuer+webseiten+optimieren+webp+alt+text) · empfohlen: Kevin Powell / web.dev
- [MDN: Bilder in HTML](https://developer.mozilla.org/de/docs/Learn_web_development/Core/Structuring_content/HTML_images)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF02-K2 Text, Listen und Links]] · [[LF02 HTML Struktur und Bedeutung|Lernfeld]] · [[00 Start]] · [[LF02-K4 Tabellen, und wann sie richtig sind]] →
