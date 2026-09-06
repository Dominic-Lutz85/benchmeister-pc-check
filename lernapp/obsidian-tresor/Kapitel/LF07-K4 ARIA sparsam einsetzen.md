---
lernfeld: 7
kapitel: 4
dauer: 30
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-07
---

# ARIA sparsam einsetzen

**Lernfeld 07, Kapitel 4** · 30 Minuten · [[LF07 Barrierefreiheit|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du benutzt ARIA nur dort, wo HTML wirklich nicht ausreicht.

## Das musst du wissen

- Die erste Regel von ARIA lautet: **benutze kein ARIA**, wenn ein normales HTML-Element dasselbe kann. Ein `<button>` ist immer besser als ein `div role="button"`.
- Falsches ARIA ist schlimmer als gar keins. Es uebersteuert, was der Browser von allein richtig meldet, und macht funktionierende Dinge kaputt.
- Sinnvolle Faelle: `aria-label` fuer Knoepfe, die nur ein Symbol zeigen, `aria-expanded` an Auf- und Zuklappern, `aria-live` fuer Meldungen, die ohne Seitenwechsel erscheinen.
- Selbstgebaute Bedienelemente (Registerkarten, Menues, Dialoge) haben festgelegte Tastaturmuster. Die stehen im ARIA Authoring Practices Guide, mit fertigen Beispielen. Nicht raten, nachschlagen.

## Aufgabe

> [!todo] Selber machen
> Baue ein Auf- und Zuklapp-Element mit `aria-expanded` und pruefe es mit einer Vorlesefunktion (Windows: Sprachausgabe mit Strg+Windows+Eingabe).

## Selbstpruefung

> [!question]
> Warum ist falsches ARIA schlechter als gar keins?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: aria basics erklaert wann benutzen](https://www.youtube.com/results?search_query=aria+basics+erklaert+wann+benutzen) · empfohlen: Deque / Kevin Powell
- [ARIA Authoring Practices](https://www.w3.org/WAI/ARIA/apg/)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF07-K3 Sehen, Lesen, Bewegung]] · [[LF07 Barrierefreiheit|Lernfeld]] · [[00 Start]] · [[LF07-K5 Pruefen Werkzeuge und Handproben]] →
