#!/usr/bin/env bash

set -e
mkdir -p build
echo "" > build/coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -v -race -coverprofile=build/profile.out -covermode=atomic $d
    if [ -f build/profile.out ]; then
        cat build/profile.out >> build/coverage.txt
        rm build/profile.out
    fi
done
