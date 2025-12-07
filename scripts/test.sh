#!/bin/bash

go test -v -coverprofile .local/coverage/.coverage.out "$@"
go tool cover -html .local/coverage/.coverage.out -o .local/coverage/.coverage.html