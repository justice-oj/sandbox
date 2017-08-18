package main

import (
	"flag"
	"fmt"
	"bytes"
	"os/exec"
	"os"
)

func main() {
	compiler := flag.String("compiler", "NOT NULL", "executable compiler")
	dir := flag.String("dir", "NOT NULL", "compiler working directory")
	file := flag.String("file", "NOT NULL", "the original source file")
	flag.Parse()

	fmt.Println(*compiler, "|", *dir, "|", *file)

	err, stdout, stderr := compile(*compiler, *dir, *file)

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

func compile(compiler, dir, file string) (error, string, string) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(compiler, file, "-o", "Main")
	cmd.Dir = dir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return err, stdout.String(), stderr.String()
}
