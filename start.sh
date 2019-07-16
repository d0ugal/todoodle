#!/bin/bash
set -eux
set -o pipefail

docker kill todoodle || true;
docker rm todoodle || true;

docker run -d \
  --name todoodle --network host --restart always \
  -e DB="todos:todos@tcp(192.168.1.200:3306)/todos?charset=utf8mb4&parseTime=True&loc=Local" \
  -e MQTT_BROKER="tcp://192.168.1.200:1883" \
  -e MQTT_USERNAME="hassio" \
  -e MQTT_PASSWORD="hassio" \
  todoodle;