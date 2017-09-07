package main

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"syscall"
	"time"
	"github.com/getsentry/raven-go"
	"../../config"
)

func main() {
	compiler := flag.String("compiler", "/usr/bin/g++", "CPP compiler with abs path")
	basedir := flag.String("basedir", "/tmp", "basedir of tmp CPP code snippet")
	filename := flag.String("filename", "Main.cpp", "name of file to be compiled")
	timeout := flag.Int("timeout", 5000, "compile timeout in milliseconds")
	flag.Parse()

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(*compiler, *filename, "-save-temps", "-std=gnu++11", "-fmax-errors=10", "-static", "-o", "Main")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = *basedir

	time.AfterFunc(time.Duration(*timeout)*time.Millisecond, func() {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})
	err := cmd.Run()

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString(stderr.String())
		if err.Error() == "signal: killed" {
			raven.SetDSN(config.SENTRY_DSN)
			raven.CaptureMessageAndWait(*basedir, map[string]string{"error": "CompileTimeExceeded"})
		}
		return
	}

	os.Stdout.WriteString("Compile OK")
}
