# 13-ml-analytics.md

## Proposito
Definir el servicio `ml-analytics` como extension de observabilidad inteligente.
El servicio observa el flujo Edge -> MQTT y expone senales como metricas Prometheus.

No crea dashboards, no modifica el backend, no publica eventos MQTT.

## Contratos que respeta
- `METRICS.md`: define el contrato de metricas Prometheus. El ML solo agrega:
  - `ml_anomaly_score` (gauge, score 0-1).
  - `ml_state` (gauge, enum 0=OK,1=WARN,2=CRIT).
- `EVENTS.md`: define eventos MQTT. El ML no publica eventos.

## Input del ML (fuente de verdad)
- MQTT topics:
  - `drone/<id>/telemetry`
  - `drone/<id>/event`

El ML no usa Prometheus como input.

## Output del ML
- Metrics Prometheus expuestas por HTTP:
  - `ml_anomaly_score`
  - `ml_state`

Prometheus las scrapea y Grafana las visualiza.

## Modelo (restricciones)
- Unsupervised (sin labels).
- Implementado con Isolation Forest.
- Ventanas temporales y latencia en segundos.
- Sin deep learning, sin LLMs.

## Estructura del servicio
```
ml-analytics/
  README.md
  requirements.txt
  Makefile
  config/
  ingestion/
  processing/
  models/
  inference/
  export/
  main.py
```

## Setup (Python)
Pandas no tiene wheels estables para Python 3.14.
Para evitar errores de build, usa Python 3.11 o 3.12.

### Opcion recomendada (Makefile)
Desde `ml-analytics/`:
```bash
make deps
make run
```

Si necesitas una version especifica:
```bash
make PYTHON=python3.12 deps
```

### Instalacion de Python (Arch + pyenv)
Instalar pyenv y dependencias:
```bash
sudo pacman -S pyenv base-devel openssl zlib xz tk
```

Instalar Python 3.12 y fijarlo en el repo:
```bash
pyenv install 3.12.4
pyenv local 3.12.4
```

Crear entorno y ejecutar con pyenv:
```bash
cd ml-analytics
make PYTHON=$(pyenv which python) run
```

### Entorno virtual (manual)
Si no usas Makefile:
```bash
python3.12 -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt
python main.py
```

## Salida actual (ejemplo)
Prometheus expone las metricas en:
```
http://localhost:9108/metrics
```

Ejemplo de salida real:
```
# HELP ml_anomaly_score Anomaly score from ML (0-1)
# TYPE ml_anomaly_score gauge
ml_anomaly_score 0.9213899649406001
# HELP ml_state Operational state derived from anomaly score (0=OK,1=WARN,2=CRIT)
# TYPE ml_state gauge
ml_state 2.0
```

## Configuracion por entorno
Variables principales:
- `MQTT_HOST` (default: `mqtt`)
- `MQTT_PORT` (default: `1883`)
- `MQTT_BASE_TOPIC` (default: `drone/alpha`)
- `PROMETHEUS_PORT` (default: `9108`)
- `WINDOW_SIZE` (default: `60`)
- `MIN_SAMPLES` (default: `20`)
- `ANOMALY_WARN` (default: `0.6`)
- `ANOMALY_CRIT` (default: `0.85`)
- `BATTERY_LOW_GRACE_SEC` (default: `120`)
- `BATTERY_LOW_BATTERY_WEIGHT` (default: `0.2`)

## No-goals
- No control de drones.
- No eventos MQTT nuevos.
- No labels de alta cardinalidad.
- No cambios en backend/edge/dashboards.
