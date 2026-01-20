# 11-cli-observability-tool.md

## Proposito
`drone-observe` es un CLI determinista para validar el pipeline de observabilidad sin depender de UI web. Su foco es operativo: confirmar que el Control Plane y el Data Plane cumplen contratos y que el flujo basico de metricas funciona.

No sustituye Grafana, Prometheus ni sistemas SOC (secure_core, Elastic/Fleet). Tampoco edita dashboards ni introduce metricas nuevas.

## Alcance
- Validar reachability de componentes clave (MQTT, backend, Prometheus).
- Visualizar en vivo el Data Plane con metricas existentes.
- Auditar el contrato de `METRICS.md` contra Prometheus y backend.

No hace:
- Diagnostico avanzado de fallas.
- Correlacion SOC/SIEM.
- Edicion de dashboards o metricas.

## Ubicacion y estructura
El CLI vive en `tools/drone-observe/` y se mantiene separado del runtime de observabilidad.

Estructura base:
```
tools/
  drone-observe/
    main.go
    cmd/
    internal/
```

## Comandos
### 1) health
Valida el Control Plane con un checklist visual:
- MQTT reachable
- Backend /metrics accesible
- Prometheus accesible
- Flujo de metricas (rate(mqtt_messages_total[1m]) > 0)

Uso:
```bash
drone-observe health
```

### 2) telemetry
Visualiza el Data Plane en vivo:
- Ultima bateria (drone_battery_last_pct)
- Tasa de mensajes por segundo

Uso:
```bash
drone-observe telemetry
```

### 3) validate
Audita el contrato de `METRICS.md`:
- Todas las metricas del contrato existen
- No hay metricas inesperadas en el backend

Uso:
```bash
drone-observe validate
```

### 4) topology
Muestra la topologia efectiva del sistema:
- Edge -> MQTT -> Backend -> Prometheus -> Grafana
- Componentes OK vs mudos

Uso:
```bash
drone-observe topology
```

### 5) freshness
Evalua recencia de datos observados:
- Tiempo desde ultima muestra
- Semaforo temporal

Uso:
```bash
drone-observe freshness
```

### 6) drift
Detecta deriva entre docs y estado real:
- Metricas documentadas vs reales
- Dashboards versionados vs docs

Uso:
```bash
drone-observe drift
```

### 7) limits
Expone limites tecnicos observados:
- Frecuencia de mensajes
- Cadencia de scrape observada
- Conteo de metricas y cardinalidad

Uso:
```bash
drone-observe limits
```

## Ayuda multi-idioma
La ayuda es bilingue y explicita (sin auto-deteccion):
```bash
drone-observe --help --es
drone-observe --help --en
drone-observe health --help --es
```

## Variables de entorno
- `MQTT_HOST` (default: `mqtt`)
- `MQTT_PORT` (default: `1883`)
- `BACKEND_HTTP_PORT` (default: `8080`)
- `PROMETHEUS_URL` (default: `http://localhost:9090`)
- `GRAFANA_URL` (default: `http://localhost:3000`)
- `METRICS_DOC` (default: `METRICS.md`)
- `FRESHNESS_WARN_SEC` (default: `30`)
- `FRESHNESS_FAIL_SEC` (default: `120`)

Nota: para `validate`, ejecutar desde la raiz del repo o ajustar `METRICS_DOC`.

## Prerrequisitos (Windows 11)
### Instalar Go
Opcion MSI (oficial):
1) Descargar el instalador desde `https://go.dev/dl/`.
2) Ejecutar el MSI y aceptar valores por defecto (`C:\Program Files\Go`).
3) Abrir una nueva terminal y verificar:
```bash
go version
```

Opcion CLI (winget):
```bash
winget install --id GoLang.Go -e
```

Opcion CLI (Chocolatey, si ya esta instalado):
```bash
choco install golang -y
```

## Build y ejecucion (Windows 11)
Desde `tools/drone-observe`:

1) Resolver dependencias:
```bash
go mod tidy
```

2) Compilar (evita VCS stamping en Windows):
```bash
go build -buildvcs=false .
```

3) Ejecutar binario local:
```bash
.\drone-observe.exe health
```

### Ejecutar sin `.\` (usar PATH)
Para usar `drone-observe` sin prefijo:

1) Instalar binario en GOPATH/bin:
```bash
go install -buildvcs=false .
```

2) Agregar GOPATH/bin al PATH (CMD):
```bat
for /f "delims=" %i in ('go env GOPATH') do set "GOPATH=%i"
set PATH=%PATH%;%GOPATH%\bin
```

3) Agregar GOPATH/bin al PATH (PowerShell):
```powershell
$env:GOPATH = (go env GOPATH)
$env:PATH = "$env:PATH;$env:GOPATH\bin"
```

Nota: estos cambios de PATH son para la sesion actual. Para persistir, usar Variables de entorno del sistema.

## Principios de diseno
- Determinismo: resultados reproducibles y sin heuristicas ocultas.
- Minimalismo: solo metricas actuales, sin labels nuevos.
- GitOps: no modifica dashboards ni configura Grafana.
- Separacion de dominios: observabilidad tecnica ≠ SOC/SIEM.

## Limitaciones actuales
- No hay soporte de autenticacion para endpoints (entorno controlado).
- No hay soporte de TLS en MQTT (segun setup actual).
- No valida payloads MQTT ni contratos de eventos (solo metricas).

## Evolucion futura (FUTURO)
- Soporte opcional de auth para Prometheus.
- Validacion de EVENTS.md con un subscriber de solo lectura.
- Reportes exportables en formato texto/JSON (sin reemplazar UI).
