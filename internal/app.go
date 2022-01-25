package coffeetray

import (
	"fmt"

	"github.com/getlantern/systray"
)

type OptionSwitch struct {
	arg      string
	menuItem *systray.MenuItem
}

func (sw *OptionSwitch) Toggle() {
	if sw.menuItem.Checked() {
		sw.menuItem.Uncheck()
	} else {
		sw.menuItem.Check()
	}
}

func (sw *OptionSwitch) On() bool {
	return sw.menuItem.Checked()
}

type Application struct {
	optionSystem    *OptionSwitch
	optionUser      *OptionSwitch
	optionDisplay   *OptionSwitch
	controlEnable   *systray.MenuItem
	controlDisable  *systray.MenuItem
	controlSchedule *systray.MenuItem
	status          *systray.MenuItem
	quit            *systray.MenuItem
}

const (
	APPLICATION_STATUS_ENABLED = iota
	APPLICATION_STATUS_DISABLED
	APPLICATION_SCHEDULER_ENABLED
	APPLICATION_SCHEDULER_DISABLED
	APPLICATION_STATUS
	APPLICATION_QUIT
)

type ApplicationEvent struct {
	Event int
	Args  []string
}

func NewApplication() *Application {
	options := systray.AddMenuItem("Options", "Options")
	control := systray.AddMenuItem("Control", "Control")
	systray.AddSeparator()
	return &Application{
		optionSystem: &OptionSwitch{
			"-s",
			options.AddSubMenuItem("System", "Prevent system from sleeping"),
		},
		optionUser: &OptionSwitch{
			"-u",
			options.AddSubMenuItem("User", "Prevent inactive user"),
		},
		optionDisplay: &OptionSwitch{
			"-d",
			options.AddSubMenuItem("Display", "Prevent display from sleeping"),
		},
		controlEnable:   control.AddSubMenuItem("Enable", "Enable"),
		controlDisable:  control.AddSubMenuItem("Disable", "Disable"),
		controlSchedule: control.AddSubMenuItem("Scheduled", "Scheduled"),
		status:          systray.AddMenuItem("Stopped", "Status"),
		quit:            systray.AddMenuItem("Quit", "Quit Coffee Tray"),
	}
}

func (app *Application) InitDefaults() {
	app.optionSystem.menuItem.Check()
	app.optionUser.menuItem.Check()
	app.optionDisplay.menuItem.Check()
	app.controlEnable.Check()
	app.status.Disable()
}

func (app *Application) Loop() chan Message {
	ch := make(chan Message)
	cmdLoop := NewCommandSupervisor("caffeinate", app.BuildArgs())
	cmdLoopCh := cmdLoop.Run()
	app.status.SetTitle("Init")
	cmdLoopCh <- &ControlMessage{"START", nil}
	app.status.SetTitle((<-cmdLoopCh).Subject())

	go func() {
		for doRun := true; doRun; {
			select {
			case <-app.optionDisplay.menuItem.ClickedCh:
				fmt.Println("Display")
				app.optionDisplay.Toggle()
				cmdLoopCh <- &ControlMessage{"APPLY", app.BuildArgs()}

			case <-app.optionUser.menuItem.ClickedCh:
				fmt.Println("User")
				app.optionUser.Toggle()
				cmdLoopCh <- &ControlMessage{"APPLY", app.BuildArgs()}

			case <-app.optionSystem.menuItem.ClickedCh:
				fmt.Println("System")
				app.optionSystem.Toggle()
				cmdLoopCh <- &ControlMessage{"APPLY", app.BuildArgs()}

			case <-app.controlEnable.ClickedCh:
				app.controlEnable.Check()
				app.controlDisable.Uncheck()
				app.controlSchedule.Uncheck()
				cmdLoopCh <- &ControlMessage{"START", app.BuildArgs()}

			case <-app.controlDisable.ClickedCh:
				app.controlDisable.Check()
				app.controlEnable.Uncheck()
				app.controlSchedule.Uncheck()
				cmdLoopCh <- &ControlMessage{"STOP", app.BuildArgs()}

			case <-app.controlSchedule.ClickedCh:
				app.controlSchedule.Check()
				app.controlEnable.Uncheck()
				app.controlDisable.Uncheck()
				cmdLoopCh <- &ControlMessage{"SCHEDULE", nil}

			case <-app.quit.ClickedCh:
				app.status.SetTitle("Stopping...")
				cmdLoopCh <- &ControlMessage{"QUIT", nil}
				doRun = false
				ch <- &ControlMessage{"QUIT", nil}

			case msg := <-cmdLoopCh:
				app.status.SetTitle(msg.Subject())
				ch <- msg.(*StatusMessage)
			}
		}
	}()

	return ch
}

func (app *Application) BuildArgs() []string {
	var args []string
	if app.optionSystem.On() {
		args = append(args, app.optionSystem.arg)
	}
	if app.optionUser.On() {
		args = append(args, app.optionUser.arg)
	}
	if app.optionDisplay.On() {
		args = append(args, app.optionDisplay.arg)
	}

	return args
}
