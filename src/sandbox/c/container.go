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
)

func init() {
	reexec.Register("justice_init", justice_init)
	if reexec.Init() {
		os.Exit(0)
	}
}

func pivot_root(new_root string) error {
	put_old := filepath.Join(new_root, "/.pivot_root")

	// bind mount new_root to itself - this is a slight hack needed to satisfy requirement (2)
	//
	// The following restrictions apply to new_root and put_old:
	// 1.  They must be directories.
	// 2.  new_root and put_old must not be on the same filesystem as the current root.
	// 3.  put_old must be underneath new_root, that is, adding a nonzero
	//     number of /.. to the string pointed to by put_old must yield the same directory as new_root.
	// 4.  No other filesystem may be mounted on put_old.
	if err := syscall.Mount(new_root, new_root, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}

	// create put_old directory
	if err := os.MkdirAll(put_old, 0700); err != nil {
		return err
	}

	// call pivot_root
	if err := syscall.PivotRoot(new_root, put_old); err != nil {
		return err
	}

	// Note that this also applies to the calling process: pivot_root() may
	// or may not affect its current working directory.  It is therefore
	// recommended to call chdir("/") immediately after pivot_root().
	if err := os.Chdir("/"); err != nil {
		return err
	}

	// umount put_old, which now lives at /.pivot_root
	put_old = "/.pivot_root"
	if err := syscall.Unmount(put_old, syscall.MNT_DETACH); err != nil {
		return err
	}

	// remove put_old
	if err := os.RemoveAll(put_old); err != nil {
		return err
	}

	return nil
}

func mount_proc(new_root string) error {
	target := filepath.Join(new_root, "/proc")
	os.MkdirAll(target, 0755)
	return syscall.Mount("proc", target, "proc", uintptr(0), "")
}

func justice_init() {
	new_root_path := os.Args[1]
	input := os.Args[2]
	expected := os.Args[3]

	raven.SetDSN(config.SENTRY_DSN)

	if err := mount_proc(new_root_path); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(models.CODE_INIT_CONTAINER_FAILED)
	}

	if err := pivot_root(new_root_path); err != nil {
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

	justice_run(input, expected)
}

func justice_run(input, expected string) {
	raven.SetDSN(config.SENTRY_DSN)

	// for c programs, compiled binary with name [Main] will be located in "/"
	var o bytes.Buffer
	cmd := exec.Command("/Main")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &o
	cmd.Stderr = os.Stderr
	cmd.Env = []string{"PS1=[justice] # "}

	if err := cmd.Run(); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerRunTimeError"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(models.CODE_CONTAINER_RUNTIME_ERROR)
	}

	output := o.String()
	if output == expected {
		result, _ := json.Marshal(models.GetAccepptedTaskResult(13,456))
		os.Stdout.Write(result)
	} else {
		result, _ := json.Marshal(models.GetWrongAnswerTaskResult(input, output, expected))
		os.Stdout.Write(result)
	}
	os.Exit(models.CODE_OK)
}

func main() {
	basedir := flag.String("basedir", "/tmp", "basedir of tmp C binary")
	input := flag.String("input", "", "test case input")
	expected := flag.String("expected", "", "test case expected")
	flag.Parse()

	raven.SetDSN(config.SENTRY_DSN)

	cmd := reexec.Command("justice_init", *basedir, *input, *expected)
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

	if err := cmd.Start(); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerRunTimeError"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(models.CODE_CONTAINER_RUNTIME_ERROR)
	}

	if err := cmd.Wait(); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerRunTimeError"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(models.CODE_CONTAINER_RUNTIME_ERROR)
	}
}
