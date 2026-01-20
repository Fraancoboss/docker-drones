// Archivo: tools/drone-observe/internal/ui/drift.go
// Rol: TUI para detectar deriva vs docs/dashboards.
// No hace: correcciones ni mutaciones.
package ui

import (
	"fmt"
	"strings"

	"drone-observe/internal/audit"
	"drone-observe/internal/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
)

type driftMsg struct {
	Findings []audit.Finding
	Err      error
}

type driftModel struct {
	cfg      config.Config
	spinner  spinner.Model
	findings []audit.Finding
	done     bool
}

func RunDrift(cfg config.Config) error {
	m := newDriftModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newDriftModel(cfg config.Config) driftModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return driftModel{cfg: cfg, spinner: s}
}

func (m driftModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, driftCmd(m.cfg))
}

// PARTE CRITICA **********************
// Drift compara solo artefactos versionados con el estado observable.
// Si se agregan heuristicas, se pierde la trazabilidad del gobierno tecnico.
// No introducir reglas SOC en esta capa.
// FIN DE PARTE CRITICA ****************
func driftCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		findings, err := audit.Drift(cfg)
		return driftMsg{Findings: findings, Err: err}
	}
}

func (m driftModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case driftMsg:
		m.findings = v.Findings
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

func (m driftModel) View() string {
	title := TitleStyle.Render("drone-observe drift")
	sub := WarnStyle.Render("Desviaciones tecnicas")

	var body strings.Builder
	body.WriteString(fmt.Sprintf("%s\n%s\n%s\n", title, sub, strings.Repeat("─", 44)))

	if !m.done {
		body.WriteString("\n" + m.spinner.View())
		return BoxStyle.Render(body.String())
	}

	if len(m.findings) == 0 {
		body.WriteString(OKStyle.Render("Sin drift detectado") + "\n")
		body.WriteString("\nPresiona 'q' para salir.\n")
		return BoxStyle.Render(body.String())
	}

	for _, f := range m.findings {
		line := fmt.Sprintf("%s %s - %s", formatSeverity(f.Severity), f.Item, f.Detail)
		body.WriteString(line + "\n")
	}

	body.WriteString("\nPresiona 'q' para salir.\n")
	return BoxStyle.Render(body.String())
}

func formatSeverity(s audit.Severity) string {
	switch s {
	case audit.SeverityHigh:
		return FailStyle.Render("ALTA")
	case audit.SeverityMed:
		return WarnStyle.Render("MEDIA")
	default:
		return OKStyle.Render("BAJA")
	}
}
