#!/bin/bash
set -eux
set -o pipefail

docker run -d --name todoodle --network host --restart always -v /usr/share/hassio/share/todoodle:/share -it todoodle;