from dataclasses import dataclass
@dataclass
class BoilerInputs:
    # Control actions
    fuel_valve: float = 0.0  # Feul valve
    feedwater_valve: float = 0.0  # Feed water valve
    steam_valve: float = 0.0  # Steam valve (out; to customers)
@dataclass
class BoilerOutputs:
    # Sensor readings
    furnace_temp: float = 20.0  # Temperature in the furnace (°C)
    steam_pressure: float = 0.0  # Steam pressure (Bar)
    drum_level: float = 500.0  # Drum level (mm)
    steam_flow: float = 0.0  # Steam flow rate (kg/s or arbitrary units)
@dataclass
class BoilerState:
    # The complete state of the system
    timestamp: float  # Time simulation
    inputs: BoilerInputs
    outputs: BoilerOutputs
