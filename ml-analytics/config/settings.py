"""
Settings are sourced from environment to keep the service stateless and
Docker-friendly. No secrets are stored here.
"""

import os


def _get_env(name: str, default: str) -> str:
    return os.getenv(name, default)


MQTT_HOST = _get_env("MQTT_HOST", "mqtt")
MQTT_PORT = int(_get_env("MQTT_PORT", "1883"))
MQTT_BASE_TOPIC = _get_env("MQTT_BASE_TOPIC", "drone/alpha")

PROMETHEUS_PORT = int(_get_env("PROMETHEUS_PORT", "9108"))

WINDOW_SIZE = int(_get_env("WINDOW_SIZE", "60"))
MIN_SAMPLES = int(_get_env("MIN_SAMPLES", "20"))

ANOMALY_WARN = float(_get_env("ANOMALY_WARN", "0.6"))
ANOMALY_CRIT = float(_get_env("ANOMALY_CRIT", "0.85"))

BATTERY_LOW_GRACE_SEC = int(_get_env("BATTERY_LOW_GRACE_SEC", "120"))
BATTERY_LOW_BATTERY_WEIGHT = float(_get_env("BATTERY_LOW_BATTERY_WEIGHT", "0.2"))
