package main

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	compiler := flag.String("compiler", "g++", "CPP compiler with abs path")
	basedir := flag.String("basedir", "/tmp", "basedir of tmp CPP code snippet")
	filename := flag.String("filename", "Main.cpp", "name of file to be compiled")
	timeout := flag.Int("timeout", 10, "timeout in seconds")
	flag.Parse()

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
		return
	}

	os.Stdout.WriteString("Compile OK")
}
