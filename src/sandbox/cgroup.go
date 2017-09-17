package sandbox

import (
	"path/filepath"
	"os"
	"io/ioutil"
)

const (
	cgCPUPathPrefix    = "/sys/fs/cgroup/cpu/"
	cgPidPathPrefix    = "/sys/fs/cgroup/pids/"
	cgMemoryPathPrefix = "/sys/fs/cgroup/memory/"
)

func InitCGroup(pid, containerID, memory string) error {
	os.Stderr.WriteString("InitCGroup starting...\n")

	dirs := []string{
		filepath.Join(cgCPUPathPrefix, containerID),
		filepath.Join(cgPidPathPrefix, containerID),
		filepath.Join(cgMemoryPathPrefix, containerID),
	}

	for _, dir := range dirs {
		os.Stderr.WriteString("os.MkdirAll " + dir + "\n")
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Stderr.WriteString("os.MkdirAll " + dir + " failed \n")
			return err
		}
	}

	if err := cpuCGroup(pid, containerID); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("cpuCGroup failed...\n")
		return err
	}

	if err := pidCGroup(pid, containerID); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("pidCGroup failed...\n")
		return err
	}

	if err := memoryCGroup(pid, containerID, memory); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("memoryCGroup failed...\n")
		return err
	}

	os.Stderr.WriteString("InitCGroup done...\n")
	return nil
}

// https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt
func cpuCGroup(pid, containerID string) error {
	cgCPUPath := filepath.Join(cgCPUPathPrefix, containerID)

	configs := map[string]string{
		"tasks":            pid,
		"cpu.cfs_quota_us": "2000",
	}

	for file, content := range configs {
		path := filepath.Join(cgCPUPath, file)
		os.Stderr.WriteString("writing [" + content + "] to file: " + path + "\n")
		if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
			os.Stderr.WriteString("write failed \n")
			return err
		}
		c, _ := ioutil.ReadFile(path)
		os.Stderr.WriteString("content of " + path + ": " + string(c))
	}

	os.Stderr.WriteString("cpuCGroup done\n")
	return nil
}

// https://www.kernel.org/doc/Documentation/cgroup-v1/pids.txt
func pidCGroup(pid, containerID string) error {
	cgPidPath := filepath.Join(cgPidPathPrefix, containerID)

	configs := map[string]string{
		"cgroup.procs": pid,
		"pids.max":     "2",
	}

	for file, content := range configs {
		path := filepath.Join(cgPidPath, file)
		os.Stderr.WriteString("writing [" + content + "] to file: " + path + "\n")
		if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
			os.Stderr.WriteString("write failed \n")
			return err
		}
		c, _ := ioutil.ReadFile(path)
		os.Stderr.WriteString("content of " + path + ": " + string(c))
	}

	os.Stderr.WriteString("pidCGroup done \n")
	return nil
}

// https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
func memoryCGroup(pid, containerID, memory string) error {
	cgMemoryPath := filepath.Join(cgMemoryPathPrefix, containerID)

	configs := map[string]string{
		"tasks":                       pid,
		"memory.limit_in_bytes":       memory + "m",
		"memory.memsw.limit_in_bytes": memory + "m",
		"memory.kmem.limit_in_bytes":  "8m",
		"memory.swappiness":           "0",
	}

	for file, content := range configs {
		path := filepath.Join(cgMemoryPath, file)
		os.Stderr.WriteString("writing [" + content + "] to file: " + path + "\n")
		if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
			os.Stderr.WriteString("write failed \n")
			return err
		}
		c, _ := ioutil.ReadFile(path)
		os.Stderr.WriteString("content of " + path + ": " + string(c))
	}

	os.Stderr.WriteString("memoryCGroup done \n")
	return nil
}

func CleanupCGroup(containerID string) error {
	os.Stderr.WriteString("CleanupCGroup starting...\n")

	dirs := []string{
		filepath.Join(cgCPUPathPrefix, containerID),
		filepath.Join(cgPidPathPrefix, containerID),
		filepath.Join(cgMemoryPathPrefix, containerID),
	}

	for _, dir := range dirs {
		os.Stderr.WriteString("os.Remove " + dir + "\n")
		if err := os.Remove(dir); err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Stderr.WriteString("os.Remove " + dir + " failed \n")
			return err
		}
	}

	os.Stderr.WriteString("CleanupCGroup done \n")
	return nil
}
