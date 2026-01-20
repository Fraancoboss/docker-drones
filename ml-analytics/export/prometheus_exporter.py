"""
Prometheus exporter for ML metrics.

Only exports approved metrics from METRICS.md extension.
"""

from prometheus_client import Gauge, start_http_server


class PrometheusExporter:
    def __init__(self) -> None:
        self.ml_anomaly_score = Gauge(
            "ml_anomaly_score",
            "Anomaly score from ML (0-1)",
        )
        self.ml_state = Gauge(
            "ml_state",
            "Operational state derived from anomaly score (0=OK,1=WARN,2=CRIT)",
        )

    def start(self, port: int) -> None:
        start_http_server(port)

    def set_metrics(self, score: float, state: int) -> None:
        self.ml_anomaly_score.set(score)
        self.ml_state.set(state)
