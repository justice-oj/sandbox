#!/bin/bash

echo "Download go packages..."
go get "github.com/getsentry/raven-go"
go get "github.com/docker/docker/pkg/reexec"
go get "github.com/satori/go.uuid"

echo "Compile binaries..."
go build -o bin/c_compiler src/cmd/c/compiler.go
go build -o bin/cpp_compiler src/cmd/cpp/compiler.go

go build -o bin/c_container src/cmd/c/container.go
go build -o bin/cpp_container src/cmd/cpp/container.go

echo "Enable automatically removing empty cgroups..."
echo 1 > /sys/fs/cgroup/cpu/notify_on_release
echo 1 > /sys/fs/cgroup/memory/notify_on_release
echo 1 > /sys/fs/cgroup/pids/notify_on_release

echo "/opt/justice-sandbox/bin/clean_cpu_cgroup.sh" > /sys/fs/cgroup/cpu/release_agent
echo "/opt/justice-sandbox/bin/clean_memory_cgroup.sh" > /sys/fs/cgroup/memory/release_agent
echo "/opt/justice-sandbox/bin/clean_pids_cgroup.sh" > /sys/fs/cgroup/pids/release_agent

echo "Done"