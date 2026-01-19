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

# MQTT simplifica el puente edge->backend sin depender de HTTP directo; es tolerante a
# desconexiones intermitentes y permite agregar consumidores sin tocar el edge.
client = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
client.connect(MQTT_HOST, MQTT_PORT, keepalive=30)

seq = 0
while True:
    seq += 1
    battery = max(0, min(100, 100 - (seq % 120)))
    altitude_m = 10 + 2 * random.random()

    # Telemetria frecuente: en sistemas reales puede ser alta tasa, por eso QoS 0
    # evita bloqueos por reintentos y no satura el broker si hay perdida temporal.
    telemetry = {
        "seq": seq,
        "ts": time.time(),
        "battery_pct": battery,
        "altitude_m": altitude_m,
    }
    client.publish(topic("telemetry"), json.dumps(telemetry), qos=0)

    # Eventos puntuales usan QoS 1 porque perderlos dificulta alertas y post-mortem.
    # Error comun: mezclar eventos y telemetria en el mismo topic y perder semantica.
    if battery == 25:
        event = {
            "ts": time.time(),
            "type": "BATTERY_LOW",
            "severity": "warning",
            "battery_pct": battery
        }
        client.publish(topic("event"), json.dumps(event), qos=1)

    time.sleep(1)
