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

	RamTotalGb int `json:"ram_total_gb"`
	// Zahl der Speicher-Steckplaetze auf dem Board, 0 wenn unbekannt.
	// Ergaenzt am 30.08.2026: "zwei Riegel verbaut" sagt ohne diese Zahl
	// nichts ueber den Ausbau. Zwei von zwei ist voll, zwei von vier
	// ist halb.
	// NICHT UEBERTRAGEN, wie die anderen rein oertlichen Angaben. Der
	// Upload zaehlt seine Felder in baueAnfrage() einzeln auf, ein neues
	// Feld landet also ohnehin nicht automatisch dort. Das `json:"-"`
	// macht die Absicht trotzdem sichtbar, damit niemand beim Lesen
	// raten muss.
	RamSteckplaetze int  `json:"-"`
	RamSpeedMhz     *int `json:"ram_speed_mhz"`

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

// IstTaktMhz liefert den Takt, mit dem dieser Riegel tatsaechlich laeuft.
//
// DIE REGEL, an einer Stelle, damit sie nicht auseinanderlaufen kann:
// ConfiguredClockSpeed gilt, wenn es gefuellt ist, sonst Speed.
//
//	Speed                = hoechster nativer Takt OHNE Profil.
//	                       Bei DDR5 also fast immer 4800.
//	ConfiguredClockSpeed = was das BIOS eingestellt hat (XMP, EXPO,
//	                       manuell). Das ist der Wert, der zaehlt.
//
// Ein Unterschied zwischen beiden ist kein Widerspruch, sondern der
// Normalfall bei aktivem Profil. Erklaert von NullPointerEx im
// PCGH-Forum am 29.08.2026.
//
// WARUM DIESE FUNKTION UEBERHAUPT EXISTIERT: Die Regel stand bis zum
// 29.08.2026 nur in internal/pruefung. Die Ist-Geschwindigkeit, die in
// der Hardware-Uebersicht ANGEZEIGT und an BenchMeister UEBERTRAGEN
// wird, entstand daneben aus einer zweiten, aelteren Rechnung, die nur
// Speed kannte. Auf einem Rechner mit aktivem Profil sagte deshalb der
// Befund 5600 und die Tabelle direkt darueber 4800, auf derselben Seite.
// Uebertragen wurde ebenfalls die 4800, also ausgerechnet in der
// Statistik, die spaeter etwas belegen soll.
//
// Verlaesslich ist auch dieser Wert nicht: Was in SMBIOS Typ 17 landet,
// entscheidet das BIOS des Boards. Deshalb behauptet die Oberflaeche
// nichts, sondern sagt "Windows meldet".
func (r RiegelInfo) IstTaktMhz() uint32 {
	if r.TaktMhzZweiter > 0 {
		return r.TaktMhzZweiter
	}
	return r.TaktMhz
}

// LaufwerkInfo ist ein Datentraeger, nur zur lokalen Anzeige.
type LaufwerkInfo struct {
	Name string
	// MSFT_PhysicalDisk.MediaType: 0 = unbekannt, 3 = HDD, 4 = SSD.
	// NICHT direkt auswerten, sondern ueber Medienart(), siehe dort.
	MedienArt uint16
	// MSFT_PhysicalDisk.BusType: 7 = USB, 11 = SATA, 17 = NVMe.
	BusArt uint16
	Bytes  uint64
}

// Medienart liefert die Art des Datentraegers und repariert dabei einen
// haeufigen Fall, in dem Windows sie gar nicht meldet.
//
// DAS PROBLEM: MSFT_PhysicalDisk.MediaType ist auf vielen Rechnern 0,
// also "unspecified". Das passiert vor allem bei NVMe-SSDs hinter
// bestimmten Treibern und Controllern und ist in Microsofts eigenen
// Foren seit Jahren dokumentiert. Ohne die Regel unten stand in der
// Uebersicht dann "Art unbekannt im M.2-Steckplatz (NVMe)", und in der
// Uebertragung landete "unbekannt" als Laufwerksart.
//
// DIE REGEL: Was per NVMe angebunden ist, IST eine SSD. NVMe ist ein
// Protokoll fuer Flash-Speicher, eine Magnetplatte spricht es nicht.
// Das ist keine Schaetzung, sondern folgt aus dem Anschluss selbst.
//
// Umgekehrt wird NICHT geraten: Ein Datentraeger am SATA-Anschluss ohne
// MediaType bleibt unbekannt, dort haengen Platten wie SSDs.
//
// Ergaenzt am 29.08.2026 bei der Durchsicht vor 1.0.6.
func (l LaufwerkInfo) Medienart() uint16 {
	const (
		unbekannt = 0
		ssd       = 4
		nvme      = 17
	)
	if l.MedienArt == unbekannt && l.BusArt == nvme {
		return ssd
	}
	return l.MedienArt
}
