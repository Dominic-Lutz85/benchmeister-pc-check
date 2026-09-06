---
lernfeld: 10
kapitel: 2
dauer: 35
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-10
---

# Bilder und Schriften optimieren

**Lernfeld 10, Kapitel 2** · 35 Minuten · [[LF10 Tempo, Auffindbarkeit, Qualitaet|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du halbierst die Ladezeit typischer Seiten mit wenigen Handgriffen.

## Das musst du wissen

- Bilder sind meist siebzig bis achtzig Prozent des Gewichts einer Seite. Hier liegt der groesste Hebel, und er kostet keine Gestaltung.
- Richtige Groesse liefern, nicht ein 4000-Pixel-Foto per CSS auf 400 schrumpfen. Mit `srcset` und `sizes` laesst du den Browser die passende Fassung waehlen.
- Schriften: nur die Schnitte laden, die du wirklich benutzt. Vier Schnitte einer Familie sind schnell ein halbes Megabyte. Variable Fonts koennen guenstiger sein, wenn man mehrere Gewichte braucht.
- `font-display:swap` und `preload` fuer die eine Schrift, die oben sichtbar ist. Alles andere darf warten.

## Aufgabe

> [!todo] Selber machen
> Nimm eine bestehende Seite, komprimiere alle Bilder, setz srcset ein und reduziere die Schriftschnitte. Miss vorher und nachher.

## Selbstpruefung

> [!question]
> Warum reicht es nicht, ein grosses Bild per CSS kleiner darzustellen?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: bilder optimieren srcset webp schriften laden](https://www.youtube.com/results?search_query=bilder+optimieren+srcset+webp+schriften+laden) · empfohlen: web.dev / Kevin Powell
- [MDN: Responsive Bilder](https://developer.mozilla.org/de/docs/Web/HTML/Guides/Responsive_images)
- [Squoosh (Kompression)](https://squoosh.app/)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF10-K1 Ladezeit verstehen]] · [[LF10 Tempo, Auffindbarkeit, Qualitaet|Lernfeld]] · [[00 Start]] · [[LF10-K3 SEO fuer Gestalter]] →
