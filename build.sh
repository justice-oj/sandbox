#!/bin/bash

rm bin/*

echo "Download go packages..."
go get "github.com/getsentry/raven-go"
go get "github.com/docker/docker/pkg/reexec"
go get "github.com/satori/go.uuid"

echo "Compile binaries..."
go build -o bin/c_compiler src/cmd/c/compiler.go
go build -o bin/cpp_compiler src/cmd/cpp/compiler.go

go build -o bin/c_container src/cmd/c/container.go
go build -o bin/cpp_container src/cmd/cpp/container.go

echo "Done"