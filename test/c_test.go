package test

import (
	"bytes"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var (
	baseDir    string
	projectDir string
)

// copy test source file `*.c` to tmp dir
func copyCSourceFile(name string, t *testing.T) {
	t.Logf("Copying file %s ...", name)
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		t.Errorf("Invoke mkdir(%s) err: %v", baseDir, err.Error())
	}

	args := []string{
		projectDir + "/resources/c/" + name,
		baseDir + "/Main.c",
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

func TestMain(m *testing.M) {
	projectDir, _ = os.Getwd()
	baseDir = projectDir + "/tmp"

	os.Exit(m.Run())
}

func TestCAC(t *testing.T) {
	name := "ac.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		So(runC(baseDir, "64", "1000", t), ShouldContainSubstring, `"status":0`)
	})
}

func TestCCompilerBomb0(t *testing.T) {
	name := "compiler_bomb_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestCCompilerBomb1(t *testing.T) {
	name := "compiler_bomb_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestCCompilerBomb2(t *testing.T) {
	name := "compiler_bomb_2.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestCCompilerBomb3(t *testing.T) {
	name := "compiler_bomb_3.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestCCoreDump0(t *testing.T) {
	name := "core_dump_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		So(runC(baseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestCCoreDump1(t *testing.T) {
	name := "core_dump_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		// warning: division by zero [-Wdiv-by-zero]
		So(compileC(name, baseDir, t), ShouldBeEmpty)
		So(runC(baseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestCCoreDump2(t *testing.T) {
	name := "core_dump_2.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		// *** stack smashing detected ***: terminated
		So(runC(baseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestCForkBomb0(t *testing.T) {
	name := "fork_bomb_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		// got `signal: killed`
		So(runC(baseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestCForkBomb1(t *testing.T) {
	name := "fork_bomb_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		// got `signal: killed`
		So(runC(baseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestCGetHostByName(t *testing.T) {
	name := "get_host_by_name.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		// Main.c:(.text+0x28): warning: Using 'gethostbyname' in statically linked applications
		// requires at runtime the shared libraries from the glibc version used for linking
		// got `exit status 1`
		So(runC(baseDir, "64", "1000", t), ShouldContainSubstring, `"status":2`)
	})
}

func TestCIncludeLeaks(t *testing.T) {
	name := "include_leaks.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldContainSubstring, "/etc/shadow")
	})
}

func TestCInfiniteLoop(t *testing.T) {
	name := "infinite_loop.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		// got `signal: killed`
		So(runC(baseDir, "64", "1000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestCMemoryAllocation(t *testing.T) {
	name := "memory_allocation.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		So(runC(baseDir, "8", "5000", t), ShouldContainSubstring, "Runtime Error")
	})
}

func TestCPlainText(t *testing.T) {
	name := "plain_text.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldContainSubstring, "error")
	})
}

func TestCRunCommandLine0(t *testing.T) {
	name := "run_command_line_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		So(runC(baseDir, "64", "1000", t), ShouldContainSubstring, `"status":5`)
	})
}

func TestCRunCommandLine1(t *testing.T) {
	name := "run_command_line_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		So(runC(baseDir, "64", "1000", t), ShouldContainSubstring, `"status":5`)
	})
}

func aTestCSyscall0(t *testing.T) {
	name := "syscall_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		So(runC(baseDir, "16", "1000", t), ShouldContainSubstring, `"status":5`)
	})
}

func TestCTCPClient(t *testing.T) {
	name := "tcp_client.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(baseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", baseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, baseDir, t), ShouldBeEmpty)
		So(runC(baseDir, "16", "5000", t), ShouldContainSubstring, "Runtime Error")
	})
}
