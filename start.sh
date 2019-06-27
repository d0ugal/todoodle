#!/bin/bash
set -eux
set -o pipefail

docker run -d --name todos --network host --restart always -v /usr/share/hassio/share/todos:/share -it todos;