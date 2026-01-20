// Archivo: tools/drone-observe/internal/ui/llm.go
// Rol: TUI para estado ML (ml_anomaly_score y ml_state).
// No hace: inferencia ML ni mutaciones de configuracion.
package ui

import (
	"context"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"drone-observe/internal/config"
	"drone-observe/internal/prometheus"

	tea "github.com/charmbracelet/bubbletea"
)

const llmRefresh = 2 * time.Second

type llmMsg struct {
	Score      string
	State      string
	AlertLabel string
	UpdatedAt  time.Time
	HasError   bool
	ErrorLabel string
}

type llmModel struct {
	cfg        config.Config
	last       llmMsg
	lastUpdate time.Time
}

func RunLLM(cfg config.Config) error {
	m := llmModel{cfg: cfg}
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m llmModel) Init() tea.Cmd {
	return tea.Batch(fetchLLMCmd(m.cfg), llmTickCmd())
}

// PARTE CRITICA **********************
// El estado ML se consulta via Prometheus para mantener consistencia con Grafana.
// No agregar queries fuera de METRICS.md.
// FIN DE PARTE CRITICA ****************
func fetchLLMCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		scoreVal, scoreOK, err := prometheus.QueryInstant(ctx, cfg.PrometheusURL, "ml_anomaly_score")
		if err != nil || !scoreOK {
			return llmMsg{
				HasError:   true,
				ErrorLabel: "sin score",
				UpdatedAt:  time.Now(),
			}
		}

		stateVal, stateOK, err := prometheus.QueryInstant(ctx, cfg.PrometheusURL, "ml_state")
		if err != nil || !stateOK {
			return llmMsg{
				HasError:   true,
				ErrorLabel: "sin estado",
				UpdatedAt:  time.Now(),
			}
		}

		stateLabel, alertLabel := llmStateLabels(stateVal)
		return llmMsg{
			Score:      fmt.Sprintf("%.3f", scoreVal),
			State:      stateLabel,
			AlertLabel: alertLabel,
			UpdatedAt:  time.Now(),
		}
	}
}

func llmTickCmd() tea.Cmd {
	return tea.Tick(llmRefresh, func(t time.Time) tea.Msg { return t })
}

func (m llmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case llmMsg:
		m.last = v
		m.lastUpdate = v.UpdatedAt
		return m, nil
	case time.Time:
		return m, tea.Batch(fetchLLMCmd(m.cfg), llmTickCmd())
	case tea.KeyMsg:
		if v.String() == "q" || v.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m llmModel) View() string {
	title := TitleStyle.Render("drone-observe llm")
	refresh := WarnStyle.Render(fmt.Sprintf("Refresh: %s", llmRefresh))
	ts := ""
	if !m.lastUpdate.IsZero() {
		ts = fmt.Sprintf("Ultima actualizacion: %s", m.lastUpdate.Format(time.RFC3339))
	}

	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, HeaderStyle.Render("Parametro")+"\t"+HeaderStyle.Render("Valor"))
	_, _ = fmt.Fprintln(tw, "---------\t-----")
	if m.last.HasError {
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", "ml_anomaly_score", FailStyle.Render("N/A ("+m.last.ErrorLabel+")"))
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", "ml_state", FailStyle.Render("N/A"))
	} else {
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", "ml_anomaly_score", OKStyle.Render(m.last.Score))
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", "ml_state", m.last.State)
	}
	_ = tw.Flush()

	alertLine := formatAlertLine(m.last)
	body := fmt.Sprintf("%s\n%s\n%s\n\n%s\n\n%s\nPresiona 'q' para salir.\n", title, refresh, SubtitleStyle.Render(ts), sb.String(), alertLine)
	return BoxStyle.Render(body)
}

func llmStateLabels(state float64) (string, string) {
	switch int(state + 0.5) {
	case 0:
		return OKStyle.Render("OK (0)"), OKStyle.Render("Sin alerta")
	case 1:
		return WarnStyle.Render("WARN (1)"), WarnStyle.Render("Alerta: observar y confirmar")
	case 2:
		return FailStyle.Render("CRIT (2)"), FailStyle.Render("Alerta: critica y activa")
	default:
		return WarnStyle.Render("UNKNOWN"), WarnStyle.Render("Alerta: estado desconocido")
	}
}

func formatAlertLine(msg llmMsg) string {
	if msg.HasError {
		return FailStyle.Render("Alerta: sin datos ML (ver Prometheus).")
	}
	return msg.AlertLabel
}
