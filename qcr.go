package qcr

import (
	"os/exec"
	"strings"
)

// tm for "Treasure Map"
var tm = make(map[string]chan *exec.Cmd)
var run = func(in chan *exec.Cmd) {
	for {
		cmd := <-in
		_ = cmd.Run()
	}
}

//Run runs your command.
func Run(cmd *exec.Cmd) error {
	key := strings.Join(cmd.Args, " ")
	if tm[key] == nil {
		tm[key] = make(chan *exec.Cmd)
		go run(tm[key])
	}

	tm[key] <- cmd
	return nil
}
