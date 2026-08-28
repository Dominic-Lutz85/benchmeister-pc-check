// Package scan liest die verbaute Hardware aus. Ausschließlich lesend,
// ohne Administratorrechte, ohne Belastungstests, ohne Eingriff ins System.
package scan

// ScanResult ist die vollständige und abschließende Liste dessen, was
// dieses Programm über den Rechner erfasst und (nach ausdrücklicher
// Zustimmung) überträgt. Was hier nicht steht, wird auch nicht erhoben.
//
// AUSDRÜCKLICH NICHT ERFASST, und zwar bewusst:
//
//	MAC-Adresse                 (Win32_NetworkAdapter.MACAddress)
//	Festplatten-Seriennummer    (Win32_DiskDrive.SerialNumber)
//	Mainboard-Seriennummer      (Win32_BaseBoard.SerialNumber)
//	BIOS-/System-UUID           (Win32_ComputerSystemProduct.UUID)
//	Windows-Benutzername        (Win32_ComputerSystem.UserName)
//	Rechnername                 (os.Hostname)
//	Installierte Software       (Win32_Product wird gar nicht erst abgefragt)
//	IP-Adresse
//
// Das sind alles Angaben, über die sich ein bestimmter Rechner oder eine
// Person wiedererkennen ließe. Genau deshalb fehlen sie. Aus einem
// übertragenen Ergebnis lässt sich ablesen, dass irgendwer irgendwo eine
// bestimmte Hardware-Kombination nutzt, mehr nicht.
//
// ZUR IP-ADRESSE, damit die Liste oben nicht mehr verspricht, als sie
// halten kann: Dieses Programm liest sie nicht aus und sendet sie nicht
// mit. Aber jede Verbindung ins Internet trägt sie zwangsläufig, das
// gilt für jede Website und auch hier. Seit 1.0.3 verrechnet der Server
// von benchmeister.de sie kurz im Arbeitsspeicher zu einer Prüfsumme,
// um zu begrenzen, wie viele Ergebnisse in kurzer Zeit von derselben
// Stelle ankommen. Ohne diese Bremse könnte jemand die Statistik mit
// erfundenen Rechnern fluten. Die Adresse wird dabei nicht gespeichert,
// nicht protokolliert und landet in keiner Datenbank, und nach zehn
// Minuten ist auch die Prüfsumme weg. Siehe Abschnitt 3g der
// Datenschutzerklärung auf benchmeister.de.
//
// Die Modellbezeichnungen werden absichtlich UNVERÄNDERT so übernommen,
// wie Windows sie meldet (z.B. "AMD Ryzen 7 7800X3D 8-Core Processor").
// Die Übersetzung auf den BenchMeister-Katalog passiert erst auf der
// Website. Grund: Diese Zuordnungslogik muss bei jedem neuen Bauteil
// gepflegt werden, und sie soll nur an einer einzigen Stelle existieren,
// nicht zusätzlich hier im Programm.
type ScanResult struct {
	CPUName    string `json:"cpu_name"`
	CPUCores   int    `json:"cpu_cores"`
	CPUThreads int    `json:"cpu_threads"`

	GPUName string `json:"gpu_name"`
	// Zeiger, damit "unbekannt" (nil) von "0 GB" unterscheidbar bleibt.
	// Windows meldet den Grafikspeicher oberhalb von 4 GB notorisch falsch,
	// siehe gpu.go. Lieber gar kein Wert als ein falscher: die Website
	// nimmt dann den Wert aus dem eigenen Katalog.
	GPUVramGb *int `json:"gpu_vram_gb"`

	RamTotalGb  int  `json:"ram_total_gb"`
	RamSpeedMhz *int `json:"ram_speed_mhz"`

	// "SSD", "HDD" oder "unbekannt".
	StorageType       string `json:"storage_type"`
	StorageCapacityGb int    `json:"storage_capacity_gb"`

	ResolutionWidth  int `json:"resolution_width"`
	ResolutionHeight int `json:"resolution_height"`

	// Nur zur Anzeige im Programm selbst, wird NICHT übertragen (kein
	// json-Feld in der Upload-Struktur, siehe upload/client.go). Hilft
	// beim Einordnen der eigenen Hardware, ist für die Auswertung aber
	// ohne Belang und gehört deshalb nicht in die Datenbank.
	MainboardName string `json:"-"`

	// Ab 27.08.2026 fuer die Plausibilitaetspruefung (internal/pruefung).
	//
	// BEIDE LISTEN WERDEN NICHT UEBERTRAGEN, erkennbar am `json:"-"` und
	// daran, dass die Upload-Struktur in upload/client.go sie nicht kennt.
	// Sie dienen ausschliesslich der lokalen Anzeige und Auswertung.
	//
	// Die Teilenummer eines Speicherriegels ist eine Modellbezeichnung
	// wie "CMK32GX4M2E3200C16", keine Seriennummer. Sie sagt, welches
	// Produkt verbaut ist, nicht welches Exemplar. Sie bleibt trotzdem
	// lokal, weil sie fuer die Auswertung auf der Website nicht gebraucht
	// wird und die Regel hier lautet: was nicht gebraucht wird, wird nicht
	// uebertragen.
	Riegel    []RiegelInfo   `json:"-"`
	Laufwerke []LaufwerkInfo `json:"-"`
}

// RiegelInfo ist ein einzelner Speicherbaustein, nur zur lokalen Anzeige.
type RiegelInfo struct {
	KapazitaetBytes uint64
	// Win32_PhysicalMemory.Speed.
	TaktMhz uint32
	// Win32_PhysicalMemory.ConfiguredClockSpeed. Zweiter Taktwert aus
	// derselben SMBIOS-Tabelle, seit 28.08.2026 dabei. Welches der beiden
	// Felder der Ist-Takt ist, haengt vom BIOS ab, siehe die Erklaerung an
	// win32PhysicalMemory in scan.go. Beide zusammen sind belastbarer als
	// eines allein, sicher ist keines.
	TaktMhzZweiter uint32
	Teilenummer    string
	Kanal          string
}

// LaufwerkInfo ist ein Datentraeger, nur zur lokalen Anzeige.
type LaufwerkInfo struct {
	Name      string
	MedienArt uint16
	BusArt    uint16
	Bytes     uint64
}
