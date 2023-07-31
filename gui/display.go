package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Display struct {
	*tview.Grid
	view  *tview.TextView
	input *tview.InputField
}

func NewDisplay() *Display {
	inputField := tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0).
		SetAcceptanceFunc(tview.InputFieldMaxLength(200)).
		SetFieldStyle(tcell.StyleDefault).SetPlaceholder("help").SetPlaceholderStyle(tcell.StyleDefault)

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true)

	grid := tview.NewGrid().SetRows(-1, 1).SetColumns(-1)
	grid.AddItem(inputField, 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(textView, 0, 0, 1, 1, 0, 0, false)

	return &Display{
		Grid:  grid,
		view:  textView,
		input: inputField,
	}
}
