#!/usr/bin/env bash

set -euf -o pipefail

# Get Data : full path
gnmic --config car.yml get --path car:/

# Get Data : use model
gnmic --config car.yml get --model car --path ""

# Set Data
gnmic --config car.yml set --update  car:/:::json:::'{"speed":300}'

# Subscribe
timeout 5s \
  gnmic --config car.yml sub --model car --path "" --sample-interval 1s --heartbeat-interval 2s || true

# Subscribe to just tire metrics : use model
timeout 5s \
  gnmic --config car.yml sub --mode once --model car --path "/tire" || true

# Subscribe to just tire metrics : full path
timeout 5s \
  gnmic --config car.yml sub --mode once --path "car:/tire" || true
