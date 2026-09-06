---
lernfeld: 4
kapitel: 4
dauer: 40
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-04
---

# Responsive: Media und Container Queries

**Lernfeld 04, Kapitel 4** · 40 Minuten · [[LF04 Layout|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Deine Seite funktioniert vom Handy bis zum grossen Monitor, ohne Sonderfaelle zu sammeln.

## Das musst du wissen

- **Zuerst schmal bauen** (mobile first), dann per `@media (min-width:...)` erweitern. Andersherum sammelt man Ausnahmen, die man einzeln zuruecknehmen muss.
- Umbruchpunkte richten sich nach dem **Inhalt**, nicht nach Geraetenamen. Zieh das Fenster langsam breiter und setze den Punkt dort, wo es haesslich wird. Eine Liste von iPhone-Breiten veraltet jedes Jahr.
- **Container Queries** (`@container`) fragen die Breite des Behaelters statt des Fensters. Damit funktioniert dieselbe Karte in der schmalen Seitenleiste und im breiten Hauptbereich, ohne Sonderklassen. Das ist die groesste Layout-Neuerung der letzten Jahre.
- Vergiss die Extreme nicht: sehr schmal (320px) und sehr breit. Ein `max-width` auf dem Inhaltsbereich verhindert Zeilen ueber die ganze Wand.

## Aufgabe

> [!todo] Selber machen
> Nimm deine Kartenkomponente und setz sie einmal in eine schmale Seitenleiste und einmal in den breiten Bereich. Loese den Unterschied mit `@container` statt mit einer zweiten Klasse.

## Selbstpruefung

> [!question]
> Warum sind Umbruchpunkte nach Geraetemodellen eine schlechte Idee?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: css container queries tutorial deutsch](https://www.youtube.com/results?search_query=css+container+queries+tutorial+deutsch) · empfohlen: Kevin Powell
- [MDN: Media Queries](https://developer.mozilla.org/de/docs/Web/CSS/CSS_media_queries/Using_media_queries)
- [MDN: Container Queries](https://developer.mozilla.org/de/docs/Web/CSS/CSS_containment/Container_queries)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF04-K3 Grid]] · [[LF04 Layout|Lernfeld]] · [[00 Start]] · [[LF04-K5 Fliessende Groessen clamp, min, max]] →
