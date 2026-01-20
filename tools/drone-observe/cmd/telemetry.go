// Archivo: tools/drone-observe/cmd/telemetry.go
// Rol: comando telemetry para visualizar Data Plane en vivo.
// No hace: analitica avanzada ni alerting SOC.
package cmd

import (
	"drone-observe/internal/config"
	"drone-observe/internal/ui"
)

func runTelemetry(cfg config.Config) int {
	if err := ui.RunTelemetry(cfg); err != nil {
		return 1
	}
	return 0
}
