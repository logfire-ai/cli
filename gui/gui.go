package gui

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/livetail"
	"github.com/logfire-sh/cli/pkg/cmd/factory"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/filters"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app     *tview.Application
	Display *Display
	logs    string

	f *cmdutil.Factory

	Config config.Config

	Livetail *livetail.Livetail

	StartDateTimeFilter       time.Time
	EndDateTimeFilter         time.Time
	SourceFilter              []string
	SearchFilter              []string
	FieldBasedFilterName      string
	FieldBasedFilterValue     string
	FieldBasedFilterCondition string

	StartDateUnParsed string
	EndDateUnParsed   string

	Ctx context.Context
}

type LivetailStatus struct {
	LivetailEnabled bool
}

var livetailStatus = &LivetailStatus{
	LivetailEnabled: false,
}

func NewUI(cfg config.Config) *UI {

	displayInstance := NewDisplay(cfg)
	ui := &UI{
		Config:  cfg,
		Display: displayInstance,
		app:     displayInstance.App,
	}
	ui.app.EnableMouse(true)
	ui.SetDisplayCapture()
	ui.Display.Livetail = true

	ui.f = factory.New()

	err := errors.New("")
	ui.Livetail, err = livetail.NewLivetail()
	if err != nil {
		ui.Display.View.SetText(fmt.Sprintf("Error while initiating livetail:\n%s", err.Error()))
		ui.app.Stop()
	}

	// go checkWaitingForLogs(ui, ui.Livetail)

	ui.Ctx = context.Background()

	ui.Livetail.CreateConnection()

	time.Sleep(200 * time.Millisecond)

	RunLivetail(ui, livetailStatus)
	return ui
}

func display(u *UI, l *livetail.Livetail) {
	var numDots int

	for {
		select {
		case <-u.Ctx.Done():
			return
		default:
			u.app.QueueUpdateDraw(func() {
				// Case 1: If l.Logs is empty, show "Waiting for logs..." with progress dots
				if len(l.Logs) == 0 {
					numDots = (numDots + 1) % 4
					var waitingMessage string
					if u.Config.Get().Theme == "dark" {
						waitingMessage = "[white]" + "Waiting for logs" + strings.Repeat(".", numDots)
					} else {
						waitingMessage = "[black]" + "Waiting for logs" + strings.Repeat(".", numDots)
					}
					// Pad the message with spaces to keep it a constant length
					paddedMessage := fmt.Sprintf("%-40s", waitingMessage)
					u.Display.View.SetTextAlign(tview.AlignCenter)
					u.Display.View.SetText(paddedMessage)
				} else {
					// Reset the numDots counter
					numDots = 0

					// Case 4: If l.Logs is not empty and only contains "Waiting for logs...", do nothing
					if l.Logs == "Waiting for logs..." {
						return
					}

					// Case 2: If l.Logs is not empty and contains "Waiting for logs..." along with other text, remove "Waiting for logs..." and keep the rest.
					if strings.Contains(l.Logs, "Waiting for logs") {
						updatedLogs := strings.ReplaceAll(l.Logs, "Waiting for logs...", "")
						updatedLogs = strings.ReplaceAll(updatedLogs, "Waiting for logs..", "")
						updatedLogs = strings.ReplaceAll(updatedLogs, "Waiting for logs.", "")
						updatedLogs = strings.ReplaceAll(updatedLogs, "Waiting for logs", "")
						u.logs = strings.TrimSpace(updatedLogs) // remove extra spaces if any
						u.Display.View.SetText(u.logs)
						u.Display.View.SetTextAlign(tview.AlignLeft)
					} else {
						// Case 3: If l.Logs is not empty and doesn't contain "Waiting for logs...", do nothing.
						u.logs = l.Logs
						u.Display.View.SetText(l.Logs)
						u.Display.View.SetTextAlign(tview.AlignLeft)
					}
				}
			})

			time.Sleep(500 * time.Millisecond)
		}
	}

}

func splitFieldFilterValue(input string) (field, operator, value string) {
	filterValue := strings.SplitN(input, "field-filter=", 2)[1]

	// Regular expression to match the operators
	re := regexp.MustCompile(`(=|!=|:|!:|>|<|>=|<=)`)

	// Split the filterValue using the regex
	parts := re.Split(filterValue, -1)
	if len(parts) < 2 {
		return filterValue, "", "" // Only the field is present
	}

	// Extract the operator using regex
	op := re.FindString(filterValue)

	return parts[0], op, parts[1]
}

var sourceNamesList []string
var sourceIds []string

func (u *UI) runQuitCmd() {
	StopLivetail(u, livetailStatus)
	u.app.Stop()
}

func (u *UI) SetDisplayCapture() {
	u.Display.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			input := u.Display.input.GetText()
			u.Display.input.SetText("")

			if len(strings.Split(input, "=")) > 1 {
				if u.Display.Livetail {
					if strings.Split(input, "=")[0] == "source" {
						sourceNames := strings.Split(input, "=")[1]
						sourceNamesList = strings.Split(sourceNames, ",")

						for _, source := range sourceNamesList {
							for _, sourceO := range u.Display.SourceList {
								if source == sourceO.Name {
									sourceIds = append(sourceIds, sourceO.ID)
								}
							}
						}

						StopLivetail(u, livetailStatus)

						u.SourceFilter = sourceIds

						time.Sleep(500 * time.Millisecond)

						RunLivetail(u, livetailStatus)

					} else if strings.Split(input, "=")[0] == "start-date" {
						StopLivetail(u, livetailStatus)

						u.StartDateTimeFilter = filters.ShortDateTimeToGoDate(strings.Split(input, "=")[1])
						u.StartDateUnParsed = strings.Split(input, "=")[1]

						time.Sleep(500 * time.Millisecond)

						RunLivetail(u, livetailStatus)

					} else if strings.Split(input, "=")[0] == "end-date" {
						StopLivetail(u, livetailStatus)

						u.EndDateTimeFilter = filters.ShortDateTimeToGoDate(strings.Split(input, "=")[1])
						u.EndDateUnParsed = strings.Split(input, "=")[1]

						time.Sleep(500 * time.Millisecond)

						RunLivetail(u, livetailStatus)

					} else if strings.Split(input, "=")[0] == "field-filter" {
						StopLivetail(u, livetailStatus)

						field, operator, value := splitFieldFilterValue(input)

						u.FieldBasedFilterName = field
						u.FieldBasedFilterCondition = operator
						u.FieldBasedFilterValue = value

						time.Sleep(500 * time.Millisecond)

						RunLivetail(u, livetailStatus)

					} else if strings.Split(input, "=")[0] == "save-view" {
						name := strings.Split(input, "=")[1]

						var selectedSource []models.Source

						for _, source := range u.Display.SourceList {
							for _, sourceid := range sourceIds {
								if sourceid == source.ID {
									selectedSource = append(selectedSource, source)
								}
							}
						}

						err := APICalls.CreateView(u.Config.Get().Token, u.Config.Get().EndPoint, u.Config.Get().TeamId, selectedSource, []string{}, u.FieldBasedFilterName, u.FieldBasedFilterValue, u.FieldBasedFilterCondition, u.StartDateUnParsed, u.EndDateUnParsed, name)
						if err != nil {
							u.Display.input.SetFieldTextColor(tcell.ColorRed)
							input = "Failed to create view"

							go func() {
								time.Sleep(1000 * time.Millisecond)
								u.Display.input.SetFieldTextColor(tcell.ColorWhite)
								input = ""
							}()
						}
					}
					break
				} else {
					if strings.Split(input, "=")[0] == "view" {
						name := strings.Split(input, "=")[1]

						for _, view := range u.Display.ViewsList {
							if view.Name == name {
								u.StartDateTimeFilter = view.DateFilter.StartDate
								u.EndDateTimeFilter = view.DateFilter.EndDate

								if len(view.SourcesFilter) > 0 {
									for _, source := range view.SourcesFilter {
										u.SourceFilter = append(u.SourceFilter, source.ID)
									}
								}

								u.SearchFilter = view.TextFilter
								u.FieldBasedFilterName = view.SearchFilter[0].Key
								u.FieldBasedFilterValue = view.SearchFilter[0].Value
								u.FieldBasedFilterCondition = view.SearchFilter[0].Condition
							}
						}

						if livetailStatus.LivetailEnabled {
							StopLivetail(u, livetailStatus)

							time.Sleep(200 * time.Millisecond)
						}

						RunLivetail(u, livetailStatus)
					}
				}
				return nil
			}

			switch u.Display.Livetail {
			case true:
				switch input {
				case "start":
					RunLivetail(u, livetailStatus)
				case "stop":
					StopLivetail(u, livetailStatus)
				case "q":
					u.runQuitCmd()
				case "quit":
					u.runQuitCmd()
				case "exit":
					u.runQuitCmd()
				case "1":
				case "2":
					StopLivetail(u, livetailStatus)

					ResetFilters(u)

					u.logs = "\n Select a view below"
					u.Display.View.SetText(u.logs)

					u.Display.input.SetText("view=")
					u.Display.input.Autocomplete()

					u.Display.Livetail = false
					u.Display.TopHelp.SetPlaceholder("  Stream > View | 1. livetail 2. view 9.QUIT [q | quit | exit]")
					u.Display.BottomHelp.SetPlaceholder("  3.view [view=view-name]")
				case "3":
					u.Display.input.SetText("source=")
					u.Display.input.Autocomplete()
				case "4":
					u.Display.input.SetText("start-date=")
				case "5":
					u.Display.input.SetText("end-date=")
				case "6":
					u.Display.input.SetText("field-filter=")
					u.Display.input.Autocomplete()
				case "7":
					u.Display.input.SetText("save-view=")
				case "9":
					u.runQuitCmd()
				default:
					u.Display.BottomHelp.SetPlaceholder("  Invalid command").SetPlaceholderTextColor(tcell.ColorRed)

					go func() {
						time.Sleep(200 * time.Millisecond)

						u.Display.BottomHelp.SetPlaceholder("  3.source [source=source-name,source-name,source-name...] 4.start-date [start-date=now-2d] 5.end-date [end-date=now] 6.field-filter [field-filter=level=info] 7.save-view [save-view=name] 8.QUIT [q | quit | exit]").
							SetPlaceholderTextColor(tcell.ColorGray)
					}()
				}
			case false:
				switch input {
				case "start":
					RunLivetail(u, livetailStatus)
				case "stop":
					StopLivetail(u, livetailStatus)
				case "q":
					u.runQuitCmd()
				case "quit":
					u.runQuitCmd()
				case "exit":
					u.runQuitCmd()
				case "1":
					StopLivetail(u, livetailStatus)

					ResetFilters(u)

					time.Sleep(200 * time.Millisecond)

					RunLivetail(u, livetailStatus)

					u.Display.Livetail = true
					u.Display.TopHelp.SetPlaceholder("  Stream > Livetail | 1. livetail 2. view 9.QUIT [q | quit | exit]").
						SetPlaceholderTextColor(tcell.ColorGray)
					u.Display.BottomHelp.SetPlaceholder("  3.source [source=source-name,source-name,source-name...] 4.start-date [start-date=now-2d] 5.end-date [end-date=now] 6.field-filter [field-filter=level=info] 7.save-view [save-view=name]").
						SetPlaceholderTextColor(tcell.ColorGray)
				case "2":
				case "3":
					u.Display.input.SetText("view=")
					u.Display.input.Autocomplete()
				case "9":
					u.runQuitCmd()
				default:
					u.Display.BottomHelp.SetPlaceholder("  Invalid command").SetPlaceholderTextColor(tcell.ColorRed)

					go func() {
						time.Sleep(200 * time.Millisecond)

						u.Display.BottomHelp.SetPlaceholder("  3.view [view=view-name] 4.QUIT [q | quit | exit]").
							SetPlaceholderTextColor(tcell.ColorGray)
					}()
				}
			}

			return nil
		}
		return event
	})
}

var mu sync.Mutex

func RunLivetail(u *UI, livetailStatus *LivetailStatus) {
	mu.Lock()
	defer mu.Unlock()

	if !livetailStatus.LivetailEnabled {
		livetailStatus.LivetailEnabled = true

		u.logs = ""
		u.Livetail.Logs = ""
		u.Display.View.SetText(u.logs)

		u.Display.View.SetTextAlign(tview.AlignLeft)

		u.Livetail.ApplyFilter(u.Config, u.SourceFilter, u.StartDateTimeFilter, u.EndDateTimeFilter, u.FieldBasedFilterName, u.FieldBasedFilterValue, u.FieldBasedFilterCondition)
		go u.Livetail.GenerateLogs(u.Ctx, u.Config)
		go display(u, u.Livetail)
		u.Display.View.ScrollToEnd()
	}
}

func StopLivetail(u *UI, livetailStatus *LivetailStatus) {
	mu.Lock()
	defer mu.Unlock()

	if livetailStatus.LivetailEnabled {
		livetailStatus.LivetailEnabled = false
		_, cancel := context.WithCancel(u.Ctx)
		defer cancel()
	}
}

func ResetFilters(u *UI) {
	u.SourceFilter = nil
	u.StartDateTimeFilter = time.Time{}
	u.EndDateTimeFilter = time.Time{}
	u.FieldBasedFilterName = ""
	u.FieldBasedFilterValue = ""
	u.FieldBasedFilterCondition = ""
}

func (u *UI) Run() error {
	u.app.SetRoot(u.Display.Grid, true)
	if err := u.app.Run(); err != nil {
		fmt.Println("Failed to run application:", err)
		return err
	}
	return nil
}
