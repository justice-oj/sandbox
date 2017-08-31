#!/bin/bash

rm bin/*
go build -o bin/c_compiler src/c/compiler.go
go build -o bin/cpp_compiler src/cpp/compiler.go