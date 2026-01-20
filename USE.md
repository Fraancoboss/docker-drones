# USE.md

Comandos de observabilidad (operacion diaria).

## Docker / servicios
Estado de contenedores:
```bash
docker compose ps
```

Logs en vivo:
```bash
docker compose logs -f edge
docker compose logs -f backend
docker compose logs -f mqtt
docker compose logs -f prometheus
docker compose logs -f grafana
```

## MQTT (telemetria/eventos)
Escuchar todo el bus:
```bash
docker exec -it mqtt mosquitto_sub -t 'drone/#' -v
```

Escuchar solo telemetria:
```bash
docker exec -it mqtt mosquitto_sub -t 'drone/+/telemetry' -v
```

Escuchar solo eventos:
```bash
docker exec -it mqtt mosquitto_sub -t 'drone/+/event' -v
```

## Backend (metricas)
Health:
```bash
curl http://localhost:8080/healthz
```

Metrics:
```bash
curl http://localhost:8080/metrics
```

## Prometheus
UI:
```
http://localhost:9090
```

Query rapido:
```bash
curl 'http://localhost:9090/api/v1/query?query=rate(mqtt_messages_total[1m])'
```

## Grafana
UI:
```
http://localhost:3000
```

## ML Analytics (metrics)
Exporter:
```
http://localhost:9108/metrics
```

Ver metricas ML:
```bash
curl http://localhost:9108/metrics | rg 'ml_anomaly_score|ml_state'
```

Si `curl` falla, el servicio ML no esta corriendo. Levantalo:
```bash
cd ml-analytics
MQTT_HOST=localhost make PYTHON=$(pyenv which python) run
```

Si `drone-observe llm` muestra N/A, valida el target en Prometheus:
```
http://localhost:9090/targets
```

Si el target esta DOWN, recrea Prometheus para aplicar `extra_hosts`:
```bash
docker compose up -d --force-recreate prometheus
```

## CLI drone-observe
Control Plane:
```bash
drone-observe health
```

Data Plane:
```bash
drone-observe telemetry
```

Estado ML:
```bash
drone-observe llm
```

Contratos de metricas:
```bash
drone-observe validate
```

Topologia:
```bash
drone-observe topology
```

Recencia:
```bash
drone-observe freshness
```

Deriva:
```bash
drone-observe drift
```

Limites:
```bash
drone-observe limits
```

Ayuda ES/EN:
```bash
drone-observe --help --es
drone-observe --help --en
```
