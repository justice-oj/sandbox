package gotest

import (
	"testing"
	"os"
	"os/exec"
	"bytes"
	"strings"
)

// copy test source file to tmp dir
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

func Test_C_Accepted(t *testing.T) {
	baseDir, projectDir := copySourceFile("0", t)
	var compilerStderr bytes.Buffer

	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
	}

	if compilerStderr.Len() > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String())
	}

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-timeout=3000", "-input=10:10:23AM", "-expected=10:10:23"}
	containerCmd := exec.Command(projectDir + "/bin/c_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerErr := containerCmd.Run()

	if containerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr.Error())
	}

	result := containerStdout.String()
	if !strings.Contains(result, "\"status\":0") {
		os.RemoveAll(baseDir + "/")
		t.Error(result + " => status != 0")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Include_Leaks(t *testing.T) {
	baseDir, projectDir := copySourceFile("1", t)
	var compilerStderr bytes.Buffer

	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
	}

	if !strings.Contains(compilerStderr.String(), "/etc/shadow") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String() + " => Compile error does not contain string `/etc/shadow`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Error(t *testing.T) {
	baseDir, projectDir := copySourceFile("2", t)
	var compilerStderr bytes.Buffer

	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
	}

	if !strings.Contains(compilerStderr.String(), "error") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String() + " => Compile error does not contain string `error`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_0(t *testing.T) {
	baseDir, projectDir := copySourceFile("3", t)
	var compilerStderr bytes.Buffer

	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
	}

	if !strings.Contains(compilerStderr.String(), "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String() + " => Compile error does not contain string `signal: killed`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_1(t *testing.T) {
	baseDir, projectDir := copySourceFile("4", t)
	var compilerStderr bytes.Buffer

	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
	}

	if !strings.Contains(compilerStderr.String(), "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String() + " => Compile error does not contain string `signal: killed`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_2(t *testing.T) {
	baseDir, projectDir := copySourceFile("5", t)
	var compilerStderr bytes.Buffer

	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
	}

	if !strings.Contains(compilerStderr.String(), "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String() + " => Compile error does not contain string `signal: killed`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Run_Infinite_Loop(t *testing.T) {
	baseDir, projectDir := copySourceFile("6", t)
	var compilerStderr bytes.Buffer

	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
	}

	if compilerStderr.Len() > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String())
	}

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-timeout=3000", "-input=10:10:23AM", "-expected=10:10:23"}
	containerCmd := exec.Command(projectDir + "/bin/c_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerErr := containerCmd.Run()

	if containerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr.Error())
	}

	result := containerStdout.String()
	if !strings.Contains(result, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(result)
	}

	os.RemoveAll(baseDir + "/")
}

/*func Test_C_Run_Fork_Bomb(t *testing.T) {
	baseDir, projectDir := copySourceFile("7", t)
	var compilerStderr bytes.Buffer

	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerCmd.Run()
	if compilerStderr.Len() > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String())
	}

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-timeout=3000", "-input=10:10:23AM", "-expected=10:10:23"}
	containerCmd := exec.Command(projectDir + "/bin/c_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerCmd.Run()

	result := containerStdout.String()
	if !strings.Contains(result, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(result)
	}

	os.RemoveAll(baseDir + "/")
}*/

func Test_C_Run_Cli(t *testing.T) {
	baseDir, projectDir := copySourceFile("8", t)
	var compilerStderr bytes.Buffer

	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
	}

	if compilerStderr.Len() > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String())
	}

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-timeout=3000", "-input=10:10:23AM", "-expected=10:10:23"}
	containerCmd := exec.Command(projectDir + "/bin/c_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerErr := containerCmd.Run()

	if containerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr.Error())
	}

	result := containerStdout.String()
	if !strings.Contains(result, "\"status\":5") {
		os.RemoveAll(baseDir + "/")
		t.Error(result)
	}

	os.RemoveAll(baseDir + "/")
}