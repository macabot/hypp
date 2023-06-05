#!/bin/bash

set -e

cwd=$(dirname "$0")
for example in "$cwd"/*/cmd/html; do
    echo "$example"
    go run "${example}/main.go"
done
