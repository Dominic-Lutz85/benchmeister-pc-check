package scan

import (
	"fmt"
	"strings"

	"github.com/yusufpapurcu/wmi"
)

// Die WMI-Strukturen unten enthalten bewusst NUR die Felder, die wir
// tatsächlich brauchen. WMI würde bei jeder dieser Klassen Dutzende
// weitere Felder liefern, darunter Seriennummern. Was hier nicht als Feld
// steht, wird gar nicht erst abgefragt, siehe die Abfragen mit expliziter
// Spaltenliste weiter unten.

type win32Processor struct {
	Name                      string
	NumberOfCores             uint32
	NumberOfLogicalProcessors uint32
}

type win32VideoController struct {
	Name                        string
	AdapterRAM                  uint32
	CurrentHorizontalResolution uint32
	CurrentVerticalResolution   uint32
}

// PartNumber, DeviceLocator und BankLabel sind am 27.08.2026 dazugekommen,
// fuer die Plausibilitaetspruefung (internal/pruefung). WICHTIG: Diese
// Felder werden NUR lokal angezeigt und ausgewertet, sie werden NICHT
// uebertragen. Die Upload-Struktur in upload/client.go kennt sie nicht.
//
// ConfiguredClockSpeed stand bis 28.08.2026 bewusst NICHT hier, mit der
// Begruendung: "Auf dem Testrechner meldete es denselben Wert wie Speed."
// Das war eine Stichprobe von genau einem Rechner, und sie hat getrogen.
//
// Ein Moderator im PCGH-Forum hat den Fehler gefunden: Sein Speicher lief
// nachweislich mit 5600 MT/s (per HWiNFO belegt), Windows meldete in
// Speed aber 4800, also den JEDEC-Grundtakt fuer DDR5. Das Programm hat
// daraufhin behauptet, sein Speicher liefe zu langsam, obwohl er sogar
// ueber dem Profil lief. Falschalarm im wichtigsten Befund ueberhaupt.
//
// Beide Felder werden deshalb jetzt gelesen und in der Pruefung
// gegeneinander gehalten. Widersprechen sie sich, gibt es keinen Befund,
// siehe pruefung.Speicher. Verlaesslich ist keines von beiden: Was in
// SMBIOS Typ 17 landet, entscheidet das BIOS des jeweiligen Boards, und
// manche tragen dort schlicht den SPD-Grundwert ein.
type win32PhysicalMemory struct {
	Capacity      uint64
	Speed         uint32
	PartNumber    string
	DeviceLocator string
	BankLabel     string
	// Zweiter Taktwert aus derselben SMBIOS-Tabelle. Kann 0 sein, wenn das
	// Board das Feld gar nicht fuellt.
	ConfiguredClockSpeed uint32
}

type win32DiskDrive struct {
	Size uint64
}

type msftPhysicalDisk struct {
	MediaType    uint16
	Size         uint64
	BusType      uint16
	FriendlyName string
}

type win32BaseBoard struct {
	Manufacturer string
	Product      string
}

// Alles auslesen. Einzelne Fehler sind nicht tödlich: Wenn z.B. der
// Arbeitsspeicher-Takt nicht ermittelbar ist, bleibt das Feld leer und der
// Rest wird trotzdem angezeigt. Ein hart abbrechendes Programm wäre hier
// die schlechtere Wahl, weil dann eine einzige eigenwillige
// Hardware-Konfiguration die ganze Auswertung verhindern würde.
func Auslesen() (*ScanResult, error) {
	ergebnis := &ScanResult{StorageType: "unbekannt"}

	if err := prozessor(ergebnis); err != nil {
		// Ohne Prozessor und Grafikkarte ist die Auswertung sinnlos,
		// deshalb sind das die einzigen beiden harten Fehler.
		return nil, fmt.Errorf("Prozessor konnte nicht ausgelesen werden: %w", err)
	}
	if err := grafikkarte(ergebnis); err != nil {
		return nil, fmt.Errorf("Grafikkarte konnte nicht ausgelesen werden: %w", err)
	}

	arbeitsspeicher(ergebnis)
	laufwerk(ergebnis)
	mainboard(ergebnis)

	return ergebnis, nil
}

func prozessor(e *ScanResult) error {
	var liste []win32Processor
	q := "select Name, NumberOfCores, NumberOfLogicalProcessors from Win32_Processor"
	if err := wmi.Query(q, &liste); err != nil {
		return err
	}
	if len(liste) == 0 {
		return fmt.Errorf("kein Prozessor gemeldet")
	}

	// Bei mehreren Sockeln nehmen wir den ersten. Doppel-Sockel-Systeme
	// sind bei Spiele-/Bastel-PCs praktisch nicht vorhanden, und die
	// Auswertung auf der Website rechnet ohnehin mit einem Prozessor.
	p := liste[0]
	e.CPUName = strings.TrimSpace(p.Name)
	e.CPUCores = int(p.NumberOfCores)
	e.CPUThreads = int(p.NumberOfLogicalProcessors)
	return nil
}

func grafikkarte(e *ScanResult) error {
	var liste []win32VideoController
	q := "select Name, AdapterRAM, CurrentHorizontalResolution, CurrentVerticalResolution from Win32_VideoController"
	if err := wmi.Query(q, &liste); err != nil {
		return err
	}
	if len(liste) == 0 {
		return fmt.Errorf("keine Grafikkarte gemeldet")
	}

	// Viele Rechner melden mehrere Grafikeinheiten: die im Prozessor
	// eingebaute UND die eingesteckte Karte. Für die Auswertung zählt die
	// tatsächlich benutzte, und das ist die, an der ein Bildschirm hängt
	// (die also eine Auflösung meldet). Nur falls keine eine Auflösung
	// meldet, nehmen wir die erste.
	gewaehlt := liste[0]
	for _, g := range liste {
		if g.CurrentHorizontalResolution > 0 && g.CurrentVerticalResolution > 0 {
			gewaehlt = g
			break
		}
	}

	e.GPUName = strings.TrimSpace(gewaehlt.Name)
	e.ResolutionWidth = int(gewaehlt.CurrentHorizontalResolution)
	e.ResolutionHeight = int(gewaehlt.CurrentVerticalResolution)

	// AdapterRAM ist ein 32-Bit-Feld und kann deshalb gar keine 4 GB
	// fassen (der größtmögliche Wert liegt eine Zahl darunter). Bei jeder
	// Karte mit 4 GB oder mehr meldet Windows deshalb einen gedeckelten
	// Wert, typischerweise exakt 4095 MB, auch wenn 12 oder 24 GB verbaut
	// sind. Ein bekannter, seit Jahren unbehobener Windows-Fehler.
	//
	// Der Compiler hat diese Klemme beim Bauen prompt selbst vorgeführt:
	// Eine Vergleichsgrenze von 4 GB ließ sich hier gar nicht erst
	// hinschreiben, weil sie überläuft. Deshalb wird stattdessen gegen den
	// gedeckelten Wert selbst geprüft: Alles ab 4095 MB gilt als
	// unglaubwürdig und wird verworfen. Die Website nimmt dann den Wert
	// aus dem eigenen Katalog. Lieber keine Angabe als eine falsche.
	const gedeckelt = uint32(4095) * 1024 * 1024
	if gewaehlt.AdapterRAM > 0 && gewaehlt.AdapterRAM < gedeckelt {
		gb := int(gewaehlt.AdapterRAM / (1024 * 1024 * 1024))
		if gb > 0 {
			e.GPUVramGb = &gb
		}
	}

	return nil
}

func arbeitsspeicher(e *ScanResult) {
	var liste []win32PhysicalMemory
	q := "select Capacity, Speed, ConfiguredClockSpeed, PartNumber, DeviceLocator, BankLabel from Win32_PhysicalMemory"
	if err := wmi.Query(q, &liste); err != nil {
		return
	}

	var summe uint64
	var takt uint32
	for _, m := range liste {
		summe += m.Capacity
		// Nur fuer die lokale Pruefung und Anzeige, wird nicht uebertragen.
		e.Riegel = append(e.Riegel, RiegelInfo{
			KapazitaetBytes: m.Capacity,
			TaktMhz:         m.Speed,
			TaktMhzZweiter:  m.ConfiguredClockSpeed,
			Teilenummer:     strings.TrimSpace(m.PartNumber),
			Kanal:           strings.TrimSpace(m.BankLabel),
		})
		// Bei gemischten Riegeln zählt der langsamste, denn genau mit dem
		// läuft das ganze System.
		if m.Speed > 0 && (takt == 0 || m.Speed < takt) {
			takt = m.Speed
		}
	}

	e.RamTotalGb = int(summe / (1024 * 1024 * 1024))
	if takt > 0 {
		t := int(takt)
		e.RamSpeedMhz = &t
	}
}

func laufwerk(e *ScanResult) {
	// Der naheliegende Weg (Win32_DiskDrive.MediaType) meldet auch bei
	// SSDs meist nur "Fixed hard disk media" und ist damit wertlos.
	// MSFT_PhysicalDisk im Storage-Namensraum unterscheidet dagegen
	// zuverlässig, und zwar ebenfalls ohne Administratorrechte.
	var physisch []msftPhysicalDisk
	q := "select MediaType, Size, BusType, FriendlyName from MSFT_PhysicalDisk"
	err := wmi.QueryNamespace(q, &physisch, `root\Microsoft\Windows\Storage`)

	if err == nil {
		for _, d := range physisch {
			e.Laufwerke = append(e.Laufwerke, LaufwerkInfo{
				Name:      strings.TrimSpace(d.FriendlyName),
				MedienArt: d.MediaType,
				BusArt:    d.BusType,
				Bytes:     d.Size,
			})
		}
	}

	if err == nil && len(physisch) > 0 {
		// Das größte Laufwerk ist in aller Regel das, auf dem Spiele und
		// Programme liegen, also das für die Einschätzung relevante.
		var groesstes msftPhysicalDisk
		for _, d := range physisch {
			if d.Size > groesstes.Size {
				groesstes = d
			}
		}

		switch groesstes.MediaType {
		case 3:
			e.StorageType = "HDD"
		case 4:
			e.StorageType = "SSD"
		default:
			e.StorageType = "unbekannt"
		}
		e.StorageCapacityGb = int(groesstes.Size / (1000 * 1000 * 1000))
		return
	}

	// Rückfallebene: Größe wenigstens grob ermitteln, Art bleibt unbekannt.
	// Ehrlicher, als aus der Größe eine Art zu raten.
	var einfach []win32DiskDrive
	if err := wmi.Query("select Size from Win32_DiskDrive", &einfach); err != nil {
		return
	}
	var groesste uint64
	for _, d := range einfach {
		if d.Size > groesste {
			groesste = d.Size
		}
	}
	e.StorageCapacityGb = int(groesste / (1000 * 1000 * 1000))
}

func mainboard(e *ScanResult) {
	var liste []win32BaseBoard
	// Ausdrücklich OHNE SerialNumber, obwohl WMI sie liefern würde.
	q := "select Manufacturer, Product from Win32_BaseBoard"
	if err := wmi.Query(q, &liste); err != nil || len(liste) == 0 {
		return
	}
	e.MainboardName = strings.TrimSpace(
		liste[0].Manufacturer + " " + liste[0].Product,
	)
}
