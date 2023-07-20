package gui

import (
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"time"

	"github.com/logfire-sh/cli/livetail"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app     *tview.Application
	display *Display
	token   string
	teamId  string
	logs    string

	Config func() (config.Config, error)
}

type LivetailStatus struct {
	LivetailEnabled bool
}

func NewUI(token, teamId string) *UI {
	ui := &UI{
		app:     tview.NewApplication(),
		display: NewDisplay(),
		token:   token,
		teamId:  teamId,
	}
	ui.app.EnableMouse(true)
	ui.SetDisplayCapture()
	return ui
}

func display(u *UI, l *livetail.Livetail, stop chan bool) {
	for {
		u.app.QueueUpdateDraw(func() {
			u.display.view.SetText(fmt.Sprintf("Livetail Started:\n%s", l.Logs))
			u.logs = l.Logs
		})
		time.Sleep(500 * time.Millisecond)
		select {
		case <-stop:
			return
		default:
		}
	}
}

func (u *UI) SetDisplayCapture() {
	stop := make(chan bool)

	livetailStatus := LivetailStatus{
		LivetailEnabled: false,
	}

	cfg, _ := u.Config()

	u.display.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			input := u.display.input.GetText()
			u.display.input.SetText("")
			switch input {
			case "start":
				if !livetailStatus.LivetailEnabled {
					u.display.view.SetTextAlign(tview.AlignLeft)
					livetail, err := livetail.NewLivetail(u.token, u.teamId, cfg.Get().EndPoint)
					if err != nil {
						u.display.view.SetText(fmt.Sprintf("Errow while initiating livetail:\n%s", err.Error()))
						return nil
					}
					go livetail.GenerateLogs()
					go display(u, livetail, stop)
					u.display.view.ScrollToEnd()
					livetailStatus.LivetailEnabled = true
				}
				return nil

			case "stop":
				if livetailStatus.LivetailEnabled {
					livetailStatus.LivetailEnabled = false
					stop <- true
					u.logs += "\nLivetail stopped!"
				} else {
					u.logs += "\nLivetail is not running!"
				}
				u.display.view.SetText(u.logs)
				return nil

			// Todo: add a pop to show invalid command
			default:
				u.logs += fmt.Sprintf("\nInvalid command: %s. Please use command [green]start [white]for starting livetail and [green]stop [white]to stop it.", input)
				u.display.view.SetText(u.logs)
			}
			return nil
		}
		return event
	})
}

func (u *UI) Run() error {
	u.app.SetRoot(u.display.Grid, true)
	if err := u.app.Run(); err != nil {
		fmt.Println("Failed to run application:", err)
		return err
	}
	return nil
}
