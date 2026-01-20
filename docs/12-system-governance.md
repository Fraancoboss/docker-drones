# 12-system-governance.md

## Filosofia de gobernanza
La gobernanza en `drone-observe` busca responder con evidencia: "puedo confiar en el sistema ahora y a medio plazo". No agrega features ni cambia contratos; valida el estado real y el alineamiento con lo documentado.

Principios:
- Determinismo: solo se reporta lo observable (HTTP/TCP/Prometheus).
- Minimalismo: si no aporta gobierno tecnico, no entra.
- Contratos primero: METRICS.md y dashboards versionados son la fuente de verdad.
- Sin deuda: no hay heuristicas ocultas ni auto-remediacion.

## Relacion con SOC / secure_core
Esta capa no es SOC ni SIEM. No analiza amenazas ni correlaciona eventos. Su valor es tecnico: asegurar que la observabilidad funciona y que el sistema no deriva de sus contratos. secure_core consume señales de seguridad desde Elastic/Fleet; `drone-observe` valida la consistencia del pipeline tecnico.

## Observabilidad vs Control
- Observabilidad: ver el estado y la tendencia del sistema (Grafana/Prometheus).
- Control: confirmar que el sistema esta alineado con sus contratos y limites (CLI).
El CLI no sustituye los dashboards; los valida y los contextualiza.

## Ejes de gobernanza (v2.x)
### 1) Topology
Confirma la topologia efectiva:
Edge -> MQTT -> Backend -> Prometheus -> Grafana.
No hace discovery; solo checks explicitos.

### 2) Freshness
Evalua recencia de datos usando timestamps de Prometheus:
- `drone_battery_last_pct`
- `mqtt_messages_total`
Semaforo temporal configurable por variables de entorno.

### 3) Drift
Detecta desviaciones entre:
- METRICS.md vs metricas reales
- Dashboards JSON vs docs
No corrige; solo reporta.

### 4) Limits
Expone limites observados sin benchmarks:
- Frecuencia de mensajes
- Cadencia de scrape observada
- Conteo de metricas y cardinalidad

## Como usar el CLI como auditor tecnico
1) `drone-observe health` para Control Plane.
2) `drone-observe topology` para confirmar cadena completa.
3) `drone-observe freshness` para detectar datos viejos.
4) `drone-observe drift` para garantizar alineamiento con docs.
5) `drone-observe limits` para entender limites actuales sin especular.

## Parametros de gobernanza
Variables relevantes:
- `FRESHNESS_WARN_SEC` (default: 30)
- `FRESHNESS_FAIL_SEC` (default: 120)

Estos umbrales son tecnicos y deben revisarse con el operador del sistema.
