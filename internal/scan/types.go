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

	// PCIe-Anbindung von Grafikkarte und NVMe-SSDs, ab 01.09.2026.
	// BLEIBT EBENFALLS LOKAL, `json:"-"` wie die beiden Listen darueber.
	Pcie []PcieGeraet `json:"-"`
}

/*
PcieGeraet ist ein PCIe-Geraet mit der Anbindung, die es GERADE hat.

NUR DER IST-WERT, kein gemeldeter Hoechstwert. Windows liefert auch
MaxLinkWidth und MaxLinkSpeed, und die Versuchung ist gross, beides
nebeneinanderzustellen und bei Abweichung zu warnen. Das waere ein
Fehlalarmautomat, aus zwei voneinander unabhaengigen Gruenden.

ERSTENS stimmt der gemeldete Hoechstwert nicht. Nachgemessen am
01.09.2026 auf dem Entwicklungsrechner: Eine GeForce RTX 5060 Ti meldet
MaxLinkWidth 16, obwohl diese Karte von Haus aus mit acht Leitungen
gebaut ist und die andere Haelfte der Kontakte gar nicht angeschlossen
hat. Eine Regel "aktuell kleiner als maximal, also Warnung" haette also
schon auf dem eigenen Rechner Alarm geschlagen, wo alles in Ordnung ist.

ZWEITENS sind auch die Ist-Werte im Leerlauf niedriger als unter Last.
Stromsparen senkt Generation und Breite, und das ist gewolltes
Verhalten. Werkzeuge wie GPU-Z haben deshalb einen Last-Test, der die
Verbindung erst hochzwingt, bevor sie ablesen. Dieses Programm erzeugt
ausdruecklich keine Last (siehe CLAUDE.md), es kann also nicht
unterscheiden zwischen "haengt dauerhaft an zu wenigen Leitungen" und
"doest gerade". Deshalb gibt es hierzu KEINEN Befund, sondern nur eine
Zeile in der Hardware-Uebersicht mit der Einordnung daneben.

Was die Angabe trotzdem wert ist: Wer weiss, was er verbaut hat, sieht
hier, in welchem Modus es laeuft. Genau danach war gefragt.
*/
type PcieGeraet struct {
	// Geraetename, wie Windows ihn meldet, z.B. "NVIDIA GeForce RTX 5060
	// Ti" oder "Standardmaessiger NVM Express-Controller".
	Name string
	// Zahl der ausgehandelten Leitungen, also 16 bei x16.
	Breite uint32
	// Ausgehandelte PCIe-Generation, also 4 bei PCIe 4.0.
	Generation uint32
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

// SystemzustandInfo sind drei Angaben zum laufenden System, die nichts
// ueber die verbaute Hardware sagen, sondern ueber ihren Zustand.
// Angelegt am 31.08.2026 nach Vorschlaegen von "cryon1c" im
// PCGH-Forum.
//
// ALLE DREI BLEIBEN AUF DEM RECHNER. Sie stehen bewusst nicht in
// ScanResult, sondern daneben: Wie voll eine Platte ist und welche
// Schutzprogramme laufen, geht niemanden etwas an. Das Feld hier hat
// deshalb gar keine json-Markierung, und baueAnfrage() zaehlt seine
// Felder ohnehin einzeln auf.
type SystemzustandInfo struct {
	// Datentraeger mit Windows darauf, nur die belegte und die
	// gesamte Groesse. Kein Laufwerksbuchstabe, kein Datentraegername.
	SystemplatteFreiGb   int
	SystemplatteGesamtGb int
	// Ob die Werte ueberhaupt ermittelt werden konnten. Ohne dieses
	// Feld waere "0 von 0 GB frei" nicht von "nicht auslesbar" zu
	// unterscheiden, und daraus wuerde ein Fehlalarm.
	SystemplatteErkannt bool

	// Namen der aktiven Echtzeit-Schutzprogramme. Windows Defender
	// zaehlt mit, er ist der Normalfall und allein kein Befund.
	Virenschutz []string
	// Ob die Abfrage lief. Auf Windows-Server-Ausgaben fehlt der
	// Namensraum SecurityCenter2 vollstaendig, dann ist eine leere
	// Liste kein Beweis fuer "kein Schutz".
	VirenschutzErkannt bool

	// Aushandlungsgeschwindigkeit der aktiven kabelgebundenen
	// Netzwerkverbindung in Mbit/s. KEINE MAC-Adresse, kein
	// Adaptername, keine IP.
	NetzwerkMbit int
	// Ob eine kabelgebundene Verbindung gefunden wurde. Wer nur per
	// WLAN online ist, bekommt keinen Befund: Dort ist eine niedrigere
	// Aushandlung normal und kein Zeichen fuer ein Problem.
	NetzwerkErkannt bool
}
