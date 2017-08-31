package main

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"syscall"
	"time"
	"github.com/getsentry/raven-go"
)

func main() {
	compiler := flag.String("compiler", "gcc", "C compiler with abs path")
	basedir := flag.String("basedir", "/tmp", "basedir of tmp C code snippet")
	filename := flag.String("filename", "Main.c", "name of file to be compiled")
	timeout := flag.Int("timeout", 10, "timeout in seconds")
	flag.Parse()

	raven.SetDSN("http://e79ebf76a31a43d18ef7bdfa7381537e:5b21a25106584b39ac22ebf0752412db@sentry.justice.plus/3")

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(*compiler, *filename, "-save-temps", "-fmax-errors=10", "-o", "Main")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = *basedir

	time.AfterFunc(time.Duration(*timeout)*time.Second, func() {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})
	err := cmd.Run()

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString(stderr.String())
		raven.SetDSN("http://e79ebf76a31a43d18ef7bdfa7381537e:5b21a25106584b39ac22ebf0752412db@sentry.justice.plus/3")
		raven.CaptureMessage(stderr.String(), nil)
		raven.CaptureError(err, nil)
		return
	}

	os.Stdout.WriteString("Compile OK")
}
