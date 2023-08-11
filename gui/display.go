package gui

import (
	"fmt"
	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/logfire-sh/cli/internal/config"
	"github.com/logfire-sh/cli/pkg/cmd/sources/models"
	"github.com/logfire-sh/cli/pkg/cmdutil/APICalls"
	"github.com/rivo/tview"
	"log"
	"net/http"
	"strings"
)

type Display struct {
	*tview.Grid
	View             *tview.TextView
	input            *tview.InputField
	PlaceholderField *tview.InputField
	//SelectedSourceIDs *[]string
	List       *tview.List
	Window     *winman.WindowBase
	App        *tview.Application
	SourceList []models.Source
}

//type listItem struct {
//	MainText      string // The main text of the list item.
//	SecondaryText string // A secondary text to be shown underneath the main text.
//	Shortcut      rune   // The key to select the list item directly, 0 if there is no shortcut.
//	Selected      func() // The optional function which is called when the item is selected.
//}

type Task struct {
	Title string `json:"text"`
	Id    string
}

//type company struct {
//	Name string `json:"name"`
//}

var wordList []string
var schemaList []string

func NewDisplay(cfg config.Config) *Display {
	//var sourcesTask []Task
	//selectedSources := make(map[string]bool)
	var sourceIds []string

	app := tview.NewApplication()

	client := &http.Client{}

	sourcesList, err := APICalls.GetAllSources(client, cfg.Get().Token, cfg.Get().EndPoint, cfg.Get().TeamId)
	if err != nil {
		log.Fatalln(fmt.Sprint(err))
	}

	for _, source := range sourcesList {
		wordList = append(wordList, source.Name)
		sourceIds = append(sourceIds, source.ID)
		//task := Task{Title: source.Name, Id: source.ID}
		//sourcesTask = append(sourcesTask, task)
		//s := makeTaskListingTitle(task, selectedSources[source.Name])
		//list.AddItem(s, "", 0, func() {})
	}

	schemaMap, err := APICalls.GetSchema(cfg.Get().Token, cfg.Get().EndPoint, cfg.Get().TeamId, sourceIds)
	if err != nil {
		log.Fatalln(fmt.Sprint(err))
	}

	// Assuming the schemaMap is of type []map[string]string
	fieldTypeMap := make(map[string]string)
	for _, item := range schemaMap {
		for key, value := range item {
			fieldTypeMap[key] = value
		}
	}

	// Define the operator options
	var stringOptions = []string{"=", "!=", ":", "!:"}
	var integerOptions = []string{"=", "!=", ">", "<", ">=", "<="}
	var booleanOptions = []string{"="}

	for _, item := range schemaMap {
		for key := range item {
			schemaList = append(schemaList, key)
		}
	}

	//list := tview.NewList().ShowSecondaryText(false).SetSelectedBackgroundColor(tcell.ColorDarkBlue)

	//var selectedSourceIDs []string

	inputField := tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0).
		SetAcceptanceFunc(tview.InputFieldMaxLength(200)).
		SetFieldStyle(tcell.StyleDefault)

	// Set up autocomplete function.
	var typedText string
	inputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}

		typedText = currentText

		if strings.HasPrefix(typedText, "field-filter=") {
			field := strings.TrimPrefix(typedText, "field-filter=")
			if contains(schemaList, field) {
				// If a field is already selected, suggest the appropriate operators
				fieldType, exists := fieldTypeMap[field]
				if exists {
					switch fieldType {
					case "string":
						entries = append(entries, stringOptions...)
					case "int":
						entries = append(entries, integerOptions...)
					case "bool":
						entries = append(entries, booleanOptions...)
					}
				}
			} else {
				// If the current text is "field-filter=", show the schemaList suggestions
				for _, word := range schemaList {
					if strings.HasPrefix(strings.ToLower(word), strings.ToLower(field)) && !strings.Contains(currentText, word) {
						entries = append(entries, word)
					}
				}
			}
		} else {
			for _, word := range wordList {
				parts := strings.Split(currentText, "=")

				// If there's no "=", use the entire currentText as the prefix
				if len(parts) < 2 {
					if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) && !strings.Contains(currentText, word) {
						entries = append(entries, word)
					}
					continue
				}

				// Split the right side of "=" by commas to get the sources
				sources := strings.Split(parts[1], ",")
				prefix := sources[len(sources)-1] // Use the last source as the prefix

				// Check if the word starts with the prefix (case-insensitive)
				// and is not already in the currentText
				if strings.HasPrefix(strings.ToLower(word), strings.ToLower(prefix)) && !strings.Contains(currentText, word) {
					entries = append(entries, word)
				}
			}
		}

		if len(entries) < 1 {
			entries = nil
		}
		return
	})
	inputField.SetAutocompletedFunc(func(text string, index, source int) bool {
		if source != tview.AutocompletedNavigate {
			if strings.Contains(typedText, "source=") {
				// Split the typedText by commas
				parts := strings.Split(typedText, ",")

				// Check if the last part is "source="
				if parts[len(parts)-1] == "source=" {
					parts[len(parts)-1] = "source=" + text
				} else {
					// Replace the last part with the selected word
					parts[len(parts)-1] = text
				}

				// Join the parts back together
				updatedText := strings.Join(parts, ",")

				inputField.SetText(updatedText)
			} else if strings.HasPrefix(typedText, "field-filter=") {
				// If the current text is "field-filter=", show the schemaList suggestions
				if contains(schemaList, text) {
					inputField.SetText("field-filter=" + text)
				} else {
					// If a field is already selected, suggest the appropriate operators
					field := strings.TrimPrefix(typedText, "field-filter=")
					if fieldType, exists := fieldTypeMap[field]; exists {
						switch fieldType {
						case "string":
							if contains(stringOptions, text) {
								inputField.SetText("field-filter=" + field + text)
							}
						case "int":
							if contains(integerOptions, text) {
								inputField.SetText("field-filter=" + field + text)
							}
						case "bool":
							if contains(booleanOptions, text) {
								inputField.SetText("field-filter=" + field + text)
							}
						}
					}
				}
			} else {
				inputField.SetText(text)
			}
		}
		return source == tview.AutocompletedTab || source == tview.AutocompletedClick
	})

	PlaceholderField := tview.NewInputField().
		SetFieldWidth(0).
		SetFieldStyle(tcell.StyleDefault).
		SetPlaceholder("  1.source [source=source-name,source-name,source-name...] 2.start-date [start-date=now-2d] 3.end-date [end-date=now] 4.field-filter [field-filter=level=info] 5.save-view [save-view=name] 6.QUIT [q | quit | exit]").
		SetPlaceholderTextColor(tcell.ColorGray)

	PlaceholderField.SetDisabled(true)

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	//var items []listItem
	//
	//var item = listItem{
	//	MainText: "List1",
	//}

	//items = append(items, item)

	// Initialize the WindowManager with the tview application.
	//wm := winman.NewWindowManager()

	// Create the winman window and set its root to the list.
	//window := wm.NewWindow().
	//	Show().
	//	SetRoot(list).
	//	SetDraggable(false).
	//	SetResizable(false).
	//	SetTitle("Select sources to filter by:")
	//
	//window.SetRect(5, 5, 30, 10)

	// Create the Grid and add items to it.
	grid := tview.NewGrid().SetRows(-1, 1, 1).SetColumns(-1)
	grid.AddItem(textView, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(inputField, 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(PlaceholderField, 2, 0, 1, 1, 0, 0, false)

	//window.
	//	AddButton(&winman.Button{
	//		Symbol: 'X',
	//		OnClick: func() {
	//			grid.SetRows(-1, 1, 1).SetColumns(-1)
	//			grid.Clear()
	//			grid.AddItem(textView, 0, 0, 1, 1, 0, 0, false)
	//			grid.AddItem(inputField, 1, 0, 1, 1, 0, 0, true)
	//			grid.AddItem(PlaceholderField, 2, 0, 1, 1, 0, 0, false)
	//			app.SetFocus(inputField)
	//		},
	//	})
	//
	//list.AddItem("Quit", "Press to exit", 'q', func() {
	//	grid.SetRows(-1, 1, 1).SetColumns(-1)
	//	grid.Clear()
	//	grid.AddItem(textView, 0, 0, 1, 1, 0, 0, false)
	//	grid.AddItem(inputField, 1, 0, 1, 1, 0, 0, true)
	//	grid.AddItem(PlaceholderField, 2, 0, 1, 1, 0, 0, false)
	//	app.SetFocus(inputField)
	//})

	//list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	//	switch event.Key() {
	//	case tcell.KeyEnter:
	//		// Handle Enter key
	//		itemIndex := list.GetCurrentItem()
	//		sourceID := sourcesTask[itemIndex].Id
	//
	//		// Toggle the selection status of the source
	//		if selectedSources[sourceID] {
	//			selectedSources[sourceID] = false
	//			for i, id := range selectedSourceIDs {
	//				if id == sourceID {
	//					selectedSourceIDs = append(selectedSourceIDs[:i], selectedSourceIDs[i+1:]...)
	//					break
	//				}
	//			}
	//		} else {
	//			selectedSources[sourceID] = true
	//			selectedSourceIDs = append(selectedSourceIDs, sourceID)
	//		}
	//
	//		// Update the list item's title with the new checkbox status
	//		newTitle := makeTaskListingTitle(sourcesTask[itemIndex], selectedSources[sourceID])
	//		list.SetItemText(itemIndex, newTitle, "")
	//
	//	case tcell.KeyRune:
	//		switch event.Rune() {
	//		case 'q':
	//			grid.RemoveItem(list)
	//		}
	//	}
	//	return event
	//})

	return &Display{
		Grid:             grid,
		View:             textView,
		input:            inputField,
		PlaceholderField: PlaceholderField,
		//List:              list,
		//SelectedSourceIDs: &selectedSourceIDs,
		//Window:            window,
		App:        app,
		SourceList: sourcesList,
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

//func makeTaskListingTitle(task Task, selected bool) string {
//	checkbox := "[ ]"
//	if selected {
//		checkbox = "[x[]"
//	}
//
//	return fmt.Sprintf("[%s]%s %s [%s]", "smokewhite", checkbox, task.Title, task.Id[:4])
//}
