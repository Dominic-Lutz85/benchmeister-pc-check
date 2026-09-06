---
lernfeld: 9
kapitel: 2
dauer: 40
status: offen
erledigt_am:
tags:
  - webdesign/kapitel
  - webdesign/lernfeld-09
---

# DOM und Ereignisse

**Lernfeld 09, Kapitel 2** · 40 Minuten · [[LF09 JavaScript fuer Gestalter|Lernfeld oeffnen]]

- [ ] Kapitel gelesen und Aufgabe gemacht

> [!abstract] Ziel
> Du reagierst auf Klicks und aenderst die Seite gezielt.

## Das musst du wissen

- `document.querySelector` holt ein Element ueber denselben Selektor wie in CSS. Mehr Auswahlmethoden braucht man selten.
- `element.addEventListener('click', fn)` ist der Standardweg. Klick-Attribute direkt im HTML vermischen Struktur und Verhalten und lassen sich schlecht wieder loesen.
- Aendere bevorzugt **Klassen**, nicht Einzelstile: `el.classList.toggle('offen')`. Das Aussehen bleibt in CSS, wo es hingehoert, und JavaScript schaltet nur um.
- Bei vielen gleichartigen Elementen einen Zuhoerer am gemeinsamen Elternelement anbringen und im Ereignis `event.target` auswerten. Spart hunderte Zuhoerer.

## Aufgabe

> [!todo] Selber machen
> Baue ein Aufklapp-Menue und eine Bilderliste, in der ein Klick auf ein kleines Bild das grosse austauscht. Nur ueber Klassen, kein direktes Style-Setzen.

## Selbstpruefung

> [!question]
> Warum aendert man lieber Klassen als einzelne Stileigenschaften?

Antworte laut, in eigenen Worten. Wo du stockst, liegt die Luecke.

## Video und Quellen

- [Videosuche: javascript dom manipulation events tutorial deutsch](https://www.youtube.com/results?search_query=javascript+dom+manipulation+events+tutorial+deutsch) · empfohlen: Web Dev Simplified / Programmieren lernen
- [MDN: DOM-Skripting](https://developer.mozilla.org/de/docs/Learn_web_development/Core/Scripting/DOM_scripting)

## Meine Notiz

*Was habe ich gebaut? Was hat nicht funktioniert? Wo will ich nochmal ran?*


---
← [[LF09-K1 Grundlagen der Sprache]] · [[LF09 JavaScript fuer Gestalter|Lernfeld]] · [[00 Start]] · [[LF09-K3 Formulare, Pruefung, Zustaende]] →
