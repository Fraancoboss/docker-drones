// Archivo: tools/drone-observe/cmd/limits.go
// Rol: comando limits para exponer limites tecnicos observados.
// No hace: benchmarks ni pruebas de carga.
package cmd

import (
	"drone-observe/internal/config"
	"drone-observe/internal/ui"
)

func runLimits(cfg config.Config) int {
	if err := ui.RunLimits(cfg); err != nil {
		return 1
	}
	return 0
}
