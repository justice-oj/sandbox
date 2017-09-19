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

	keys := []string{"tasks", "cpu.cfs_quota_us"}
	values := []string{pid, "10000"}
	for k, v := range keys {
		path := filepath.Join(cgCPUPath, v)
		os.Stderr.WriteString("writing [" + values[k] + "] to file: " + path + "\n")
		if err := ioutil.WriteFile(path, []byte(values[k]), 0644); err != nil {
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

	keys := []string{"cgroup.procs", "pids.max"}
	values := []string{pid, "64"}
	for k, v := range keys {
		path := filepath.Join(cgPidPath, v)
		os.Stderr.WriteString("writing [" + values[k] + "] to file: " + path + "\n")
		if err := ioutil.WriteFile(path, []byte(values[k]), 0644); err != nil {
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

	keys := []string{"memory.kmem.limit_in_bytes", "tasks", "memory.limit_in_bytes", "memory.memsw.limit_in_bytes"}
	values := []string{"256m", pid, memory + "m", memory + "m"}
	for k, v := range keys {
		path := filepath.Join(cgMemoryPath, v)
		os.Stderr.WriteString("writing [" + values[k] + "] to file: " + path + "\n")
		if err := ioutil.WriteFile(path, []byte(values[k]), 0644); err != nil {
			os.Stderr.WriteString("write failed \n")
			return err
		}
		c, _ := ioutil.ReadFile(path)
		os.Stderr.WriteString("content of " + path + ": " + string(c))
	}

	os.Stderr.WriteString("memoryCGroup done \n")
	return nil
}
