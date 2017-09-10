package cgroup

import (
	"io/ioutil"
	"path/filepath"
	"os"
	"os/exec"
)

const (
	cgCPUPathPrefix    = "/sys/fs/cgroup/cpu/"
	cgMemoryPathPrefix = "/sys/fs/cgroup/memory/"
)

func InitCGroup(pid, containerID, memory string) error {
	if err := cpuCGroup(pid, containerID); err != nil {
		return err
	}

	if err := memoryCGroup(pid, containerID, memory); err != nil {
		return err
	}

	return nil
}

func cpuCGroup(pid, containerID string) error {
	cgCPUPath := filepath.Join(cgCPUPathPrefix, containerID)

	// add sub cgroup system
	if err := os.Mkdir(cgCPUPath, 0755); err != nil {
		return err
	}

	// add current pid to cgroup cpu
	if err := ioutil.WriteFile(filepath.Join(cgCPUPath, "/tasks"), []byte(pid), 0755); err != nil {
		return err
	}

	// cpu usage max up to 2%
	if err := ioutil.WriteFile(filepath.Join(cgCPUPath, "/cpu.cfs_quota_us"), []byte("2000"), 0755); err != nil {
		return err
	}

	return nil
}

func memoryCGroup(pid, containerID, memory string) error {
	cgMemoryPath := filepath.Join(cgMemoryPathPrefix, containerID)

	// add sub cgroup system
	if err := os.Mkdir(cgMemoryPath, 0755); err != nil {
		return err
	}

	// add current pid to cgroup memory
	if err := ioutil.WriteFile(filepath.Join(cgMemoryPath, "/tasks"), []byte(string(pid)), 0755); err != nil {
		return err
	}

	// set memory usage limitation
	if err := ioutil.WriteFile(filepath.Join(cgMemoryPath, "/memory.limit_in_bytes"), []byte(memory+"m"), 0755); err != nil {
		return err
	}

	return nil
}

func Cleanup(containerID string) error {
	cleanCPUCommand := exec.Command("rmdir", cgCPUPathPrefix + containerID)
	if cpuErr := cleanCPUCommand.Run(); cpuErr != nil {
		return cpuErr
	}

	cleanMemoryCommand := exec.Command("rmdir", cgMemoryPathPrefix + containerID)
	if cpuErr := cleanMemoryCommand.Run(); cpuErr != nil {
		return cpuErr
	}

	return nil
}
