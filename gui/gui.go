package gui

import (
	"fmt"
	"logfire/livetail"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app     *tview.Application
	display *Display
	token   string
	teamId  string
}

type LivetailStatus struct {
	LivetailEnabled bool
	LivetailStopped bool
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
			u.display.view.SetText(fmt.Sprintf("Livetail:\n%s", l.Logs))
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
		LivetailStopped: true,
	}

	u.display.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			input := u.display.input.GetText()
			u.display.input.SetText("")
			switch input {
			case "livetail":
				if !livetailStatus.LivetailEnabled {
					livetail, err := livetail.NewLivetail(u.token, u.teamId)
					if err != nil {
						u.display.view.SetText(fmt.Sprintf("Errow while initiating livetail:\n%s", err.Error()))
						return nil
					}
					go livetail.GenerateLogs()
					go display(u, livetail, stop)
					u.display.view.ScrollToEnd()
					livetailStatus.LivetailEnabled = true
					livetailStatus.LivetailStopped = false
				}
				return nil

			case "stop":
				if !livetailStatus.LivetailStopped {
					livetailStatus.LivetailStopped = true
					stop <- true
				}
				return nil

			// Todo: add a pop to show invalid command
			default:
				// u.display.view.SetText(fmt.Sprintf("Invalid command: %s. Please use command [blue]livetail [white]for starting livetail and [blue]stop [white]to stop it.", input))
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
