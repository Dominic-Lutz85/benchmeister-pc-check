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
