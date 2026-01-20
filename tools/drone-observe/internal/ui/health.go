// Archivo: tools/drone-observe/internal/ui/health.go
// Rol: TUI para el comando health (Control Plane).
// No hace: acciones correctivas ni configuracion.
package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"drone-observe/internal/config"
	"drone-observe/internal/mqtt"
	"drone-observe/internal/prometheus"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
)

type healthStatus int

const (
	statusPending healthStatus = iota
	statusOK
	statusFail
)

type healthItem struct {
	Name   string
	Status healthStatus
	Detail string
}

type healthResultMsg struct {
	Items []healthItem
	OK    bool
}

type healthModel struct {
	spinner spinner.Model
	items   []healthItem
	cfg     config.Config
	done    bool
	ok      bool
	err     error
}

func RunHealth(cfg config.Config) error {
	m := newHealthModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newHealthModel(cfg config.Config) healthModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return healthModel{
		spinner: s,
		cfg:     cfg,
		items: []healthItem{
			{Name: "MQTT reachable", Status: statusPending},
			{Name: "Backend /metrics", Status: statusPending},
			{Name: "Prometheus accesible", Status: statusPending},
			{Name: "Flujo de metricas", Status: statusPending},
		},
	}
}

func (m healthModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, healthChecksCmd(m.cfg))
}

// PARTE CRITICA **********************
// Health checks deben ser simples y deterministas para evitar falsos positivos.
// Si se agregan chequeos pesados, el CLI deja de ser usable en incidentes.
// No incorporar logica SOC ni correlacion aqui.
// FIN DE PARTE CRITICA ****************
func healthChecksCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items := []healthItem{
		{Name: "MQTT reachable"},
		{Name: "Backend /metrics"},
		{Name: "Prometheus accesible"},
		{Name: "Flujo de metricas"},
	}

	if err := mqtt.CheckReachable(cfg.MQTTHost, cfg.MQTTPort); err != nil {
		items[0].Status = statusFail
		items[0].Detail = err.Error()
	} else {
		items[0].Status = statusOK
	}

	if err := httpGetOK(ctx, cfg.BackendMetricsURL); err != nil {
		items[1].Status = statusFail
		items[1].Detail = err.Error()
	} else {
		items[1].Status = statusOK
	}

	if err := prometheus.CheckReady(ctx, cfg.PrometheusURL); err != nil {
		items[2].Status = statusFail
		items[2].Detail = err.Error()
	} else {
		items[2].Status = statusOK
	}

	val, ok, err := prometheus.QueryInstant(ctx, cfg.PrometheusURL, "rate(mqtt_messages_total[1m])")
	if err != nil || !ok || val <= 0 {
		items[3].Status = statusFail
		if err != nil {
			items[3].Detail = err.Error()
		} else {
			items[3].Detail = "rate=0 o sin datos"
		}
	} else {
		items[3].Status = statusOK
	}

	allOK := true
	for _, it := range items {
		if it.Status != statusOK {
			allOK = false
			break
		}
	}

		return healthResultMsg{Items: items, OK: allOK}
	}
}

func (m healthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case healthResultMsg:
		m.items = v.Items
		m.done = true
		m.ok = v.OK
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

func (m healthModel) View() string {
	title := TitleStyle.Render("drone-observe health")
	statusLine := WarnStyle.Render("Ejecutando checks...")
	if m.done {
		if m.ok {
			statusLine = OKStyle.Render("Resultado: OK")
		} else {
			statusLine = FailStyle.Render("Resultado: FAIL")
		}
	}

	var body strings.Builder
	body.WriteString(fmt.Sprintf("%s\n%s\n%s\n", title, statusLine, strings.Repeat("─", 44)))
	for _, it := range m.items {
		line := fmt.Sprintf("%s  %s", statusIcon(it.Status), it.Name)
		if it.Detail != "" {
			line += fmt.Sprintf(" (%s)", it.Detail)
		}
		body.WriteString(line + "\n")
	}
	if !m.done {
		body.WriteString("\n" + m.spinner.View())
	} else {
		body.WriteString("\nPresiona 'q' para salir.\n")
	}
	return BoxStyle.Render(body.String())
}

func statusIcon(s healthStatus) string {
	switch s {
	case statusOK:
		return OKStyle.Render("✔")
	case statusFail:
		return FailStyle.Render("✖")
	default:
		return WarnStyle.Render("…")
	}
}

func httpGetOK(ctx context.Context, url string) error {
	return simpleGet(ctx, url)
}
