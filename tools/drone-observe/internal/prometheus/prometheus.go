// Archivo: tools/drone-observe/internal/prometheus/prometheus.go
// Rol: cliente minimo para consultas Prometheus HTTP API.
// No hace: autodescubrimiento ni mutaciones de dashboards.
package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const httpTimeout = 3 * time.Second

type queryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func CheckReady(ctx context.Context, baseURL string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/-/ready", nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("prometheus not ready: %s", resp.Status)
	}
	return nil
}

// PARTE CRITICA **********************
// Las consultas se hacen via /api/v1/query para mantener compatibilidad Prometheus.
// Si se cambia a endpoints no estables, se rompe la validacion de contratos.
// No agregar queries que impliquen alta cardinalidad o labels variables.
// FIN DE PARTE CRITICA ****************
func QueryInstant(ctx context.Context, baseURL, expr string) (float64, bool, error) {
	u := fmt.Sprintf("%s/api/v1/query", baseURL)
	q := url.Values{}
	q.Set("query", expr)
	u = u + "?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, false, err
	}

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return 0, false, err
	}
	defer resp.Body.Close()

	var payload queryResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, false, err
	}
	if payload.Status != "success" || len(payload.Data.Result) == 0 {
		return 0, false, nil
	}
	if len(payload.Data.Result[0].Value) < 2 {
		return 0, false, nil
	}
	valStr, ok := payload.Data.Result[0].Value[1].(string)
	if !ok {
		return 0, false, nil
	}
	var val float64
	if _, err := fmt.Sscanf(valStr, "%f", &val); err != nil {
		return 0, false, nil
	}
	return val, true, nil
}
