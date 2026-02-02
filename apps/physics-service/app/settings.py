from pydantic_settings import BaseSettings


# TODO: comments and docs
class Settings(BaseSettings):
    MAX_FURNACE_TEMP: float = 1200.0
    HEATING_RATE: float = 0.05
    COOLING_RATE: float = 0.02
    MAX_PRESSURE: float = 100.0
    PRESSURE_DROP_RATE: float = 0.5
    MAX_DRUM_LEVEL: float = 1000.0
    EVAPORATION_RATE: float = 2.0
    FEEDWATER_RATE: float = 2.5

    class Config:
        env_prefix = "BOILER_"


settings = Settings()
