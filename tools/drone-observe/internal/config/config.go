// Archivo: tools/drone-observe/internal/config/config.go
// Rol: leer configuracion determinista desde variables de entorno.
// No hace: lectura de archivos .env ni auto-deteccion de entorno.
package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	MQTTHost          string
	MQTTPort          int
	BackendMetricsURL string
	PrometheusURL     string
	MetricsDocPath    string
}

const (
	defaultMQTTHost      = "mqtt"
	defaultMQTTPort      = 1883
	defaultBackendPort   = 8080
	defaultPrometheusURL = "http://localhost:9090"
	defaultMetricsDoc    = "METRICS.md"
)

// PARTE CRITICA **********************
// Las rutas y defaults deben mantenerse estables para garantizar ejecucion reproducible.
// Si se cambian sin documentar, se rompen los contratos de uso en CLI/Docs.
// No agregar logica que intente "adivinar" paths fuera del repo.
// FIN DE PARTE CRITICA ****************
func FromEnv() Config {
	mqttHost := getenv("MQTT_HOST", defaultMQTTHost)
	mqttPort := getenvInt("MQTT_PORT", defaultMQTTPort)
	backendPort := getenvInt("BACKEND_HTTP_PORT", defaultBackendPort)
	backendURL := fmt.Sprintf("http://localhost:%d/metrics", backendPort)
	promURL := getenv("PROMETHEUS_URL", defaultPrometheusURL)

	return Config{
		MQTTHost:          mqttHost,
		MQTTPort:          mqttPort,
		BackendMetricsURL: backendURL,
		PrometheusURL:     promURL,
		MetricsDocPath:    getenv("METRICS_DOC", defaultMetricsDoc),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
