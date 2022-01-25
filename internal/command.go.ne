package coffeetray

import (
	"fmt"
	"os/exec"
	"syscall"
)

const (
	COMMAND_START = iota
	COMMAND_STOP  = iota
	COMMAND_QUIT  = iota
)

type CommandLoop struct {
	Path    string
	command *exec.Cmd
}

type CommandEvent struct {
	Event int
	Args  []string
}

func NewCommandLoop(name string) *CommandLoop {
	path, err := exec.LookPath(name)
	AssertNoError(err)

	return &CommandLoop{path, nil}
}

func (cl *CommandLoop) Run() chan CommandEvent {

	ch := make(chan CommandEvent)

	go func() {
		for {
			defer cl.stop(true)

			event := <-ch
			switch event.Event {
			case COMMAND_STOP:
				fmt.Println("Stopping")
				cl.stop(false)
			case COMMAND_START:
				fmt.Println("Starting")
				cl.start(event.Args)
			case COMMAND_QUIT:
				cl.stop(true)
				return
			}
		}
	}()

	return ch
}

func (cl *CommandLoop) running() bool {
	return !(cl.command == nil || cl.command.Process == nil || (cl.command.ProcessState != nil && cl.command.ProcessState.Exited()))
}

func (cl *CommandLoop) stop(kill bool) {
	if !cl.running() {
		fmt.Println("Returning")
		return
	}
	if kill {
		AssertNoError(cl.command.Process.Signal(syscall.SIGKILL))
	} else {
		AssertNoError(cl.command.Process.Signal(syscall.SIGTERM))
	}
	_, err := cl.command.Process.Wait()
	cl.command = nil
	AssertNoError(err)
}

func (cl *CommandLoop) start(args []string) {
	if cl.running() {
		fmt.Println("Returning")
		return
	}
	cl.command = exec.Command(cl.Path, args...)
	err := cl.command.Start()
	AssertNoError(err)
}
