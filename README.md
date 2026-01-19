# Drone Observability (MQTT + Rust SOC + Prometheus/Grafana)

Repositorio base para un sistema de observabilidad de eventos en drones con arquitectura event-driven:
Edge → MQTT → Backend (Rust) → Observabilidad (Prometheus/Grafana).

## 1. Objetivo
- Recibir telemetría y eventos (eventos discretos + estado).
- Normalizar y enrutar por MQTT.
- Consumir y procesar en backend Rust (SOC / correlación / reglas).
- Exponer métricas y salud para observabilidad.
- Preparar base para simulación MAVLink/PX4 y futuro gemelo digital.

## 2. Requisitos
- Host: Arch Linux (solo Docker)
- Docker Engine + Docker Compose plugin
- No se instala Mosquitto/Grafana/Prometheus/Rust en el host.

## 3. Arquitectura
### 3.1 Capas
- Edge: publica `telemetry` + `event` en MQTT
- MQTT Broker: transporte y desacoplo
- Backend (Rust): consumo, métricas, API, reglas
- Observabilidad: Prometheus scraping + Grafana dashboards

### 3.2 Contratos (Topics y Payloads)
- Base topic: `drone/<id>`
- Topics:
  - `drone/<id>/telemetry` (JSON)
  - `drone/<id>/event` (JSON)
- Reglas:
  - Telemetry: alta frecuencia, QoS 0
  - Event: baja frecuencia, QoS 1

## 4. Arranque rápido
```bash
docker compose up --build

Endpoints:

Grafana: http://localhost:3000
 (admin/admin)

Prometheus: http://localhost:9090

Backend: http://localhost:8080/healthz
 y /metrics

5. Desarrollo
5.1 Logs útiles
docker compose logs -f edge
docker compose logs -f backend
docker compose logs -f mqtt

5.2 Pruebas MQTT (manuales)
docker exec -it mqtt mosquitto_sub -t 'drone/#' -v

6. Observabilidad
Qué queda “críticamente asentado” con esto

Contrato MQTT (drone/<id>/telemetry y drone/<id>/event) listo para evolucionar.

Backend Rust ya consume y mide (métricas reales).

Observabilidad “de verdad”: Prometheus + Grafana provisionado (GitOps).

Todo dockerizado con runtime Ubuntu 24.04 en tus servicios propios.