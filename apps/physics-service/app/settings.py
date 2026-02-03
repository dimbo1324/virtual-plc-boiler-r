from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    # Application-wide configuration settings loaded from environment variables or default values.

    # Environment / ambient conditions
    AMBIENT_TEMP: float = 20.0  # Base ambient temperature (°C)
    MAX_FURNACE_TEMP: float = 1200.0  # Maximum achievable furnace temperature (°C)

    # Heating/cooling dynamics
    HEATING_RATE: float = 0.05  # Temperature increase factor per second
    COOLING_RATE: float = 0.02  # Temperature decrease factor per second

    # Pressure related parameters
    MAX_PRESSURE: float = 100.0  # Maximum possible steam pressure (bar)
    PRESSURE_DROP_RATE: float = (
        0.5  # Pressure loss per % of open steam valve per second
    )

    # Drum / water level parameters
    MAX_DRUM_LEVEL: float = 1000.0  # Maximum drum water level (mm)
    EVAPORATION_RATE: float = 2.0  # Water evaporation rate (mm/s at full power)
    FEEDWATER_RATE: float = 2.5  # Feedwater inflow rate (mm/s at 100% valve)

    class Config:
        env_prefix = "BOILER_"  # All env variables start with BOILER_


# Global singleton instance of settings (used everywhere in the app)
settings = Settings()
