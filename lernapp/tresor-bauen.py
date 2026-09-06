#!/usr/bin/env python3
# Baut aus dem Lehrplan in index.html einen Obsidian-Tresor.
# Die Daten werden aus der HTML-Datei gezogen und nicht doppelt gepflegt:
# eine zweite Kopie des Lehrplans waere die naechste Stelle, die
# auseinanderlaeuft.
import io, os, re, json, html, shutil, subprocess, sys

HIER = os.path.dirname(os.path.abspath(__file__))
QUELLE = os.path.join(HIER, "index.html")
ZIEL = os.path.join(HIER, "obsidian-tresor")

def daten():
    s = io.open(QUELLE, encoding="utf-8").read()
    block = re.findall(r"<script>(.*?)</script>", s, re.S)[0]
    js = block + "\nprocess.stdout.write(JSON.stringify({LEHRPLAN:LEHRPLAN,PROJEKTE:PROJEKTE}));"
    r = subprocess.run(["node", "-e", js], capture_output=True, text=True)
    if r.returncode:
        sys.exit("Lehrplan liess sich nicht auslesen:\n" + r.stderr)
    return json.loads(r.stdout)

def text(s):
    """HTML-Schnipsel aus dem Lehrplan in Markdown umsetzen."""
    s = re.sub(r"</?b>", "**", s)
    s = re.sub(r"<code>(.*?)</code>", lambda m: "`" + m.group(1) + "`", s, flags=re.S)
    s = re.sub(r"<[^>]+>", "", s)
    return html.unescape(s)

def dateiname(s):
    for zeichen in ':/\\?*"<>|':
        s = s.replace(zeichen, "")
    return re.sub(r"\s+", " ", s).strip()

def schreib(pfad, inhalt):
    os.makedirs(os.path.dirname(pfad), exist_ok=True)
    io.open(pfad, "w", encoding="utf-8", newline="\n").write(inhalt.rstrip() + "\n")

D = daten()
LEHRPLAN, PROJEKTE = D["LEHRPLAN"], D["PROJEKTE"]

if os.path.isdir(ZIEL):
    shutil.rmtree(ZIEL)

# ---------- Kapitel flach machen ----------
ALLE = []
for lf in LEHRPLAN:
    for i, k in enumerate(lf["k"]):
        ALLE.append(dict(k, lfNr=lf["nr"], lfName=lf["name"], kNr=i + 1))

def lfDatei(lf):   return dateiname("LF%02d %s" % (lf["nr"], lf["name"]))
def kapDatei(k):   return dateiname("LF%02d-K%d %s" % (k["lfNr"], k["kNr"], k["t"]))

# ---------- Kapitelnotizen ----------
for pos, k in enumerate(ALLE):
    lf = next(x for x in LEHRPLAN if x["nr"] == k["lfNr"])
    vor = ALLE[pos - 1] if pos > 0 else None
    nach = ALLE[pos + 1] if pos + 1 < len(ALLE) else None

    t = []
    t.append("---")
    t.append("lernfeld: %d" % k["lfNr"])
    t.append("kapitel: %d" % k["kNr"])
    t.append("dauer: %d" % k["min"])
    t.append("status: offen")
    t.append("erledigt_am:")
    t.append("tags:")
    t.append("  - webdesign/kapitel")
    t.append("  - webdesign/lernfeld-%02d" % k["lfNr"])
    t.append("---")
    t.append("")
    t.append("# %s" % text(k["t"]))
    t.append("")
    t.append("**Lernfeld %02d, Kapitel %d** · %d Minuten · [[%s|Lernfeld oeffnen]]"
             % (k["lfNr"], k["kNr"], k["min"], lfDatei(lf)))
    t.append("")
    t.append("- [ ] Kapitel gelesen und Aufgabe gemacht")
    t.append("")
    t.append("> [!abstract] Ziel")
    t.append("> " + text(k["ziel"]))
    t.append("")
    t.append("## Das musst du wissen")
    t.append("")
    for p in k["kern"]:
        t.append("- " + text(p))
    t.append("")
    t.append("## Aufgabe")
    t.append("")
    t.append("> [!todo] Selber machen")
    t.append("> " + text(k["praxis"]))
    t.append("")
    t.append("## Selbstpruefung")
    t.append("")
    t.append("> [!question]")
    t.append("> " + text(k["check"]))
    t.append("")
    t.append("Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.")
    t.append("")
    t.append("## Video und Quellen")
    t.append("")
    such = "https://www.youtube.com/results?search_query=" + k["v"][0].replace(" ", "+")
    t.append("- [Videosuche: %s](%s) · empfohlen: %s" % (k["v"][0], such, k["v"][1]))
    for q in k["q"]:
        t.append("- [%s](%s)" % (q[0], q[1]))
    t.append("")
    t.append("## Meine Notiz")
    t.append("")
    t.append("*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*")
    t.append("")
    t.append("")
    t.append("---")
    weg = []
    if vor:  weg.append("← [[%s]]" % kapDatei(vor))
    weg.append("[[%s|Lernfeld]]" % lfDatei(lf))
    weg.append("[[00 Start]]")
    if nach: weg.append("[[%s]] →" % kapDatei(nach))
    t.append(" · ".join(weg))
    schreib(os.path.join(ZIEL, "Kapitel", kapDatei(k) + ".md"), "\n".join(t))

# ---------- Lernfeld-Uebersichten ----------
for lf in LEHRPLAN:
    kap = [k for k in ALLE if k["lfNr"] == lf["nr"]]
    dauer = sum(k["min"] for k in kap)
    t = ["---", "tags:", "  - webdesign/lernfeld", "lernfeld: %d" % lf["nr"], "---", "",
         "# Lernfeld %02d: %s" % (lf["nr"], lf["name"]), "",
         "*%s*" % text(lf["kurz"]), "",
         "%d Kapitel, zusammen etwa %d Minuten." % (len(kap), dauer), "",
         "## Kapitel", ""]
    for k in kap:
        t.append("- [ ] [[%s]] · %d Min · %s" % (kapDatei(k), k["min"], text(k["ziel"])))
    frei = [p for p in PROJEKTE if p["nach"] == lf["nr"]]
    if frei:
        t += ["", "## Danach freigeschaltet", ""]
        for p in frei:
            t.append("- [[%s]]" % dateiname("%s %s" % (p["stufe"], p["name"])))
    t += ["", "---", "[[00 Start]] · [[01 Lernplan]]"]
    schreib(os.path.join(ZIEL, "Lernfelder", lfDatei(lf) + ".md"), "\n".join(t))

# ---------- Projekte ----------
for p in PROJEKTE:
    name = dateiname("%s %s" % (p["stufe"], p["name"]))
    t = ["---", "tags:", "  - webdesign/projekt", "frei_nach_lernfeld: %d" % p["nach"],
         "status: offen", "---", "",
         "# %s: %s" % (p["stufe"], text(p["name"])), "",
         "Frei nach Lernfeld %d · Aufwand: %s" % (p["nach"], text(p["dauer"])), "",
         "> [!abstract] Ziel", "> " + text(p["ziel"]), "",
         "## Das muss drin sein", ""]
    for x in p["punkte"]:
        t.append("- [ ] " + text(x))
    t += ["", "## Protokoll", "",
          "*Ausgangslage, Entscheidungen, was schiefging. Daraus wird spaeter die Fallgeschichte.*", "",
          "", "## Ergebnis", "",
          "- Link:", "- Screenshot:", "- Zahl, die sich verbessert hat:", "",
          "---", "[[00 Start]] · [[02 Projekte]] · Vorlage: [[Vorlage Fallgeschichte]]"]
    schreib(os.path.join(ZIEL, "Projekte", name + ".md"), "\n".join(t))

print("Kapitel, Lernfelder und Projekte geschrieben nach", ZIEL)

# ---------- Lernplan als eine Liste ----------
t = ["---", "tags:", "  - webdesign/uebersicht", "---", "",
     "# Lernplan", "",
     "Alle %d Kapitel in der Reihenfolge, in der sie aufeinander aufbauen." % len(ALLE),
     "Hak ab, was du **gebaut** hast, nicht was du gelesen hast.", ""]
for lf in LEHRPLAN:
    kap = [k for k in ALLE if k["lfNr"] == lf["nr"]]
    t.append("## Lernfeld %02d · [[%s|%s]]" % (lf["nr"], lfDatei(lf), text(lf["name"])))
    t.append("")
    for k in kap:
        t.append("- [ ] [[%s|%s]] · %d Min" % (kapDatei(k), text(k["t"]), k["min"]))
    frei = [p for p in PROJEKTE if p["nach"] == lf["nr"]]
    for p in frei:
        t.append("- [ ] **[[%s]]** bauen" % dateiname("%s %s" % (p["stufe"], p["name"])))
    t.append("")
t += ["---", "[[00 Start]]"]
schreib(os.path.join(ZIEL, "01 Lernplan.md"), "\n".join(t))

# ---------- Startseite ----------
gesamt = sum(k["min"] for k in ALLE)
t = ["---", "tags:", "  - webdesign/start", "---", "",
     "# Werkbank Webdesign", "",
     "Der komplette Weg vom ersten HTML-Tag bis zum abgerechneten Kundenprojekt.",
     "%d Lernfelder, %d Kapitel, etwa %d Stunden Lernzeit, dazu %d Praxisprojekte."
     % (len(LEHRPLAN), len(ALLE), round(gesamt / 60), len(PROJEKTE)), "",
     "## Taeglich", "",
     "1. [[01 Lernplan]] oeffnen und das naechste offene Kapitel nehmen.",
     "2. Lesen, **Aufgabe bauen**, Selbstpruefung laut beantworten.",
     "3. Notiz im Kapitel ergaenzen, Haken setzen, Eintrag ins [[Tagebuch]].", "",
     "Ein Kapitel am Tag reicht. Bei 30 Minuten taeglich bist du in etwa %d Tagen durch."
     % round(gesamt / 30), "",
     "## Einstiege", "",
     "- [[01 Lernplan]] · alle Kapitel als eine Liste zum Abhaken",
     "- [[02 Projekte]] · die sechs Sachen, die am Ende dein Portfolio sind",
     "- [[Tagebuch]] · was du an welchem Tag gemacht hast",
     "- [[Merksaetze]] · die Zahlen und Regeln, die immer wieder gebraucht werden",
     "- [[Abnahmeliste]] · vor jedem Livegang durchgehen",
     "- [[Alle Quellen]] · jede verlinkte Seite aus allen Kapiteln", "",
     "## Die zwoelf Lernfelder", ""]
for lf in LEHRPLAN:
    kap = [k for k in ALLE if k["lfNr"] == lf["nr"]]
    t.append("%d. [[%s|%s]] · %d Kapitel · *%s*"
             % (lf["nr"], lfDatei(lf), text(lf["name"]), len(kap), text(lf["kurz"])))
t += ["", "## Wie dieser Tresor entstanden ist", "",
      "Siehe [[Sitzungsnotiz]]. Dort steht auch, wo die App liegt und wie",
      "dieser Tresor neu gebaut wird, wenn sich der Lehrplan aendert.", "",
      "> [!tip] Wenn du Dataview installiert hast",
      "> Dann zeigt dir diese Abfrage den Stand:",
      "> ```dataview",
      "> TABLE dauer AS Minuten, status FROM #webdesign/kapitel WHERE status = \"offen\" SORT lernfeld, kapitel LIMIT 5",
      "> ```",
      "> Ohne das Plugin funktioniert alles andere trotzdem."]
schreib(os.path.join(ZIEL, "00 Start.md"), "\n".join(t))

# ---------- Projektuebersicht ----------
t = ["---", "tags:", "  - webdesign/uebersicht", "---", "",
     "# Projekte", "",
     "Wissen, das nie gebaut wurde, verschwindet. Diese sechs Projekte sind",
     "der eigentliche Ertrag des Kurses: am Ende hast du ein Portfolio,",
     "keine Urkunde.", ""]
for p in PROJEKTE:
    name = dateiname("%s %s" % (p["stufe"], p["name"]))
    t.append("### [[%s|%s: %s]]" % (name, p["stufe"], text(p["name"])))
    t.append("")
    t.append("Frei nach Lernfeld %d · %s" % (p["nach"], text(p["dauer"])))
    t.append("")
    t.append(text(p["ziel"]))
    t.append("")
t += ["---", "[[00 Start]]"]
schreib(os.path.join(ZIEL, "02 Projekte.md"), "\n".join(t))

# ---------- Alle Quellen ----------
t = ["---", "tags:", "  - webdesign/referenz", "---", "",
     "# Alle Quellen", "",
     "Jede in den Kapiteln verlinkte Seite, nach Lernfeld sortiert. Die",
     "Videolinks sind Suchen und keine festen Videos: einzelne Videos",
     "verschwinden, ein toter Link kostet mehr Vertrauen als er einbringt.", ""]
gesehen = set()
for lf in LEHRPLAN:
    t.append("## Lernfeld %02d: %s" % (lf["nr"], text(lf["name"])))
    t.append("")
    for k in [x for x in ALLE if x["lfNr"] == lf["nr"]]:
        for q in k["q"]:
            if q[1] in gesehen:
                continue
            gesehen.add(q[1])
            t.append("- [%s](%s)" % (q[0], q[1]))
    t.append("")
t += ["---", "[[00 Start]]"]
schreib(os.path.join(ZIEL, "Referenz", "Alle Quellen.md"), "\n".join(t))
print("Uebersichten und Quellenliste geschrieben.")

# ---------- Feste Notizen ----------
# Diese Texte stehen hier und nicht im Tresor, weil das Skript den Ordner
# bei jedem Lauf neu baut. Was nur im Tresor liegt, waere beim naechsten
# Bauen weg.
FEST = {}

FEST["LIESMICH.md"] = """# Werkbank Webdesign · Obsidian-Tresor

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
"""

FEST["Referenz/Merksaetze.md"] = """---
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
"""

FEST["Referenz/Abnahmeliste.md"] = """---
tags:
  - webdesign/referenz
---

# Abnahmeliste vor dem Livegang

Einmal schreiben, immer wieder benutzen. Verhindert genau die peinlichen
Fehler, die man am Tag der Übergabe macht.

## Inhalt

- [ ] Alle Links geprüft, keiner läuft ins Leere
- [ ] Rechtschreibung gelesen, am besten von jemand anderem
- [ ] Echte Texte, keine Platzhalter mehr im Quelltext

## Technik

- [ ] 404-Seite vorhanden und gestaltet
- [ ] Favicon gesetzt
- [ ] Vorschaubild und Beschreibung für soziale Netze
- [ ] HTTPS erzwungen
- [ ] `www` und Nicht-`www` auf eine Fassung umgeleitet
- [ ] Formular einmal echt abgesendet und Empfang geprüft

## Qualität

- [ ] Lighthouse gelaufen, Befunde abgearbeitet
- [ ] Nur mit Tastatur bedienbar, Fokus überall sichtbar
- [ ] Bei 200 Prozent Zoom benutzbar
- [ ] Auf einem echten Handy angesehen, nicht nur in der Simulation
- [ ] In zwei Browserfamilien geprüft

## Recht

- [ ] Impressum in höchstens zwei Klicks erreichbar
- [ ] Datenschutzerklärung beschreibt, was tatsächlich passiert
- [ ] Schriften selbst gehostet
- [ ] Externe Dienste erst nach Einwilligung geladen
- [ ] Bildlizenzen dokumentiert

## Übergabe

- [ ] Sicherung angelegt
- [ ] Zugangsdaten liegen beim Kunden
- [ ] Kurze Anleitung, wie er selbst etwas ändern kann

---
[[00 Start]]
"""

FEST["Tagebuch.md"] = """---
tags:
  - webdesign/tagebuch
---

# Tagebuch

Ein Eintrag pro Lerntag. Drei Zeilen reichen. Der Sinn ist nicht die
Dokumentation, sondern dass du in drei Monaten siehst, wie weit du
gekommen bist.

Muster steht in [[Vorlage Tageseintrag]].

## Einträge

<!-- Neueste oben -->
"""

FEST["Vorlagen/Vorlage Tageseintrag.md"] = """## JJJJ-MM-TT

**Kapitel:** [[ ]]
**Gebaut:** 
**Hat geklickt:** 
**Hängt noch:** 
"""

FEST["Vorlagen/Vorlage Fallgeschichte.md"] = """---
tags:
  - webdesign/portfolio
---

# Fallgeschichte: 

## Ausgangslage

*Wer, welches Geschäft, welche Seite gab es vorher. Zwei bis drei Sätze.*

## Aufgabe

*Was sollte besser werden, woran wird das gemessen. Ein Satz.*

## Mein Weg

*Die drei wichtigsten Entscheidungen und warum. Keine Aufzählung aller
Arbeitsschritte, sondern die Stellen, an denen es hätte anders laufen
können.*

1. 
2. 
3. 

## Ergebnis

*Mindestens eine Zahl. Ladezeit, Anfragen, Absprungrate, Kontrastwert,
irgendetwas Messbares.*

- Vorher: 
- Nachher: 

## Was ich daraus gelernt habe

*Ein ehrlicher Absatz. Der überzeugt Kunden mehr als jedes Lob.*
"""

FEST["Vorlagen/Vorlage Kundenbriefing.md"] = """---
tags:
  - webdesign/kunde
---

# Briefing: 

**Datum:** 
**Ansprechpartner:** 

## Die zehn Fragen

1. Was soll die Seite erreichen?
2. Wer soll darauf landen?
3. Was soll dieser Mensch danach getan haben?
4. Wie viele Seiten ungefähr?
5. Wer liefert die Texte, bis wann?
6. Wer liefert Bilder, gibt es Rechte daran?
7. Gibt es ein Logo, Farben, Vorgaben?
8. Welche Seiten gefallen dir und warum?
9. Wunschtermin?
10. Welcher Rahmen ist für das Budget vorgesehen?

## Wer pflegt die Seite in zwei Jahren?

*Die Frage entscheidet über die Wahl des Systems.*

## Nicht enthalten

*Was ausdrücklich nicht Teil des Auftrags ist. Dieser Punkt spart die
meisten Streitigkeiten.*

- 

## Nächster Schritt

- [ ] Angebot mit Umfang, Korrekturschleifen und Zahlungsplan
"""

for pfad, inhalt in FEST.items():
    schreib(os.path.join(ZIEL, pfad), inhalt)
print("Feste Notizen geschrieben:", len(FEST))

# ---------- Sitzungsnotiz und Obsidian-Grundeinstellungen ----------
schreib(os.path.join(ZIEL, "Referenz", "Sitzungsnotiz.md"), """---
tags:
  - webdesign/referenz
---

# Sitzungsnotiz: wie das hier entstanden ist

Angelegt am 06.09.2026. Damit in einem halben Jahr noch nachvollziehbar
ist, was wo liegt und warum es so aussieht.

## Was gebaut wurde

Eine Lern-App für Webdesign und dieser Tresor als Ablage dazu. Beides
speist sich aus **einer** Quelle: dem Lehrplan in `lernapp/index.html`.

| Was | Wo |
| --- | --- |
| Die App | `lernapp/index.html`, eine einzelne Datei ohne Abhängigkeiten |
| Dieser Tresor | `lernapp/obsidian-tresor/`, erzeugt aus der App |
| Das Bauskript | `lernapp/tresor-bauen.py` |
| Kurzbeschreibung | `lernapp/README.md` |
| Zweig | `claude/webdesign-learning-app-3j6ycg` |
| Veröffentlicht | https://claude.ai/code/artifact/e06882ac-1952-49ab-9d56-7195f13a4472 |

Der Tresor wird bei jedem Lauf des Skripts **neu gebaut**. Eigene Notizen
in den Kapiteldateien sind danach weg. Wer sammelt, sammelt in eigenen
Dateien.

## Entscheidungen und ihre Gründe

**Videos sind Suchlinks, keine festen Video-IDs.** Einzelne YouTube-Videos
werden gelöscht oder auf privat gestellt. Ein toter Link kostet mehr
Vertrauen, als ein Direktlink an Bequemlichkeit einbringt. Stattdessen
steht in jedem Kapitel ein vorbereiteter Suchbegriff plus der Kanal, dem
man trauen kann.

**Der Lehrplan steht nur in der App.** Eine zweite gepflegte Kopie im
Tresor wäre die nächste Stelle, die auseinanderläuft. Deshalb zieht das
Skript die Daten aus `index.html`.

**Die Reihenfolge der Lernfelder ist keine Geschmacksfrage.** Es wird nie
etwas vorausgesetzt, was noch nicht dran war. Lernfeld 5 (Gestaltung)
kommt nach 3 und 4, weil man Gestaltungsregeln erst anwenden kann, wenn
man Layout bauen kann.

**Fortschritt liegt im Browser**, zusätzlich in der Artefakt-Datenbank,
damit der Stand auf mehreren Geräten gleich ist. Ohne beides läuft die
App weiter, nur ohne Gedächtnis.

## Offener Punkt

Der sichtbare Text der App benutzt durchgehend umschriebene Umlaute
("Uebersicht", "oeffnen"). Die Regel in `CLAUDE.md` verlangt für
sichtbaren Text echte Umlaute; nur Kommentare sollen ASCII sein. Das
gehört bei Gelegenheit korrigiert, aber nicht mit einer pauschalen
Ersetzung: Wörter wie "neue" oder "Quelle" enthalten dieselben
Buchstabenfolgen und dürfen nicht angefasst werden. Es braucht eine
geprüfte Wortliste.

## Ideen, die noch nicht gebaut sind

- Karteikarten aus den 65 Selbstprüfungsfragen, Wiederholung nach
  Vergessenskurve
- Werkbuch: gelöste Probleme als eigene durchsuchbare Sammlung
- Aus den Projektnotizen automatisch eine Fallgeschichte erzeugen
- Code-Spielwiese direkt im Kapitel
- Wochenrückblick am Freitag

## Was nicht ging

Die Sitzung konnte sich nicht beim Artefakt-Dienst anmelden (HTTP 403).
Kommentare an der veröffentlichten App und Neuveröffentlichungen von
anderer Stelle kommen deshalb nicht automatisch an.

---
[[00 Start]]
""")

schreib(os.path.join(ZIEL, ".obsidian", "app.json"), """{
  "alwaysUpdateLinks": true,
  "attachmentFolderPath": "Anhaenge",
  "readableLineLength": true,
  "newLinkFormat": "shortest"
}""")

schreib(os.path.join(ZIEL, ".obsidian", "appearance.json"), """{
  "accentColor": "#0C6B5E"
}""")

print("Sitzungsnotiz und Obsidian-Einstellungen geschrieben.")
