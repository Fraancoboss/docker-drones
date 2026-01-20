// Archivo: tools/drone-observe/internal/mqtt/mqtt.go
// Rol: chequeo basico de reachability MQTT via TCP.
// No hace: suscripciones ni validacion de payloads.
package mqtt

import (
	"fmt"
	"net"
	"time"
)

const dialTimeout = 2 * time.Second

// PARTE CRITICA **********************
// Se usa TCP simple para confirmar reachability sin introducir dependencias MQTT.
// Si se cambia a un handshake complejo, se pierde determinismo y aumenta fragilidad.
// No implementar publish/subscribe aqui; eso seria mezclar responsabilidades.
// FIN DE PARTE CRITICA ****************
func CheckReachable(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, dialTimeout)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}
