package gotest

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// HELPER
// copy test source file `*.cpp` to tmp dir
func copyCppSourceFile(name string, t *testing.T) (string, string) {
	t.Logf("Copying file %s ...", name)

	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath+"/tmp", absPath+"/../.."
	os.MkdirAll(baseDir, os.ModePerm)

	cmd := exec.Command("cp", projectDir+"/src/test/resources/cpp/"+name, baseDir+"/Main.cpp")
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	return baseDir, projectDir
}

// HELPER
// compile CPP source file
func compileCpp(name, baseDir, projectDir string, t *testing.T) string {
	t.Logf("Compiling file %s ...", name)

	var stderr bytes.Buffer
	args := []string{"-compiler=/usr/bin/g++", "-basedir=" + baseDir, "-filename=Main.cpp", "-timeout=3000", "-std=gnu++14"}
	cmd := exec.Command(projectDir+"/bin/clike_compiler", args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	return stderr.String()
}

// HELPER
// run CPP binary in our container
func runCpp(baseDir, projectDir, memory, timeout string, t *testing.T) string {
	t.Log("Running file /Main ...")

	var stdout, stderr bytes.Buffer
	args := []string{"-basedir=" + baseDir, "-input=10:10:23AM", "-expected=10:10:23", "-memory=" + memory, "-timeout=" + timeout}
	cmd := exec.Command(projectDir+"/bin/clike_container", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	t.Log(stderr.String())
	return stdout.String()
}

func TestCppAC(t *testing.T) {
	name := "ac.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runCpp(baseDir, projectDir, "16", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":0") {
		t.Error(containerOutput + " => status != 0")
	}
}

func TestCppCompilerBomb0(t *testing.T) {
	name := "compiler_bomb_0.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func TestCppCompilerBomb1(t *testing.T) {
	name := "compiler_bomb_1.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func TestCppCompilerBomb2(t *testing.T) {
	name := "compiler_bomb_2.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "compilation terminated due to -fmax-errors=") {
		t.Error(compilerStderr + " => Compile error does not contain string `fmax-errors`")
	}
}

func TestCppCompilerBomb3(t *testing.T) {
	name := "compiler_bomb_3.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "template instantiation depth exceeds maximum of") {
		t.Error(compilerStderr + " => Compile error does not contain string `template instantiation depth exceeds`")
	}
}

func TestCppCompilerBomb4(t *testing.T) {
	name := "compiler_bomb_4.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func TestCppCoreDump0(t *testing.T) {
	name := "core_dump_0.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	// terminate called after throwing an instance of 'char const*'
	containerOutput := runCpp(baseDir, projectDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCppForkBomb(t *testing.T) {
	name := "fork_bomb.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runCpp(baseDir, projectDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCppIncludeLeaks(t *testing.T) {
	name := "include_leaks.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "/etc/shadow") {
		t.Error(compilerStderr + " => Compile error does not contain string `/etc/shadow`")
	}
}

func TestCppInfiniteLoop(t *testing.T) {
	name := "infinite_loop.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runCpp(baseDir, projectDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCppMemoryAllocation(t *testing.T) {
	name := "memory_allocation.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runCpp(baseDir, projectDir, "64", "1000", t)
	if !strings.Contains(containerOutput, "Runtime Error") {
		t.Error(containerOutput)
	}
}

func TestCppPlainText(t *testing.T) {
	name := "plain_text.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "error") {
		t.Error(compilerStderr + " => Compile error does not contain string `error`")
	}
}

func TestCppRunCommandLine0(t *testing.T) {
	name := "run_command_line_0.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if len(compilerStderr) > 0 {
		t.Error(compilerStderr)
		return
	}

	containerOutput := runCpp(baseDir, projectDir, "16", "1000", t)
	if !strings.Contains(containerOutput, "\"status\":5") {
		t.Error(containerOutput)
	}
}
