package main

import (
	"fmt"
	"log/syslog"
	"os"

	coffeetray "durold.de/coffeetray/internal"
	"github.com/getlantern/systray"
)

func main() {
	logwriter, err := syslog.New(syslog.LOG_ERR|syslog.LOG_SYSLOG, "coffeetray")
	coffeetray.AssertNoError(err)
	logwriter.Info("Starting Coffeetray")
	fmt.Println(os.Getwd())
	systray.Run(onReady, onExit)
}

func onReady() {
	openIcon, err := coffeetray.GetPngIconBuffer("/Users/patrick.durold/Pictures/EyeConOpen.png")
	coffeetray.AssertNoError(err)
	openIconBytes := openIcon.Bytes()
	closedIcon, err := coffeetray.GetPngIconBuffer("/Users/patrick.durold/Pictures/EyeConClosed.png")
	coffeetray.AssertNoError(err)
	closedIconBytes := closedIcon.Bytes()

	schedulerClosedIcon, err := coffeetray.GetPngIconBuffer("/Users/patrick.durold/Pictures/EyeConScheduleClosed.png")
	coffeetray.AssertNoError(err)
	schedulerClosedIconBytes := schedulerClosedIcon.Bytes()

	schedulerOpenIcon, err := coffeetray.GetPngIconBuffer("/Users/patrick.durold/Pictures/EyeConScheduleOpen.png")
	coffeetray.AssertNoError(err)
	schedulerOpenIconBytes := schedulerOpenIcon.Bytes()

	systray.SetTemplateIcon(openIconBytes, openIconBytes)
	systray.SetTooltip("Coffee Tray")

	app := coffeetray.NewApplication()
	app.InitDefaults()

	ch := app.Loop()
	for {
		message := <-ch
		switch typedMessage := message.(type) {
		case *coffeetray.ControlMessage:
			if typedMessage.Subject() == "QUIT" {
				systray.Quit()
				return
			}
		case *coffeetray.StatusMessage:
			switch typedMessage.Status() {
			case coffeetray.Running:
				systray.SetTemplateIcon(openIconBytes, openIconBytes)
			case coffeetray.Stopped:
				systray.SetTemplateIcon(closedIconBytes, closedIconBytes)
			case coffeetray.ScheduleRunning:
				systray.SetTemplateIcon(schedulerOpenIconBytes, schedulerOpenIconBytes)
			case coffeetray.ScheduleIdle:
				systray.SetTemplateIcon(schedulerClosedIconBytes, schedulerClosedIconBytes)
			}
		}
	}
}

func onExit() {
	os.Exit(0)
}
