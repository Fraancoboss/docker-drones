// Archivo: tools/drone-observe/internal/ui/style.go
// Rol: estilos TUI consistentes usando lipgloss.
// No hace: temas dinamicos ni deteccion avanzada de capacidades del terminal.
package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7DD3FC"))

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A3A3A3"))

	OKStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#22C55E"))

	FailStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EF4444"))

	WarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B"))

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#334155")).
			Padding(1, 2).
			Width(94)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E2E8F0"))
)
