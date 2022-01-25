package coffeetray

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

type Message interface {
	Subject() string
}

type Payload []string

type Payloaded interface {
	Message
	Payload() Payload
}

type Status int

const (
	Running = iota
	Stopped
	ScheduleRunning
	ScheduleIdle
)

type Statused interface {
	Message
	Status() Status
}

type CommandSupervisor struct {
	Path    string
	Args    []string
	Should  Assertion
	command *exec.Cmd
	ticker  *time.Ticker
}

type ControlMessage struct {
	subject string
	payload Payload
}

func (pl *ControlMessage) Subject() string {
	return pl.subject
}

func (pl *ControlMessage) Payload() []string {
	return pl.payload
}

type StatusMessage struct {
	subject string
	status  Status
}

func (pl *StatusMessage) Subject() string {
	return pl.subject
}

func (pl *StatusMessage) Status() Status {
	return pl.status
}

type Assertion int

const (
	Stop Assertion = iota
	Run
	Scheduled
)

func NewCommandSupervisor(name string, initialArgs []string) *CommandSupervisor {
	path, err := exec.LookPath(name)
	AssertNoError(err)
	return &CommandSupervisor{Path: path, Args: initialArgs, Should: Stop}
}

func (cs *CommandSupervisor) Start() {
	cs.Should = Run
	cs.Apply()
}

func (cs *CommandSupervisor) Stop() {
	cs.Should = Stop
	cs.Apply()
}

func (cs *CommandSupervisor) Schedule() {
	cs.Should = Scheduled
	cs.Apply()
}

func (cs *CommandSupervisor) doStart() {
	cs.command = exec.Command(cs.Path, cs.Args...)
	err := cs.command.Start()
	AssertNoError(err)
}

func (cs *CommandSupervisor) doStop(kill bool) {
	if kill {
		AssertNoError(cs.command.Process.Signal(syscall.SIGKILL))
	} else {
		AssertNoError(cs.command.Process.Signal(syscall.SIGTERM))
	}
	_, err := cs.command.Process.Wait()
	cs.command = nil
	AssertNoError(err)
}

func (cs *CommandSupervisor) Apply() {
	switch cs.Should {
	case Run:
		if !cs.isRunning() {
			cs.doStart()
		}
	case Stop:
		if cs.isRunning() {
			cs.doStop(false)
		}
	case Scheduled:
		shouldRun := cs.checkSchedule()
		if shouldRun && !cs.isRunning() {
			cs.doStart()
		} else if !shouldRun && cs.isRunning() {
			cs.doStop(false)
		}
	}
}

func (cs *CommandSupervisor) Run() chan Message {
	cs.ticker = time.NewTicker(time.Second * 30)
	ch := make(chan Message)

	go func() {
		defer func() {
			cs.ticker.Stop()
			cs.doStop(true)
		}()
		for doRun := true; doRun; {
			select {
			case t := <-cs.ticker.C:
				fmt.Println("Debug: Tick", t)
				cs.Apply()
			case msg := <-ch:
				controlMessage, ok := msg.(*ControlMessage)
				if !ok {
					break
				}
				switch controlMessage.Subject() {
				case "STOP":
					cs.Stop()
				case "APPLY":
					cs.doStop(false)
					cs.Args = controlMessage.Payload()
					cs.Apply()
				case "SCHEDULE":
					cs.Schedule()
				case "START":
					cs.Start()
				case "QUIT":
					cs.ticker.Stop()
					cs.Stop()
					doRun = false
				}
			}
			if cs.isRunning() {
				if cs.Should == Scheduled {
					ch <- &StatusMessage{"Scheduled (running)", ScheduleRunning}
				} else {
					ch <- &StatusMessage{"Running", Running}
				}
			} else {
				if cs.Should == Scheduled {
					ch <- &StatusMessage{"Scheduled (idle)", ScheduleIdle}
				} else {
					ch <- &StatusMessage{"Stopped", Stopped}
				}
			}
		}
	}()

	return ch
}

func (cs *CommandSupervisor) isRunning() bool {
	if cs.command != nil && cs.command.ProcessState != nil {
		fmt.Println(cs.command.ProcessState.Exited())
	}
	return !(cs.command == nil || cs.command.Process == nil || (cs.command.ProcessState != nil && cs.command.ProcessState.Exited()))
}

// Hardcoded schedule 9-19 on weekdays
func (cs *CommandSupervisor) checkSchedule() bool {
	now := time.Now()

	if now.Weekday() == time.Sunday || now.Weekday() == time.Saturday {
		return false
	}

	frm := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	if frm.After(now) {
		frm = frm.AddDate(0, 0, -1)
	}

	to := frm.Add(time.Hour * 10)

	return frm.Before(now) && to.After(now)
}
