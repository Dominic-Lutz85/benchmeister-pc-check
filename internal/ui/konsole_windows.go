package ui

import "syscall"

// Windows-Konsolen zeigen standardmäßig eine alte Zeichentabelle an
// (meist Codepage 850). Go schreibt aber UTF-8, weshalb aus "geöffnet"
// ohne diese Umstellung "geÃ¶ffnet" wird.
//
// 65001 ist die Kennung für UTF-8. Die Umstellung gilt nur für dieses eine
// Konsolenfenster und endet mit dem Programm, es wird also nichts am
// System verändert, siehe die Zusage im Kopfkommentar von main.go.
func init() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleOutputCP := kernel32.NewProc("SetConsoleOutputCP")
	// Rückgabewert bewusst ignoriert: Klappt es nicht (z.B. weil die
	// Ausgabe in eine Datei umgeleitet wird), sind lediglich die Umlaute
	// unschön. Das ist kein Grund, das Programm abbrechen zu lassen.
	_, _, _ = setConsoleOutputCP.Call(65001)
}
