// Archivo: tools/drone-observe/internal/ui/telemetry.go
// Rol: TUI para visualizar Data Plane en vivo sin UI web.
// No hace: graficos complejos ni historicos.
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

const telemetryRefresh = 2 * time.Second

type telemetryMsg struct {
	Battery    string
	MsgRate    string
	UpdatedAt  time.Time
	HasError   bool
	ErrorLabel string
}

type telemetryModel struct {
	cfg        config.Config
	last       telemetryMsg
	lastUpdate time.Time
}

func RunTelemetry(cfg config.Config) error {
	m := telemetryModel{cfg: cfg}
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m telemetryModel) Init() tea.Cmd {
	return tea.Batch(fetchTelemetryCmd(m.cfg), tickCmd())
}

// PARTE CRITICA **********************
// La telemetria se consulta via Prometheus para mantener consistencia con Grafana.
// Si se consulta directo al backend, se puede ocultar discrepancias de scraping.
// No agregar nuevos queries fuera de METRICS.md.
// FIN DE PARTE CRITICA ****************
func fetchTelemetryCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		batteryVal, batteryOK, err := prometheus.QueryInstant(ctx, cfg.PrometheusURL, "drone_battery_last_pct")
		if err != nil || !batteryOK {
			return telemetryMsg{
				HasError:   true,
				ErrorLabel: "sin bateria",
				UpdatedAt:  time.Now(),
			}
		}

		msgRateVal, msgRateOK, err := prometheus.QueryInstant(ctx, cfg.PrometheusURL, "rate(mqtt_messages_total[1m])")
		if err != nil || !msgRateOK {
			return telemetryMsg{
				HasError:   true,
				ErrorLabel: "sin rate",
				UpdatedAt:  time.Now(),
			}
		}

		return telemetryMsg{
			Battery:   fmt.Sprintf("%.0f%%", batteryVal),
			MsgRate:   fmt.Sprintf("%.2f msg/s", msgRateVal),
			UpdatedAt: time.Now(),
		}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(telemetryRefresh, func(t time.Time) tea.Msg { return t })
}

func (m telemetryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case telemetryMsg:
		m.last = v
		m.lastUpdate = v.UpdatedAt
		return m, nil
	case time.Time:
		return m, tea.Batch(fetchTelemetryCmd(m.cfg), tickCmd())
	case tea.KeyMsg:
		if v.String() == "q" || v.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m telemetryModel) View() string {
	title := TitleStyle.Render("drone-observe telemetry")
	refresh := WarnStyle.Render(fmt.Sprintf("Refresh: %s", telemetryRefresh))
	ts := ""
	if !m.lastUpdate.IsZero() {
		ts = fmt.Sprintf("Ultima actualizacion: %s", m.lastUpdate.Format(time.RFC3339))
	}

	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, HeaderStyle.Render("Bateria actual (%)")+"\t"+HeaderStyle.Render("Mensajes por segundo"))
	_, _ = fmt.Fprintln(tw, "------------------\t-------------------")
	if m.last.HasError {
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", FailStyle.Render("N/A"), FailStyle.Render("N/A ("+m.last.ErrorLabel+")"))
	} else {
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", OKStyle.Render(m.last.Battery), OKStyle.Render(m.last.MsgRate))
	}
	_ = tw.Flush()

	body := fmt.Sprintf("%s\n%s\n%s\n\n%s\nPresiona 'q' para salir.\n", title, refresh, SubtitleStyle.Render(ts), sb.String())
	return BoxStyle.Render(body)
}
