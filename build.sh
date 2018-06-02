#!/bin/bash

echo -e "[DEBUG] GOPATH is \e[34m\e[1m${GOPATH}\033[0m"

echo "Download go packages..."
go get "github.com/getsentry/raven-go"
go get "github.com/docker/docker/pkg/reexec"
go get "github.com/satori/go.uuid"

echo "Compile binaries..."
mkdir -p "/opt/justice-sandbox/bin"
go build -o /opt/justice-sandbox/bin/clike_compiler ${GOPATH}/src/github.com/justice-oj/sandbox/src/cmd/clike/compiler.go
go build -o /opt/justice-sandbox/bin/clike_container ${GOPATH}/src/github.com/justice-oj/sandbox/src/cmd/clike/container.go

echo "Enable automatically removing empty cgroups..."
echo 1 > /sys/fs/cgroup/cpu/notify_on_release
echo 1 > /sys/fs/cgroup/memory/notify_on_release
echo 1 > /sys/fs/cgroup/pids/notify_on_release

echo "${GOPATH}/src/github.com/justice-oj/sandbox/bin/clean_cpu_cgroup.sh" > /sys/fs/cgroup/cpu/release_agent
echo "${GOPATH}/src/github.com/justice-oj/sandbox/bin/clean_memory_cgroup.sh" > /sys/fs/cgroup/memory/release_agent
echo "${GOPATH}/src/github.com/justice-oj/sandbox/bin/clean_pids_cgroup.sh" > /sys/fs/cgroup/pids/release_agent

echo "Done!"