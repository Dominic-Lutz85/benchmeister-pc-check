# Frankenstein-Build: Stand

Persönliche Projektnotiz, **kein Teil des Programms**. Sie liegt hier,
weil eine Notiz, die nur in einer Unterhaltung steht, beim nächsten Mal
weg ist. Siehe die Regel "Gespeichert ist es erst, wenn die Datei da
ist" in `CLAUDE.md`.

Stand vom 03.09.2026 aus einer lokalen Sitzung, hier eingetragen am
05.09.2026.

## Was feststeht

**Grafikkarte: gebrauchte RTX 3050 8 GB, rund 150 Euro.**

Zwei Gründe, und beide sind der Kern der Auswahl:

- Sie braucht kein Resizable BAR, läuft also in jedem Unterbau. Das
  hält die Suche nach dem Rechner darunter offen, statt sie auf
  Plattformen mit ReBAR-BIOS einzuengen.
- Sie hat RT-Einheiten. Ohne die startet DOOM The Dark Ages gar nicht,
  siehe die dritte Sperre unten.

## Was gesucht wird

| Teil | Vorgabe | Preis |
|---|---|---|
| Komplettrechner | 6 Kerne, 16 GB, SSD, Ziel i5-10400F | bis 175 Euro |
| Netzteil | 550 W, neu, mit PCIe-Stecker | 53 Euro |
| Grafikkarte | steht fest, siehe oben | 150 Euro |
| **Summe** | | **rund 359 Euro** |

## Die drei Sperren

Daran scheitern billige Gebrauchtrechner, und zwar hart. Nicht
"langsam", sondern "startet nicht".

1. **SSE4.2** fehlt bei Phenom II und Athlon II.
2. **AVX2** fehlt bei AMD FX und bei Intel vor der 4. Generation.
3. **RT-Einheiten** fehlen bei allen GTX und allen RX 500.

Die ersten beiden sind der Grund, warum ein 60-Euro-Bürorechner aus dem
Kleinanzeigenmarkt keine Option ist, auch wenn die Kernzahl stimmt.

## Zwei Zahlen, die noch nicht nachgeprüft sind

Nach der ersten Regel dieses Projekts gehört dazugeschrieben, was
gesetzt, aber nicht belegt ist:

- **"rund 60 fps Medium in DOOM The Dark Ages"** auf einer RTX 3050
  ist optimistisch. Die Karte liegt nah an der Mindestanforderung des
  Spiels. Vor dem Kauf einen aktuellen Test genau dieser Karte in genau
  diesem Spiel ansehen, mit Auflösung und Upscaling-Einstellung dabei.
- **"ReBAR ab der 10. Generation offiziell"** stimmt für die CPU,
  hängt in der Praxis aber am BIOS des jeweiligen Mainboards. Bei einem
  Fertigrechner ist genau das die Stelle, an der es fehlen kann. Für
  diesen Build ist es kein Ausschlusskriterium, weil die 3050 ohne
  ReBAR auskommt.

## Wo die Details liegen

Beides liegt lokal auf dem eigenen Rechner und ist aus einer
Cloud-Sitzung nicht erreichbar:

- Obsidian-Vault, unterhalb von `01 Projekte/BenchMeister/Konzepte/`
- Teile-Katalog in der lokalen Skill `frankenstein-hardware`,
  Datei `TEILE-KATALOG.md`

Wer daran weiterarbeiten will, startet Claude Code **lokal** in dem
Ordner. Aus der Cloud geht es nicht, und eine Sitzung, die etwas
anderes behauptet, irrt sich.

## Bastel-Trends 2026 und was davon hier taugt

Eingeordnet am 05.09.2026 anhand einer Trend-Übersicht, gemessen an
diesem Budget:

- **Laptop-Mainboard-Repurposing (Cyberdeck).** Der einzige Trend, der
  Geld spart. Defekte Gaming-Notebooks sind billig. Haken: Speicher oft
  verlötet, CPU immer, also kein Aufrüstpfad. Für diesen Build nicht
  passend, weil die 3050 dann nicht hineinpasst.
- **Einbau-Displays (CYD mit ESP32).** Unter 15 Euro, reine Optik,
  unkritisch.
- **BTF-Stecker-Löten.** Ein umgelötetes Board ist bei einem Fehler
  Schrott, und gespart wird dabei nichts. Falscher Ort zum Üben.
- **Carbon-Nanotube-Pads und Delidding.** Zielen auf die letzten Grad
  bei High-End-CPUs. Bei einem i5-10400F gibt es die nicht zu holen.
- **Edge-AI-Lüftersteuerung.** Löst ein Problem, das die BIOS-Kurve
  schon löst.

## Ein Befund im Programm, der hier hängt

Ein Cyberdeck aus einer Notebook-Platine meldet über SMBIOS
"Notebook". `internal/scan/bauform.go` stuft es damit als mobil ein,
und die Speicher-Empfehlungen werden zurückgehalten, obwohl die
SODIMM-Steckplätze bei so einem Aufbau offen daliegen.

Das ist **kein Fehler**, sondern die gewollte Richtung: lieber ein
vorsichtiger Rat als ein unbefolgbarer. Notiert, damit es nicht in
einem halben Jahr als Bug missverstanden wird.

## Offen

- Konkreter Fertigrechner unter 175 Euro, der die drei Sperren nimmt.
  Preise dafür kommen von Geizhals und aus Kleinanzeigen, nicht aus
  einem Sprachmodell.
- Gegenprobe der beiden ungeprüften Zahlen oben.
