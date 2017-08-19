package main

import (
	"flag"
	"bytes"
	"os/exec"
	"os"
	"path/filepath"
)

func main() {
	dir := flag.String("dir", "NOT NULL", "sandbox working directory")
	file := flag.String("file", "NOT NULL", "binary filename")
	stdin := flag.String("stdin", "NOT NULL", "test case read from stdin")
	flag.Parse()

	//fmt.Println(*dir, "|", *file)

	err, stdout, stderr := run(*dir, *file, *stdin)

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

func run(dir, file, stdin string) (error, string, string) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(dir + string(filepath.Separator) + file)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = bytes.NewReader([]byte(stdin))
	err := cmd.Run()

	return err, stdout.String(), stderr.String()
}
