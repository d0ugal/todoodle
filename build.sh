#!/bin/bash
set -eux
set -o pipefail

docker build . -t todoodle;