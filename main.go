package main

import (
	"encoding/json"
	"fmt"
	"os"
	"bytes"
	"os/exec"
)

type Output struct {
	Runtime int32
	Memory  int32
	Status  int32
	Error   string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("one argument only")
	}

	cmd := os.Args[1]
	fmt.Println(cmd)

	err, stdout, stderr := run_command(cmd)
	if err != nil {
		fmt.Print(json.Marshal(Output{-1, -1, 4, err.Error()}))
	}

	if len(stderr) > 0 {
		fmt.Print(json.Marshal(Output{-1, -1, 4, nil}))
	} else {
		fmt.Print(json.Marshal(Output{-1, -1, 4, nil}))
	}
}

func run_command(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}
