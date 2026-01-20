// Archivo: tools/drone-observe/internal/ui/topology.go
// Rol: TUI para topologia efectiva del sistema.
// No hace: discovery dinamico ni inferencias.
package ui

import (
	"fmt"
	"strings"

	"drone-observe/internal/config"
	"drone-observe/internal/topology"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
)

type topologyMsg struct {
	Items []topology.Component
}

type topologyModel struct {
	cfg     config.Config
	spinner spinner.Model
	items   []topology.Component
	done    bool
}

func RunTopology(cfg config.Config) error {
	m := newTopologyModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newTopologyModel(cfg config.Config) topologyModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return topologyModel{cfg: cfg, spinner: s}
}

func (m topologyModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, topologyCmd(m.cfg))
}

// PARTE CRITICA **********************
// La topologia solo refleja checks observables, no inferencias.
// Si se agregan dependencias ocultas, se pierde determinismo.
// No usar discovery ni suposiciones de infraestructura.
// FIN DE PARTE CRITICA ****************
func topologyCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		items := topology.Check(cfg)
		return topologyMsg{Items: items}
	}
}

func (m topologyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case topologyMsg:
		m.items = v.Items
		m.done = true
		return m, nil
	case tea.KeyMsg:
		if v.String() == "q" || v.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m topologyModel) View() string {
	title := TitleStyle.Render("drone-observe topology")
	sub := WarnStyle.Render("System Topology")

	var body strings.Builder
	body.WriteString(fmt.Sprintf("%s\n%s\n%s\n", title, sub, strings.Repeat("─", 44)))
	if len(m.items) == 0 && !m.done {
		body.WriteString("\n" + m.spinner.View())
		return BoxStyle.Render(body.String())
	}

	for i, it := range m.items {
		prefix := "├─"
		if i == len(m.items)-1 {
			prefix = "└─"
		}
		status := formatTopoStatus(it.Status)
		line := fmt.Sprintf("%s %s ............ %s", prefix, it.Name, status)
		if it.Detail != "" {
			line += fmt.Sprintf(" (%s)", it.Detail)
		}
		body.WriteString(line + "\n")
	}

	body.WriteString("\nPresiona 'q' para salir.\n")
	return BoxStyle.Render(body.String())
}

func formatTopoStatus(s topology.Status) string {
	switch s {
	case topology.StatusOK:
		return OKStyle.Render("OK")
	case topology.StatusSilent:
		return WarnStyle.Render("MUDO")
	default:
		return FailStyle.Render("FAIL")
	}
}
