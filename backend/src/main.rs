use axum::{response::IntoResponse, routing::get, Router};
use prometheus::{Encoder, IntCounter, IntGauge, Registry, TextEncoder};
use rumqttc::{AsyncClient, Event, Incoming, MqttOptions, QoS};
use std::{env, net::SocketAddr, sync::Arc, time::Duration};
use tracing::{error, info};

#[derive(Clone)]
struct AppState {
    registry: Registry,
    mqtt_messages_total: IntCounter,
    battery_last_pct: IntGauge,
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    // La configuración por variables de entorno permite reutilizar la misma imagen
    // en distintos entornos sin recompilar ni acoplar el backend al broker.
    let mqtt_host = env::var("MQTT_HOST").unwrap_or_else(|_| "mqtt".to_string());
    let mqtt_port: u16 = env::var("MQTT_PORT").ok().and_then(|v| v.parse().ok()).unwrap_or(1883);
    let base = env::var("MQTT_BASE_TOPIC").unwrap_or_else(|_| "drone/alpha".to_string());

    // Exponemos métricas con Prometheus porque es el estándar de facto en observabilidad
    // y evita inventar protocolos propietarios difíciles de integrar con Grafana/Alerting.
    let registry = Registry::new();
    let mqtt_messages_total = IntCounter::new("mqtt_messages_total", "Total MQTT messages consumed").unwrap();
    let battery_last_pct = IntGauge::new("drone_battery_last_pct", "Last seen drone battery percentage").unwrap();

    registry.register(Box::new(mqtt_messages_total.clone())).unwrap();
    registry.register(Box::new(battery_last_pct.clone())).unwrap();

    let state = Arc::new(AppState {
        registry,
        mqtt_messages_total,
        battery_last_pct,
    });

    // MQTT como bus de eventos mantiene desacoplados edge y backend;
    // esto evita dependencias directas y permite escalar productores/consumidores por separado.
    // ⚠️ client_id fijo: válido para laboratorio.
    // En despliegues con múltiples instancias debe ser único por instancia
    // (por ejemplo añadiendo hostname o un UUID) o usar shared subscriptions.
    let mut mqtt_options = MqttOptions::new("backend", mqtt_host, mqtt_port);
    mqtt_options.set_keep_alive(Duration::from_secs(30));

    let (client, mut eventloop) = AsyncClient::new(mqtt_options, 10);

    // QoS distintos: telemetria tolera perdida (volumen alto); eventos piden entrega al menos una vez.
    // Error comun: usar QoS alto para todo y saturar el broker con reintentos.
    let t_telemetry = format!("{}/telemetry", base);
    let t_event = format!("{}/event", base);
    client.subscribe(t_telemetry.clone(), QoS::AtMostOnce).await.unwrap();
    client.subscribe(t_event.clone(), QoS::AtLeastOnce).await.unwrap();

    info!("Subscribed to {}, {}", t_telemetry, t_event);

    // Loop dedicado para consumir MQTT y no bloquear el servidor HTTP.
    // En el futuro puede aislarse en un task supervisor si se agregan mas suscripciones.
    let state_mqtt = state.clone();
    tokio::spawn(async move {
        loop {
            match eventloop.poll().await {
                Ok(Event::Incoming(Incoming::Publish(p))) => {
                    state_mqtt.mqtt_messages_total.inc();

                    // Solo parseamos el campo que nos interesa para no acoplar el backend
                    // a esquemas completos; evita romperse ante cambios de payload.
                    if p.topic.ends_with("/telemetry") {
                        if let Ok(v) = serde_json::from_slice::<serde_json::Value>(&p.payload) {
                            if let Some(b) = v.get("battery_pct").and_then(|x| x.as_i64()) {
                                state_mqtt.battery_last_pct.set(b as i64);
                            }
                        }
                    }
                }
                Ok(_) => {}
                Err(e) => {
                    // ⚠️ Deuda técnica aceptada: no exponemos métricas de desconexión ni backoff.
                    // En producción esto puede ocultar caídas del broker y generar ruido en logs.
                    error!("MQTT poll error: {e}. retrying...");
                    tokio::time::sleep(Duration::from_secs(1)).await;
                }
            }
        }
    });

    // HTTP separado de MQTT: simplifica healthchecks y scraping sin mezclar protocolos.
    let app = Router::new()
        .route("/healthz", get(|| async { "ok" }))
        .route("/metrics", get(metrics_handler))
        .with_state(state);

    let addr: SocketAddr = "0.0.0.0:8080".parse().unwrap();
    info!("HTTP listening on {}", addr);
    axum::serve(tokio::net::TcpListener::bind(addr).await.unwrap(), app).await.unwrap();
}

async fn metrics_handler(axum::extract::State(state): axum::extract::State<Arc<AppState>>) -> impl IntoResponse {
    let metric_families = state.registry.gather();
    let mut buffer = Vec::new();
    let encoder = TextEncoder::new();
    encoder.encode(&metric_families, &mut buffer).unwrap();
    String::from_utf8(buffer).unwrap()
}
