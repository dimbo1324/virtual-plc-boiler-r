from typing import Tuple


def get_demo_controls(
    elapsed: float, step_duration: float
) -> Tuple[str, float, float, float]:

    step = int(elapsed // step_duration)

    stages = [
        ("WARM-UP PHASE    ", 50.0, 10.0, 0.0),
        ("FULL POWER MODE  ", 100.0, 20.0, 0.0),
        ("STEAM BLOWDOWN   ", 100.0, 25.0, 80.0),
        ("SHUTDOWN         ", 0.0, 50.0, 100.0),
    ]

    return stages[min(step, len(stages) - 1)]
