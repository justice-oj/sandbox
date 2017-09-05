package gotest

import (
	"testing"
	"os"
	"os/exec"
	"bytes"
	"strings"
)

// HELPER
// copy test source file `*.in` to tmp dir
func copySourceFile(n string, t *testing.T) (string, string) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."
	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/c/" + n + ".in", baseDir + "/Main.c")
	cpErr := cpCmd.Run()

	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("Copy " + n + ".in failed")
	}

	return baseDir, projectDir
}

// HELPER
// compile source file
func compile(baseDir, projectDir string, t *testing.T) (string) {
	var compilerStderr bytes.Buffer

	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
	}

	return compilerStderr.String()
}

// HELPER
// run binary in our container
func run(baseDir, projectDir string, t *testing.T) (string) {
	var containerStdout bytes.Buffer

	containerArgs := []string{"-basedir=" + baseDir, "-timeout=3000", "-input=10:10:23AM", "-expected=10:10:23"}
	containerCmd := exec.Command(projectDir + "/bin/c_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerErr := containerCmd.Run()

	if containerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr.Error())
	}

	return containerStdout.String()
}

func Test_C_Accepted(t *testing.T) {
	baseDir, projectDir := copySourceFile("0", t)

	compilerStderr := compile(baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
	}

	containerErr := run(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "\"status\":0") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr + " => status != 0")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Include_Leaks(t *testing.T) {
	baseDir, projectDir := copySourceFile("1", t)

	compilerStderr := compile(baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "/etc/shadow") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `/etc/shadow`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Error(t *testing.T) {
	baseDir, projectDir := copySourceFile("2", t)

	compilerStderr := compile(baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "error") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `error`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_0(t *testing.T) {
	baseDir, projectDir := copySourceFile("3", t)

	compilerStderr := compile(baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_1(t *testing.T) {
	baseDir, projectDir := copySourceFile("4", t)

	compilerStderr := compile(baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_2(t *testing.T) {
	baseDir, projectDir := copySourceFile("5", t)

	compilerStderr := compile(baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Run_Infinite_Loop(t *testing.T) {
	baseDir, projectDir := copySourceFile("6", t)

	compilerStderr := compile(baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
	}

	containerErr := run(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
	}

	os.RemoveAll(baseDir + "/")
}

/*func Test_C_Run_Fork_Bomb(t *testing.T) {
	baseDir, projectDir := copySourceFile("7", t)

	compilerStderr := compile(baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
	}

	containerErr := run(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
	}

	os.RemoveAll(baseDir + "/")
}*/

func Test_C_Run_Command_Line(t *testing.T) {
	baseDir, projectDir := copySourceFile("8", t)

	compilerStderr := compile(baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
	}

	containerErr := run(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "\"status\":5") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Run_System_Call(t *testing.T) {
	baseDir, projectDir := copySourceFile("9", t)

	compilerStderr := compile(baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
	}

	containerErr := run(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "File not found") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
	}

	os.RemoveAll(baseDir + "/")
}