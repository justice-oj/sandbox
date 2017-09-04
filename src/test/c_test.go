package gotest

import (
	"testing"
	"os"
	"os/exec"
	"bytes"
	"strings"
)

func Test_C_Accepted(t *testing.T) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."
	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/c/0.in", baseDir + "/Main.c")
	cpErr := cpCmd.Run()
	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("copy source file failed")
	}

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=10000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerCmd.Run()
	if compilerStderr.Len() > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String())
	}

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-timeout=10000", "-input=10:10:23AM", "-expected=10:10:23"}
	containerCmd := exec.Command(projectDir + "/bin/c_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerCmd.Run()

	result := containerStdout.String()
	if !strings.Contains(result, "\"status\":0") {
		os.RemoveAll(baseDir + "/")
		t.Error(result)
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Include_Leaks(t *testing.T) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."

	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/c/1.in", baseDir + "/Main.c")
	cpErr := cpCmd.Run()
	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("copy source file failed")
	}

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=10000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerCmd.Run()

	if compilerStderr.Len() == 0 {
		os.RemoveAll(baseDir + "/")
		t.Error("Compile error missing")
	}

	if !strings.Contains(compilerStderr.String(), "/etc/shadow") {
		os.RemoveAll(baseDir + "/")
		t.Error("Compile error does not contain string `/etc/shadow`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Error(t *testing.T) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."

	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/c/2.in", baseDir + "/Main.c")
	cpErr := cpCmd.Run()
	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("copy source file failed")
	}

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=10000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerCmd.Run()

	if compilerStderr.Len() == 0 {
		os.RemoveAll(baseDir + "/")
		t.Error("Compile error missing")
	}

	if !strings.Contains(compilerStderr.String(), "error") {
		os.RemoveAll(baseDir + "/")
		t.Error("Compile error does not contain string `error`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_0(t *testing.T) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."

	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/c/3.in", baseDir + "/Main.c")
	cpErr := cpCmd.Run()
	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("copy source file failed")
	}

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=10000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerCmd.Run()

	if compilerStderr.Len() == 0 {
		os.RemoveAll(baseDir + "/")
		t.Error("Compile error missing")
	}

	if !strings.Contains(compilerStderr.String(), "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error("Compile error does not contain string `signal: killed`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_1(t *testing.T) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."

	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/c/4.in", baseDir + "/Main.c")
	cpErr := cpCmd.Run()
	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("copy source file failed")
	}

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=10000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerCmd.Run()

	if compilerStderr.Len() == 0 {
		os.RemoveAll(baseDir + "/")
		t.Error("Compile error missing")
	}

	if !strings.Contains(compilerStderr.String(), "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error("Compile error does not contain string `signal: killed`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_2(t *testing.T) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."

	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/c/5.in", baseDir + "/Main.c")
	cpErr := cpCmd.Run()
	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("copy source file failed")
	}

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=10000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerCmd.Run()

	if compilerStderr.Len() == 0 {
		os.RemoveAll(baseDir + "/")
		t.Error("Compile error missing")
	}

	if !strings.Contains(compilerStderr.String(), "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error("Compile error does not contain string `signal: killed`")
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Run_Infinite_Loop(t *testing.T) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."
	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/c/6.in", baseDir + "/Main.c")
	cpErr := cpCmd.Run()
	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("copy source file failed")
	}

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=10000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerCmd.Run()
	if compilerStderr.Len() > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String())
	}

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-timeout=10000", "-input=10:10:23AM", "-expected=10:10:23"}
	containerCmd := exec.Command(projectDir + "/bin/c_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerCmd.Run()

	result := containerStdout.String()
	if !strings.Contains(result, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(result)
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Run_Fork_Bomb(t *testing.T) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."
	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/c/7.in", baseDir + "/Main.c")
	cpErr := cpCmd.Run()
	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("copy source file failed")
	}

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=10000"}
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
}

func Test_C_Run_Cli(t *testing.T) {
	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath + "/tmp", absPath + "/../.."
	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir + "/src/resources/c/8.in", baseDir + "/Main.c")
	cpErr := cpCmd.Run()
	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error("copy source file failed")
	}

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-timeout=10000"}
	compilerCmd := exec.Command(projectDir + "/bin/c_compiler", compilerArgs...)
	compilerCmd.Stderr = &compilerStderr
	compilerCmd.Run()
	if compilerStderr.Len() > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr.String())
	}

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-timeout=10000", "-input=10:10:23AM", "-expected=10:10:23"}
	containerCmd := exec.Command(projectDir + "/bin/c_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerCmd.Run()

	result := containerStdout.String()
	if !strings.Contains(result, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(result)
	}

	os.RemoveAll(baseDir + "/")
}