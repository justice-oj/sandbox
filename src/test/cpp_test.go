package gotest

import (
	"testing"
	"os"
	"bytes"
	"strings"
	"os/exec"
)

func Test_Cpp_Accepted(t *testing.T) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."

	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/cpp/0.in", baseDir + "/Main.cpp")
	cpErr := cpCmd.Run()
	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("copy source file failed")
	}

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/g++", "-basedir=" + baseDir, "-timeout=10000"}
	compilerCmd := exec.Command(projectDir + "/bin/cpp_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerCmd.Run()
	if compilerStderr.Len() > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String())
	}

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-timeout=10000", "-input=10:10:23AM", "-expected=10:10:23"}
	containerCmd := exec.Command(projectDir + "/bin/cpp_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerCmd.Run()

	result := containerStdout.String()
	if !strings.Contains(result, "\"status\":0") {
		os.RemoveAll(baseDir + "/")
		t.Error(result)
	}

	os.RemoveAll(baseDir + "/")
}