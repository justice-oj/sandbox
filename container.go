package sandbox

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/reexec"
	"github.com/getsentry/raven-go"
	"github.com/satori/go.uuid"

	"github.com/justice-oj/sandbox/config"
	"github.com/justice-oj/sandbox/model"
	"github.com/justice-oj/sandbox/sandbox"
)

func init() {
	raven.SetDSN(config.SENTRY_DSN)
	// register "justiceInit" => justiceInit() every time
	reexec.Register("justiceInit", justiceInit)

	/**
	* 0. `init()` adds key "justiceInit" in `map`;
	* 1. reexec.Init() seeks if key `os.Args[0]` exists in `registeredInitializers`;
	* 2. for the first time this binary is invoked, the key is os.Args[0], AKA "/path/to/clike_container",
	     which `registeredInitializers` will return `false`;
	* 3. `main()` calls binary itself by reexec.Command("justiceInit", args...);
	* 4. for the second time this binary is invoked, the key is os.Args[0], AKA "justiceInit",
	*    which exists in `registeredInitializers`;
	* 5. the value `justiceInit()` is invoked, any hooks(like set hostname) before fork() can be placed here.
	*/
	if reexec.Init() {
		os.Exit(0)
	}
}

func justiceInit() {
	basedir := os.Args[1]
	input := os.Args[2]
	expected := os.Args[3]
	timeout, _ := strconv.ParseInt(os.Args[4], 10, 32)

	r := new(model.Result)
	if err := sandbox.InitNamespace(basedir); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(r.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(0)
	}

	var o, e bytes.Buffer
	cmd := exec.Command("/Main")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &o
	cmd.Stderr = &e
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = []string{"PS1=[justice] # "}

	time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	if err := cmd.Run(); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerRunTimeError"})
		result, _ := json.Marshal(r.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Stderr.WriteString("Error: " + err.Error() + "\n")
		return
	}
	endTime := time.Now().UnixNano() / int64(time.Millisecond)

	if e.Len() > 0 {
		result, _ := json.Marshal(r.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Stderr.WriteString("stderr: " + e.String() + "\n")
		return
	}

	output := strings.TrimSpace(o.String())
	if output == expected {
		// ms, MB
		timeCost, memoryCost := endTime-startTime, cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss/1024
		// timeCost value 0 will be omitted
		if timeCost == 0 {
			timeCost = 1
		}

		result, _ := json.Marshal(r.GetAcceptedTaskResult(timeCost, memoryCost))
		os.Stdout.Write(result)
	} else {
		result, _ := json.Marshal(r.GetWrongAnswerTaskResult(input, output, expected))
		os.Stdout.Write(result)
	}

	os.Stderr.WriteString("output: " + output + " | " + "expected: " + expected + "\n")
}

// logs will be printed to os.Stderr
func main() {
	basedir := flag.String("basedir", "/tmp", "basedir of tmp C binary")
	input := flag.String("input", "<input>", "test case input")
	expected := flag.String("expected", "<expected>", "test case expected")
	timeout := flag.String("timeout", "2000", "timeout in milliseconds")
	memory := flag.String("memory", "256", "memory limitation in MB")
	flag.Parse()

	result := new(model.Result)
	u, _ := uuid.NewV4()

	if err := sandbox.InitCGroup(strconv.Itoa(os.Getpid()), u.String(), *memory); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(result.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(0)
	}

	cmd := reexec.Command("justiceInit", *basedir, *input, *expected, *timeout, *memory)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Run(); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "ContainerRunTimeError"})
		result, _ := json.Marshal(result.GetRuntimeErrorTaskResult())
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stdout.Write(result)
	}

	os.Exit(0)
}
