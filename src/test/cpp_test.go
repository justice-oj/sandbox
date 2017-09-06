package gotest

import (
	"testing"
	"os"
	"bytes"
	"strings"
	"os/exec"
)

// HELPER
// copy test source file `*.in` to tmp dir
func copyCppSourceFile(name string, t *testing.T) (string, string) {
	t.Logf("Copying file %s ...", name)

	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath+"/tmp", absPath+"/../.."
	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir+"/src/resources/cpp/"+name, baseDir+"/Main.cpp")
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
// compileC source file
func compileCpp(name, baseDir, projectDir string, t *testing.T) (string) {
	t.Logf("Compiling file %s ...", name)

	var compilerStderr bytes.Buffer
	compilerArgs := []string{"-compiler=/usr/bin/g++", "-basedir=" + baseDir, "-timeout=3000"}
	compilerCmd := exec.Command(projectDir+"/bin/cpp_compiler", compilerArgs...)
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
// runC binary in our container
func runCpp(name, baseDir, projectDir string, t *testing.T) (string) {
	t.Logf("Running file %s ...", name)

	var containerStdout bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-timeout=3000", "-input=10:10:23AM", "-expected=10:10:23"}
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

func Test_Cpp_Accepted(t *testing.T) {
	name := "ac.cpp"
	baseDir, projectDir := copyCppSourceFile(name, t)
	compilerStderr := compileCpp(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runCpp(name, baseDir, projectDir, t)
	if !strings.Contains(containerErr, "\"status\":0") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr + " => status != 0")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}
