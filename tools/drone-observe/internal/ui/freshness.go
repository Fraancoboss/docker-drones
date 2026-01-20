// Archivo: tools/drone-observe/internal/ui/freshness.go
// Rol: TUI para recencia de datos usando timestamps reales.
// No hace: alertas ni auto-remediacion.
package ui

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"drone-observe/internal/config"
	"drone-observe/internal/freshness"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
)

type freshnessMsg struct {
	Signals []freshness.Signal
}

type freshnessModel struct {
	cfg     config.Config
	spinner spinner.Model
	signals []freshness.Signal
	done    bool
}

func RunFreshness(cfg config.Config) error {
	m := newFreshnessModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newFreshnessModel(cfg config.Config) freshnessModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return freshnessModel{cfg: cfg, spinner: s}
}

func (m freshnessModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, freshnessCmd(m.cfg))
}

// PARTE CRITICA **********************
// Freshness se basa en timestamps de Prometheus, no en valores del payload.
// Si se cambia a heuristicas, se pierde capacidad de detectar datos viejos.
// No agregar nuevas metricas aqui.
// FIN DE PARTE CRITICA ****************
func freshnessCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		signals := freshness.Check(cfg)
		return freshnessMsg{Signals: signals}
	}
}

func (m freshnessModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case freshnessMsg:
		m.signals = v.Signals
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

func (m freshnessModel) View() string {
	title := TitleStyle.Render("drone-observe freshness")
	sub := WarnStyle.Render(fmt.Sprintf("Umbrales: warn=%ds, fail=%ds", m.cfg.FreshnessWarnSec, m.cfg.FreshnessFailSec))

	var body strings.Builder
	body.WriteString(fmt.Sprintf("%s\n%s\n%s\n", title, sub, strings.Repeat("─", 44)))

	if !m.done {
		body.WriteString("\n" + m.spinner.View())
		return BoxStyle.Render(body.String())
	}

	tw := tabwriter.NewWriter(&body, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, HeaderStyle.Render("Signal")+"\t"+HeaderStyle.Render("Ultima muestra")+"\t"+HeaderStyle.Render("Estado"))
	_, _ = fmt.Fprintln(tw, "------\t--------------\t------")
	for _, s := range m.signals {
		age := "sin datos"
		if s.AgeSeconds >= 0 {
			age = fmt.Sprintf("hace %ds", s.AgeSeconds)
		}
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\n", s.Name, age, formatFreshStatus(s.Status))
	}
	_ = tw.Flush()

	body.WriteString("\nPresiona 'q' para salir.\n")
	return BoxStyle.Render(body.String())
}

func formatFreshStatus(s freshness.Status) string {
	switch s {
	case freshness.StatusOK:
		return OKStyle.Render("OK")
	case freshness.StatusWarn:
		return WarnStyle.Render("WARN")
	default:
		return FailStyle.Render("FAIL")
	}
}
