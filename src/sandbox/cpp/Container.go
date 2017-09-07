package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"syscall"
	"github.com/docker/docker/pkg/reexec"
	"path/filepath"
	"github.com/getsentry/raven-go"
	"../../models"
	"../../config"
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

func pivotRoot(newRoot string) error {
	putOld := filepath.Join(newRoot, "/.pivot_root")

	if err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}

	if err := os.MkdirAll(putOld, 0700); err != nil {
		return err
	}

	if err := syscall.PivotRoot(newRoot, putOld); err != nil {
		return err
	}

	if err := os.Chdir("/"); err != nil {
		return err
	}

	putOld = "/.pivot_root"
	if err := syscall.Unmount(putOld, syscall.MNT_DETACH); err != nil {
		return err
	}

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
		result, _ := json.Marshal(models.GetAccepptedTaskResult(endTime - startTime, memory))
		os.Stdout.Write(result)
	} else {
		result, _ := json.Marshal(models.GetWrongAnswerTaskResult(input, output, expected))
		os.Stdout.Write(result)
	}
}

func main() {
	basedir := flag.String("basedir", "/tmp", "basedir of tmp CPP binary")
	input := flag.String("input", "<input>", "test case input")
	expected := flag.String("expected", "<expected>", "test case expected")
	timeout := flag.String("timeout", "2000", "timeout in milliseconds")
	flag.Parse()

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
		os.Exit(models.CODE_CONTAINER_RUNTIME_ERROR)
	}

	os.Exit(models.CODE_OK)
}
