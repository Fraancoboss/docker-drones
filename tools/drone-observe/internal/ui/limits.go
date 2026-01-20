// Archivo: tools/drone-observe/internal/ui/limits.go
// Rol: TUI para exponer limites tecnicos observados.
// No hace: benchmarks ni pruebas de carga.
package ui

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"drone-observe/internal/config"
	"drone-observe/internal/limits"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
)

type limitsMsg struct {
	Snapshot limits.Snapshot
	Err      error
}

type limitsModel struct {
	cfg     config.Config
	spinner spinner.Model
	data    limits.Snapshot
	done    bool
}

func RunLimits(cfg config.Config) error {
	m := newLimitsModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newLimitsModel(cfg config.Config) limitsModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return limitsModel{cfg: cfg, spinner: s}
}

func (m limitsModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, limitsCmd(m.cfg))
}

// PARTE CRITICA **********************
// Limits muestra solo datos observables sin extrapolar futuro.
// Si se agregan benchmarks, se rompe el enfoque de gobernanza determinista.
// No introducir pruebas de carga ni supuestos de escala.
// FIN DE PARTE CRITICA ****************
func limitsCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		snap, err := limits.Observe(cfg)
		return limitsMsg{Snapshot: snap, Err: err}
	}
}

func (m limitsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case limitsMsg:
		m.data = v.Snapshot
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

func (m limitsModel) View() string {
	title := TitleStyle.Render("drone-observe limits")
	sub := WarnStyle.Render("Limites observados (sin benchmark)")

	var body strings.Builder
	body.WriteString(fmt.Sprintf("%s\n%s\n%s\n", title, sub, strings.Repeat("─", 44)))

	if !m.done {
		body.WriteString("\n" + m.spinner.View())
		return BoxStyle.Render(body.String())
	}

	tw := tabwriter.NewWriter(&body, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, HeaderStyle.Render("Metrica")+"\t"+HeaderStyle.Render("Valor"))
	_, _ = fmt.Fprintln(tw, "-------\t-----")
	_, _ = fmt.Fprintf(tw, "Mensajes por segundo\t%.2f\n", m.data.MessageRate)
	_, _ = fmt.Fprintf(tw, "Series observadas\t%.0f\n", m.data.SeriesCount)
	_, _ = fmt.Fprintf(tw, "Nombres de metricas\t%.0f\n", m.data.MetricNameCount)
	if m.data.ScrapeAgeSeconds >= 0 {
		_, _ = fmt.Fprintf(tw, "Ultimo scrape (age)\t%ds\n", m.data.ScrapeAgeSeconds)
	} else {
		_, _ = fmt.Fprintf(tw, "Ultimo scrape (age)\tN/A\n")
	}
	_ = tw.Flush()

	body.WriteString("\nPresiona 'q' para salir.\n")
	return BoxStyle.Render(body.String())
}
