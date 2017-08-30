package main

import (
	"fmt"
	"flag"
	"os/exec"
	"time"
	"bytes"
	"os"
	"syscall"
)

func main() {
	compiler := flag.String("compiler", "gcc", "C compiler with abs path")
	basedir := flag.String("basedir", "/tmp", "basedir of tmp C code snippet")
	filename := flag.String("filename", "Main.c", "name of file to be compiled")
	timeout := flag.Int("timeout", 10, "timeout in seconds")
	flag.Parse()
	fmt.Println(*compiler, " | ", *basedir, " | ", *filename, " | ", *timeout)

	cmd := exec.Command(*compiler, *filename, "-o", "Main")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = *basedir
	cmd.Start()

	err, isTimeout := CmdRunWithTimeout(cmd, time.Duration(*timeout)*time.Second)

	if isTimeout {
		os.Stderr.WriteString("Compile Time Exceeded")
		return
	}

	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}

	os.Stdout.WriteString("Compile OK")
}

func CmdRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		cmd.Process.Signal(syscall.SIGINT)
		err = cmd.Process.Kill()
		go func() {
			<-done
		}()
		return err, true
	case err = <-done:
		return err, false
	}
}
