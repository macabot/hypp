#!/bin/bash

set -e

for example in */cmd/html; do
    echo "$example"
    gotip run "${example}/main.go"
done
