package cgroup

import (
	"io/ioutil"
	"path/filepath"
	"os"
	"os/exec"
)

func CPUInit(pid, cgCPUPath string) error {
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

func MemoryInit(pid, cgMemoryPath, memory string) error {
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

func Cleanup(path string) error {
	cmd := exec.Command("rmdir", path)
	return cmd.Run()
}
