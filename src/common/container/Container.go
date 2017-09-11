package container

import (
	"bytes"
	"strings"
	"syscall"
	"github.com/getsentry/raven-go"
	"os/exec"
	"os"
	"time"
	"encoding/json"
	"../../models"
	"../../common/namespace"
	"../../common/cgroup"
)

func Run(timeout int32, memory, pid, containerID, basedir, input, expected, cmdName string, cmdArgs ...string) {
	// Init Namespace
	if err := namespace.InitNamespace(basedir); err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		os.Exit(models.CODE_INIT_CONTAINER_FAILED)
	}

	// Init CGroup
	if err := cgroup.InitCGroup(string(pid), containerID, memory); err != nil {
		cgroup.Cleanup(containerID)
		raven.CaptureErrorAndWait(err, map[string]string{"error": "InitContainerFailed"})
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		return
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
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		return
	}
	endTime := time.Now().UnixNano() / int64(time.Millisecond)

	if e.Len() > 0 {
		result, _ := json.Marshal(models.GetRuntimeErrorTaskResult())
		os.Stdout.Write(result)
		return
	}

	output := strings.TrimSpace(o.String())
	if output == expected {
		// ms, MB
		timeCost, memoryCost := endTime-startTime, cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss/1024
		result, _ := json.Marshal(models.GetAccepptedTaskResult(timeCost, memoryCost))
		os.Stdout.Write(result)
	} else {
		result, _ := json.Marshal(models.GetWrongAnswerTaskResult(input, output, expected))
		os.Stdout.Write(result)
	}

	cgroup.Cleanup(containerID)
}
