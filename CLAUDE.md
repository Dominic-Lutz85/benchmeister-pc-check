# Arbeitsregeln für dieses Repository

Angelegt am 30.08.2026, nachdem "Jaffech" im PCGH-Forum angemerkt hat,
dass hier keine liegt. Er hatte recht: Die Regeln standen bis dahin
verstreut in Kommentaren, und Kommentare liest man erst, wenn man
ohnehin schon an der richtigen Stelle ist.

## Worum es bei diesem Programm geht

Ein portables Windows-Werkzeug, das ausliest, was in einem Rechner
steckt, und auf verschenkte Leistung hinweist. Kein Benchmark, keine
Last, keine Änderung am System. Es liest nur, was Windows ohnehin
meldet.

Die Zielgruppe sind Leute, die solchen Programmen misstrauen. Das ist
keine Randnotiz, sondern der Maßstab für jede Entscheidung hier.

## Die drei Regeln, die alles andere überstimmen

### 1. Lieber nichts sagen als etwas Falsches

Windows meldet Hardware-Daten aus einer Tabelle, die das Mainboard
füllt. Was da steht, entscheidet dessen Hersteller, und manche tragen
Standardwerte ein. Jede Auswertung muss deshalb zwei Fälle kennen: den
erkannten und den unklaren. Wo nichts sicher ablesbar ist, wird
geschwiegen.

Das ist teuer erkauft. Am 28.08.2026 meldete das Programm einem
PCGH-Moderator zu langsamen Speicher, obwohl seiner schneller lief als
vorgesehen. Ein Fehlalarm kostet nicht nur diesen einen Befund, er
kostet das Vertrauen in alle anderen mit.

Praktisch heißt das: Funktionen wie `kanalAus()` geben den Wert **und**
ein `ok` zurück. Wer das zusammenfasst, baut den nächsten Fehlalarm.

### 2. Befunde bleiben kurz, Begründungen kommen dahinter

Ein Befund hat drei Teile:

- `Feststellung`: die Messung, ein Satz mit den Zahlen
- `Empfehlung`: die Handlung, ein bis zwei Sätze
- `Hintergrund`: warum, wie man gegenprüft, welche Ausnahmen es gibt

Die ersten beiden zusammen dürfen **220 Zeichen** nicht überschreiten,
`TestBefundeBleibenKurz` bricht sonst ab. Der Hintergrund darf lang
sein, er steht in einem zugeklappten Block.

Anlass war "Hellhammer" im PCGH-Forum am 29.08.2026: zu viel begründet,
zu ausformuliert. Nachgemessen hatte er recht, der längste Befund hatte
590 Zeichen. Der Fehler war nicht die Ehrlichkeit, sondern dass alles
auf einer Ebene stand.

Faustregel beim Schreiben: Feststellung und Empfehlung müssen allein
ausreichen, um richtig zu handeln.

### 3. Ohne ausdrückliche Zustimmung verlässt nichts den Rechner

Vor der Zustimmung gibt es null Netzwerkverkehr, und was bei Zustimmung
übertragen würde, steht vorher wörtlich auf dem Bildschirm. Keine
Seriennummern, keine MAC-Adressen, kein Benutzername.

Wer hier etwas ergänzt, ergänzt es auch in der Vorschau. Eine
Übertragung, die im Text nicht auftaucht, bricht das zentrale
Versprechen des Programms.

## Aufbau

```
main.go                     Einstieg, Versionsnummer
internal/scan/              liest die Hardware aus (WMI, Registry)
internal/pruefung/          entscheidet, was dazu gesagt wird
internal/ui/                lokaler Webserver und die Ergebnisseite
internal/upload/            Übertragung, nur nach Zustimmung
```

Die Trennung zwischen `scan` und `pruefung` ist Absicht: `pruefung`
weiß nichts über WMI und lässt sich ohne Windows testen.

**Achtung, eine Regel steht doppelt:** Welches Feld den Ist-Takt
enthält, ist sowohl in `scan.RiegelInfo.IstTaktMhz` als auch in
`pruefung.go` beschrieben. Genau diese Doppelung hat am 29.08.2026 dazu
geführt, dass der Befund 5600 sagte und die Tabelle darüber 4800. Wer
eine der beiden Stellen ändert, ändert die andere mit.
`TestScanUndPruefungRechnenGleich` hält sie zusammen.

## Veröffentlichen

Nie von Hand. `release-bauen.ps1 v1.0.9` klont den Tag frisch und baut
daraus, weil ein Bau aus dem Arbeitsordner nachweislich eine andere
Datei ergibt als ein Klon. Der ganze Ablauf steht im Kopf des Skripts.

Die exe muss `BenchMeister-PC-Check.exe` heißen. Der Download-Knopf auf
benchmeister.de zeigt fest auf diesen Namen, ebenso die Links in den
Forenbeiträgen und das winget-Manifest.

Nach dem Veröffentlichen die Gegenprobe machen: Datei über den echten
Download-Weg holen und die Prüfsumme gegen die im Release-Text halten.
Dauert eine Minute und hat schon zwei Fehler gefunden.

## Sprache

Deutsch, auch in Kommentaren und Commit-Nachrichten. Kommentare in
ASCII (also "waere" statt "wäre"), sichtbarer Text mit echten Umlauten.
Keine Gedankenstriche als Satzzeichen.

Kommentare erklären das **Warum**, nicht das Was. Ein Kommentar, der
beschreibt was die Zeile darunter tut, ist überflüssig; einer, der
sagt warum sie so aussehen muss, verhindert, dass jemand sie in einem
halben Jahr "aufräumt".
