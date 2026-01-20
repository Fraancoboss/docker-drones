// Archivo: tools/drone-observe/cmd/topology.go
// Rol: comando topology para mostrar la topologia efectiva del sistema.
// No hace: descubrimiento dinamico ni inferencias magicas.
package cmd

import (
	"drone-observe/internal/config"
	"drone-observe/internal/ui"
)

func runTopology(cfg config.Config) int {
	if err := ui.RunTopology(cfg); err != nil {
		return 1
	}
	return 0
}
