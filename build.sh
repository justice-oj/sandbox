#!/bin/bash

echo "Compile binaries..."
mkdir -p "/opt/justice-sandbox/bin"
go build -o /opt/justice-sandbox/bin/clike_compiler compiler.go
go build -o /opt/justice-sandbox/bin/clike_container container.go

echo "Enable automatically removing empty cgroups..."
echo 1 > /sys/fs/cgroup/cpu/notify_on_release
echo 1 > /sys/fs/cgroup/memory/notify_on_release
echo 1 > /sys/fs/cgroup/pids/notify_on_release

echo "${PWD}/scripts/clean_cpu_cgroup.sh" > /sys/fs/cgroup/cpu/release_agent
echo "${PWD}/scripts/clean_memory_cgroup.sh" > /sys/fs/cgroup/memory/release_agent
echo "${PWD}/scripts/clean_pids_cgroup.sh" > /sys/fs/cgroup/pids/release_agent

echo "Done!"
