// Archivo: tools/drone-observe/cmd/llm.go
// Rol: comando llm para visualizar estado ML en vivo.
// No hace: inferencia ni control de drones.
package cmd

import (
	"drone-observe/internal/config"
	"drone-observe/internal/ui"
)

func runLLM(cfg config.Config) int {
	if err := ui.RunLLM(cfg); err != nil {
		return 1
	}
	return 0
}
