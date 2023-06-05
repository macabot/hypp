#!/bin/bash

set -e

cwd=$(dirname "$0")
for example in "$cwd"/*/; do
    echo "$example"
    GOOS=js GOARCH=wasm go build -o "${example}cmd/js/public/main.wasm" "${example}cmd/js/main.go"
done
