// Archivo: tools/drone-observe/cmd/freshness.go
// Rol: comando freshness para evaluar recencia de datos.
// No hace: alarmado ni auto-remediacion.
package cmd

import (
	"drone-observe/internal/config"
	"drone-observe/internal/ui"
)

func runFreshness(cfg config.Config) int {
	if err := ui.RunFreshness(cfg); err != nil {
		return 1
	}
	return 0
}
