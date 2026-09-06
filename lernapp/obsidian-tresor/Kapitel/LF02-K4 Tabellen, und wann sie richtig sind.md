---
lernfeld: 2
kapitel: 4
dauer: 20
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-02
---

# Tabellen, und wann sie richtig sind

**Lernfeld 02, Kapitel 4** · 20 Minuten · [[LF02 HTML Struktur und Bedeutung|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du erkennst echte Tabellendaten und baust sie zugaenglich auf.

## Das musst du wissen

- Tabellen sind fuer **Daten in Zeilen und Spalten**, niemals fuer Layout. Das war in den 90ern ueblich und ist heute ein Fehler.
- `<th>` statt `<td>` fuer Kopfzellen, dazu `scope="col"` oder `scope="row"`. Erst damit weiss ein Vorleseprogramm, zu welcher Spalte ein Wert gehoert.
- Eine `<caption>` sagt in einem Satz, was die Tabelle zeigt. Kostet nichts und macht sie sofort verstaendlicher.
- Auf kleinen Bildschirmen sind breite Tabellen ein Problem. Loesung: die Tabelle in einen Behaelter mit `overflow-x:auto` setzen, damit nur sie scrollt und nicht die ganze Seite.

## Aufgabe

> [!todo] Selber machen
> Baue eine Preistabelle mit drei Spalten und Kopfzeile, richtig ausgezeichnet, und mach sie auf dem Handy seitlich scrollbar.

## Selbstpruefung

> [!question]
> Woran erkennst du, dass ein Inhalt keine Tabelle sein sollte?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: html tabellen barrierefrei th scope caption](https://www.youtube.com/results?search_query=html+tabellen+barrierefrei+th+scope+caption) · empfohlen: MDN / Kevin Powell
- [MDN: Tabellen](https://developer.mozilla.org/de/docs/Learn_web_development/Core/Structuring_content/HTML_table_basics)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF02-K3 Bilder, Medien und Formate]] · [[LF02 HTML Struktur und Bedeutung|Lernfeld]] · [[00 Start]] · [[LF02-K5 Formulare]] →
