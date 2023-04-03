#!/usr/bin/env sh
set -eu

openssl req -new -x509 -sha256 -newkey rsa:2048 -nodes \
    -days 999 -keyout server.key -out server.crt \
    -subj "/C=US/ST=Self/L=Self/O=Self/OU=Self/CN=localhost"
