#!/bin/bash

rm bin/*

go build -o bin/c_compiler src/sandbox/c/Compiler.go
go build -o bin/cpp_compiler src/sandbox/cpp/Compiler.go

go build -o bin/c_container src/sandbox/c/Container.go
go build -o bin/cpp_container src/sandbox/cpp/Container.go