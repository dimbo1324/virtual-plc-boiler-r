import time
from app import settings
from app.core.state import BoilerInputs, BoilerOutputs, BoilerState


class BoilerSimulator:
    def __init__(self):
        self.inputs = BoilerInputs()
        self.outputs = BoilerOutputs()
        self.last_tick_time = time.time()

    def set_controls(self, fuel: float, water: float, steam: float):
        # Accepts control commands (0 - 100%)
        # Limit the values to between 0 and 100 for safety reasons
        self.inputs.fuel_valve = max(0, 0, min(100.0, fuel))
        self.inputs.feedwater_valve = max(0, 0, min(100.0, water))
        self.inputs.steam_valve = max(0, 0, min(100.0, steam))

    def get_state(self) -> BoilerState:
        # Returns a snapshot of the current system state
        return BoilerState(
            timestamp=time.time(), inputs=self.inputs, outputs=self.outputs
        )

    def tick(self):
        # -----
        # Recalculates the physics for one time step
        # Called in a loop
        # -----
        current_time = time.time()
        # dt - how many seconds have passed since the last calc
        dt = current_time - self.last_tick_time
        self.last_tick_time = current_time

        # -----
        # ----- 1. Temperature calculation (inertion) -----
        # -----
        # Target temperature depends on the fuel opening (linearly)
        target_temp = 20.0 + (
            settings.MAX_FURNACE_TEMP * (self.inputs.fuel_valve / 100.0)
        )
        # If the current temperature is lower than the target temperature, we heat up; else, we cool down
        if self.outputs.furnace_temp < target_temp:
            change = (
                (target_temp - self.outputs.furnace_temp) * settings.HEATING_RATE * dt
            )
        else:
            change = (
                (target_temp - self.outputs.furnace_temp) * settings.COOLING_RATE * dt
            )
        self.outputs.furnace_temp += change

        # -----
        # ----- 2. Pressure calculation -----
        # -----
        # Base pressure depends on the temperature (at the moment P ~ T), TODO: make more real calculation of press-temperature depending
        # P_base = (Temp / Max_Temp) * Max_Pressure
        base_pressure = (
            self.outputs.furnace_temp / settings.MAX_FURNACE_TEMP
        ) * settings.MAX_PRESSURE
        pressure_loss = self.inputs.steam_valve * settings.PRESSURE_DROP_RATE * dt
        # The final pressure (cannot be less than 0)
        self.outputs.steam_pressure = max(0.0, base_pressure - pressure_loss)

        # -----
        # ----- 3. Steam flow calculation -----
        # -----
        # Steam flow is if pressure is and valve open
        self.outputs.steam_flow = (
            self.outputs.steam_pressure / settings.MAX_PRESSURE
        ) * self.inputs.steam_valve

        # -----
        # ----- 4. Water level calculation -----
        # -----
        # Water inlet (from the supply valve)
        inflow = self.inputs.feedwater_valve * settings.FEEDWATER_RATE * dt
        evaporation = (
            (self.outputs.furnace_temp / settings.MAX_FURNACE_TEMP)
            * settings.EVAPORATION_RATE
            * dt
        )
        self.outputs.drum_level += inflow - evaporation
        self.outputs.drum_level = max(
            0.0, min(settings.MAX_DRUM_LEVEL, self.outputs.drum_level)
        )
