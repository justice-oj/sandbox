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

	cpCmd := exec.Command("cp", projectDir+"/src/test/resources/cpp/"+name, baseDir+"/Main.cpp")
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
// compile CPP source file
func compileCpp(name, baseDir, projectDir string, t *testing.T) (string) {
	t.Logf("Compiling file %s ...", name)

	var compilerStderr bytes.Buffer
	compilerCmd := exec.Command(projectDir+"/bin/cpp_compiler", "-basedir=" + baseDir)
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
// run CPP binary in our container
func runCpp(baseDir, projectDir string, t *testing.T) (string) {
	t.Log("Running file /Main ...")

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-input=10:10:23AM", "-expected=10:10:23"}
	containerCmd := exec.Command(projectDir+"/bin/cpp_container", containerArgs...)
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

func Test_Cpp_AC(t *testing.T) {
	name := "ac.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runCpp(baseDir, projectDir, t)
	if !strings.Contains(containerErr, "\"status\":0") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr + " => status != 0")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_Cpp_Compiler_Bomb_0(t *testing.T) {
	name := "compiler_bomb_0.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_Cpp_Compiler_Bomb_1(t *testing.T) {
	name := "compiler_bomb_1.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_Cpp_Compiler_Bomb_2(t *testing.T) {
	name := "compiler_bomb_2.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "compilation terminated due to -fmax-errors=") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `fmax-errors`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_Cpp_Compiler_Bomb_3(t *testing.T) {
	name := "compiler_bomb_3.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "template instantiation depth exceeds maximum of") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `template instantiation depth exceeds`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_Cpp_Compiler_Bomb_4(t *testing.T) {
	name := "compiler_bomb_4.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_Cpp_Fork_Bomb(t *testing.T) {
	name := "fork_bomb.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runCpp(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_Cpp_Include_Leaks(t *testing.T) {
	name := "include_leaks.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "/etc/shadow") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `/etc/shadow`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_Cpp_Infinite_Loop(t *testing.T) {
	name := "infinite_loop.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runCpp(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_Cpp_Plain_Text(t *testing.T) {
	name := "plain_text.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)

	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "error") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `error`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_Cpp_Run_Command_Line_0(t *testing.T) {
	name := "run_command_line_0.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runCpp(baseDir, projectDir, t)

	if !strings.Contains(containerErr, "\"status\":5") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}