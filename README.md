# Drone Observability (MQTT + Rust SOC + Prometheus/Grafana)

Repositorio base para un sistema de observabilidad de eventos en drones con arquitectura event-driven:
Edge -> MQTT -> Backend (Rust) -> Observabilidad (Prometheus/Grafana).

## Objetivo
- Recibir telemetria y eventos (eventos discretos + estado).
- Normalizar y enrutar por MQTT.
- Consumir y procesar en backend Rust (SOC / correlacion / reglas).
- Exponer metricas y salud para observabilidad.
- Preparar base para simulacion MAVLink/PX4 y futuro gemelo digital.

## Requisitos
- Host con Docker Engine y Docker Compose plugin.
- No se instala Mosquitto/Grafana/Prometheus/Rust en el host.

## Arquitectura
- Edge: publica `telemetry` + `event` en MQTT.
- MQTT Broker: transporte y desacoplo.
- Backend (Rust): consumo, metricas, API, reglas.
- Observabilidad: Prometheus scraping + Grafana dashboards.

## Contratos MQTT
- Base topic: `drone/<id>`.
- Topics:
  - `drone/<id>/telemetry` (JSON).
  - `drone/<id>/event` (JSON).
- Reglas:
  - Telemetry: alta frecuencia, QoS 0.
  - Event: baja frecuencia, QoS 1.

## Uso
### Construir imagenes
```bash
docker compose build
```

### Levantar contenedores
```bash
docker compose up -d
```

Para primera vez puedes usar:
```bash
docker compose up --build
```

### Parar todo
```bash
docker compose down
```

Si quieres eliminar datos persistidos:
```bash
docker compose down -v
```

## Endpoints
- Grafana: http://localhost:3000 (admin/admin por defecto, ver `.env`).
- Prometheus: http://localhost:9090.
- Backend: http://localhost:8080/healthz y http://localhost:8080/metrics.

## Logs utiles
```bash
docker compose logs -f edge
docker compose logs -f backend
docker compose logs -f mqtt
```

## Pruebas MQTT (manuales)
```bash
docker exec -it mqtt mosquitto_sub -t 'drone/#' -v
```

## Notas
- Configuracion central en `.env` (puerto HTTP, credenciales de Grafana, base topic).
- Grafana y Prometheus quedan provisionados desde `observability/grafana` y `observability/prometheus.yml`.
