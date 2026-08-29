<#
    Baut die Veroeffentlichungs-Fassung und gibt die Pruefsumme aus.

    WARUM ES DIESES SKRIPT GIBT, angelegt am 29.08.2026 nach einem
    Fehler bei der Fassung 1.0.8:

    Zwei Dinge sind an dem Abend schiefgegangen, und beide haette dieses
    Skript verhindert.

    1. FALSCHER DATEINAME. Die exe wurde als "BenchMeisterPCCheck.exe"
       gebaut und so ins Release gehaengt. Acht Fassungen davor hiess
       sie "BenchMeister-PC-Check.exe". Der Download-Knopf auf
       benchmeister.de zeigt fest auf
       releases/latest/download/BenchMeister-PC-Check.exe und lief
       damit ins Leere. Dasselbe gilt fuer die Links in den
       Forenbeitraegen und fuer das winget-Manifest.

    2. SCHMUTZIGER ARBEITSBAUM. Die Datei wurde gebaut, BEVOR die
       Aenderungen committet waren. Go stempelt den Git-Stand ins
       Binary (vcs.revision, vcs.modified). Die hochgeladene Datei trug
       deshalb den vorherigen Commit und den Vermerk "geaendert", und
       ihre Pruefsumme liess sich von niemandem nachbauen.

       Das ist der schwerere Fehler. Auf der Seite und im Release steht
       "selbst bauen ergibt dieselbe Pruefsumme". Genau das ersetzt bei
       uns die fehlende Signatur. Stimmt es nicht, muss ein
       misstrauischer Nutzer von einer manipulierten Datei ausgehen,
       und misstrauische Nutzer sind unsere Zielgruppe.

    Pruefen laesst sich der Stempel in einer fertigen exe jederzeit mit:
        go version -m BenchMeister-PC-Check.exe

    NACHTRAG VOM SELBEN ABEND, wichtiger als die zwei Punkte oben:

    Dieses Skript baut im Arbeitsordner, und der ist fuer ein Release
    die falsche Grundlage. Nachgemessen am 29.08.2026 an Fassung 1.0.8:

        Arbeitsordner, sauber auf b2fd0d7   -> 3404A1B4...
        frischer Klon des Tags v1.0.8       -> 451C0F45...
        zweiter frischer Klon desselben Tags -> 451C0F45...

    Gleicher Commit, gleicher eingebetteter Stempel (vcs.revision und
    vcs.modified stimmen ueberein), trotzdem zwei verschiedene Dateien.
    Woran die lokale Abweichung liegt, ist ungeklaert.

    Entscheidend ist aber nicht, was hier herauskommt, sondern was ein
    Nachpruefer bekommt. Und der klont und baut, wie die beiden Klone
    oben. DIE ZAHL AUS EINEM FRISCHEN KLON IST DIE RICHTIGE.

    Fuer ein Release deshalb IMMER release-bauen.ps1 nehmen, das genau
    das tut. Dieses Skript hier ist fuer den Alltag.
#>

$ErrorActionPreference = "Stop"
Set-Location -Path $PSScriptRoot

# Der Name ist nicht verhandelbar, siehe Punkt 1 oben.
$Ziel = "BenchMeister-PC-Check.exe"

# --- Arbeitsbaum sauber? -------------------------------------------
$offen = git status --porcelain
if ($offen) {
    Write-Host ""
    Write-Host "ABBRUCH: Es gibt nicht committete Aenderungen." -ForegroundColor Red
    Write-Host ""
    Write-Host $offen
    Write-Host ""
    Write-Host "Go schreibt den Git-Stand in die exe. Aus einem geaenderten"
    Write-Host "Arbeitsbaum entsteht eine Datei, deren Pruefsumme niemand"
    Write-Host "nachbauen kann. Erst committen, dann bauen."
    exit 1
}

$commit = (git rev-parse --short HEAD).Trim()
$version = (Select-String -Path "main.go" -Pattern 'const version = "([^"]+)"').Matches[0].Groups[1].Value

Write-Host "Baue Fassung $version aus Commit $commit ..."

# Dieselben Schalter wie in der Anleitung im README und im Release-Text.
# Wer sie hier aendert, muss sie dort mitaendern, sonst stimmt die
# Pruefsumme fuer niemanden mehr.
go build -trimpath -ldflags="-s -w" -o $Ziel .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

$summe = (Get-FileHash $Ziel -Algorithm SHA256).Hash
$groesse = [math]::Round((Get-Item $Ziel).Length / 1MB, 2)

Write-Host ""
Write-Host "Fertig: $Ziel ($groesse MB)" -ForegroundColor Green
Write-Host "SHA256: $summe"
Write-Host ""
Write-Host "Diese Pruefsumme gehoert in den Release-Text. Zur Gegenprobe"
Write-Host "kann jeder das Repo klonen, dieses Skript laufen lassen und"
Write-Host "muss denselben Wert bekommen."
