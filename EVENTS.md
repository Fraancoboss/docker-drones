# EVENTS.md

## 1. Proposito y principios
Este documento define el contrato de eventos MQTT para drones.
Un evento es un hecho puntual y significativo (alerta, anomalia, cambio de estado).
La telemetria es un flujo continuo de estado y no debe mezclarse con eventos.

Principios:
- Eventos deben ser idempotentes o deduplicables.
- Eventos usan QoS 1; telemetria usa QoS 0.
- Contrato estable y ampliable (versionable).

## 2. Convenciones
- Topic base: `drone/<id>/event`.
  - Estado actual: el sistema usa `MQTT_BASE_TOPIC` (por defecto `drone/alpha`).
  - Por tanto, el topic real es `MQTT_BASE_TOPIC + "/event"`.
- `schema_version` (FUTURO): version del esquema de evento.
- `event_id` o `idempotency_key` (FUTURO): clave para deduplicacion.
- Campos obligatorios (estado actual): `ts`, `type`, `severity`.
- Campos recomendados (FUTURO): `schema_version`, `event_id`, `drone_id`, `payload`.
 - `severity` (contrato cerrado): `info`, `warning`, `critical`.
 - `type` (contrato cerrado): debe pertenecer al catalogo de este documento.
 - `payload`: especifico por evento; no debe duplicar valores ya expuestos como metricas.

## 3. Esquema base de evento (FUTURO)
Formato propuesto (JSON):
```json
{
  "schema_version": "1.0",
  "ts": 1710000000.123,
  "drone_id": "alpha",
  "type": "BATTERY_LOW",
  "severity": "warning",
  "payload": {
    "battery_pct": 25
  }
}
```

## 4. Catalogo de eventos V1 (estado actual)
Eventos realmente publicados por el edge:
- `BATTERY_LOW`
  - Trigger actual: `battery_pct == 25`
  - Campos presentes: `ts`, `type`, `severity`, `battery_pct`
  - QoS: 1

## 5. Eventos FUTUROS (no implementados)
- `OBJECT_DETECTED` (vision)
- `JAMMING_SUSPECTED`
- `SPOOFING_SUSPECTED`
- `LIDAR_OBSTACLE`

## 6. Ejemplos completos
### 6.1 Telemetria (para contrastar, estado actual)
Topic: `drone/alpha/telemetry`
```json
{
  "seq": 101,
  "ts": 1710000000.123,
  "battery_pct": 87,
  "altitude_m": 12.3
}
```

### 6.2 Evento actual (BATTERY_LOW)
Topic: `drone/alpha/event`
```json
{
  "ts": 1710000123.456,
  "type": "BATTERY_LOW",
  "severity": "warning",
  "battery_pct": 25
}
```

### 6.3 Evento FUTURO (con schema y payload)
Topic: `drone/alpha/event`
```json
{
  "schema_version": "1.0",
  "ts": 1710000123.456,
  "drone_id": "alpha",
  "event_id": "evt-0001",
  "type": "LINK_DEGRADED",
  "severity": "warning",
  "payload": {
    "rssi_dbm": -85,
    "packet_loss_pct": 12
  }
}
```

## 7. Buenas practicas SOC
- Deduplicacion: usar `event_id` o un hash de `(type, ts, drone_id, payload)` y ventana temporal.
- Correlacion: agrupar por `drone_id`, `type`, y rango de tiempo.
- Auditoria: incluir `schema_version`, `event_id`, `drone_id`, `ts`, `severity`, `source`, `firmware_version` (FUTURO).
