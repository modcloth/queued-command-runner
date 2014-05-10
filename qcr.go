package qcr

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

import (
	structures "github.com/hishboy/gocommons/lang"
)

// tm for "Treasure Map"
var tm = make(map[string]*runner)
var tmLock = &sync.Mutex{}

/*
Done is qcr's exit channel - if you use qcr, you MUST wait on Done to ensure
your commands get run.  This can be accomplished by including the following at
the bottom of main():

  <-qcr.Done
*/
var Done = make(chan bool)

type runner struct {
	queue *structures.Queue
	*sync.Mutex
	key string
}

func (r *runner) start() {
	for {
		r.Lock()
		cmd := r.queue.Poll()
		if cmd == nil {
			destroyRunner(r)
			//r.Unlock() //should this be here?
			break
		} else {
			cmd := cmd.(*exec.Cmd)

			r.Unlock()

			if err := cmd.Run(); err != nil {
				fmt.Printf("OOPS, qcr encountered an error for %q: %q\n", r.key, err)
			}
		}
	}
}

func (r *runner) enqueue(cmd *exec.Cmd) {
	r.Lock()
	defer r.Unlock()
	r.queue.Push(cmd)
}

//Run runs your command.
func Run(cmd *exec.Cmd) error {
	tmLock.Lock()
	defer tmLock.Unlock()

	key := strings.Join(cmd.Args, " ")

	if tm[key] == nil {
		tm[key] = newRunner(cmd)
		go tm[key].start()
	} else {
		tm[key].enqueue(cmd)
	}

	return nil
}

func newRunner(cmd *exec.Cmd) *runner {
	q := structures.NewQueue()
	q.Push(cmd)

	ret := &runner{
		key:   strings.Join(cmd.Args, " "),
		Mutex: &sync.Mutex{},
		queue: q,
	}
	return ret
}

func destroyRunner(r *runner) {
	tmLock.Lock()
	defer tmLock.Unlock()

	if r.queue.Len() != 0 {
		panic("HOW THE HELL DID YOU GET HERE?!?!")
	}

	delete(tm, r.key)
	if len(tm) == 0 {
		Done <- true
	}
}
