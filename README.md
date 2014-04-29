queued-command-runner
=====================

[![Build Status](https://travis-ci.org/modcloth/queued-command-runner.svg?branch=master)](https://travis-ci.org/modcloth/queued-command-runner)
[![GoDoc](https://godoc.org/github.com/modcloth/queued-command-runner?status.png)](https://godoc.org/github.com/modcloth/queued-command-runner)

Runs distinct commands parallel goroutines while queueing indistinct
commands so they're never run twice at the same time.

## Installation &amp; Usage

Install with:

```bash
go get github.com/modcloth/queued-command-runner
```

Then, in your code:

```go
package main

// package name is "qcr"
import (
  "fmt"
  "github.com/modcloth/queued-command-runner"
  "os"
  "os/exec"
)

func main() {
  fmt.Println("Running a command now.")

  pwd := os.Getenv("PWD")

  cmd := exec.Command("ls", "-la", pwd)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  if err := qcr.Run(cmd) ; err != nil {
    // print a useful message
    // exit nonzero
  }

  os.Exit(0)
}
```
