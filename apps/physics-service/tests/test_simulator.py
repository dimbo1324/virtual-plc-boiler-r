import pytest
from app.core.simulator import BoilerSimulator
from app.core.steam_table import SteamTable
from hypothesis import given, settings, HealthCheck
from hypothesis.strategies import floats
def test_set_controls_clamps(simulator):
    simulator.set_controls(fuel=-50.0, water=150.0, steam=50.0)
    assert simulator.inputs.fuel_valve == 0.0
    assert simulator.inputs.feedwater_valve == 100.0
    assert simulator.inputs.steam_valve == 50.0
def test_tick_temperature_heating_and_cooling(monkeypatch, simulator):
    simulator.outputs.furnace_temp = 20.0
    simulator.inputs.fuel_valve = 100.0
    monkeypatch.setattr(
        "time.time", lambda: 10.0
    )
    simulator.last_tick_time = 0.0
    simulator.tick()
    ambient_temp = 20.0
    max_furnace_temp = 1200.0
    heating_rate = 0.05
    cooling_rate = 0.02
    target_temp = ambient_temp + (
        max_furnace_temp * (simulator.inputs.fuel_valve / 100.0)
    )
    expected_change = (target_temp - 20.0) * heating_rate * 10.0
    assert simulator.outputs.furnace_temp == pytest.approx(20.0 + expected_change)
    simulator.outputs.furnace_temp = target_temp + 50.0
    monkeypatch.setattr("time.time", lambda: 20.0)
    simulator.last_tick_time = 10.0
    simulator.tick()
    expected_change2 = (
        (target_temp - (target_temp + 50.0)) * cooling_rate * 10.0
    )
    assert simulator.outputs.furnace_temp == pytest.approx(
        (target_temp + 50.0) + expected_change2
    )
def test_tick_pressure_and_flow_and_level(monkeypatch, simulator):
    simulator.outputs.furnace_temp = 120.0
    simulator.inputs.fuel_valve = (
        (simulator.outputs.furnace_temp - 20.0) / 1200.0
    ) * 100.0
    simulator.inputs.steam_valve = 50.0
    simulator.inputs.feedwater_valve = 10.0
    simulator.outputs.drum_level = 500.0
    monkeypatch.setattr("time.time", lambda: 2.0)
    simulator.last_tick_time = 0.0
    dt = 2.0
    ambient_temp = 20.0
    max_furnace_temp = 1200.0
    max_pressure = 100.0
    pressure_drop_rate = 0.5
    feedwater_rate = 2.5
    evaporation_rate = 2.0
    max_drum_level = 1000.0
    base_pressure = SteamTable.get_pressure(simulator.outputs.furnace_temp)
    pressure_loss = simulator.inputs.steam_valve * pressure_drop_rate * dt
    expected_pressure = max(0.0, base_pressure - pressure_loss)
    expected_flow = (expected_pressure / max_pressure) * simulator.inputs.steam_valve
    inflow = simulator.inputs.feedwater_valve * feedwater_rate * dt
    evaporation = (
        (simulator.outputs.furnace_temp / max_furnace_temp) * evaporation_rate * dt
    )
    expected_level = simulator.outputs.drum_level + inflow - evaporation
    expected_level = max(0.0, min(max_drum_level, expected_level))
    simulator.tick()
    assert simulator.outputs.steam_pressure == pytest.approx(expected_pressure)
    assert simulator.outputs.steam_flow == pytest.approx(expected_flow)
    assert simulator.outputs.drum_level == pytest.approx(expected_level)
def test_drum_level_clamped_to_bounds(monkeypatch, simulator):
    simulator.outputs.drum_level = 1.0
    simulator.outputs.furnace_temp = 1500.0
    simulator.inputs.feedwater_valve = 0.0
    monkeypatch.setattr("time.time", lambda: 1.0)
    simulator.last_tick_time = 0.0
    simulator.tick()
    assert simulator.outputs.drum_level >= 0.0
    simulator.outputs.drum_level = 1000.0 - 1.0
    simulator.inputs.feedwater_valve = 100.0
    monkeypatch.setattr("time.time", lambda: 2.0)
    simulator.last_tick_time = 1.0
    simulator.tick()
    assert simulator.outputs.drum_level <= 1000.0
def test_multiple_ticks_accumulate_changes(monkeypatch, simulator):
    simulator.outputs.furnace_temp = 20.0
    simulator.inputs.fuel_valve = 50.0
    simulator.last_tick_time = 0.0
    ambient_temp = 20.0
    max_furnace_temp = 1200.0
    heating_rate = 0.05
    for t in [1.0, 2.0, 3.0]:
        monkeypatch.setattr("time.time", lambda: t)
        simulator.tick()
    target_temp = ambient_temp + max_furnace_temp * (50.0 / 100.0)
    expected = 20.0
    for _ in range(3):
        expected += (target_temp - expected) * heating_rate * 1.0
    assert simulator.outputs.furnace_temp == pytest.approx(expected)
def test_zero_pressure_when_loss_exceeds_base(monkeypatch, simulator):
    simulator.outputs.furnace_temp = 100.0
    simulator.inputs.steam_valve = 100.0
    simulator.last_tick_time = 0.0
    monkeypatch.setattr("time.time", lambda: 10.0)
    simulator.tick()
    assert simulator.outputs.steam_pressure == 0.0
    assert simulator.outputs.steam_flow == 0.0
def test_get_state_timestamp_updates(simulator):
    initial_state = simulator.get_state()
    assert initial_state.timestamp > 0
    import time
    time.sleep(0.1)
    new_state = simulator.get_state()
    assert new_state.timestamp > initial_state.timestamp
@settings(suppress_health_check=[HealthCheck.function_scoped_fixture])
@given(fuel=floats(0, 100), water=floats(0, 100), steam=floats(0, 100))
def test_set_controls_clamps_hypothesis(simulator, fuel, water, steam):
    simulator.set_controls(fuel, water, steam)
    assert 0 <= simulator.inputs.fuel_valve <= 100
    assert 0 <= simulator.inputs.feedwater_valve <= 100
    assert 0 <= simulator.inputs.steam_valve <= 100
