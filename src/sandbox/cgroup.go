package sandbox

import (
	"path/filepath"
	"os/exec"
	"os"
)

const (
	cgCPUPathPrefix    = "/sys/fs/cgroup/cpu/"
	cgPidPathPrefix    = "/sys/fs/cgroup/pids/"
	cgMemoryPathPrefix = "/sys/fs/cgroup/memory/"
)

func InitCGroup(pid, containerID, memory string) error {
	os.Stderr.WriteString("InitCGroup starting...\n")

	if err := cpuCGroup(pid, containerID); err != nil {
		os.Stderr.WriteString("cpuCGroup failed...\n")
		return err
	}

	if err := pidCGroup(pid, containerID); err != nil {
		os.Stderr.WriteString("pidCGroup failed...\n")
		return err
	}

	if err := memoryCGroup(pid, containerID, memory); err != nil {
		os.Stderr.WriteString("memoryCGroup failed...\n")
		return err
	}

	os.Stderr.WriteString("InitCGroup done...\n")
	return nil
}

// https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt
func cpuCGroup(pid, containerID string) error {
	cgCPUPath := filepath.Join(cgCPUPathPrefix, containerID)
	os.Stderr.WriteString("cgCPUPath: " + cgCPUPath + "\n")

	// add sub cgroup system
	mkdirCmd := exec.Command("/usr/bin/mkdir", cgCPUPath)
	if err := mkdirCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("mkdirCmd failed \n")
		return err
	}

	// add current pid to cgroup cpu
	taskCmd := exec.Command("/usr/bin/echo", pid, ">", filepath.Join(cgCPUPath, "/tasks"))
	if err := taskCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("taskCmd failed \n")
		return err
	}

	// limit a group to 2% of 1 CPU
	quotaCmd := exec.Command("/usr/bin/echo", "2000", ">", filepath.Join(cgCPUPath, "/cpu.cfs_quota_us"))
	if err := quotaCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("quotaCmd failed \n")
		return err
	}

	os.Stderr.WriteString("cpuCGroup done\n")
	return nil
}

// https://www.kernel.org/doc/Documentation/cgroup-v1/pids.txt
func pidCGroup(pid, containerID string) error {
	cgPidPath := filepath.Join(cgPidPathPrefix, containerID)
	os.Stderr.WriteString("cgPidPath: " + cgPidPath + "\n")

	// add sub cgroup system
	mkdirCmd := exec.Command("/usr/bin/mkdir", cgPidPath)
	if err := mkdirCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("mkdirCmd failed \n")
		return err
	}

	// add current pid to cgroup pids
	taskCmd := exec.Command("/usr/bin/echo", pid, ">", filepath.Join(cgPidPath, "/cgroup.procs"))
	if err := taskCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("taskCmd failed \n")
		return err
	}

	// max pids limitation
	quotaCmd := exec.Command("/usr/bin/echo", "2", ">", filepath.Join(cgPidPath, "/pids.max"))
	if err := quotaCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("quotaCmd failed \n")
		return err
	}

	os.Stderr.WriteString("pidCGroup done\n")
	return nil
}

// https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
func memoryCGroup(pid, containerID, memory string) error {
	cgMemoryPath := filepath.Join(cgMemoryPathPrefix, containerID)
	os.Stderr.WriteString("cgMemoryPath: " + cgMemoryPath + "\n")

	// add sub cgroup system
	mkdirCmd := exec.Command("/usr/bin/mkdir", cgMemoryPath)
	if err := mkdirCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("mkdirCmd failed \n")
		return err
	}

	// add current pid to cgroup memory
	taskCmd := exec.Command("/usr/bin/echo", pid, ">", filepath.Join(cgMemoryPath, "/tasks"))
	if err := taskCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("taskCmd failed \n")
		return err
	}

	// set memory usage limitation
	memoryCmd := exec.Command("/usr/bin/echo", memory+"m", ">", filepath.Join(cgMemoryPath, "/memory.limit_in_bytes"))
	if err := memoryCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("memoryCmd failed \n")
		return err
	}

	memswCmd := exec.Command("/usr/bin/echo", memory+"m", ">", filepath.Join(cgMemoryPath, "/memory.memsw.limit_in_bytes"))
	if err := memswCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("memswCmd failed \n")
		return err
	}

	swappinessCmd := exec.Command("/usr/bin/echo", "0", ">", filepath.Join(cgMemoryPath, "/memory.swappiness"))
	if err := swappinessCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("swappinessCmd failed \n")
		return err
	}

	kernelMemoryCmd := exec.Command("/usr/bin/echo", memory+"m", ">", filepath.Join(cgMemoryPath, "/memory.kmem.limit_in_bytes"))
	if err := kernelMemoryCmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString("kernelMemoryCmd failed \n")
		return err
	}

	os.Stderr.WriteString("memoryCGroup done\n")
	return nil
}

func CleanupCGroup(containerID string) error {
	os.Stderr.WriteString("CleanupCGroup starting...\n")

	cleanCPUCommand := exec.Command("rmdir", filepath.Join(cgCPUPathPrefix, containerID))
	os.Stderr.WriteString("rmdir " + filepath.Join(cgCPUPathPrefix, containerID) + "\n")
	if err := cleanCPUCommand.Run(); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("cleanCPUCommand failed \n")
		return err
	}

	cleanPidCommand := exec.Command("rmdir", filepath.Join(cgPidPathPrefix, containerID))
	os.Stderr.WriteString("rmdir " + filepath.Join(cgPidPathPrefix, containerID) + "\n")
	if err := cleanPidCommand.Run(); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("cleanPidCommand failed \n")
		return err
	}

	cleanMemoryCommand := exec.Command("rmdir", filepath.Join(cgMemoryPathPrefix, containerID))
	os.Stderr.WriteString("rmdir " + filepath.Join(cgMemoryPathPrefix, containerID) + "\n")
	if err := cleanMemoryCommand.Run(); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("cleanMemoryCommand failed \n")
		return err
	}

	os.Stderr.WriteString("CleanupCGroup done\n")
	return nil
}
