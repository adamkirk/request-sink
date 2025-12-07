#!/bin/bash

REPO_ROOT=$(git rev-parse --show-toplevel)

(
    cd $REPO_ROOT
    VIOLATIONS=$(gofmt -l .)

    if [ ! -z "$VIOLATIONS" ]; then
        echo "Violations found in the following files, run gofmt (or make fmt) locally!"
        echo ""
        echo $VIOLATIONS
        exit 1
    fi
)