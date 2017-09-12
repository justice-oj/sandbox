package gotest

import (
	"testing"
	"os"
	"os/exec"
	"bytes"
	"strings"
)

// HELPER
// copy test source file `*.c` to tmp dir
func copyCSourceFile(name string, t *testing.T) (string, string) {
	t.Logf("Copying file %s ...", name)

	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath+"/tmp", absPath+"/../.."
	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir+"/src/test/resources/c/"+name, baseDir+"/Main.c")
	cpErr := cpCmd.Run()

	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(cpErr.Error())
		t.FailNow()
	}

	t.Log("Done")
	return baseDir, projectDir
}

// HELPER
// compile C source file
func compileC(name, baseDir, projectDir string, t *testing.T) (string) {
	t.Logf("Compiling file %s ...", name)

	var compilerStderr bytes.Buffer
	compilerCmd := exec.Command(projectDir+"/bin/c_compiler", "-basedir="+baseDir)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
		t.FailNow()
	}

	t.Log("Done")
	return compilerStderr.String()
}

// HELPER
// run C binary in our container
func runC(baseDir, projectDir string, t *testing.T) (string) {
	t.Log("Running binary /Main ...")

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-input=10:10:23PM", "-expected=22:10:23"}
	containerCmd := exec.Command(projectDir+"/bin/c_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerErr := containerCmd.Run()

	if containerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr.Error())
		t.FailNow()
	}

	t.Log("Done")
	return containerStdout.String()
}

func Test_C_AC(t *testing.T) {
	name := "ac.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, t)
	if !strings.Contains(containerErr, "\"status\":0") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr + " => status != 0")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_0(t *testing.T) {
	name := "compiler_bomb_0.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_1(t *testing.T) {
	name := "compiler_bomb_1.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_2(t *testing.T) {
	name := "compiler_bomb_2.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Fork_Bomb(t *testing.T) {
	name := "fork_bomb.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Include_Leaks(t *testing.T) {
	name := "include_leaks.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "/etc/shadow") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `/etc/shadow`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Infinite_Loop(t *testing.T) {
	name := "infinite_loop.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Memory_Allocation(t *testing.T) {
	name := "memory_allocation.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, t)

	// `Killed` will be sent to tty, both stdout and stderr are empty
	if !strings.Contains(containerErr, "\"status\":5") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Plain_Text(t *testing.T) {
	name := "plain_text.c"
	baseDir, projectDir := copyCSourceFile(name, t)

	compilerStderr := compileC(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "error") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `error`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Run_Command_Line_0(t *testing.T) {
	name := "run_command_line_0.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "\"status\":5") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Run_Command_Line_1(t *testing.T) {
	name := "run_command_line_1.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "\"status\":5") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Syscall_0(t *testing.T) {
	name := "syscall_0.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "\"status\":5") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}
