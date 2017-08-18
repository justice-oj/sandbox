package main

import (
	"flag"
	"fmt"
	"bytes"
	"os/exec"
	"os"
)

func main() {
	dir := flag.String("dir", "NOT NULL", "sandbox working directory")
	file := flag.String("file", "NOT NULL", "binary filename")
	flag.Parse()

	fmt.Println(*dir, "|", *file)

	err, stdout, stderr := run(*dir, *file)

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

func run(dir, file string) (error, string, string) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("./" + file)
	cmd.Dir = dir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return err, stdout.String(), stderr.String()
}