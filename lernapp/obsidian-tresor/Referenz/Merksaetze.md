---
tags:
  - webdesign/referenz
---

# Merksätze

Die Zahlen und Regeln, die in mehreren Lernfeldern wiederkommen. Zum
Nachschlagen, nicht zum Auswendiglernen.

## Text

- Zeilenlänge 45 bis 75 Zeichen (`max-width: 65ch`)
- Zeilenhöhe etwa 1,5 im Fließtext, 1,1 bis 1,25 bei großen Überschriften
- Fließtext mindestens 16px, in `rem` statt `px`
- Je größer die Schrift, desto weniger Zeilenabstand braucht sie
- Zwei Schriften reichen, eine dritte braucht eine eigene Aufgabe

## Farbe und Kontrast

- 4,5:1 für normalen Text, 3:1 für große Schrift und Bedienelemente
- Farbe nie als einziger Träger einer Information
- Palette nach Rollen benennen, nicht nach Aussehen
- Neutralgrau mit einem Hauch des Akzenttons statt reinem Grau

## Abstand und Layout

- Abstände aus einer festen Reihe: 4, 8, 12, 16, 24, 32, 48, 64
- Was zusammengehört, steht enger beieinander als zum Nachbarn
- Grid für das Geräst, Flexbox für alles darin
- `box-sizing: border-box` gehört in jedes Projekt
- Umbruchpunkte nach Inhalt setzen, nie nach Gerätenamen

## Tempo

- LCP unter 2,5 Sekunden, INP unter 200 ms, CLS unter 0,1
- Bilder sind meist 70 bis 80 Prozent des Seitengewichts
- Bilder immer mit `width` und `height`, sonst springt das Layout
- Auf dem Handy messen, nicht auf dem eigenen Rechner

## Barrierefreiheit

- `outline: none` ohne Ersatz ist der schädlichste Einzeiler in CSS
- Bei 200 Prozent Zoom muss die Seite noch benutzbar sein
- Erste ARIA-Regel: kein ARIA benutzen, wenn HTML dasselbe kann
- Automatische Prüfer finden nur 30 bis 40 Prozent der Probleme

## Kunden

- Fünf Testpersonen finden etwa 85 Prozent der Bedienprobleme
- Texte und Bilder sind der häufigste Verzögerungsgrund, schriftlich klären
- Anzahl der Korrekturschleifen gehört ins Angebot
- Vorschuss: ein Drittel bei Auftrag, bei Freigabe, bei Livegang

---
[[00 Start]]
