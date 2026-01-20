# 10-dashboard-data-plane.md

## 1. Proposito del dashboard
Este dashboard muestra el estado tecnico del dato observado (data plane), usando solo metricas actuales. Resuelve la necesidad de ver el ultimo valor operativo disponible y su tendencia, sin mezclarlo con salud del pipeline.

No intenta diagnosticar fallas del backend, ni sustituir SOC/SIEM, ni agregar detalle por drone. Su alcance se limita a las metricas actuales definidas en `METRICS.md`.

## 2. Principios de diseno
- Data Plane separado del Control Plane: se observa el dato, no la infraestructura.
- Estado actual vs tendencia: el dato presente en stat y el flujo en series.
- Minimalismo intencional: pocos paneles, alta claridad operativa.
- Sin alta cardinalidad: no se usan labels ni dimensiones por drone.

## 3. Ubicacion y ciclo de vida
- JSON fuente: `observability/grafana/dashboards/drones-data-plane.json`.
- Provisioning: `observability/grafana/provisioning/dashboards/dashboards.yml`.
- Datasource: `observability/grafana/provisioning/datasources/datasources.yml`.
- No se edita en UI; se versiona en Git para trazabilidad (GitOps).

## 4. Descripcion panel por panel
### Panel 1
- Nombre exacto: `Bateria actual (%)`.
- Tipo: Stat.
- PromQL exacto: `drone_battery_last_pct`.
- Significado operativo: ultimo porcentaje de bateria recibido en telemetria.
- No decidir solo con este panel: no indica frescura del dato ni tendencia.
- Fallos comunes: interpretar dato viejo como actual; usarlo como alerta SOC.

### Panel 2
- Nombre exacto: `MQTT mensajes por segundo (1m) - actual`.
- Tipo: Stat.
- PromQL exacto: `rate(mqtt_messages_total[1m])`.
- Significado operativo: valor actual estimado de tasa de mensajes.
- No decidir solo con este panel: no diferencia eventos vs telemetria.
- Fallos comunes: interpretar el valor como volumen absoluto sin contexto.

### Panel 3
- Nombre exacto: `MQTT mensajes por segundo (1m)`.
- Tipo: Time series.
- PromQL exacto: `rate(mqtt_messages_total[1m])`.
- Significado operativo: tendencia de la tasa de mensajes.
- No decidir solo con este panel: no diferencia eventos vs telemetria.
- Fallos comunes: tratar picos como incidentes de seguridad.

### Panel 4
- Nombre exacto: `MQTT mensajes por minuto (1m)`.
- Tipo: Time series.
- PromQL exacto: `increase(mqtt_messages_total[1m])`.
- Significado operativo: volumen por minuto, util para ver caidas o pausas.
- No decidir solo con este panel: no reemplaza healthchecks del backend.
- Fallos comunes: asumir cero como fallo SOC; puede ser falta de trafico real.

## 5. Semantica de colores y thresholds
- `Bateria actual (%)`: verde >= 20, amarillo >= 10, rojo < 10.
- Los colores son operativos, no severidad SOC.
- Un valor "verde" no implica salud completa del sistema.

## 6. Relacion con el SOC
- Aporta contexto operativo del dron (estado de bateria y flujo observado).
- No sustituye SIEM, Fleet, EDR ni correlacion en secure_core.
- Debe usarse como complemento tecnico, no como fuente de alertas SOC.

## 7. Limitaciones actuales
- No hay paneles por drone ni etiquetas de baja cardinalidad (no implementadas).
- No hay indicadores de frescura de telemetria (timestamp).
- No hay separacion entre telemetria y eventos en metricas actuales.

## 8. Evolucion futura controlada (FUTURO)
- Podrian incluirse: `telemetry_last_seen_ts`, `mqtt_connected`, `mqtt_errors_total` si existen.
- No deben incluirse aqui: sensores detallados ni analitica avanzada (van a dashboards dedicados).
- Criterio de adicion: metricas actuales, baja cardinalidad, contrato estable en `METRICS.md`.

## 9. Checklist operacional
1) Verificar `drone_battery_last_pct` en Prometheus antes de culpar a Grafana.
2) Verificar `rate(mqtt_messages_total[1m])` para confirmar flujo.
3) Confirmar datasource Prometheus (default).
4) Confirmar provisioning y presencia del JSON en `observability/grafana/dashboards/`.
5) Si falta el dashboard, reiniciar Grafana para re-provisioning.
