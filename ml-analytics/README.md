# ML Analytics (Observability)

This service consumes raw MQTT telemetry/events and emits Prometheus metrics for ML-based
observability. It does not control drones and it does not publish MQTT events.

## Scope
- Input: MQTT `drone/<id>/telemetry` and `drone/<id>/event`.
- Output: Prometheus metrics `ml_anomaly_score` and `ml_state`.
- Mode: streaming / near-real-time.

## Non-goals
- No dashboards or APIs.
- No backend/edge changes.
- No high-cardinality labels.

## Metrics (contract extension)
- `ml_anomaly_score` (gauge, score 0-1): anomaly score from ML pipeline.
- `ml_state` (gauge, enum 0=OK, 1=WARN, 2=CRIT): operational state from anomaly score.

## Run (local)
1) Create a venv and install deps:
```bash
python -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt
```

2) Set env (optional):
```bash
export MQTT_HOST=localhost
export MQTT_PORT=1883
export MQTT_BASE_TOPIC=drone/alpha
export PROMETHEUS_PORT=9108
export WINDOW_SIZE=60
export MIN_SAMPLES=20
```

3) Start:
```bash
python ml-analytics/main.py
```

Prometheus can scrape `http://localhost:9108/metrics`.
