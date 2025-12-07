#!/bin/bash

ROOT_DIR="$(cd $(dirname ${BASH_SOURCE[0]}) && cd .. && pwd)"

(
    go build -o ./bin/api ./cmd/api/*.go
)