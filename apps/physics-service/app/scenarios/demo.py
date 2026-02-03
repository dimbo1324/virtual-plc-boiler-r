from typing import Tuple


def get_demo_controls(
    elapsed: float, step_duration: float
) -> Tuple[str, float, float, float]:
    """
    Returns control parameters for the demo scenario based on elapsed time.

    Args:
        elapsed: Time since simulation start (seconds)
        step_duration: Duration of each phase (seconds)

    Returns:
        Tuple of (phase_name, fuel_percent, water_percent, steam_percent)
    """
    step = int(elapsed // step_duration)

    stages = [
        # Phase name             fuel   water  steam
        ("WARM-UP PHASE    ", 50.0, 10.0, 0.0),  # 0–10s
        ("FULL POWER MODE  ", 100.0, 20.0, 0.0),  # 10–20s
        ("STEAM BLOWDOWN   ", 100.0, 25.0, 80.0),  # 20–30s
        ("SHUTDOWN         ", 0.0, 50.0, 100.0),  # 30s+
    ]

    # Return current phase or last one if time exceeded
    return stages[min(step, len(stages) - 1)]
