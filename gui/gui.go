package gui

import (
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/livetail"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app     *tview.Application
	display *Display
	logs    string

	Config func() (config.Config, error)

	StartDateTimeFilter       string
	EndDateTimeFilter         string
	SourceFilter              []string
	SearchFilter              []string
	FieldBasedFilterName      string
	FieldBasedFilterValue     string
	FieldBasedFilterCondition string
}

type LivetailStatus struct {
	LivetailEnabled bool
}

var livetailStatus = &LivetailStatus{
	LivetailEnabled: false,
}

var stop = make(chan error)

func NewUI() *UI {
	ui := &UI{
		app:     tview.NewApplication(),
		display: NewDisplay(),
	}
	ui.app.EnableMouse(true)
	ui.SetDisplayCapture()
	RunLivetail(ui, livetailStatus, stop)
	return ui
}

func display(u *UI, l *livetail.Livetail, stop chan error) {
	for {
		select {
		case <-stop:
			return
		default:
			u.app.QueueUpdateDraw(func() {
				u.display.view.SetText(fmt.Sprintf("%s", l.Logs))
				u.logs = l.Logs
			})
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (u *UI) SetDisplayCapture() {
	u.display.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			input := u.display.input.GetText()
			u.display.input.SetText("")
			if input != "help" {
				u.display.input.SetPlaceholder("help")
			}

			//if strings.Split(input, "=")[0] == "source" {
			//	u.SourceFilter = append(u.SourceFilter, strings.Split(input, "=")[1])
			//	ClearLogs(u, livetailStatus, stop)
			//	RunLivetail(u, livetailStatus, stop)
			//} else {
			switch input {
			case "start":
				RunLivetail(u, livetailStatus, stop)
			case "stop":
				StopLivetail(u, livetailStatus, stop)
			//case "help":
			//	u.display.input.SetPlaceholder("source [source=source-id] startDate [startDate=now-2d] endDate [endDate=now-30d] fieldFilter [level=info] saveView [saveView=name]")
			default:
				u.logs += fmt.Sprintf("\nInvalid command: %s. Please use command [green]start [white]for starting livetail and [green]stop [white]to stop it.\n", input)
				StopLivetail(u, livetailStatus, stop)
			}
			//}

			return nil
		}
		return event
	})
}

func RunLivetail(u *UI, livetailStatus *LivetailStatus, stop chan error) {
	cfg, _ := config.NewConfig()

	if !livetailStatus.LivetailEnabled {
		livetailStatus.LivetailEnabled = true
		u.display.view.SetTextAlign(tview.AlignLeft)
		livetail, err := livetail.NewLivetail()
		if err != nil {
			u.display.view.SetText(fmt.Sprintf("Errow while initiating livetail:\n%s", err.Error()))
			return
		}
		livetail.ApplyFilter(cfg, u.SourceFilter, "now-30d", "")
		go livetail.GenerateLogs("", stop)
		go display(u, livetail, stop)
		u.display.view.ScrollToEnd()
	}
	return
}

func StopLivetail(u *UI, livetailStatus *LivetailStatus, stop chan error) {
	if livetailStatus.LivetailEnabled {
		livetailStatus.LivetailEnabled = false
		stop <- errors.New("")
		u.logs += "\nLivetail stopped!"
	} else {
		u.logs += "\nLivetail is not running!"
	}
	u.display.view.SetText(u.logs)
	return
}

func ClearLogs(u *UI, livetailStatus *LivetailStatus, stop chan error) {
	livetailStatus.LivetailEnabled = false
	stop <- errors.New("")
	u.logs = ""
	u.display.view.SetText(u.logs)
	return
}

func (u *UI) Run() error {
	u.app.SetRoot(u.display.Grid, true)
	if err := u.app.Run(); err != nil {
		fmt.Println("Failed to run application:", err)
		return err
	}
	return nil
}
