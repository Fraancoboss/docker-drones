// Archivo: tools/drone-observe/internal/freshness/freshness.go
// Rol: evaluacion de recencia de datos via timestamps de Prometheus.
// No hace: inferencias ni autocorrecciones.
package freshness

import (
	"context"
	"time"

	"drone-observe/internal/config"
	"drone-observe/internal/prometheus"
)

type Status int

const (
	StatusOK Status = iota
	StatusWarn
	StatusFail
)

type Signal struct {
	Name       string
	AgeSeconds int
	Status     Status
	Detail     string
}

// PARTE CRITICA **********************
// La recencia se calcula con timestamps reales de Prometheus.
// Si se cambia a valores inferidos, se rompe la capacidad de detectar datos viejos.
// No usar labels ni nuevas metricas para este calculo.
// FIN DE PARTE CRITICA ****************
func Check(cfg config.Config) []Signal {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return []Signal{
		checkMetric(ctx, cfg, "drone_battery_last_pct", "Bateria"),
		checkMetric(ctx, cfg, "mqtt_messages_total", "MQTT mensajes"),
	}
}

func checkMetric(ctx context.Context, cfg config.Config, metric, label string) Signal {
	_, ts, ok, err := prometheus.QueryInstantWithTimestamp(ctx, cfg.PrometheusURL, metric)
	if err != nil || !ok {
		return Signal{
			Name:       label,
			AgeSeconds: -1,
			Status:     StatusFail,
			Detail:     "sin datos",
		}
	}

	age := int(time.Since(time.Unix(int64(ts), 0)).Seconds())
	status := StatusOK
	if age >= cfg.FreshnessFailSec {
		status = StatusFail
	} else if age >= cfg.FreshnessWarnSec {
		status = StatusWarn
	}

	return Signal{
		Name:       label,
		AgeSeconds: age,
		Status:     status,
		Detail:     "ok",
	}
}
