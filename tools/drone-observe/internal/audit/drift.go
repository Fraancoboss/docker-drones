// Archivo: tools/drone-observe/internal/audit/drift.go
// Rol: auditoria de deriva entre estado real y documentacion.
// No hace: correccion automatica ni mutacion de config.
package audit

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"drone-observe/internal/config"
	"drone-observe/internal/prometheus"
)

type Severity string

const (
	SeverityHigh Severity = "alta"
	SeverityMed  Severity = "media"
	SeverityLow  Severity = "baja"
)

type Finding struct {
	Severity Severity
	Item     string
	Detail   string
}

// PARTE CRITICA **********************
// La deriva se evalua solo contra artefactos versionados y metricas reales.
// Si se agregan heuristicas, el resultado pierde valor como gobierno tecnico.
// No usar fuentes externas no versionadas.
// FIN DE PARTE CRITICA ****************
func Drift(cfg config.Config) ([]Finding, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	contract, err := readMetricsContract(cfg.MetricsDocPath)
	if err != nil {
		return []Finding{{Severity: SeverityHigh, Item: "METRICS.md", Detail: err.Error()}}, nil
	}

	findings := []Finding{}

	for _, name := range contract {
		_, ok, err := prometheus.QueryInstant(ctx, cfg.PrometheusURL, name)
		if err != nil || !ok {
			findings = append(findings, Finding{
				Severity: SeverityHigh,
				Item:     "Metrica documentada ausente",
				Detail:   name,
			})
		}
	}

	actual, err := readBackendMetrics(cfg.BackendMetricsURL)
	if err != nil {
		findings = append(findings, Finding{
			Severity: SeverityHigh,
			Item:     "Backend /metrics",
			Detail:   err.Error(),
		})
	} else {
		extra := diffUnexpected(contract, actual)
		for _, e := range extra {
			findings = append(findings, Finding{
				Severity: SeverityMed,
				Item:     "Metrica no documentada",
				Detail:   e,
			})
		}
	}

	dashFindings := checkDashboardDocs()
	findings = append(findings, dashFindings...)

	docFindings, err := checkDocsMetrics(contract)
	if err == nil {
		findings = append(findings, docFindings...)
	}

	order := map[Severity]int{SeverityHigh: 0, SeverityMed: 1, SeverityLow: 2}
	sort.Slice(findings, func(i, j int) bool {
		return order[findings[i].Severity] < order[findings[j].Severity]
	})
	return findings, nil
}

func checkDashboardDocs() []Finding {
	findings := []Finding{}
	dashDir := filepath.Join("observability", "grafana", "dashboards")
	files, _ := filepath.Glob(filepath.Join(dashDir, "*.json"))

	for _, f := range files {
		base := filepath.Base(f)
		doc := ""
		switch base {
		case "drones-control-plane.json":
			doc = filepath.Join("docs", "09-dashboard-control-plane.md")
		case "drones-data-plane.json":
			doc = filepath.Join("docs", "10-dashboard-data-plane.md")
		default:
			doc = ""
		}

		if doc == "" {
			findings = append(findings, Finding{
				Severity: SeverityMed,
				Item:     "Dashboard sin doc",
				Detail:   base,
			})
			continue
		}
		if _, err := os.Stat(doc); err != nil {
			findings = append(findings, Finding{
				Severity: SeverityHigh,
				Item:     "Doc faltante para dashboard",
				Detail:   base,
			})
		}
	}

	return findings
}

func checkDocsMetrics(contract []string) ([]Finding, error) {
	allowed := map[string]struct{}{"up": {}}
	for _, c := range contract {
		allowed[c] = struct{}{}
	}

	docFiles := []string{
		filepath.Join("docs", "09-dashboard-control-plane.md"),
		filepath.Join("docs", "10-dashboard-data-plane.md"),
	}

	findings := []Finding{}
	for _, doc := range docFiles {
		content, err := os.ReadFile(doc)
		if err != nil {
			continue
		}
		metrics := extractMetricTokens(string(content))
		for _, m := range metrics {
			if _, ok := allowed[m]; !ok {
				findings = append(findings, Finding{
					Severity: SeverityMed,
					Item:     "Doc refiere metrica no documentada",
					Detail:   fmt.Sprintf("%s -> %s", filepath.Base(doc), m),
				})
			}
		}
	}
	return findings, nil
}

func extractMetricTokens(s string) []string {
	re := regexp.MustCompile(`\b[a-z][a-z0-9_]*\b`)
	raw := re.FindAllString(s, -1)
	seen := map[string]struct{}{}
	var out []string
	for _, token := range raw {
		if token != "up" && !strings.Contains(token, "_") {
			continue
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		out = append(out, token)
	}
	sort.Strings(out)
	return out
}

func readMetricsContract(path string) ([]string, error) {
	f, _, err := openMetricsFile(path)
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
	return metrics, nil
}

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

func readBackendMetrics(url string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	body, err := fetchText(ctx, url)
	if err != nil {
		return nil, err
	}

	set := map[string]struct{}{}
	scanner := bufio.NewScanner(strings.NewReader(body))
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

func fetchText(ctx context.Context, url string) (string, error) {
	req, err := httpRequest(ctx, url)
	if err != nil {
		return "", err
	}
	return req, nil
}
