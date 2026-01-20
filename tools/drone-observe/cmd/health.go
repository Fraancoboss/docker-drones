// Archivo: tools/drone-observe/cmd/health.go
// Rol: comando health para validar Control Plane de forma visual y determinista.
// No hace: correcciones automaticas ni cambios de configuracion.
package cmd

import (
	"drone-observe/internal/config"
	"drone-observe/internal/ui"
)

func runHealth(cfg config.Config) int {
	if err := ui.RunHealth(cfg); err != nil {
		return 1
	}
	return 0
}
