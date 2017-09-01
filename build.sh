#!/bin/bash

rm bin/*
go build -o bin/c_compiler src/sandbox/c/compiler.go
go build -o bin/cpp_compiler src/sandbox/cpp/compiler.go