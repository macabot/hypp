#!/bin/bash

set -e

cwd=$(dirname "$0")
for example in "$cwd"/*/; do
    echo "$example"
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "${example}cmd/js/public"
    GOOS=js GOARCH=wasm go build -o "${example}cmd/js/public/main.wasm" "${example}cmd/js/main.go"
done
