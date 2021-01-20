#!/bin/bash

terser client/client.js -o static/client.min.js --compress --mangle
go run server/*.go
