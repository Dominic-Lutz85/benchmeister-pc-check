<#
    Baut die exe fuer ein Release, aus einem frischen Klon des Tags.

    WARUM NICHT EINFACH bauen.ps1, angelegt am 29.08.2026:

    Der eigene Arbeitsordner ist fuer ein Release die falsche Grundlage.
    Nachgemessen an Fassung 1.0.8:

        Arbeitsordner, sauber auf b2fd0d7    -> 3404A1B4...
        frischer Klon des Tags v1.0.8        -> 451C0F45...
        zweiter frischer Klon desselben Tags -> 451C0F45...

    Gleicher Commit, gleicher eingebetteter Git-Stempel, trotzdem zwei
    verschiedene Dateien. Woran die lokale Abweichung liegt, ist bis
    heute ungeklaert.

    Das muss es auch nicht sein. Entscheidend ist nicht, was auf dem
    Entwicklungsrechner herauskommt, sondern was jemand bekommt, der die
    Pruefsumme nachrechnet. Und der klont und baut. Deshalb wird die
    Datei fuer ein Release genau so erzeugt: geklont und gebaut.

    Zwei unabhaengige Klone ergaben bytegleich dasselbe. Das Versprechen
    "bau es selbst und vergleiche" haelt also, solange die Datei im
    Release aus einem Klon stammt und nicht von hier.

    ABLAUF EINER VEROEFFENTLICHUNG:
      1. Aenderungen committen und pushen.
      2. Auf GitHub das Release mit dem Tag anlegen, noch ohne Datei.
      3. Dieses Skript mit dem Tag aufrufen:  .\release-bauen.ps1 v1.0.9
      4. Die erzeugte Datei ins Release haengen und die ausgegebene
         Pruefsumme in den Release-Text schreiben.

    Schritt 2 laesst sich auch vorziehen, indem der Tag von Hand
    gesetzt und gepusht wird. Dann aber LEICHTGEWICHTIG, also
    `git tag v1.0.9 <commit>` OHNE -a. Ein annotierter Tag zeigt auf ein
    Tag-Objekt statt auf den Commit, und `git clone --branch --depth 1`
    scheitert daran mit "is not a commit". GitHub legt beim Anlegen
    eines Releases ebenfalls einen leichtgewichtigen an, alle bisherigen
    Fassungen sind so getaggt. Am 30.08.2026 einmal falsch gemacht und
    den Tag nachtraeglich ersetzen muessen.
#>

param(
    [Parameter(Mandatory = $true)]
    [string]$Tag
)

$ErrorActionPreference = "Stop"

$Repo = "https://github.com/Dominic-Lutz85/benchmeister-pc-check.git"
# Der Name ist nicht verhandelbar: Der Download-Knopf auf
# benchmeister.de zeigt fest auf
# releases/latest/download/BenchMeister-PC-Check.exe, ebenso die Links
# in den Forenbeitraegen und das winget-Manifest. Ein anderer Name macht
# den Download auf der Seite kaputt, ohne dass es jemand merkt.
$Ziel = "BenchMeister-PC-Check.exe"

$arbeit = Join-Path $env:TEMP "bmpc-release-$Tag"
if (Test-Path $arbeit) { Remove-Item $arbeit -Recurse -Force }

Write-Host "Klone $Tag frisch nach $arbeit ..."
# -c advice.detachedHead=false: Ohne das schreibt git seinen Hinweis
# zum losgeloesten HEAD auf die Fehlerausgabe, und PowerShell wertet
# bei ErrorActionPreference=Stop jede Zeile davon als Fehler. Das
# Skript brach dadurch ab, obwohl der Klon einwandfrei durchlief.
git -c advice.detachedHead=false clone --quiet --branch $Tag --depth 1 $Repo $arbeit
if ($LASTEXITCODE -ne 0) {
    Write-Host "ABBRUCH: Tag $Tag laesst sich nicht klonen. Ist er gepusht?" -ForegroundColor Red
    exit 1
}

Push-Location $arbeit
try {
    # Dieselben Schalter wie im README und im Release-Text. Wer sie hier
    # aendert, muss sie dort mitaendern, sonst stimmt die Pruefsumme fuer
    # niemanden mehr.
    # CGO_ENABLED=0 ausdruecklich, nicht dem Zufall ueberlassen.
    # Go schaltet CGO von selbst ein, sobald ein C-Compiler im Pfad
    # liegt, und baut dann eine andere Datei. Auf diesem Rechner ist
    # keiner installiert, auf einem anderen vielleicht schon, und der
    # GitHub-Runner hat einen. Ohne diese Zeile haengt die Pruefsumme
    # davon ab, wer gerade baut.
    $env:CGO_ENABLED = "0"
    go build -trimpath -ldflags="-s -w" -o $Ziel .
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

    $summe = (Get-FileHash $Ziel -Algorithm SHA256).Hash
    $groesse = [math]::Round((Get-Item $Ziel).Length / 1MB, 2)
    $version = (Select-String -Path "main.go" -Pattern 'const version = "([^"]+)"').Matches[0].Groups[1].Value

    # Gegenprobe: Traegt die Datei wirklich den Stand des Tags und keinen
    # Aenderungsvermerk? Ein schmutziger Baum kann im frischen Klon nicht
    # vorkommen, aber die Pruefung kostet nichts und deckt auf, wenn Go
    # den Stempel gar nicht setzen konnte.
    $stempel = (go version -m $Ziel) -join "`n"
    if ($stempel -notmatch "vcs\.revision=") {
        Write-Host "WARNUNG: Die Datei traegt keinen Git-Stempel." -ForegroundColor Yellow
        Write-Host "Die Pruefsumme ist dann nicht an einen Commit gebunden."
    }
    if ($stempel -match "vcs\.modified=true") {
        Write-Host "ABBRUCH: Die Datei traegt den Vermerk 'geaendert'." -ForegroundColor Red
        exit 1
    }

    $ablage = Join-Path ([Environment]::GetFolderPath("Desktop")) $Ziel
    Copy-Item $Ziel $ablage -Force

    Write-Host ""
    Write-Host "Fertig. Fassung $version aus Tag $Tag ($groesse MB)" -ForegroundColor Green
    Write-Host "Liegt auf dem Desktop: $ablage"
    Write-Host ""
    Write-Host "SHA256 fuer den Release-Text:" -ForegroundColor Cyan
    Write-Host $summe
}
finally {
    Pop-Location
}
