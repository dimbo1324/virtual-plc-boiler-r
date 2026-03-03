from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    INFLUX_URL: str
    INFLUX_TOKEN: str
    INFLUX_ORG: str
    MQTT_BROKER_HOST: str = "mosquitto"
    MQTT_TOPIC: str = "v1/devices/boiler/telemetry"

    class Config:
        env_file = ".env"
        extra = "ignore"

settings = Settings()