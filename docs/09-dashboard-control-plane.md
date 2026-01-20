# 09-dashboard-control-plane.md

## 1. Proposito del dashboard
Este dashboard resume el estado del "control plane" de observabilidad tecnica: si el backend esta vivo, si hay flujo de mensajes y el ultimo estado de bateria observado. Resuelve el problema de detectar interrupciones basicas sin depender de UI interactiva ni logs.

No intenta diagnosticar causas raiz, ni reemplazar un SIEM/SOC, ni mostrar telemetria completa por drone. Su alcance se limita a lo definido en `METRICS.md` y a las metricas actuales disponibles.

## 2. Principios de diseno
- Control Plane vs Data Plane: se monitorea la salud del pipeline, no el contenido detallado de la telemetria.
- Estado actual vs tendencias: combina "estado actual" (stats) con tendencia (serie temporal).
- Minimalismo intencional: solo tres paneles para reducir ruido y decisiones precipitadas.
- Evitar alta cardinalidad: no hay labels por drone ni dimensiones inestables.

## 3. Ubicacion y ciclo de vida
- JSON fuente: `observability/grafana/dashboards/drones-control-plane.json`.
- Provisioning: `observability/grafana/provisioning/dashboards/dashboards.yml`.
- Datasource: `observability/grafana/provisioning/datasources/datasources.yml` (Prometheus por defecto).
- Se evita la UI para garantizar reproducibilidad y trazabilidad (GitOps).
- Cualquier cambio debe versionarse en Git y no editarse desde Grafana.

## 4. Descripcion panel por panel
### Panel 1
- Nombre exacto: `Backend UP`.
- Tipo: Stat.
- PromQL exacto: `up{job="backend"}`.
- Significado operativo: indica si Prometheus logra hacer scrape del backend.
- No decidir solo con este panel: no confirma que MQTT procese mensajes ni que la app este funcional a nivel logico.
- Fallos comunes: asumir que "UP" implica salud completa; ignorar problemas de pipeline MQTT.

### Panel 2
- Nombre exacto: `Batería actual (%)`.
- Tipo: Stat.
- PromQL exacto: `drone_battery_last_pct`.
- Significado operativo: ultimo porcentaje de bateria visto en telemetria.
- No decidir solo con este panel: no refleja tendencia ni garantiza frescura del dato.
- Fallos comunes: interpretar caida de bateria como incidente SOC; confundir dato viejo con dato actual.

### Panel 3
- Nombre exacto: `MQTT mensajes por segundo (1m)`.
- Tipo: Time series.
- PromQL exacto: `rate(mqtt_messages_total[1m])`.
- Significado operativo: tasa de mensajes consumidos por el backend.
- No decidir solo con este panel: no distingue telemetria vs eventos, ni valida integridad de payloads.
- Fallos comunes: interpretar picos como ataques; ignorar que la tasa puede variar por entorno.

## 5. Semantica de colores y thresholds
- Verde/amarillo/rojo indican umbrales operativos, no severidad SOC.
- `Backend UP`: verde cuando el valor es 1 (UP), rojo cuando es 0 (DOWN).
- `Batería actual (%)`: verde >= 20, amarillo >= 10, rojo < 10.
- UP no implica sistema sano: solo confirma scrape exitoso.

## 6. Relacion con el SOC
- Aporta una vista rapida del estado tecnico del pipeline de observabilidad.
- No sustituye SIEM, Fleet, EDR ni correlacion en secure_core.
- Se complementa con Elastic/Fleet para eventos de seguridad y con secure_core para analitica SOC.

## 7. Limitaciones actuales
- No hay paneles por drone ni por componente especifico (baja cardinalidad intencional).
- No hay indicador de frescura de telemetria (por ejemplo, timestamp de ultimo mensaje).
- No hay metricas de errores MQTT ni estado de conexion (solo estado actual).
- Ampliar sin metricas base aumenta ruido y deuda tecnica.

## 8. Evolucion futura controlada (FUTURO)
- Encajar aqui: `mqtt_connected`, `mqtt_errors_total`, `telemetry_last_seen_ts` (si se implementan).
- No encajar aqui: metricas de sensores o analitica avanzada (van a dashboards de data plane).
- Criterios para nuevos paneles: baja cardinalidad, alto valor operacional y contractualmente definidos en `METRICS.md`.

## 9. Checklist operacional
1) Verificar que Prometheus tenga el target del backend en UP.
2) Consultar `up{job="backend"}` en Prometheus para descartar problema de Grafana.
3) Consultar `drone_battery_last_pct` y `rate(mqtt_messages_total[1m])` en Prometheus.
4) Revisar datasource en Grafana (Prometheus como default).
5) Confirmar provisioning: `observability/grafana/provisioning/dashboards/dashboards.yml`.
6) Verificar que el JSON exista en `observability/grafana/dashboards/`.
7) Si falla solo Grafana, reiniciar el contenedor grafana y revalidar.
