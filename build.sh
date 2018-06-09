#!/bin/bash

echo -e "[DEBUG] GOPATH is \e[34m\e[1m${GOPATH}\033[0m"

echo "Compile binaries..."
mkdir -p "/opt/justice-sandbox/bin"
cd ${GOPATH}/src/github.com/justice-oj/sandbox/
go build -o /opt/justice-sandbox/bin/clike_compiler compiler.go
go build -o /opt/justice-sandbox/bin/clike_container container.go

echo "Enable automatically removing empty cgroups..."
echo 1 > /sys/fs/cgroup/cpu/notify_on_release
echo 1 > /sys/fs/cgroup/memory/notify_on_release
echo 1 > /sys/fs/cgroup/pids/notify_on_release

echo "${GOPATH}/src/github.com/justice-oj/sandbox/scripts/clean_cpu_cgroup.sh" > /sys/fs/cgroup/cpu/release_agent
echo "${GOPATH}/src/github.com/justice-oj/sandbox/scripts/clean_memory_cgroup.sh" > /sys/fs/cgroup/memory/release_agent
echo "${GOPATH}/src/github.com/justice-oj/sandbox/scripts/clean_pids_cgroup.sh" > /sys/fs/cgroup/pids/release_agent

echo "Done!"
