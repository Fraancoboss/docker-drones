// Archivo: tools/drone-observe/cmd/root.go
// Rol: enrutador principal de comandos y ayuda multi-idioma (ES/EN).
// No hace: ejecucion de chequeos; eso vive en cmd/* y internal/*.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"drone-observe/internal/config"
)

type lang string

const (
	langES lang = "es"
	langEN lang = "en"
)

func Execute(args []string) int {
	cmd, flags := parseArgs(args[1:])
	language := parseLang(flags)

	if hasHelpFlag(flags) || cmd == "" {
		printHelp(cmd, language)
		return 0
	}

	cfg := config.FromEnv()

	switch cmd {
	case "health":
		return runHealth(cfg)
	case "telemetry":
		return runTelemetry(cfg)
	case "validate":
		return runValidate(cfg)
	default:
		printHelp("", language)
		return 2
	}
}

// PARTE CRITICA **********************
// El parseo manual garantiza ayuda determinista sin dependencias externas.
// Si se rompe, la UX queda ambigua y se pierde la promesa de help bilingue.
// No agregar auto-deteccion de idioma aqui; debe ser explicita por flag.
// FIN DE PARTE CRITICA ****************
func parseArgs(args []string) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}
	if strings.HasPrefix(args[0], "-") {
		return "", args
	}
	return args[0], args[1:]
}

func parseLang(flags []string) lang {
	for _, f := range flags {
		if f == "--en" {
			return langEN
		}
		if f == "--es" {
			return langES
		}
	}
	return langES
}

func hasHelpFlag(flags []string) bool {
	for _, f := range flags {
		if f == "--help" || f == "-h" {
			return true
		}
	}
	return false
}

func printHelp(cmd string, language lang) {
	if language == langEN {
		fmt.Fprint(os.Stdout, helpEN(cmd))
		return
	}
	fmt.Fprint(os.Stdout, helpES(cmd))
}

func helpES(cmd string) string {
	switch cmd {
	case "health":
		return `drone-observe health
Valida el Control Plane del pipeline.

Checks:
  - MQTT reachable
  - Backend /metrics accesible
  - Prometheus accesible
  - Flujo de metricas (rate(mqtt_messages_total[1m]) > 0)

Flags:
  --help, -h   ayuda
  --es         español (default)
  --en         english
`
	case "telemetry":
		return `drone-observe telemetry
Visualiza Data Plane en vivo (sin UI web).

Muestra:
  - Ultima bateria (drone_battery_last_pct)
  - Tasa de mensajes por segundo

Flags:
  --help, -h   ayuda
  --es         español (default)
  --en         english
`
	case "validate":
		return `drone-observe validate
Valida contratos de METRICS.md contra Prometheus.

Verifica:
  - Todas las metricas del contrato existen
  - No hay metricas inesperadas en el backend

Flags:
  --help, -h   ayuda
  --es         español (default)
  --en         english
`
	default:
		return `drone-observe
CLI de validacion del pipeline de observabilidad (GitOps).

Comandos:
  health     valida Control Plane
  telemetry  observa Data Plane en vivo
  validate   audita contratos de metricas

Flags:
  --help, -h   ayuda
  --es         español (default)
  --en         english

Variables de entorno:
  MQTT_HOST (default: mqtt)
  MQTT_PORT (default: 1883)
  BACKEND_HTTP_PORT (default: 8080)
  PROMETHEUS_URL (default: http://localhost:9090)

Nota: ejecutar desde la raiz del repo para leer METRICS.md.
`
	}
}

func helpEN(cmd string) string {
	switch cmd {
	case "health":
		return `drone-observe health
Validates the Control Plane of the pipeline.

Checks:
  - MQTT reachable
  - Backend /metrics reachable
  - Prometheus reachable
  - Metric flow (rate(mqtt_messages_total[1m]) > 0)

Flags:
  --help, -h   help
  --es         spanish (default)
  --en         english
`
	case "telemetry":
		return `drone-observe telemetry
Live Data Plane view (no web UI).

Shows:
  - Last battery (drone_battery_last_pct)
  - Messages per second rate

Flags:
  --help, -h   help
  --es         spanish (default)
  --en         english
`
	case "validate":
		return `drone-observe validate
Validates METRICS.md contract against Prometheus.

Checks:
  - All contract metrics exist
  - No unexpected backend metrics

Flags:
  --help, -h   help
  --es         spanish (default)
  --en         english
`
	default:
		return `drone-observe
Observability pipeline validation CLI (GitOps).

Commands:
  health     validate Control Plane
  telemetry  live Data Plane view
  validate   audit metric contracts

Flags:
  --help, -h   help
  --es         spanish (default)
  --en         english

Environment:
  MQTT_HOST (default: mqtt)
  MQTT_PORT (default: 1883)
  BACKEND_HTTP_PORT (default: 8080)
  PROMETHEUS_URL (default: http://localhost:9090)

Note: run from repo root to read METRICS.md.
`
	}
}
