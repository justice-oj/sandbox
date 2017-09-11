package main

import (
	"encoding/json"
	"os"
	"syscall"
	"github.com/docker/docker/pkg/reexec"
	"github.com/satori/go.uuid"
	"github.com/getsentry/raven-go"
	"../../models"
	"../../config"
	"../../common/container"
	"flag"
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
	basedir := os.Args[1]
	input := os.Args[2]
	expected := os.Args[3]
	timeout, _ := strconv.ParseInt(os.Args[4], 10, 32)
	memory := os.Args[5]
	pid := os.Args[6]
	containerID := os.Args[7]

	container.Run(int32(timeout), memory, pid, containerID, basedir, input, expected, "/Main")
}

func main() {
	basedir := flag.String("basedir", "/tmp", "basedir of tmp C binary")
	input := flag.String("input", "<input>", "test case input")
	expected := flag.String("expected", "<expected>", "test case expected")
	timeout := flag.String("timeout", "2000", "timeout in milliseconds")
	memory := flag.String("memory", "64", "memory limitation in MB")
	flag.Parse()

	cmd := reexec.Command("justiceInit", *basedir, *input, *expected, *timeout, *memory,
		strconv.Itoa(os.Getpid()), uuid.NewV4().String())
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
	os.Exit(models.CODE_OK)
}
