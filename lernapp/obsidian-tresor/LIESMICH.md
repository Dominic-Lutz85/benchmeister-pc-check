# Werkbank Webdesign · Obsidian-Tresor

## Öffnen

1. Obsidian starten
2. *Tresor öffnen* → *Ordner als Tresor öffnen*
3. Diesen Ordner auswählen
4. Mit [[00 Start]] anfangen

Es sind reine Markdown-Dateien. Sie funktionieren auch ohne Obsidian in
jedem Editor, die Verweise in doppelten Klammern sind dann eben nur Text.

## Auf den Stick

Ordner kopieren, fertig. Obsidian kann einen Tresor direkt vom Stick
öffnen. Zwei Hinweise dazu:

- Der versteckte Ordner `.obsidian` enthält deine Einstellungen. Beim
  Kopieren muss er mit, sonst sind Ansicht und Plugins wieder Standard.
  In Windows dafür *Ausgeblendete Elemente* im Explorer einschalten.
- Ein Stick ist keine Sicherung. Er geht verloren oder kaputt. Halte
  denselben Ordner zusätzlich an einer zweiten Stelle vor.

## Aufbau

| Ordner | Inhalt |
| --- | --- |
| `Kapitel/` | 65 Lerneinheiten, je eine Datei |
| `Lernfelder/` | 12 Übersichten, verlinken ihre Kapitel |
| `Projekte/` | 6 Praxisprojekte mit Prüfliste |
| `Referenz/` | Merksätze, Abnahmeliste, Quellen, Sitzungsnotiz |
| `Vorlagen/` | Muster für Tageseintrag, Fallgeschichte, Briefing |

## Neu bauen

Wenn sich der Lehrplan in `index.html` ändert:

```bash
python3 tresor-bauen.py
```

Das Skript leert den Ordner und schreibt ihn neu. **Eigene Notizen in den
Kapiteldateien gehen dabei verloren.** Wer Notizen sammeln will, legt sie
in eigenen Dateien ab oder sichert den Ordner vorher.
