"""
Unsupervised anomaly detector.

Uses IsolationForest over a sliding window. No labels, no persistence.
"""

from typing import List, Dict

import pandas as pd
from sklearn.ensemble import IsolationForest


class AnomalyDetector:
    def __init__(self, min_samples: int) -> None:
        self._min_samples = min_samples
        self._model = IsolationForest(
            n_estimators=100,
            contamination=0.1,
            random_state=42,
        )

    def score(self, window: List[Dict[str, float]]) -> float:
        if len(window) < self._min_samples:
            return 0.0
        frame = pd.DataFrame(window)
        self._model.fit(frame)
        scores = self._model.score_samples(frame)
        if scores.size == 0:
            return 0.0
        latest = scores[-1]
        min_s = scores.min()
        max_s = scores.max()
        if max_s - min_s == 0:
            return 0.0
        anomaly = (max_s - latest) / (max_s - min_s)
        if anomaly < 0:
            return 0.0
        if anomaly > 1:
            return 1.0
        return float(anomaly)
