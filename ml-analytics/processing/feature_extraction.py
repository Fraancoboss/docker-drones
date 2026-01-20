"""
Feature extraction from raw MQTT payloads.

Scope: convert telemetry/events into numeric features.
Non-goals: create new events or modify MQTT contracts.
"""

from typing import Dict, Optional


def telemetry_features(
    payload: Dict,
    prev_payload: Optional[Dict] = None,
    battery_weight: float = 1.0,
) -> Dict[str, float]:
    battery = float(payload.get("battery_pct", 0.0))
    altitude = float(payload.get("altitude_m", 0.0))

    if prev_payload:
        prev_battery = float(prev_payload.get("battery_pct", battery))
        prev_altitude = float(prev_payload.get("altitude_m", altitude))
    else:
        prev_battery = battery
        prev_altitude = altitude

    battery_delta = (battery - prev_battery) * battery_weight
    altitude_delta = altitude - prev_altitude

    return {
        "battery_delta": battery_delta,
        "altitude_delta": altitude_delta,
        "event_severity": 0.0,
        "battery_low_event": 0.0,
    }


def event_features(payload: Dict) -> Dict[str, float]:
    severity_map = {"info": 0.0, "warning": 1.0, "critical": 2.0}
    severity = severity_map.get(str(payload.get("severity", "info")), 0.0)
    event_type = str(payload.get("type", ""))
    battery_low = 1.0 if event_type == "BATTERY_LOW" else 0.0
    return {
        "battery_delta": 0.0,
        "altitude_delta": 0.0,
        "event_severity": severity,
        "battery_low_event": battery_low,
    }
