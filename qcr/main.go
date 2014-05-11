package main

import "github.com/modcloth/queued-command-runner"

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

import (
	"github.com/onsi/gocleanup"
)

func main() {
	//setup tmpdir
	tmp, err := ioutil.TempDir("", "runner")
	if err != nil {
		gocleanup.Exit(1)
	}

	gocleanup.Register(func() {
		os.RemoveAll(tmp)
	})

	script := `#!/bin/bash

echo begin >&2
sleep 1
echo end >&2
`
	script2 := `#!/bin/bash

echo begin >&2
sleep 2
echo end >&2
`
	scriptPath := fmt.Sprintf("%s/foo", tmp)
	scriptPath2 := fmt.Sprintf("%s/foo2", tmp)

	if err := ioutil.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		fmt.Printf("ERR: %q\n", err)
		gocleanup.Exit(3)
	}

	if err := ioutil.WriteFile(scriptPath2, []byte(script2), 0755); err != nil {
		fmt.Printf("ERR: %q\n", err)
		gocleanup.Exit(3)
	}

	for x := 0; x < 5; x++ {
		cmd := &exec.Cmd{
			Path:   scriptPath,
			Args:   []string{"foo"},
			Stdout: os.Stdout,
		}

		runner.Run(cmd)
	}

	for x := 0; x < 5; x++ {
		cmd := &exec.Cmd{
			Path:   scriptPath2,
			Args:   []string{"foo2"},
			Stdout: os.Stdout,
		}

		runner.Run(cmd)
	}

	// wait for commands to finish
	<-runner.Done

	gocleanup.Exit(0)
}
