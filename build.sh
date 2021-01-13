#!/bin/bash
tinygo build -o static/client.wasm -target wasm -no-debug -ldflags "-s" ./client/client.go
