package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	CBaseDir    string
	CProjectDir string
)

// copy test source file `*.c` to tmp dir
func copyCSourceFile(name string, t *testing.T) {
	t.Logf("Copying file %s ...", name)
	if err := os.MkdirAll(CBaseDir, os.ModePerm); err != nil {
		t.Errorf("Invoke mkdir(%s) err: %v", CBaseDir, err.Error())
	}

	args := []string{
		CProjectDir + "/resources/c/" + name,
		CBaseDir + "/Main.c",
	}
	cmd := exec.Command("cp", args...)
	if err := cmd.Run(); err != nil {
		t.Errorf("Invoke `cp %s` err: %v", strings.Join(args, " "), err)
	}
}

// compile C source file
func compileC(name, baseDir string, t *testing.T) string {
	t.Logf("Compiling file %s ...", name)

	var stderr bytes.Buffer
	args := []string{
		"-compiler=/usr/bin/gcc",
		"-basedir=" + baseDir,
		"-filename=Main.c",
		"-timeout=3000",
		"-std=gnu11",
	}
	cmd := exec.Command("/opt/justice-sandbox/bin/clike_compiler", args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Errorf("Invoke `/opt/justice-sandbox/bin/clike_compiler %s` err: %v", strings.Join(args, " "), err)
	}

	return stderr.String()
}

// run binary in our container
func runC(baseDir, memory, timeout string, t *testing.T) string {
	t.Log("Running binary /Main ...")

	var stdout, stderr bytes.Buffer
	args := []string{
		"-basedir=" + baseDir,
		"-input=10:10:23PM",
		"-expected=22:10:23",
		"-memory=" + memory,
		"-timeout=" + timeout,
	}
	cmd := exec.Command("/opt/justice-sandbox/bin/clike_container", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Errorf("Invoke `/opt/justice-sandbox/bin/clike_container %s` err: %v", strings.Join(args, " "), err)
	}

	t.Logf("stderr of runC: %s", stderr.String())
	return stdout.String()
}

func TestC0000Fixture(t *testing.T) {
	CProjectDir, _ = os.Getwd()
	CBaseDir = t.TempDir()
}

func TestC0001AC(t *testing.T) {
	name := "ac.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		So(runC(CBaseDir, "64", "1000", t), ShouldContainSubstring, `"status":0`)
	})
}

func TestC0002CompilerBomb0(t *testing.T) {
	name := "compiler_bomb_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestC0003CompilerBomb1(t *testing.T) {
	name := "compiler_bomb_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestC0004CompilerBomb2(t *testing.T) {
	name := "compiler_bomb_2.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestC0005CompilerBomb3(t *testing.T) {
	name := "compiler_bomb_3.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestC0006CoreDump0(t *testing.T) {
	name := "core_dump_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		So(runC(CBaseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestC0007CoreDump1(t *testing.T) {
	name := "core_dump_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		// warning: division by zero [-Wdiv-by-zero]
		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		So(runC(CBaseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestC0008CoreDump2(t *testing.T) {
	name := "core_dump_2.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		// *** stack smashing detected ***: terminated
		So(runC(CBaseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestC0009ForkBomb0(t *testing.T) {
	name := "fork_bomb_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		// got `signal: killed`
		So(runC(CBaseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestC0010ForkBomb1(t *testing.T) {
	name := "fork_bomb_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		// got `signal: killed`
		So(runC(CBaseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestC0011GetHostByName(t *testing.T) {
	name := "get_host_by_name.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		// Main.c:(.text+0x28): warning: Using 'gethostbyname' in statically linked applications
		// requires at runtime the shared libraries from the glibc version used for linking
		// got `exit status 1`
		So(runC(CBaseDir, "64", "1000", t), ShouldContainSubstring, `"status":2`)
	})
}

func TestC0012IncludeLeaks(t *testing.T) {
	name := "include_leaks.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "/etc/shadow")
	})
}

func TestC0013InfiniteLoop(t *testing.T) {
	name := "infinite_loop.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		// got `signal: killed`
		So(runC(CBaseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestC0014MemoryAllocation(t *testing.T) {
	name := "memory_allocation.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		So(runC(CBaseDir, "8", "5000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestC0015PlainText(t *testing.T) {
	name := "plain_text.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "error")
	})
}

func TestC0016RunCommandLine0(t *testing.T) {
	name := "run_command_line_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		So(runC(CBaseDir, "64", "1000", t), ShouldContainSubstring, `"status":5`)
	})
}

func TestC0017RunCommandLine1(t *testing.T) {
	name := "run_command_line_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		So(runC(CBaseDir, "64", "1000", t), ShouldContainSubstring, `"status":5`)
	})
}

func TestC0018Syscall0(t *testing.T) {
	name := "syscall_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		So(runC(CBaseDir, "16", "1000", t), ShouldContainSubstring, `"status":5`)
	})
}

func TestC0019TCPClient(t *testing.T) {
	name := "tcp_client.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		So(runC(CBaseDir, "16", "5000", t), ShouldContainSubstring, "Runtime Error")
	})
}
