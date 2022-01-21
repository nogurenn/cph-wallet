#!/bin/bash -e

rm -rf bin/wallet

ldflags=""
if [ "$STAGE" == "production" ]; then
  ldflags="-s -w"
fi

GOOS=linux CGO_ENABLED=0 go build -ldflags="$ldflags" -o bin/wallet main.go