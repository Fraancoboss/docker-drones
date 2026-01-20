// Archivo: tools/drone-observe/cmd/drift.go
// Rol: comando drift para detectar desviaciones respecto a docs.
// No hace: normalizacion ni correccion automatica.
package cmd

import (
	"drone-observe/internal/config"
	"drone-observe/internal/ui"
)

func runDrift(cfg config.Config) int {
	if err := ui.RunDrift(cfg); err != nil {
		return 1
	}
	return 0
}
