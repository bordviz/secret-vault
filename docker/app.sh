#!/bin/bash

go run cmd/migrations/main.go --action=up

go build -o build/main cmd/vault/main.go

./build/main