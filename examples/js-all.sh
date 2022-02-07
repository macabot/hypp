#!/bin/bash

set -e

for example in */; do
    echo "$example"
    GOOS=js GOARCH=wasm gotip build -o "${example}cmd/js/public/main.wasm" "${example}cmd/js/main.go"
done
