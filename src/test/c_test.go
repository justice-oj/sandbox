package gotest

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// HELPER
// copy test source file `*.c` to tmp dir
func copyCSourceFile(name string, t *testing.T) string {
	t.Logf("Copying file %s ...", name)

	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath+"/tmp", absPath
	os.MkdirAll(baseDir, os.ModePerm)

	cmd := exec.Command("cp", projectDir+"/resources/c/"+name, baseDir+"/Main.c")
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	return baseDir
}

// HELPER
// compile C source file
func compileC(name, baseDir string, t *testing.T) string {
	t.Logf("Compiling file %s ...", name)

	var stderr bytes.Buffer
	args := []string{"-compiler=/usr/bin/gcc", "-basedir=" + baseDir, "-filename=Main.c", "-timeout=3000", "-std=gnu11"}
	cmd := exec.Command("/opt/justice-sandbox/bin/clike_compiler", args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	return stderr.String()
}

// HELPER
// run C binary in our container
func runC(baseDir, memory, timeout string, t *testing.T) string {
	t.Log("Running binary /Main ...")

	var stdout, stderr bytes.Buffer
	args := []string{"-basedir=" + baseDir, "-input=10:10:23PM", "-expected=22:10:23", "-memory=" + memory, "-timeout=" + timeout}
	cmd := exec.Command("/opt/justice-sandbox/bin/clike_container", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	t.Log(stderr.String())
	return stdout.String()
}

func TestCAC(t *testing.T) {
	name := "ac.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":0") {
		t.Error(containerOutput + " => status != 0")
	}
}

func TestCCompilerBomb0(t *testing.T) {
	name := "compiler_bomb_0.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func TestCCompilerBomb1(t *testing.T) {
	name := "compiler_bomb_1.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func TestCCompilerBomb2(t *testing.T) {
	name := "compiler_bomb_2.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func TestCCompilerBomb3(t *testing.T) {
	name := "compiler_bomb_3.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func TestCCoreDump0(t *testing.T) {
	name := "core_dump_0.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// function 'foo' recurses infinitely
	containerOutput := runC(baseDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCCoreDump1(t *testing.T) {
	name := "core_dump_1.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	// warning: division by zero [-Wdiv-by-zero]
	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCCoreDump2(t *testing.T) {
	name := "core_dump_2.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// *** stack smashing detected ***: terminated
	containerOutput := runC(baseDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCForkBomb0(t *testing.T) {
	name := "fork_bomb_0.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// got `signal: killed`
	containerOutput := runC(baseDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCForkBomb1(t *testing.T) {
	name := "fork_bomb_1.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// got `signal: killed`
	containerOutput := runC(baseDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCGetHostByName(t *testing.T) {
	name := "get_host_by_name.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// Main.c:(.text+0x28): warning: Using 'gethostbyname' in statically linked applications
	// requires at runtime the shared libraries from the glibc version used for linking
	// got `exit status 1`
	containerOutput := runC(baseDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":2") {
		t.Error(containerOutput)
	}
}

func TestCIncludeLeaks(t *testing.T) {
	name := "include_leaks.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if !strings.Contains(compilerStderr, "/etc/shadow") {
		t.Error(compilerStderr + " => Compile error does not contain string `/etc/shadow`")
	}
}

func TestCInfiniteLoop(t *testing.T) {
	name := "infinite_loop.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// got `signal: killed`
	containerOutput := runC(baseDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCMemoryAllocation(t *testing.T) {
	name := "memory_allocation.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// `Killed` is sent to tty by kernel (and record will also be kept in /var/log/message)
	// both stdout and stderr are empty which will lead to status WA
	// OR...
	// just running out of time
	containerOutput := runC(baseDir, "8", "5000", t)
	if !strings.ContainsAny(containerOutput, "\"status\":5 & Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCPlainText(t *testing.T) {
	name := "plain_text.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if !strings.Contains(compilerStderr, "error") {
		t.Error(compilerStderr + " => Compile error does not contain string `error`")
	}
}

func TestCRunCommandLine0(t *testing.T) {
	name := "run_command_line_0.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, "16", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":5") {
		t.Error(containerOutput)
	}
}

func TestCRunCommandLine1(t *testing.T) {
	name := "run_command_line_1.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, "16", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":5") {
		t.Error(containerOutput)
	}
}

func TestCSyscall0(t *testing.T) {
	name := "syscall_0.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, "16", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":5") {
		t.Error(containerOutput)
	}
}

func TestCTCPClient(t *testing.T) {
	name := "tcp_client.c"
	baseDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, "16", "5000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}
