import configparser
import socket
import os

KAFKA_HOST = os.environ.get("KAFKA_HOST", "localhost:9092")

config = configparser.ConfigParser(interpolation=None)

config.read("config.ini")

KAFKA_CONF = {
    "bootstrap.servers": KAFKA_HOST,
}

TOPIC_NAME = config.get("KAFKA", "topic_name")
SLEEP_INTERVAL = config.getfloat("PRODUCER", "sleep_interval")

logging_config = {
    "filename": config.get("LOGGING", "file_name"),
    "level": config.get("LOGGING", "log_level"),
    "format": config.get("LOGGING", "format"),
}
