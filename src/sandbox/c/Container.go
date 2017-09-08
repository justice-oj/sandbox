package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"syscall"
	"github.com/docker/docker/pkg/reexec"
	"github.com/satori/go.uuid"
	"path/filepath"
	"github.com/getsentry/raven-go"
	"../../models"
	"../../config"
	"../../common/cgroup"
	"../../common/namespace"
	"flag"
	"strings"
	"bytes"
	"time"
	"strconv"
)

func init() {
	raven.SetDSN(config.SENTRY_DSN)
	reexec.Register("justiceInit", justiceInit)
	if reexec.Init() {
		os.Exit(models.CODE_OK)
	}
}

func justiceInit() {
	newRootPath := os.Args[1]
	input := os.Args[2]
	expected := os.Args[3]
	timeout, _ := strconv.ParseInt(os.Args[4], 10, 32)

	if err := &namespace.MountProc(newRootPath); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(models.CODE_INIT_CONTAINER_FAILED)
	}

	if err := &namespace.PivotRoot(newRootPath); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(models.CODE_INIT_CONTAINER_FAILED)
	}

	if err := syscall.Sethostname([]byte("justice")); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(models.CODE_INIT_CONTAINER_FAILED)
	}

	justiceRun(input, expected, int32(timeout))
}

func justiceRun(input, expected string, timeout int32) {
	// for c programs, compiled binary with name [Main] will be located in "/"
	var o, e bytes.Buffer
	cmd := exec.Command("/Main")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &o
	cmd.Stderr = &e
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = []string{"PS1=[justice] # "}

	time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	// ms
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	if err := cmd.Run(); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerRunTimeError"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		return
	}
	endTime := time.Now().UnixNano() / int64(time.Millisecond)

	if e.Len() > 0 {
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		return
	}

	// MB
	memory := cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss / 1024

	output := strings.TrimSpace(o.String())
	if output == expected {
		result, _ := json.Marshal(models.GetAccepptedTaskResult(endTime-startTime, memory))
		os.Stdout.Write(result)
	} else {
		result, _ := json.Marshal(models.GetWrongAnswerTaskResult(input, output, expected))
		os.Stdout.Write(result)
	}
}

func main() {
	basedir := flag.String("basedir", "/tmp", "basedir of tmp C binary")
	input := flag.String("input", "<input>", "test case input")
	expected := flag.String("expected", "<expected>", "test case expected")
	timeout := flag.String("timeout", "2000", "timeout in milliseconds")
	memory := flag.String("memory", "64", "memory limitation in MB")
	flag.Parse()

	pid, containerID := os.Getpid(), uuid.NewV4().String()
	cgCPUPath := filepath.Join("/sys/fs/cgroup/cpu/", containerID)
	cgMemoryPath := filepath.Join("/sys/fs/cgroup/memory/", containerID)

	// CPU
	if err := &cgroup.CPUInit(string(pid), cgCPUPath); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		return
	}

	// MEMORY
	if err := &cgroup.MemoryInit(string(pid), cgMemoryPath, *memory); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		return
	}

	cmd := reexec.Command("justiceInit", *basedir, *input, *expected, *timeout)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Run(); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerRunTimeError"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
	}

	if err := &cgroup.Cleanup(cgCPUPath); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerCleanupError"})
	}

	if err := &cgroup.Cleanup(cgMemoryPath); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerCleanupError"})
	}

	os.Exit(models.CODE_OK)
}
