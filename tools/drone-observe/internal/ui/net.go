// Archivo: tools/drone-observe/internal/ui/net.go
// Rol: helpers de red compartidos por UIs, sin dependencias extra.
// No hace: retries avanzados ni backoff.
package ui

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const httpTimeout = 3 * time.Second

// PARTE CRITICA **********************
// Helper minimo para evitar variaciones por cliente HTTP global.
// Si se agrega logica compleja aqui, todos los comandos se vuelven mas lentos.
// No implementar autenticacion; este CLI asume entorno controlado.
// FIN DE PARTE CRITICA ****************
func simpleGet(ctx context.Context, url string) error {
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

func getBody(ctx context.Context, url string) (string, error) {
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
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
