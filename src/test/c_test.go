package gotest

import (
	"testing"
	"os"
	"os/exec"
)

func Test_C_0_Accepted(t *testing.T) {
	absPath, _ := os.Getwd()

	cpCmd := exec.Command("cp", absPath + "/../resources/c/0.in", absPath + "/tmp/0/Main.c")
	cpCmd.Run()

	cmd := exec.Command(absPath + "/../../bin/c_compiler", "-basedir=" + absPath + "/tmp/0/")
	cmd.Run()
}