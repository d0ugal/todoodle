#!/bin/bash
set -eux
set -o pipefail

docker run -d \
  --name todoodle --network host --restart always \
  -e MYSQL_URL="todos:todos@tcp(192.168.1.200:3306)/todos?charset=utf8mb4&parseTime=True&loc=Local" \
  todoodle;