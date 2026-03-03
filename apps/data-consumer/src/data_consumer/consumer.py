import json
import logging
from paho.mqtt import client as mqtt_client
from influxdb_client import InfluxDBClient, Point
from influxdb_client.client.write_api import SYNCHRONOUS

from .settings import settings

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class BoilerDataConsumer:
    def __init__(self):
        self.influx_client = InfluxDBClient(url=settings.INFLUX_URL, token=settings.INFLUX_TOKEN, org=settings.INFLUX_ORG)
        self.write_api = self.influx_client.write_api(write_options=SYNCHRONOUS)
        
        self.mqtt_client = mqtt_client.Client(client_id="data-consumer-01")
        self.mqtt_client.on_connect = self.on_connect
        self.mqtt_client.on_message = self.on_message

    def on_connect(self, client, userdata, flags, rc):
        if rc == 0:
            client.subscribe(settings.MQTT_TOPIC)
            logger.info(f"Подписка на {settings.MQTT_TOPIC}")
        else:
            logger.error(f"MQTT ошибка {rc}")

    def on_message(self, client, userdata, msg):
        try:
            payload = json.loads(msg.payload.decode())
            point = Point("boiler_telemetry") \
                .tag("asset_id", payload["asset_id"]) \
                .field("temperature", payload["tags"]["temperature"]) \
                .field("pressure", payload["tags"]["pressure"]) \
                .field("fuel", payload["tags"]["fuel"]) \
                .field("drum_level", payload["tags"].get("drum_level", 500.0)) \
                .field("steam_flow", payload["tags"].get("steam_flow", 0.0)) \
                .time(payload["timestamp"])

            self.write_api.write(bucket="telemetry", record=point)
            logger.info(f"Записано: T={payload['tags']['temperature']:.1f}°C | P={payload['tags']['pressure']:.1f} bar | Fuel={payload['tags']['fuel']:.1f}%")
        except Exception as e:
            logger.error(f"Ошибка: {e}")

    def start(self):
        self.mqtt_client.connect(settings.MQTT_BROKER_HOST, 1883)
        self.mqtt_client.loop_forever()

def main():
    consumer = BoilerDataConsumer()
    try:
        consumer.start()
    except KeyboardInterrupt:
        logger.info("Остановлено")
    finally:
        consumer.influx_client.close()