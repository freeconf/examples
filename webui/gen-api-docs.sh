#!/usr/bin/env sh
set -eu

FC_YANG="go run github.com/freeconf/yang/cmd/fc-yang"
YPATH=../car
MODULE=car

# HTML API Docs
${FC_YANG} doc -f dot -module ${MODULE} -ypath ${YPATH} > ${MODULE}.dot
dot -Tsvg ${MODULE}.dot -o ${MODULE}.svg
${FC_YANG} doc -f html -module ${MODULE} -title "Car REST API" -ypath ${YPATH} > ${MODULE}-api.html

# Markdown API Docs
${FC_YANG} doc -f md -module ${MODULE} -ypath ${YPATH} > ${MODULE}-api.md
