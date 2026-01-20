#!/usr/bin/env python3
# Archivo: edge/mavlink_to_mqtt/main.py
# Rol: simulador de edge que publica telemetria y eventos en MQTT.
# No hace: validacion estricta de esquema ni reintentos complejos.
# Decisiones intencionadas: payloads simples y consistentes para pruebas reproducibles.
import os
import json
import time
import random
import paho.mqtt.client as mqtt

MQTT_HOST = os.getenv("MQTT_HOST", "mqtt")
MQTT_PORT = int(os.getenv("MQTT_PORT", "1883"))
BASE_TOPIC = os.getenv("MQTT_BASE_TOPIC", "drone/alpha")

def topic(suffix: str) -> str:
    return f"{BASE_TOPIC}/{suffix}"

def encode_payload(data: dict) -> str:
    # PARTE CRITICA **********************
    # Helper aislado para no mezclar formatos en esta iteracion.
    # Cambiar aqui impacta el contrato de EVENTS.md y METRICS.md.
    # Suposicion: JSON es el unico formato activo.
    # Evolucion esperada: soporte TOON sin auto-deteccion ambigua (FUTURO).
    # FIN DE PARTE CRITICA ****************
    # Solo JSON por ahora; dejamos el punto de extension para TOON en el futuro.
    return json.dumps(data)

# MQTT simplifica el puente edge->backend sin depender de HTTP directo; es tolerante a
# desconexiones intermitentes y permite agregar consumidores sin tocar el edge.
client = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
client.connect(MQTT_HOST, MQTT_PORT, keepalive=30)

seq = 0
while True:
    # PARTE CRITICA **********************
    # Loop principal de emision; la cadencia define volumen de telemetria.
    # Cambios de frecuencia afectan carga del broker y la estabilidad del backend.
    # Suposicion: 1 Hz es suficiente para laboratorio.
    # Evolucion esperada: tasa configurable por entorno (FUTURO).
    # FIN DE PARTE CRITICA ****************
    # Validacion manual: docker exec -it mqtt mosquitto_sub -t "drone/#" -v
    seq += 1
    battery = max(0, min(100, 100 - (seq % 120)))
    altitude_m = 10 + 2 * random.random()

    # Telemetria frecuente: en sistemas reales puede ser alta tasa, por eso QoS 0
    # evita bloqueos por reintentos y no satura el broker si hay perdida temporal.
    # PARTE CRITICA **********************
    # Telemetria usa QoS 0 para evitar bloqueo por reintentos.
    # Mezclar telemetria y eventos rompe la semantica de observabilidad.
    # Suposicion: perdida puntual no compromete el monitoreo.
    # Evolucion esperada: enriquecimiento de campos sin romper el contrato actual.
    # FIN DE PARTE CRITICA ****************
    telemetry = {
        "seq": seq,
        "ts": time.time(),
        "battery_pct": battery,
        "altitude_m": altitude_m,
    }
    client.publish(topic("telemetry"), encode_payload(telemetry), qos=0)

    # Eventos puntuales usan QoS 1 porque perderlos dificulta alertas y post-mortem.
    # Error comun: mezclar eventos y telemetria en el mismo topic y perder semantica.
    if battery == 25:
        # PARTE CRITICA **********************
        # Evento puntual con QoS 1 para reducir perdida.
        # Cambiar la condicion o el tipo rompe el contrato de EVENTS.md.
        # Suposicion: battery_pct == 25 es el umbral de BATTERY_LOW.
        # Evolucion esperada: umbral configurable por politica SOC (FUTURO).
        # FIN DE PARTE CRITICA ****************
        event = {
            "ts": time.time(),
            "type": "BATTERY_LOW",
            "severity": "warning",
            "battery_pct": battery
        }
        client.publish(topic("event"), encode_payload(event), qos=1)

    time.sleep(1)
