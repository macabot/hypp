#!/bin/bash

set -e

for example in */cmd/html; do
    echo "$example"
    go run "${example}/main.go"
done
