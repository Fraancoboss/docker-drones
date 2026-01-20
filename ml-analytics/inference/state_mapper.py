"""
Map anomaly score to operational state.

State enum: 0=OK, 1=WARN, 2=CRIT.
"""


def score_to_state(score: float, warn: float, crit: float, in_grace: bool = False) -> int:
    if score >= crit:
        if in_grace:
            return 1
        return 2
    if score >= warn:
        return 1
    return 0
