# Despliegue y operación

## Flujo de despliegue con Docker
- `docker compose build` construye imágenes.
- `docker compose up -d` crea y arranca contenedores.
- El sistema vive en Docker aunque el host sea Arch Linux.

## Variables de entorno
- `MQTT_HOST`: hostname del broker.
- `MQTT_PORT`: puerto MQTT.
- `MQTT_BASE_TOPIC`: raíz de topics por drone.

## Operación básica
- Ver estado: `docker compose ps`.
- Ver logs: `docker compose logs -f <servicio>`.
- Ver métricas: `curl http://localhost:8080/metrics`.

## Debugging y observabilidad operativa
- Confirmar suscripciones en backend por logs.
- Validar publicación con `mosquitto_sub`.
- Identificar cortes de broker por errores repetidos.

## Errores comunes en producción
- Construir imágenes sin hacer `up`.
- Reutilizar `client_id` entre instancias.
- Publicar en topics no esperados por el backend.

## Estado actual vs evolución
- Estado actual: operación manual y simple.
- Evolución futura: healthchecks más ricos y métricas de conectividad.
