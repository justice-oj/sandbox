package sandbox

import (
	"bytes"
	"strings"
	"syscall"
	"os/exec"
	"os"
	"time"
	"encoding/json"

	"github.com/getsentry/raven-go"

	"../model"
)

func Run(timeout int32, basedir, input, expected, cmdName string, cmdArgs ...string) {
	r := new(model.Result)

	if err := InitNamespace(basedir); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(r.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(0)
	}

	var o, e bytes.Buffer
	cmd := exec.Command(cmdName, cmdArgs...)
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
		result, _ := json.Marshal(r.GetAcceptedTaskResult(timeCost, memoryCost))
		os.Stdout.Write(result)
	} else {
		result, _ := json.Marshal(r.GetWrongAnswerTaskResult(input, output, expected))
		os.Stdout.Write(result)
	}

	os.Stderr.WriteString("output: " + output + " | " + "expected: " + expected + "\n")
}
