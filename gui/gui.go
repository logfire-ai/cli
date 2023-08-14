package gui

import (
	"errors"
	"fmt"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/livetail"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/logfire-sh/cli/pkg/cmdutil/filters"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app     *tview.Application
	Display *Display
	logs    string

	Config config.Config

	StartDateTimeFilter       time.Time
	EndDateTimeFilter         time.Time
	SourceFilter              []string
	SearchFilter              []string
	FieldBasedFilterName      string
	FieldBasedFilterValue     string
	FieldBasedFilterCondition string

	StartDateUnParsed string
	EndDateUnParsed   string
}

type LivetailStatus struct {
	LivetailEnabled bool
}

var livetailStatus = &LivetailStatus{
	LivetailEnabled: false,
}

var stop = make(chan error)

func NewUI() *UI {
	cfg, _ := config.NewConfig()

	displayInstance := NewDisplay(cfg)
	ui := &UI{
		Config:  cfg,
		Display: displayInstance,
		app:     displayInstance.App,
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
				u.Display.View.SetText(fmt.Sprintf("%s", l.Logs))
				u.logs = l.Logs
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

func (u *UI) SetDisplayCapture() {
	u.Display.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			input := u.Display.input.GetText()
			u.Display.input.SetText("")

			if len(strings.Split(input, "=")) > 1 {
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

					u.SourceFilter = sourceIds

					ClearLogs(u, livetailStatus, stop)

					time.Sleep(500 * time.Millisecond)

					RunLivetail(u, livetailStatus, stop)

				} else if strings.Split(input, "=")[0] == "start-date" {
					u.StartDateTimeFilter = filters.ShortDateTimeToGoDate(strings.Split(input, "=")[1])
					u.StartDateUnParsed = strings.Split(input, "=")[1]

					ClearLogs(u, livetailStatus, stop)

					time.Sleep(500 * time.Millisecond)

					RunLivetail(u, livetailStatus, stop)

				} else if strings.Split(input, "=")[0] == "end-date" {
					u.EndDateTimeFilter = filters.ShortDateTimeToGoDate(strings.Split(input, "=")[1])
					u.EndDateUnParsed = strings.Split(input, "=")[1]

					ClearLogs(u, livetailStatus, stop)

					time.Sleep(500 * time.Millisecond)

					RunLivetail(u, livetailStatus, stop)

				} else if strings.Split(input, "=")[0] == "field-filter" {
					field, operator, value := splitFieldFilterValue(input)

					u.FieldBasedFilterName = field
					u.FieldBasedFilterCondition = operator
					u.FieldBasedFilterValue = value

					ClearLogs(u, livetailStatus, stop)
					ClearLogs(u, livetailStatus, stop)

					time.Sleep(500 * time.Millisecond)

					RunLivetail(u, livetailStatus, stop)

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
				} else if strings.Split(input, "=")[0] == "stream-view" {
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

					ClearLogs(u, livetailStatus, stop)

					time.Sleep(500 * time.Millisecond)

					RunLivetail(u, livetailStatus, stop)
				}

				break
			}

			switch input {
			case "start":
				RunLivetail(u, livetailStatus, stop)
			case "stop":
				StopLivetail(u, livetailStatus, stop)
			case "q":
				ClearLogs(u, livetailStatus, stop)
				os.Exit(0)
			case "quit":
				ClearLogs(u, livetailStatus, stop)
				os.Exit(0)
			case "exit":
				ClearLogs(u, livetailStatus, stop)
				os.Exit(0)
			case "1":
				u.Display.input.SetText("source=")
				u.Display.input.Autocomplete()
			case "2":
				u.Display.input.SetText("start-date=")
			case "3":
				u.Display.input.SetText("end-date=")
			case "4":
				u.Display.input.SetText("field-filter=")
				u.Display.input.Autocomplete()
			case "5":
				u.Display.input.SetText("save-view=")
			case "6":
				u.Display.input.SetText("stream-view=")
				u.Display.input.Autocomplete()
			case "7":
				ClearLogs(u, livetailStatus, stop)
				os.Exit(0)
			default:
				u.Display.PlaceholderField.SetPlaceholder("  Invalid command").SetPlaceholderTextColor(tcell.ColorRed)

				go func() {
					time.Sleep(1000 * time.Millisecond)

					u.Display.PlaceholderField.SetPlaceholder("  1.source [source=source-name,source-name,source-name...] 2.start-date [start-date=now-2d] 3.end-date [end-date=now] 4.field-filter [field-filter=level=info] 5.save-view [save-view=name] 6.stream-view [stream-view=view-name] 7.QUIT [q | quit | exit]").
						SetPlaceholderTextColor(tcell.ColorGray)
				}()
			}

			return nil
		}
		return event
	})
}

var mu sync.Mutex

func RunLivetail(u *UI, livetailStatus *LivetailStatus, stop chan error) {
	mu.Lock()
	defer mu.Unlock()

	if !livetailStatus.LivetailEnabled {
		livetailStatus.LivetailEnabled = true
		u.Display.View.SetTextAlign(tview.AlignLeft)

		livetail, err := livetail.NewLivetail()
		if err != nil {
			u.Display.View.SetText(fmt.Sprintf("Error while initiating livetail:\n%s", err.Error()))
			return
		}
		livetail.ApplyFilter(u.Config, u.SourceFilter, u.StartDateTimeFilter, u.EndDateTimeFilter, u.FieldBasedFilterName, u.FieldBasedFilterValue, u.FieldBasedFilterCondition)
		go livetail.GenerateLogs(stop)
		go display(u, livetail, stop)
		u.Display.View.ScrollToEnd()
	}
}

func StopLivetail(u *UI, livetailStatus *LivetailStatus, stop chan error) {
	mu.Lock()
	defer mu.Unlock()

	if livetailStatus.LivetailEnabled {
		livetailStatus.LivetailEnabled = false
		stop <- errors.New("")
		u.logs += "\nLivetail stopped!"
	} else {
		u.logs += "\nLivetail is not running!"
	}
	u.Display.View.SetText(u.logs)
	return
}

func ClearLogs(u *UI, livetailStatus *LivetailStatus, stop chan error) {
	livetailStatus.LivetailEnabled = false
	stop <- errors.New("stop")
	stop <- errors.New("stop")
	u.logs = ""

	u.Display.View.SetText(u.logs)
}

func (u *UI) Run() error {
	u.app.SetRoot(u.Display.Grid, true)
	if err := u.app.Run(); err != nil {
		fmt.Println("Failed to run application:", err)
		return err
	}
	return nil
}
