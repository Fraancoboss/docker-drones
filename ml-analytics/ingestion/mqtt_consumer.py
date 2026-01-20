"""
MQTT consumer for raw telemetry and events.

Scope: subscribe to MQTT topics and forward payloads to processing.
Non-goals: persistence, transformation beyond JSON parsing.
"""

import json
import queue
from dataclasses import dataclass
from typing import Any, Dict

import paho.mqtt.client as mqtt

from config import settings


@dataclass
class MqttMessage:
    kind: str  # telemetry | event
    payload: Dict[str, Any]


class MqttConsumer:
    def __init__(self, out_queue: queue.Queue) -> None:
        self._queue = out_queue
        self._client = mqtt.Client()
        self._client.on_connect = self._on_connect
        self._client.on_message = self._on_message

    def connect(self) -> None:
        self._client.connect(settings.MQTT_HOST, settings.MQTT_PORT, keepalive=60)
        self._client.loop_start()

    def stop(self) -> None:
        self._client.loop_stop()
        self._client.disconnect()

    def _on_connect(self, client: mqtt.Client, userdata: Any, flags: Dict[str, Any], rc: int) -> None:
        base = settings.MQTT_BASE_TOPIC
        client.subscribe(f"{base}/telemetry", qos=0)
        client.subscribe(f"{base}/event", qos=1)

    def _on_message(self, client: mqtt.Client, userdata: Any, msg: mqtt.MQTTMessage) -> None:
        kind = "telemetry" if msg.topic.endswith("/telemetry") else "event"
        try:
            payload = json.loads(msg.payload.decode("utf-8"))
        except (ValueError, json.JSONDecodeError):
            return
        self._queue.put(MqttMessage(kind=kind, payload=payload))
