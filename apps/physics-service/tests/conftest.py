# app/tests/conftest.py
import os
import sys
import pytest

# make sure app package is importable when running pytest from repo root
ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
if ROOT not in sys.path:
    sys.path.insert(0, ROOT)

from app.core.simulator import BoilerSimulator


@pytest.fixture
def simulator():
    sim = BoilerSimulator()
    # set deterministic initial last_tick_time for tests
    sim.last_tick_time = 0.0
    return sim
