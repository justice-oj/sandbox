package main

import (
	"flag"
	"fmt"
	"bytes"
	"os/exec"
	"os"
)

type Output struct {
	Runtime int32
	Memory  int32
	Status  int32
	Error   string
}

func main() {
	dir := flag.String("dir", "/tmp", "compiler working directory")
	flag.Parse()

	fmt.Println(*dir)

	err, stdout, stderr := compile(*dir)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}

	if len(stderr) > 0 {
		os.Stderr.WriteString(stderr)
	} else {
		os.Stdout.WriteString(stdout)
	}
}

func compile(dir string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("gcc", "Main.c", "-o", "Main")
	cmd.Dir = dir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return err, stdout.String(), stderr.String()
}
