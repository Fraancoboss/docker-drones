"""
Entry point for ML Analytics service.

Consumes MQTT telemetry/events and exports Prometheus metrics.
"""

import queue
import signal
import sys
import time

from config import settings
from export.prometheus_exporter import PrometheusExporter
from ingestion.mqtt_consumer import MqttConsumer
from inference.state_mapper import score_to_state
from models.anomaly_detector import AnomalyDetector
from processing.feature_extraction import telemetry_features, event_features
from processing.windowing import WindowBuffer


def main() -> int:
    stop_flag = {"stop": False}
    last_battery_low_ts = None
    last_telemetry = None

    def _handle_stop(signum, frame):
        stop_flag["stop"] = True

    signal.signal(signal.SIGINT, _handle_stop)
    signal.signal(signal.SIGTERM, _handle_stop)

    msg_queue: queue.Queue = queue.Queue(maxsize=1000)
    consumer = MqttConsumer(msg_queue)
    exporter = PrometheusExporter()
    exporter.start(settings.PROMETHEUS_PORT)

    window = WindowBuffer(settings.WINDOW_SIZE)
    detector = AnomalyDetector(settings.MIN_SAMPLES)

    consumer.connect()

    try:
        while not stop_flag["stop"]:
            try:
                msg = msg_queue.get(timeout=1.0)
            except queue.Empty:
                continue

            now = time.time()
            in_grace = False
            if last_battery_low_ts is not None:
                in_grace = (now - last_battery_low_ts) <= settings.BATTERY_LOW_GRACE_SEC

            battery_weight = settings.BATTERY_LOW_BATTERY_WEIGHT if in_grace else 1.0

            if msg.kind == "telemetry":
                features = telemetry_features(
                    msg.payload,
                    prev_payload=last_telemetry,
                    battery_weight=battery_weight,
                )
                last_telemetry = msg.payload
            else:
                if str(msg.payload.get("type", "")) == "BATTERY_LOW":
                    try:
                        event_ts = float(msg.payload.get("ts", now))
                    except (TypeError, ValueError):
                        event_ts = now
                    last_battery_low_ts = event_ts
                features = event_features(msg.payload)

            window.add(features)
            score = detector.score(window.as_list())
            state = score_to_state(score, settings.ANOMALY_WARN, settings.ANOMALY_CRIT, in_grace=in_grace)
            exporter.set_metrics(score, state)
    finally:
        consumer.stop()

    return 0


if __name__ == "__main__":
    sys.exit(main())
