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
	"flag"
	"strings"
	"bytes"
	"time"
	"strconv"
	"io/ioutil"
)

func init() {
	raven.SetDSN(config.SENTRY_DSN)
	reexec.Register("justiceInit", justiceInit)
	if reexec.Init() {
		os.Exit(models.CODE_OK)
	}
}

func pivotRoot(newRoot string) error {
	putOld := filepath.Join(newRoot, "/.pivot_root")

	// bind mount new_root to itself - this is a slight hack needed to satisfy requirement (2)
	//
	// The following restrictions apply to new_root and put_old:
	// 1.  They must be directories.
	// 2.  new_root and put_old must not be on the same filesystem as the current root.
	// 3.  put_old must be underneath new_root, that is, adding a nonzero
	//     number of /.. to the string pointed to by put_old must yield the same directory as new_root.
	// 4.  No other filesystem may be mounted on put_old.
	if err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}

	// create put_old directory
	if err := os.MkdirAll(putOld, 0700); err != nil {
		return err
	}

	// call pivotRoot
	if err := syscall.PivotRoot(newRoot, putOld); err != nil {
		return err
	}

	// Note that this also applies to the calling process: pivotRoot() may
	// or may not affect its current working directory.  It is therefore
	// recommended to call chdir("/") immediately after pivotRoot().
	if err := os.Chdir("/"); err != nil {
		return err
	}

	// umount put_old, which now lives at /.pivot_root
	putOld = "/.pivot_root"
	if err := syscall.Unmount(putOld, syscall.MNT_DETACH); err != nil {
		return err
	}

	// remove put_old
	if err := os.RemoveAll(putOld); err != nil {
		return err
	}

	return nil
}

func mountProc(newRoot string) error {
	target := filepath.Join(newRoot, "/proc")
	os.MkdirAll(target, 0755)
	return syscall.Mount("proc", target, "proc", uintptr(0), "")
}

func justiceInit() {
	newRootPath := os.Args[1]
	input := os.Args[2]
	expected := os.Args[3]
	timeout, _ := strconv.ParseInt(os.Args[4], 10, 32)

	if err := mountProc(newRootPath); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(models.CODE_INIT_CONTAINER_FAILED)
	}

	if err := pivotRoot(newRootPath); err != nil {
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

func cgroupCPUInit(pid, cgCPUPath string) error {
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

func cgroupMemoryInit(pid, cgMemoryPath, memory string) error {
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

func cgroupCleanup(path string) error {
	cmd := exec.Command("rmdir", path)
	return cmd.Run()
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
	if err := cgroupCPUInit(string(pid), cgCPUPath); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		return
	}

	// MEMORY
	if err := cgroupMemoryInit(string(pid), cgMemoryPath, *memory); err != nil {
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

	if err := cgroupCleanup(cgCPUPath); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerCleanupError"})
	}

	if err := cgroupCleanup(cgMemoryPath); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerCleanupError"})
	}

	os.Exit(models.CODE_OK)
}
