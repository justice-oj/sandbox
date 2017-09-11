package cgroup

import (
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
	taskCmd := exec.Command("/usr/bin/echo", pid, ">", filepath.Join(cgCPUPath, "/tasks"))
	if err := taskCmd.Run(); err != nil {
		return err
	}

	// cpu usage max up to 2%
	quotaCmd := exec.Command("/usr/bin/echo", "2000", ">", filepath.Join(cgCPUPath, "/cpu.cfs_quota_us"))
	if err := quotaCmd.Run(); err != nil {
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
	taskCmd := exec.Command("/usr/bin/echo", pid, ">", filepath.Join(cgMemoryPath, "/tasks"))
	if err := taskCmd.Run(); err != nil {
		return err
	}

	// set memory usage limitation
	quotaCmd := exec.Command("/usr/bin/echo", memory+"m", ">", filepath.Join(cgMemoryPath, "/memory.limit_in_bytes"))
	if err := quotaCmd.Run(); err != nil {
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
