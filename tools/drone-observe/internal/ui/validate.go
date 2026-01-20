// Archivo: tools/drone-observe/internal/ui/validate.go
// Rol: TUI para validar el contrato METRICS.md contra Prometheus y backend.
// No hace: inferencias SOC ni normalizacion de metricas.
package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"drone-observe/internal/config"
	"drone-observe/internal/prometheus"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
)

type validateItem struct {
	Name   string
	Status healthStatus
	Detail string
}

type validateResultMsg struct {
	Items []validateItem
	OK    bool
}

type validateModel struct {
	cfg     config.Config
	spinner spinner.Model
	items   []validateItem
	done    bool
	ok      bool
}

func RunValidate(cfg config.Config) error {
	m := newValidateModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newValidateModel(cfg config.Config) validateModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return validateModel{
		cfg:     cfg,
		spinner: s,
	}
}

func (m validateModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, validateCmd(m.cfg))
}

// PARTE CRITICA **********************
// Validacion debe seguir el contrato de METRICS.md y no inventar reglas.
// Si se flexibiliza, se pierde el valor de auditoria y control de deuda tecnica.
// No mezclar con reglas SOC ni heuristicas operativas.
// FIN DE PARTE CRITICA ****************
func validateCmd(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()

		contractMetrics, err := readMetricsContract(cfg.MetricsDocPath)
		if err != nil {
			return validateResultMsg{
				Items: []validateItem{{Name: "Leer METRICS.md", Status: statusFail, Detail: err.Error()}},
				OK:    false,
			}
		}

		items := make([]validateItem, 0, len(contractMetrics)+1)
		for _, name := range contractMetrics {
			val, ok, err := prometheus.QueryInstant(ctx, cfg.PrometheusURL, name)
			if err != nil || !ok {
				items = append(items, validateItem{
					Name:   fmt.Sprintf("Metrica %s", name),
					Status: statusFail,
					Detail: "no visible en Prometheus",
				})
				continue
			}
			items = append(items, validateItem{
				Name:   fmt.Sprintf("Metrica %s", name),
				Status: statusOK,
				Detail: fmt.Sprintf("valor=%.2f", val),
			})
		}

		unexpected, err := readBackendMetrics(cfg.BackendMetricsURL)
		if err != nil {
			items = append(items, validateItem{
				Name:   "Metricas inesperadas en backend",
				Status: statusFail,
				Detail: err.Error(),
			})
		} else {
			extra := diffUnexpected(contractMetrics, unexpected)
			if len(extra) > 0 {
				items = append(items, validateItem{
					Name:   "Metricas inesperadas en backend",
					Status: statusFail,
					Detail: strings.Join(extra, ", "),
				})
			} else {
				items = append(items, validateItem{
					Name:   "Metricas inesperadas en backend",
					Status: statusOK,
					Detail: "ninguna",
				})
			}
		}

		allOK := true
		for _, it := range items {
			if it.Status != statusOK {
				allOK = false
				break
			}
		}
		return validateResultMsg{Items: items, OK: allOK}
	}
}

func (m validateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case validateResultMsg:
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

func (m validateModel) View() string {
	title := TitleStyle.Render("drone-observe validate")
	statusLine := WarnStyle.Render("Ejecutando validacion...")
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
		line := fmt.Sprintf("%s\t%s", statusIcon(it.Status), it.Name)
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

func readMetricsContract(path string) ([]string, error) {
	f, usedPath, err := openMetricsFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var metrics []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cols := strings.Split(line, "|")
		if len(cols) < 2 {
			continue
		}
		name := strings.TrimSpace(cols[1])
		if name == "" || name == "nombre" || strings.HasPrefix(name, "---") {
			continue
		}
		metrics = append(metrics, name)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	_ = usedPath
	return metrics, nil
}

// PARTE CRITICA **********************
// Se buscan rutas relativas controladas para evitar fallos al ejecutar desde tools/drone-observe.
// Si se expanden rutas arbitrarias, se pierde determinismo y trazabilidad del contrato.
// No usar paths absolutos hardcodeados aqui.
// FIN DE PARTE CRITICA ****************
func openMetricsFile(path string) (*os.File, string, error) {
	candidates := []string{
		path,
		filepath.Join("..", path),
		filepath.Join("..", "..", path),
	}
	for _, p := range candidates {
		if f, err := os.Open(p); err == nil {
			return f, p, nil
		}
	}
	return nil, "", fmt.Errorf("no se encontro %s en rutas conocidas", path)
}

// PARTE CRITICA **********************
// Se usa /metrics del backend para detectar metricas no contractuales.
// Prometheus agrega metricas propias; por eso se evita usar label __name__.
// No filtrar ni suprimir nombres aqui: se debe exponer el drift.
// FIN DE PARTE CRITICA ****************
func readBackendMetrics(url string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := getBody(ctx, url)
	if err != nil {
		return nil, err
	}

	set := map[string]struct{}{}
	scanner := bufio.NewScanner(strings.NewReader(resp))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name := line
		if idx := strings.IndexAny(line, " {"); idx > 0 {
			name = line[:idx]
		}
		set[name] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var names []string
	for k := range set {
		names = append(names, k)
	}
	sort.Strings(names)
	return names, nil
}

func diffUnexpected(contract, actual []string) []string {
	allowed := map[string]struct{}{}
	for _, c := range contract {
		allowed[c] = struct{}{}
	}
	var extra []string
	for _, a := range actual {
		if _, ok := allowed[a]; !ok {
			extra = append(extra, a)
		}
	}
	sort.Strings(extra)
	return extra
}
