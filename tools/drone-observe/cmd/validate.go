// Archivo: tools/drone-observe/cmd/validate.go
// Rol: comando validate para auditar el contrato de METRICS.md.
// No hace: inferencias SOC ni cambios en metricas.
package cmd

import (
	"drone-observe/internal/config"
	"drone-observe/internal/ui"
)

func runValidate(cfg config.Config) int {
	if err := ui.RunValidate(cfg); err != nil {
		return 1
	}
	return 0
}
