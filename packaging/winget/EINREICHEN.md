# winget-Einreichung

Die drei YAML-Dateien im Unterordner 1.0.3 sind das fertige Manifest fuer
Version 1.0.3. Einreichen geht so:

1. Öffne https://github.com/microsoft/winget-pkgs und klicke oben rechts
   auf **Fork**.
2. Im eigenen Fork: **Add file → Upload files** im Ordner
   `manifests/d/DominicLutz/BenchMeisterPCCheck/1.0.3/`
   (den Pfad bei "Name your file" einfach vor den Dateinamen tippen,
   GitHub legt die Ordner an).
3. Alle drei YAML-Dateien aus dem Ordner 1.0.3 hochladen.
4. **Commit changes**, dann den vorgeschlagenen **Pull Request** eröffnen.
5. Warten: Ein automatischer Prüflauf validiert das Manifest, danach
   schaut ein Mensch von Microsoft drüber. Dauert meist wenige Tage.

Vorab lokal testen (optional, PowerShell):

    winget validate packaging/winget/1.0.3
    winget install --manifest packaging/winget/1.0.3

## Bei jeder neuen Version

1. Release mit versioniertem Tag veröffentlichen (v1.0.4 usw.).
2. In den drei Dateien `PackageVersion` erhöhen.
3. In der installer-Datei die `InstallerUrl` auf den neuen Tag zeigen
   lassen und `InstallerSha256` neu berechnen:

    certutil -hashfile BenchMeister-PC-Check.exe SHA256

4. Wieder als PR einreichen, gleicher Weg, neuer Versionsordner.
