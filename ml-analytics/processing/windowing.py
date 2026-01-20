"""
Windowing for streaming telemetry.

Scope: maintain a fixed-size window of feature vectors.
Non-goals: persistence, cross-drone state.
"""

from collections import deque
from typing import Deque, Dict, List


class WindowBuffer:
    def __init__(self, size: int) -> None:
        self._size = size
        self._data: Deque[Dict[str, float]] = deque(maxlen=size)

    def add(self, features: Dict[str, float]) -> None:
        self._data.append(features)

    def as_list(self) -> List[Dict[str, float]]:
        return list(self._data)

    def count(self) -> int:
        return len(self._data)
