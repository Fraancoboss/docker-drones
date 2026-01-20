# METRICS.md

## 1. Proposito y principios
Este documento define el contrato de metricas Prometheus del sistema de drones.
Objetivos principales:
- Observar salud y comportamiento del flujo Edge -> MQTT -> Backend.
- Detectar fallos sin explotar cardinalidad (sin labels de alta variacion).
- Mantener nombres estables y compatibles con Prometheus/Grafana.

Principios:
- No usar lat/long ni etiquetas de alta cardinalidad.
- Preferir gauges/counters con unidades explicitas.
- Separar "estado actual" vs "FUTURO" para evolucionar sin romper.

## 2. Convenciones
- Nombres: `snake_case` y en ingles tecnico si aplica.
- Unidades: usar sufijos estandar (`_total`, `_pct`, `_ms`, `_dbm`, `_celsius`).
- Tipos:
  - `counter`: incrementa monotonico (eventos, errores).
  - `gauge`: valor actual (bateria, rssi, temperatura).
- Etiquetas permitidas (estado actual): ninguna.
- Etiquetas permitidas (FUTURO, baja cardinalidad):
  - `drone_id` (numero reducido y controlado de drones).
  - `component` (valores fijos: `edge`, `backend`, `mqtt`).

## 3. Catalogo de metricas (estado actual)
| nombre | tipo | unidad | descripcion | labels | fuente | frecuencia esperada / notas |
|---|---|---|---|---|---|---|
| mqtt_messages_total | counter | mensajes | Total de mensajes MQTT consumidos por el backend (telemetria + eventos). | - | backend | Incrementa por cada publish recibido. |
| drone_battery_last_pct | gauge | pct | Ultimo porcentaje de bateria visto en telemetria. | - | backend | Actualiza cuando llega telemetria con `battery_pct`. |
| ml_anomaly_score | gauge | score | Anomaly score derivado del modelo ML (0-1). | - | ml-analytics | Calculado sobre ventana de telemetria/eventos. |
| ml_state | gauge | enum | Estado operacional derivado del anomaly score (0=OK,1=WARN,2=CRIT). | - | ml-analytics | Calculado sobre ventana de telemetria/eventos. |

## 4. Metricas minimas V1 (actuales)
- `mqtt_messages_total` (counter, mensajes).
- `drone_battery_last_pct` (gauge, pct).
- `ml_anomaly_score` (gauge, score).
- `ml_state` (gauge, enum).

## 5. Metricas recomendadas FUTURO
Marcadas como FUTURO y no implementadas aun.
- `mqtt_connected` (gauge, 0/1): estado de conexion al broker MQTT.
- `mqtt_errors_total` (counter, errores): errores al consumir MQTT.
- `drone_rtt_ms` (gauge, ms): latencia de red aproximada.
- `drone_packet_loss_pct` (gauge, pct): perdida de paquetes.
- `drone_rssi_dbm` (gauge, dbm): intensidad de senal.
- `drone_temperature_celsius` (gauge, celsius): temperatura de sensores.
- `drone_pressure_hpa` (gauge, hpa): presion atmosferica.
- `drone_light_lux` (gauge, lux): luminosidad ambiental.
- `vision_fps` (gauge, fps): tasa de frames procesados.
- `vision_inference_latency_ms` (FUTURO, histograma): latencia de inferencia.
- `vision_inference_latency_ms_bucket` (histograma): buckets para percentiles.
- `vision_inference_latency_ms_sum` (histograma): suma de observaciones.
- `vision_inference_latency_ms_count` (histograma): total de observaciones.
- `vision_detections_total` (counter, detecciones): total de detecciones.
- `lidar_obstacle_distance_m` (gauge, m): distancia a obstaculo.
- `lidar_ground_distance_m` (gauge, m): distancia al suelo.
- `telemetry_last_seen_ts` (gauge, unix): ultimo timestamp de telemetria visto.

## 6. Dashboards minimos (PromQL exacto)
Paneles con PromQL real. Los marcados como FUTURO requieren metricas no implementadas.
- Mensajes MQTT por segundo: `rate(mqtt_messages_total[1m])`
- Mensajes MQTT por minuto: `increase(mqtt_messages_total[1m])`
- Bateria actual (pct): `drone_battery_last_pct`
- Estado de conexion MQTT (FUTURO): `mqtt_connected`
- Errores MQTT por minuto (FUTURO): `increase(mqtt_errors_total[1m])`
- RTT promedio (FUTURO): `avg_over_time(drone_rtt_ms[5m])`
- Packet loss promedio (FUTURO): `avg_over_time(drone_packet_loss_pct[5m])`
- RSSI promedio (FUTURO): `avg_over_time(drone_rssi_dbm[5m])`
- Latencia de inferencia P95 (FUTURO): `histogram_quantile(0.95, rate(vision_inference_latency_ms_bucket[5m]))`
- Detecciones por minuto (FUTURO): `increase(vision_detections_total[1m])`

## 7. Alertas recomendadas (propuesta)
No se implementan aqui; se listan como guia.
- Silencio de telemetria: `increase(mqtt_messages_total[5m]) == 0`
- Bateria baja sostenida: `drone_battery_last_pct < 20`
- Errores MQTT altos (FUTURO): `rate(mqtt_errors_total[5m]) > 1`
