# BenchMeister PC-Check

Kleines Windows-Programm, das ausliest, welche Hardware in einem Rechner
steckt, und das Ergebnis auf [benchmeister.de](https://benchmeister.de)
auswertet: Gesamtwertung, Vergleich mit kuratierten Komplettsystemen und
Hinweise, wo sich Aufrüsten lohnt.

Keine Installation. Eine Datei, Doppelklick, danach kann man sie löschen.

## Warum dieser Quelltext offen ist

Dieses Programm bittet darum, eine fremde ausführbare Datei zu starten,
und es überträgt (nach Zustimmung) Daten. Bei so etwas ist
"vertrau uns" kein Argument. Deshalb kann hier jede:r nachlesen, was
tatsächlich passiert, und das Programm selbst aus dem Quelltext bauen,
statt der fertigen Datei zu vertrauen.

Die BenchMeister-Website selbst liegt in einem eigenen, privaten
Repository. Offen ist genau das, was auf fremden Rechnern läuft.

## Was das Programm tut

1. Fragt Windows über WMI nach der verbauten Hardware. Dieselbe Auskunft,
   die auch der Geräte-Manager anzeigt.
2. Zeigt das Ergebnis lokal im Standardbrowser an, inklusive der Rohdaten,
   die übertragen würden. **Bis hierhin gibt es keinerlei Netzwerkverkehr.**
3. Fragt zwei getrennte Einwilligungen ab, beide standardmäßig leer.
4. Überträgt die Daten nur nach ausdrücklicher Zustimmung und öffnet die
   Ergebnisseite.

## Was das Programm ausdrücklich nicht tut

- nichts installieren, keine Registry-Einträge, kein Autostart, kein
  Dienst im Hintergrund
- keine Administratorrechte anfordern
- keine Belastungstests fahren (kein Hochtakten, keine Volllast)
- **vor der Zustimmung keine einzige Netzwerkverbindung aufbauen**

Der letzte Punkt lässt sich mit einer Firewall nachprüfen, und dazu wird
ausdrücklich eingeladen.

## Welche Daten erfasst werden

| Erfasst | Nicht erfasst |
|---|---|
| Prozessor-Modell, Kerne, Threads | Seriennummern (Festplatte, Mainboard) |
| Grafikkarten-Modell, Speicher | MAC-Adresse, BIOS-Kennung |
| Arbeitsspeicher: Größe, Takt | Windows-Benutzername, Rechnername |
| Laufwerk: SSD/HDD, Größe | Installierte Programme, Dateien |
| Bildschirmauflösung | IP-Adresse |

Die verbindliche Liste steht in
[`internal/scan/types.go`](internal/scan/types.go). Was übertragen wird,
steht in [`internal/upload/client.go`](internal/upload/client.go), der
einzigen Datei im Projekt, die überhaupt eine Netzwerkverbindung nach
außen aufbaut.

Das Mainboard wird ausgelesen und lokal angezeigt, aber **nicht**
übertragen.

Aus einem übertragenen Ergebnis lässt sich weder eine Person noch ein
bestimmter Rechner wiedererkennen.

## Die zwei Einwilligungen

Bewusst getrennt, beide standardmäßig leer:

- **(a) Ergebnis ansehen.** Nötig, damit die Daten überhaupt hochgeladen
  werden und die Auswertung angezeigt werden kann.
- **(b) Für Marktforschung freigeben.** Freiwillig und unabhängig davon.
  Erlaubt, die anonymen Angaben in eine Gesamtstatistik einfließen zu
  lassen, die an Hersteller und Händler weitergegeben werden kann. Damit
  finanziert sich BenchMeister mit.

**Wer nur seinen eigenen Score sehen will, muss (b) nie zustimmen.**

Ohne (a) wird gar nichts gesendet, und das ist an drei Stellen
abgesichert: in der Oberfläche, im lokalen Server
([`internal/ui/server.go`](internal/ui/server.go)) und in der Datenbank
selbst, die einen Upload ohne diese Einwilligung ablehnt.

Details siehe [Datenschutzerklärung](https://benchmeister.de/datenschutz),
Abschnitt 3d.

## Selbst bauen

Voraussetzung: [Go](https://go.dev/dl/) ab Version 1.21.

```
go build -ldflags="-s -w" -o BenchMeister-PC-Check.exe .
```

Ergibt eine eigenständige Datei von etwa 9 MB, ohne weitere
Abhängigkeiten.

## Hinweis zur Windows-Warnung

Die veröffentlichte Datei ist nicht signiert, deshalb zeigt Windows beim
ersten Start eine Warnung mit "Unbekannter Herausgeber". Grund: Ein
Signaturzertifikat kostet mehrere hundert Euro pro Jahr, und BenchMeister
hat bisher keinen Umsatz. Das wird nachgeholt, sobald es sich trägt.

Wem das zu unsicher ist: Der Quelltext liegt hier, das Programm lässt sich
mit dem Befehl oben in wenigen Sekunden selbst bauen.

## Aufbau

```
main.go                    Ablauf: auslesen, anzeigen, ggf. senden
internal/scan/types.go     Abschließende Liste der erfassten Felder
internal/scan/scan.go      WMI-Abfragen
internal/ui/server.go      Lokale Anzeige (nur 127.0.0.1) und Zustimmung
internal/ui/assets/        Die HTML-Seite, die im Browser erscheint
internal/upload/client.go  Der einzige Netzwerkaufruf nach außen
```

## Lizenz

MIT
