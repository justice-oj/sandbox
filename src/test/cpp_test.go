package gotest

import (
	"testing"
	"os"
	"bytes"
	"strings"
	"os/exec"
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
func compileCpp(name, baseDir, projectDir string, t *testing.T) (string) {
	t.Logf("Compiling file %s ...", name)

	var stderr bytes.Buffer
	cmd := exec.Command(projectDir+"/bin/cpp_compiler", "-basedir="+baseDir)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	return stderr.String()
}

// HELPER
// run CPP binary in our container
func runCpp(baseDir, projectDir, memory, timeout string, t *testing.T) (string) {
	t.Log("Running file /Main ...")

	var stdout, stderr bytes.Buffer
	args := []string{"-basedir=" + baseDir, "-input=10:10:23AM", "-expected=10:10:23", "-memory=" + memory, "-timeout=" + timeout}
	cmd := exec.Command(projectDir+"/bin/cpp_container", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

	t.Log(stderr.String())
	return stdout.String()
}

func Test_Cpp_AC(t *testing.T) {
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

func Test_Cpp_Compiler_Bomb_0(t *testing.T) {
	name := "compiler_bomb_0.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func Test_Cpp_Compiler_Bomb_1(t *testing.T) {
	name := "compiler_bomb_1.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func Test_Cpp_Compiler_Bomb_2(t *testing.T) {
	name := "compiler_bomb_2.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "compilation terminated due to -fmax-errors=") {
		t.Error(compilerStderr + " => Compile error does not contain string `fmax-errors`")
	}
}

func Test_Cpp_Compiler_Bomb_3(t *testing.T) {
	name := "compiler_bomb_3.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "template instantiation depth exceeds maximum of") {
		t.Error(compilerStderr + " => Compile error does not contain string `template instantiation depth exceeds`")
	}
}

func Test_Cpp_Compiler_Bomb_4(t *testing.T) {
	name := "compiler_bomb_4.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "signal: killed") {
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
	}
}

func Test_Cpp_Fork_Bomb(t *testing.T) {
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

func Test_Cpp_Include_Leaks(t *testing.T) {
	name := "include_leaks.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "/etc/shadow") {
		t.Error(compilerStderr + " => Compile error does not contain string `/etc/shadow`")
	}
}

func Test_Cpp_Infinite_Loop(t *testing.T) {
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

func Test_Cpp_Plain_Text(t *testing.T) {
	name := "plain_text.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	defer os.RemoveAll(baseDir)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)
	if !strings.Contains(compilerStderr, "error") {
		t.Error(compilerStderr + " => Compile error does not contain string `error`")
	}
}

func Test_Cpp_Run_Command_Line_0(t *testing.T) {
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
