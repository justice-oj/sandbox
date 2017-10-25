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

// compiler wrapper with timeout limitation
// os.Stderr will not be empty if any error occurred
func main() {
	compiler := flag.String("compiler", "/usr/bin/gcc", "C/CPP compiler with abs path")
	basedir := flag.String("basedir", "/tmp", "basedir of tmp C/CPP code snippet")
	filename := flag.String("filename", "Main.c", "name of file to be compiled")
	timeout := flag.Int("timeout", 5000, "compile timeout in milliseconds")
	std := flag.String("std", "gnu11", "language standards supported by gcc")
	flag.Parse()

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(*compiler, *filename, "-save-temps", "-std=" + *std, "-fmax-errors=10", "-static", "-o", "Main")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = *basedir

	time.AfterFunc(time.Duration(*timeout)*time.Millisecond, func() {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	if err := cmd.Run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString(stderr.String())
		// err.Error() == "signal: killed" means compiler is killed by our timer.
		if err.Error() == "signal: killed" {
			raven.SetDSN(config.SENTRY_DSN)
			raven.CaptureMessageAndWait(*basedir, map[string]string{"error": "CompileTimeExceeded"})
		}
		return
	}

	os.Stdout.WriteString("Compile OK")
}
