---
tags:
  - webdesign/referenz
---

# Sitzungsnotiz: wie das hier entstanden ist

Angelegt am 06.09.2026. Damit in einem halben Jahr noch nachvollziehbar
ist, was wo liegt und warum es so aussieht.

## Was gebaut wurde

Eine Lern-App für Webdesign und dieser Tresor als Ablage dazu. Beides
speist sich aus **einer** Quelle: dem Lehrplan in `lernapp/index.html`.

| Was | Wo |
| --- | --- |
| Die App | `lernapp/index.html`, eine einzelne Datei ohne Abhängigkeiten |
| Dieser Tresor | `lernapp/obsidian-tresor/`, erzeugt aus der App |
| Das Bauskript | `lernapp/tresor-bauen.py` |
| Kurzbeschreibung | `lernapp/README.md` |
| Zweig | `claude/webdesign-learning-app-3j6ycg` |
| Veröffentlicht | https://claude.ai/code/artifact/e06882ac-1952-49ab-9d56-7195f13a4472 |

Der Tresor wird bei jedem Lauf des Skripts **neu gebaut**. Eigene Notizen
in den Kapiteldateien sind danach weg. Wer sammelt, sammelt in eigenen
Dateien.

## Entscheidungen und ihre Gründe

**Videos sind Suchlinks, keine festen Video-IDs.** Einzelne YouTube-Videos
werden gelöscht oder auf privat gestellt. Ein toter Link kostet mehr
Vertrauen, als ein Direktlink an Bequemlichkeit einbringt. Stattdessen
steht in jedem Kapitel ein vorbereiteter Suchbegriff plus der Kanal, dem
man trauen kann.

**Der Lehrplan steht nur in der App.** Eine zweite gepflegte Kopie im
Tresor wäre die nächste Stelle, die auseinanderläuft. Deshalb zieht das
Skript die Daten aus `index.html`.

**Die Reihenfolge der Lernfelder ist keine Geschmacksfrage.** Es wird nie
etwas vorausgesetzt, was noch nicht dran war. Lernfeld 5 (Gestaltung)
kommt nach 3 und 4, weil man Gestaltungsregeln erst anwenden kann, wenn
man Layout bauen kann.

**Fortschritt liegt im Browser**, zusätzlich in der Artefakt-Datenbank,
damit der Stand auf mehreren Geräten gleich ist. Ohne beides läuft die
App weiter, nur ohne Gedächtnis.

## Offener Punkt

Der sichtbare Text der App benutzt durchgehend umschriebene Umlaute
("Uebersicht", "oeffnen"). Die Regel in `CLAUDE.md` verlangt für
sichtbaren Text echte Umlaute; nur Kommentare sollen ASCII sein. Das
gehört bei Gelegenheit korrigiert, aber nicht mit einer pauschalen
Ersetzung: Wörter wie "neue" oder "Quelle" enthalten dieselben
Buchstabenfolgen und dürfen nicht angefasst werden. Es braucht eine
geprüfte Wortliste.

## Ideen, die noch nicht gebaut sind

- Karteikarten aus den 65 Selbstprüfungsfragen, Wiederholung nach
  Vergessenskurve
- Werkbuch: gelöste Probleme als eigene durchsuchbare Sammlung
- Aus den Projektnotizen automatisch eine Fallgeschichte erzeugen
- Code-Spielwiese direkt im Kapitel
- Wochenrückblick am Freitag

## Was nicht ging

Die Sitzung konnte sich nicht beim Artefakt-Dienst anmelden (HTTP 403).
Kommentare an der veröffentlichten App und Neuveröffentlichungen von
anderer Stelle kommen deshalb nicht automatisch an.

---
[[00 Start]]
