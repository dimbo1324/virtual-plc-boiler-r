import pytest
from app.scenarios.demo import get_demo_controls
@pytest.mark.parametrize(
    "elapsed,step_duration,expected_index",
    [
        (0.0, 10.0, 0),
        (9.9, 10.0, 0),
        (10.0, 10.0, 1),
        (19.9, 10.0, 1),
        (20.0, 10.0, 2),
        (30.0, 10.0, 3),
        (100.0, 10.0, 3),
    ],
)
def test_get_demo_controls_stages(elapsed, step_duration, expected_index):
    stage_name, fuel, water, steam = get_demo_controls(elapsed, step_duration)
    assert isinstance(stage_name, str)
    assert 0.0 <= fuel <= 100.0
    assert 0.0 <= water <= 100.0
    assert 0.0 <= steam <= 100.0
    stages = [
        ("WARM-UP PHASE    ", 50.0, 10.0, 0.0),
        ("FULL POWER MODE  ", 100.0, 20.0, 0.0),
        ("STEAM BLOWDOWN   ", 100.0, 25.0, 80.0),
        ("SHUTDOWN         ", 0.0, 50.0, 100.0),
    ]
    exp = stages[min(int(elapsed // step_duration), len(stages) - 1)]
    assert (stage_name, fuel, water, steam) == exp
@pytest.mark.asyncio
async def test_run_demo_console_simulation_brief(monkeypatch, simulator):
    from app.runners.console import run_demo_console_simulation
    import asyncio
    times = [
        0.0,
        1.0,
        1.0,
        1.0,
    ]
    def mock_time():
        return times.pop(0)
    monkeypatch.setattr("time.time", mock_time)
    monkeypatch.setattr("builtins.print", lambda *args, **kwargs: None)
    task = asyncio.create_task(
        run_demo_console_simulation(simulator, update_interval=0.1, step_duration=10.0)
    )
    await asyncio.sleep(0.01)
    task.cancel()
    try:
        await task
    except asyncio.CancelledError:
        pass
    assert simulator.inputs.fuel_valve == 50.0
    assert simulator.inputs.feedwater_valve == 10.0
    assert simulator.inputs.steam_valve == 0.0
