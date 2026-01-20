// Archivo: tools/drone-observe/internal/audit/net.go
// Rol: helpers HTTP para auditoria sin dependencias externas.
// No hace: retries ni autenticacion.
package audit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const httpTimeout = 3 * time.Second

func httpRequest(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("http status %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
