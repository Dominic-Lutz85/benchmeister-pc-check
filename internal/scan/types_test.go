package scan

import "testing"

/*
 * Bis zum 29.08.2026 hatte dieses Paket gar keine Tests, weil das Meiste
 * darin WMI abfragt und sich ohne echten Windows-Rechner nicht pruefen
 * laesst. Die drei Regeln hier sind aber reine Rechnerei und damit sehr
 * wohl pruefbar. Zwei von ihnen haben zuvor Fehler verursacht.
 */

func TestIstTaktMhz(t *testing.T) {
	faelle := []struct {
		name    string
		speed   uint32
		konfig  uint32
		erwartt uint32
	}{
		// Der Fall aus dem PCGH-Forum: Speed ist der Grundtakt, der
		// konfigurierte Wert zaehlt. Genau hier lag der Fehler, der in
		// 1.0.5 als behoben galt und es nicht war.
		{"Profil aktiv", 4800, 5600, 5600},
		{"kein Profil", 4800, 4800, 4800},
		{"Board fuellt das zweite Feld nicht", 2133, 0, 2133},
		{"Board fuellt nur das zweite Feld", 0, 3200, 3200},
		{"gar nichts gemeldet", 0, 0, 0},
		// Untertaktet: Der konfigurierte Wert gilt auch dann, wenn er
		// NIEDRIGER ist. Sonst wuerde ein bewusst gedrosselter Riegel
		// als schnell gemeldet.
		{"untertaktet", 3200, 2400, 2400},
	}

	for _, f := range faelle {
		t.Run(f.name, func(t *testing.T) {
			r := RiegelInfo{TaktMhz: f.speed, TaktMhzZweiter: f.konfig}
			if got := r.IstTaktMhz(); got != f.erwartt {
				t.Errorf("Speed %d, ConfiguredClockSpeed %d: %d erwartet, %d bekommen",
					f.speed, f.konfig, f.erwartt, got)
			}
		})
	}
}

func TestMedienart(t *testing.T) {
	const (
		unbekannt = 0
		hdd       = 3
		ssd       = 4
		usb       = 7
		sata      = 11
		nvme      = 17
	)

	faelle := []struct {
		name    string
		art     uint16
		bus     uint16
		erwartt uint16
	}{
		// Der eigentliche Anlass: Windows meldet bei NVMe-SSDs haeufig
		// gar keine Art. Der Anschluss beweist sie trotzdem.
		{"NVMe ohne Angabe ist eine SSD", unbekannt, nvme, ssd},
		{"NVMe mit Angabe bleibt SSD", ssd, nvme, ssd},
		// Am SATA-Anschluss haengen Platten wie SSDs. Hier wird bewusst
		// nicht geraten.
		{"SATA ohne Angabe bleibt unbekannt", unbekannt, sata, unbekannt},
		{"SATA-SSD bleibt SSD", ssd, sata, ssd},
		{"Festplatte bleibt Festplatte", hdd, sata, hdd},
		{"USB ohne Angabe bleibt unbekannt", unbekannt, usb, unbekannt},
	}

	for _, f := range faelle {
		t.Run(f.name, func(t *testing.T) {
			l := LaufwerkInfo{MedienArt: f.art, BusArt: f.bus}
			if got := l.Medienart(); got != f.erwartt {
				t.Errorf("MediaType %d am Bus %d: %d erwartet, %d bekommen",
					f.art, f.bus, f.erwartt, got)
			}
		})
	}
}

func TestAufGanzeGb(t *testing.T) {
	faelle := []struct {
		name    string
		bytes   uint64
		erwartt int
	}{
		// Der echte Registry-Wert einer RTX 5060 Ti mit 16 GB, auf dem
		// Entwicklungsrechner ausgelesen. Abschneiden haette 15 ergeben.
		{"16-GB-Karte, echter Wert", 17103323136, 16},
		// Typischer Wert einer 8-GB-Karte, ebenfalls knapp darunter.
		{"8-GB-Karte, knapp drunter", 8585740288, 8},
		{"glatte 4 GB", 4 * 1024 * 1024 * 1024, 4},
		{"glatte 2 GB", 2 * 1024 * 1024 * 1024, 2},
		{"nichts gemeldet", 0, 0},
	}

	for _, f := range faelle {
		t.Run(f.name, func(t *testing.T) {
			if got := aufGanzeGb(f.bytes); got != f.erwartt {
				t.Errorf("%d Bytes: %d GB erwartet, %d bekommen", f.bytes, f.erwartt, got)
			}
		})
	}
}
