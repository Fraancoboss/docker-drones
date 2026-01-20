// Archivo: tools/drone-observe/internal/topology/topology.go
// Rol: checks de topologia efectiva sin discovery ni heuristicas.
// No hace: inferencias de negocio ni correlacion SOC.
package topology

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"drone-observe/internal/config"
	"drone-observe/internal/mqtt"
	"drone-observe/internal/prometheus"
)

type Status int

const (
	StatusOK Status = iota
	StatusSilent
	StatusFail
)

type Component struct {
	Name   string
	Status Status
	Detail string
}

const httpTimeout = 3 * time.Second

// PARTE CRITICA **********************
// La topologia se construye solo con checks explicitos y observables.
// Si se agregan supuestos ocultos, se degrada la gobernanza y la trazabilidad.
// No usar discovery dinamico ni inferencias de infraestructura.
// FIN DE PARTE CRITICA ****************
func Check(cfg config.Config) []Component {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var out []Component

	edge := Component{Name: "Edge"}
	rate, ok, err := prometheus.QueryInstant(ctx, cfg.PrometheusURL, "rate(mqtt_messages_total[1m])")
	if err != nil || !ok {
		edge.Status = StatusFail
		edge.Detail = "telemetria no observable"
	} else if rate <= 0 {
		edge.Status = StatusSilent
		edge.Detail = "telemetria silenciosa"
	} else {
		edge.Status = StatusOK
	}
	out = append(out, edge)

	mqttC := Component{Name: "MQTT Broker"}
	if err := mqtt.CheckReachable(cfg.MQTTHost, cfg.MQTTPort); err != nil {
		mqttC.Status = StatusFail
		mqttC.Detail = err.Error()
	} else {
		mqttC.Status = StatusOK
	}
	out = append(out, mqttC)

	backend := Component{Name: "Backend Rust"}
	if err := httpGetOK(ctx, cfg.BackendMetricsURL); err != nil {
		backend.Status = StatusFail
		backend.Detail = err.Error()
	} else {
		backend.Status = StatusOK
	}
	out = append(out, backend)

	prom := Component{Name: "Prometheus"}
	if err := prometheus.CheckReady(ctx, cfg.PrometheusURL); err != nil {
		prom.Status = StatusFail
		prom.Detail = err.Error()
	} else {
		prom.Status = StatusOK
	}
	out = append(out, prom)

	graf := Component{Name: "Grafana"}
	if err := httpGetOK(ctx, cfg.GrafanaURL+"/api/health"); err != nil {
		graf.Status = StatusFail
		graf.Detail = err.Error()
	} else {
		graf.Status = StatusOK
	}
	out = append(out, graf)

	return out
}

func httpGetOK(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http status %s", resp.Status)
	}
	return nil
}
