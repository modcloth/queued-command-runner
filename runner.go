/*
Package runner runs distinct commands in parallel goroutines while queueing
indistinct commands so they're never run twice at the same time (but all get run eventually).

Simply send your command to runner and let it do the rest!


Example usage:

  package main

  import "github.com/modcloth/queued-command-runner"

  import (
  	"fmt"
  	"os"
  	"os/exec"
  )

  func main() {
  	fmt.Println("Running a command now.")

  	pwd := os.Getenv("PWD")

	ls := exec.Command("ls", "-la", pwd)
	ls.Stdout = os.Stdout
	ls.Stderr = os.Stderr
	cmd := &runner.Command{
		Cmd: ls,
	}

	runner.Run(cmd)

  	WaitOnRunner:
  	for {
  		select {
  		case <-runner.Done:
  			break WaitOnRunner
		case err := <-runner.Errors:
			fmt.Printf("Uh oh, got an error: %q\n", err)
  		}
  	}

  	os.Exit(0)
  }
*/
package runner

import (
	"os/exec"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	structures "github.com/hishboy/gocommons/lang"
)

// tm for "Treasure Map"
var tm = make(map[string]*runner)
var tmLock = &sync.Mutex{}
var logger = logrus.New()

func init() {
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.Info // default to Info
}

/*
SetLogLevel sets the log level (from Sirupsen/logrus) on queued-command-runner's internal logger
*/
func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

/*
Done is qcr's exit channel - if you use qcr, you MUST wait on Done to ensure
your commands get run.  This can be accomplished by including the following at
the bottom of main():

  <-runner.Done

*/
var Done = make(chan bool)

/*
Error is the channel that qcr will use to report any errors that occur.
*/
var Errors = make(chan *QCRError)

/*
QCRError is a custom error type that includes CommandStr, the command args of
the command that failed.
*/
type QCRError struct {
	CommandStr string
	Key        string
	error
}

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
			logger.WithFields(logrus.Fields{
				"key": r.key,
			}).Debug("no command available for runner, destroying runner")

			destroyRunner(r)
			break
		} else {
			cmd := cmd.(*exec.Cmd)

			r.Unlock()

			if err := cmd.Run(); err != nil {
				logger.WithFields(logrus.Fields{
					"key":   r.key,
					"error": err,
				}).Error("error running command, notifying on Errors channel")

				Errors <- &QCRError{
					error:      err,
					CommandStr: r.key,
				}
			}
		}
	}
}

/*
Command is a small wrapper for *exec.Cmd so that a custom key may be specified.  If no key
is specified (i.e. Key == ""), key is defaulted to the following:

	key = strings.Join(cmd.Cmd.Args, " ")
*/
type Command struct {
	Key string
	Cmd *exec.Cmd
}

func (r *runner) enqueue(cmd *exec.Cmd) {
	r.Lock()
	defer r.Unlock()
	r.queue.Push(cmd)
}

//Run runs your command.
func Run(cmd *Command) {
	tmLock.Lock()
	defer tmLock.Unlock()

	if cmd.Key == "" {
		cmd.Key = strings.Join(cmd.Cmd.Args, " ")
	}

	key := cmd.Key

	logger.WithFields(logrus.Fields{
		"key": key,
	}).Info("runner received run request")

	if tm[key] == nil {
		logger.WithFields(logrus.Fields{
			"key": key,
		}).Debug("creating new runner for key")

		tm[key] = newRunner(cmd)
		go tm[key].start()
	} else {
		logger.WithFields(logrus.Fields{
			"key": key,
		}).Debug("runner for key already exists, enqueueing")

		tm[key].enqueue(cmd.Cmd)
	}
}

func newRunner(cmd *Command) *runner {
	q := structures.NewQueue()
	q.Push(cmd.Cmd)

	ret := &runner{
		key:   cmd.Key,
		Mutex: &sync.Mutex{},
		queue: q,
	}
	return ret
}

func destroyRunner(r *runner) {
	tmLock.Lock()
	defer tmLock.Unlock()

	if r.queue.Len() != 0 {
		logger.WithFields(logrus.Fields{
			"key": r.key,
		}).Panic("HOW THE HELL DID YOU GET HERE?!?!")
	}

	delete(tm, r.key)
	if len(tm) == 0 {
		logger.WithFields(logrus.Fields{
			"key": r.key,
		}).Debug("after deleting key, map is empty, notifying on Done channel")

		Done <- true
	}
}
