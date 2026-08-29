# winget-Einreichung

Die drei YAML-Dateien im Unterordner 1.0.6 sind das fertige Manifest fuer
Version 1.0.6. Einreichen geht so:

1. Öffne https://github.com/microsoft/winget-pkgs und klicke oben rechts
   auf **Fork**.
2. Im eigenen Fork: **Add file → Upload files** im Ordner
   `manifests/d/DominicLutz/BenchMeisterPCCheck/1.0.6/`
   (den Pfad bei "Name your file" einfach vor den Dateinamen tippen,
   GitHub legt die Ordner an).
3. Alle drei YAML-Dateien aus dem Ordner 1.0.6 hochladen.
4. **Commit changes**, dann den vorgeschlagenen **Pull Request** eröffnen.
5. Warten: Ein automatischer Prüflauf validiert das Manifest, danach
   schaut ein Mensch von Microsoft drüber. Dauert meist wenige Tage.

Vorab lokal testen (optional, PowerShell):

    winget validate packaging/winget/1.0.6
    winget install --manifest packaging/winget/1.0.6

## Bei jeder neuen Version

1. Release mit versioniertem Tag veröffentlichen (v1.0.4 usw.).
2. In den drei Dateien `PackageVersion` erhöhen.
3. In der installer-Datei die `InstallerUrl` auf den neuen Tag zeigen
   lassen und `InstallerSha256` neu berechnen:

    certutil -hashfile BenchMeister-PC-Check.exe SHA256

4. Wieder als PR einreichen, gleicher Weg, neuer Versionsordner.

## WICHTIG: nie eine ueberholte Fassung durchlaufen lassen

Am 29.08.2026 stand der PR microsoft/winget-pkgs#425633 offen auf 1.0.4
und war merge-bereit, waehrend 1.0.6 laengst veroeffentlicht war. Die
1.0.4 enthaelt aber genau den Fehler, bei dem der Speichertakt falsch
gemeldet wird: Die Reparatur war drin, kam aber wegen einer vergessenen
Kopierstelle nie beim Nutzer an.

Eine bekannte fehlerhafte Fassung in einen Paketmanager zu schieben ist
schlimmer als ein paar Tage Verzoegerung. Solange ein PR offen ist und
eine neue Version erscheint, wird der PR auf die neue Version umgestellt
oder geschlossen. Nie durchlaufen lassen.

## Den offenen PR auf eine neue Version umstellen

Im eigenen Fork (Dominic-Lutz85/winget-pkgs, Zweig master), ueber die
Weboberflaeche:

1. Die drei Dateien im alten Versionsordner loeschen.
2. Die drei Dateien aus dem neuen Versionsordner hier hochladen, unter
   `manifests/d/DominicLutz/BenchMeisterPCCheck/<neue Version>/`.
3. Den PR-Titel auf die neue Version aendern.
4. Einen Kommentar posten, warum umgestellt wurde.

Der PR aktualisiert sich von selbst, er zeigt auf den Zweig, nicht auf
einen festen Stand. Die automatische Pruefung laeuft danach neu.
