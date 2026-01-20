# DEPLOY_GUIDE.md

Guia de despliegue desde cero en Arch Linux (suponiendo dependencias instaladas).
Incluye: stack Docker, ML analytics en Python, y CLI en Go.

## 0) Ubicacion y pre-requisitos
Trabaja desde la raiz del repo:
```bash
cd /home/fran/developer/Cibersecurity/BlueTeam/Proyects/drone-observavility
```

## 1) Levantar el stack principal (Docker)
Esto levanta: MQTT, backend Rust, Prometheus y Grafana.
```bash
docker compose build
docker compose up -d
```

Verifica contenedores:
```bash
docker compose ps
```

Endpoints:
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Backend: http://localhost:8080/healthz y http://localhost:8080/metrics

## 2) Levantar ML Analytics (Python)
Esto levanta: consumidor MQTT + exporter Prometheus de `ml_anomaly_score` y `ml_state`.

### 2.1 Opcion recomendada (pyenv + Makefile)
```bash
cd ml-analytics
pyenv install 3.12.4
pyenv local 3.12.4
make PYTHON=$(pyenv which python) run
```

### 2.2 Opcion manual (venv)
```bash
cd ml-analytics
python3.12 -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt
python main.py
```

Verifica metricas ML:
```bash
curl http://localhost:9108/metrics | rg 'ml_anomaly_score|ml_state'
```

## 3) Compilar y usar el CLI (Go)
Esto levanta: CLI `drone-observe` para validaciones y vistas en tiempo real.

```bash
cd tools/drone-observe
go mod tidy
go build -buildvcs=false .
./drone-observe health
```

Instalar en PATH (opcional):
```bash
go install -buildvcs=false .
export GOPATH="$(go env GOPATH)"
export PATH="$PATH:$GOPATH/bin"
drone-observe telemetry
drone-observe llm
```

## 4) Comprobaciones rapidas (todo arriba)
- MQTT: `docker exec -it mqtt mosquitto_sub -t 'drone/#' -v`
- Backend metrics: `curl http://localhost:8080/metrics`
- Prometheus: `http://localhost:9090`
- ML metrics: `http://localhost:9108/metrics`
- CLI: `drone-observe health`, `drone-observe telemetry`, `drone-observe llm`

## 5) Bajar todo
### 5.1 Bajar contenedores Docker
```bash
docker compose down
```

### 5.2 Borrar datos persistidos (opcional)
```bash
docker compose down -v
```

### 5.3 Detener ML Analytics
- Si corre en foreground: `Ctrl+C`
- Si corre en otra terminal: cerrar la sesion o enviar `Ctrl+C`

### 5.4 Limpiar entorno Python (opcional)
```bash
cd ml-analytics
rm -rf .venv
```

## 6) Notas de entorno
- `ml-analytics` usa MQTT como input y expone metricas en `9108`.
- `drone-observe` consulta Prometheus (no usa MQTT directo).
- Variables utiles (si necesitas overrides):
  - `MQTT_HOST`, `MQTT_PORT`, `MQTT_BASE_TOPIC`
  - `PROMETHEUS_URL` (CLI)
  - `PROMETHEUS_PORT` (ML exporter)
  - `ANOMALY_WARN`, `ANOMALY_CRIT` (ML thresholds)
