// Archivo: tools/drone-observe/internal/limits/limits.go
// Rol: exponer limites tecnicos observados sin benchmarks.
// No hace: predicciones ni pruebas de carga.
package limits

import (
	"context"
	"time"

	"drone-observe/internal/config"
	"drone-observe/internal/prometheus"
)

type Snapshot struct {
	MessageRate      float64
	SeriesCount      float64
	MetricNameCount  float64
	ScrapeAgeSeconds int
}

// PARTE CRITICA **********************
// Se usan solo metricas observables en Prometheus para evitar heuristicas.
// Si se agregan benchmarks, se rompe el enfoque determinista del CLI.
// No extrapolar limites futuros; solo mostrar estado actual.
// FIN DE PARTE CRITICA ****************
func Observe(cfg config.Config) (Snapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rate, _, _ := prometheus.QueryInstant(ctx, cfg.PrometheusURL, "rate(mqtt_messages_total[1m])")
	series, _, _ := prometheus.QueryInstant(ctx, cfg.PrometheusURL, "count({job=\"backend\"})")
	names, _, _ := prometheus.QueryInstant(ctx, cfg.PrometheusURL, "count(count by(__name__) ({job=\"backend\"}))")

	_, ts, ok, _ := prometheus.QueryInstantWithTimestamp(ctx, cfg.PrometheusURL, "up{job=\"backend\"}")
	age := -1
	if ok {
		age = int(time.Since(time.Unix(int64(ts), 0)).Seconds())
	}

	return Snapshot{
		MessageRate:      rate,
		SeriesCount:      series,
		MetricNameCount:  names,
		ScrapeAgeSeconds: age,
	}, nil
}
