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

## Gespeichert ist es erst, wenn die Datei da ist

Angelegt am 05.09.2026, nachdem ein langes Gespräch über einen
Budget-Rechner verloren schien. Es war darum gebeten worden, das
Ergebnis nach Obsidian zu schreiben. Dafür gab es gar keinen Weg: Ein
Vault liegt auf dem eigenen Rechner, eine Unterhaltung in der Claude-App
nicht, und ein Connector dorthin war nie eingerichtet.

Ob damals eine falsche Zusage fiel oder nur eine missverstandene, lässt
sich nicht mehr klären. Genau das ist der Punkt. Deshalb gilt ab jetzt:

**Eine Zusage, etwas sei gespeichert, zählt nicht. Nur die Datei zählt.**

Praktisch heißt das drei Dinge:

- Was festgehalten werden soll, wird als Datei angelegt, mit Pfad und
  Inhalt zurückgemeldet. "Ich habe es notiert" ist keine Ablage.
- Was eine Sitzung überleben soll, wird committet und gepusht.
  Sitzungen enden, Container werden weggeräumt, Git bleibt.
- Fehlt das Werkzeug für den gewünschten Ablageort, wird das gesagt,
  **bevor** gearbeitet wird. Lieber "das geht hier nicht" als eine
  Zusage, die niemand einlösen kann.

Das ist Regel 1 eine Ebene höher. Ein Fehlalarm kostet das Vertrauen in
alle anderen Befunde; eine erfundene Speicherung kostet das Vertrauen in
die gesamte Zusammenarbeit.

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

Nie von Hand, und seit dem 31.08.2026 auch nicht mehr auf diesem
Rechner. Diese Beschreibung stand bis zum 02.09.2026 falsch hier und
beschrieb noch den Ablauf davor.

**So geht es heute**, drei Schritte:

1. Versionsnummer in `main.go` hochsetzen (`const version`), committen,
   pushen. Ohne das meldet das Programm später die alte Nummer.
2. Auf GitHub ein Release mit dem neuen Tag anlegen und
   **veröffentlichen**.
3. Fertig. Alles Weitere macht `.github/workflows/release.yml`: Tests,
   Bau mit `CGO_ENABLED=0`, Prüfung des Git-Stempels, Beglaubigung der
   Herkunft, Anhängen der Datei ans Release und Eintragen der
   Prüfsumme in den Release-Text.

`release-bauen.ps1` ist damit NICHT mehr der Weg zum Release, sondern
das Werkzeug für die Gegenprobe: Es klont den Tag frisch und baut
daraus, und das Ergebnis muss bytegleich zu dem sein, was am Release
hängt. Genau das ist das Versprechen "bau es selbst und vergleiche".

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

## Vor jedem Push an .github/workflows/

Ein kaputter Workflow scheitert, **bevor** ein einziger Job entsteht.
In der Actions-Übersicht steht dann nur ein rotes Kreuz ohne Schritte,
und der Grund ist nirgends zu sehen. Am 31.08.2026 zweimal passiert.

Beide Sprachen deshalb vorher prüfen, das dauert zusammen zehn
Sekunden:

```bash
npx --yes js-yaml .github/workflows/release.yml > /dev/null && echo "YAML ok"
```

```powershell
$errs = $null
[System.Management.Automation.Language.Parser]::ParseFile($datei, [ref]$null, [ref]$errs)
```

**Die Falle, die beide Male zuschlug:** PowerShell-Here-Strings
(`@" ... "@`) und YAML-Blöcke vertragen sich nicht. Der Inhalt eines
Here-Strings darf nicht eingerückt sein, ein YAML-Block verlangt genau
das. Eine Zeile `---` im Text beendet dann das YAML-Dokument.

Mehrzeiligen Text stattdessen als Array einfach zitierter Zeilen bauen:
eingerückt, ohne Variablenersetzung, ohne Escape-Wirkung des Backticks.
