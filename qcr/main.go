package main

import "github.com/modcloth/queued-command-runner"

import (
	"os"
	"os/exec"
)

func main() {

	for x := 0; x < 50; x++ {
		cmd := exec.Command("ls", "/tmp")
		qcr.Run(cmd)
	}
	os.Exit(0)
}
