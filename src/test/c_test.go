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

	cmd := exec.Command("cp", projectDir+"/src/test/resources/c/"+name, baseDir+"/Main.c")
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	return baseDir, projectDir
}

// HELPER
// compile C source file
func compileC(name, baseDir, projectDir string, t *testing.T) (string) {
	t.Logf("Compiling file %s ...", name)

	var stderr bytes.Buffer
	cmd := exec.Command(projectDir+"/bin/c_compiler", "-basedir="+baseDir)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	return stderr.String()
}

// HELPER
// run C binary in our container
func runC(baseDir, projectDir, memory, timeout string, t *testing.T) (string) {
	t.Log("Running binary /Main ...")

	var stdout, stderr bytes.Buffer
	args := []string{"-basedir=" + baseDir, "-input=10:10:23PM", "-expected=22:10:23", "-memory=" + memory, "-timeout=" + timeout}
	cmd := exec.Command(projectDir+"/bin/c_container", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	t.Log(stderr.String())
	return stdout.String()
}

func Test_C_AC(t *testing.T) {
	name := "ac.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, projectDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":0") {
		t.Error(containerOutput + " => status != 0")
	}
}

func Test_C_Compiler_Bomb_0(t *testing.T) {
	name := "compiler_bomb_0.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func Test_C_Compiler_Bomb_1(t *testing.T) {
	name := "compiler_bomb_1.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func Test_C_Compiler_Bomb_2(t *testing.T) {
	name := "compiler_bomb_2.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func Test_C_Fork_Bomb(t *testing.T) {
	name := "fork_bomb.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// got `signal: killed`
	containerOutput := runC(baseDir, projectDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func Test_C_Get_Host_By_Name(t *testing.T) {
	name := "get_host_by_name.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// Main.c:(.text+0x28): warning: Using 'gethostbyname' in statically linked applications
	// requires at runtime the shared libraries from the glibc version used for linking
	// got `exit status 1`
	containerOutput := runC(baseDir, projectDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":2") {
		t.Error(containerOutput)
	}
}

func Test_C_Include_Leaks(t *testing.T) {
	name := "include_leaks.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "/etc/shadow") {
		t.Error(compilerStderr + " => Compile error does not contain string `/etc/shadow`")
	}
}

func Test_C_Infinite_Loop(t *testing.T) {
	name := "infinite_loop.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// got `signal: killed`
	containerOutput := runC(baseDir, projectDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func Test_C_Memory_Allocation(t *testing.T) {
	name := "memory_allocation.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// `Killed` is sent to tty by kernel (and record will also be kept in /var/log/message)
	// both stdout and stderr are empty which will lead to status WA
	// OR...
	// just running out of time
	containerOutput := runC(baseDir, projectDir, "8", "5000", t)
	if !strings.ContainsAny(containerOutput, "\"status\":5 & Runtime Error") {
		t.Error(containerOutput)
	}
}

func Test_C_Plain_Text(t *testing.T) {
	name := "plain_text.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "error") {
		t.Error(compilerStderr + " => Compile error does not contain string `error`")
	}
}

func Test_C_Run_Command_Line_0(t *testing.T) {
	name := "run_command_line_0.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, projectDir, "16", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":5") {
		t.Error(containerOutput)
	}
}

func Test_C_Run_Command_Line_1(t *testing.T) {
	name := "run_command_line_1.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, projectDir, "16", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":5") {
		t.Error(containerOutput)
	}
}

func Test_C_Syscall_0(t *testing.T) {
	name := "syscall_0.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, projectDir, "16", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":5") {
		t.Error(containerOutput)
	}
}

func Test_C_TCP_Client(t *testing.T) {
	name := "tcp_client.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileC(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runC(baseDir, projectDir, "16", "5000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}
