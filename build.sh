#!/bin/bash

rm bin/*

echo "Download go packages..."
go get "github.com/getsentry/raven-go"
go get "github.com/docker/docker/pkg/reexec"
go get "github.com/satori/go.uuid"

echo "Compile binaries..."
go build -o bin/c_compiler src/sandbox/c/Compiler.go
go build -o bin/cpp_compiler src/sandbox/cpp/Compiler.go

go build -o bin/c_container src/sandbox/c/Container.go
go build -o bin/cpp_container src/sandbox/cpp/Container.go

echo "Done"