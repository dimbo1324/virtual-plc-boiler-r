from app.core.state import BoilerInputs, BoilerOutputs, BoilerState
from app.settings import Settings
import time
def test_dataclasses_default_values():
    inp = BoilerInputs()
    out = BoilerOutputs()
    assert inp.fuel_valve == 0.0
    assert out.furnace_temp == 20.0
    state = BoilerState(timestamp=time.time(), inputs=inp, outputs=out)
    assert isinstance(state.timestamp, float)
    assert state.inputs is inp
    assert state.outputs is out
def test_settings_from_env(monkeypatch):
    monkeypatch.setenv("BOILER_HEATING_RATE", "0.5")
    settings = Settings()
    assert settings.HEATING_RATE == 0.5
