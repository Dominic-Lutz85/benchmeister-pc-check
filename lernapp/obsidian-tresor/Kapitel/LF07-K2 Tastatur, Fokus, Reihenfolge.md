---
lernfeld: 7
kapitel: 2
dauer: 35
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-07
---

# Tastatur, Fokus, Reihenfolge

**Lernfeld 07, Kapitel 2** · 35 Minuten · [[LF07 Barrierefreiheit|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Deine Seite ist vollstaendig ohne Maus bedienbar, und man sieht immer, wo man ist.

## Das musst du wissen

- Alles Bedienbare muss mit Tabulator erreichbar sein, in einer sinnvollen Reihenfolge. Diese Reihenfolge kommt aus dem HTML, nicht aus dem CSS. Wer per Grid Dinge optisch umsortiert, zerreisst sie.
- `outline:none` ohne Ersatz ist der schaedlichste Einzeiler in CSS. Wenn dir der Standardring nicht gefaellt, gestalte ihn ueber `:focus-visible`, aber entferne ihn nie.
- Ein **Sprunglink** ganz oben ('Zum Inhalt springen') erspart taeglichen Nutzern hunderte Tastendruecke. Fuenf Zeilen Arbeit.
- Bei Dialogen muss der Fokus hinein wandern, drinnen bleiben und beim Schliessen zum ausloesenden Knopf zurueck. Das ist der haeufigste Fehler bei selbstgebauten Overlays.

## Aufgabe

> [!todo] Selber machen
> Bedien eine deiner Seiten komplett mit Tab, Enter und Escape. Ergaenze einen Sprunglink und einen sichtbaren Fokusring.

## Selbstpruefung

> [!question]
> Was ist an `outline:none` ohne Ersatz so schlimm?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: keyboard accessibility focus visible skip link tutorial](https://www.youtube.com/results?search_query=keyboard+accessibility+focus+visible+skip+link+tutorial) · empfohlen: Kevin Powell / Deque
- [MDN: :focus-visible](https://developer.mozilla.org/de/docs/Web/CSS/:focus-visible)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF07-K1 Warum, und fuer wen]] · [[LF07 Barrierefreiheit|Lernfeld]] · [[00 Start]] · [[LF07-K3 Sehen, Lesen, Bewegung]] →
