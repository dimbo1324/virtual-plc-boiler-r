import os
import sys
import pytest
from datetime import datetime
ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
if ROOT not in sys.path:
    sys.path.insert(0, ROOT)
from app.core.simulator import BoilerSimulator
@pytest.hookimpl(tryfirst=True)
def pytest_configure(config):
    if not config.option.htmlpath and config.pluginmanager.hasplugin("html"):
        timestamp = datetime.now().strftime("%Y-%m-%d_%H-%M-%S")
        report_name = f"report_{timestamp}.html"
        base_dir = config.rootpath
        out_dir = base_dir / "out" / "tests"
        out_dir.mkdir(parents=True, exist_ok=True)
        report_path = out_dir / report_name
        config.option.htmlpath = str(report_path)
        config.option.self_contained_html = True
@pytest.fixture
def simulator():
    sim = BoilerSimulator()
    sim.last_tick_time = 0.0
    return sim
