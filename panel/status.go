package panel

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/gravallvar/game_of_life/game"
)

// The following strings are labels for the status panel.
const (
	LabelPaused    = "          Paused"
	LabelSlow      = "           Slow"
	LabelFast      = "           Fast"
	LabelUnlimited = "        Unlimited"
)

var pausedStyle = ui.NewStyle(ui.ColorYellow)
var slowStyle = ui.NewStyle(ui.ColorGreen)
var fastStyle = ui.NewStyle(ui.ColorBlue)
var unlimitedStyle = ui.NewStyle(ui.ColorMagenta)

// Status - side panel displaying game status.
// Wraps widgets.List to allow easy rendering.
type Status struct {
	*widgets.List
	width   int
	height  int
	message string
}

// GetSpeedLabel returns the appropriate label for a given game speed.
func GetSpeedLabel(speed int) string {
	switch speed {
	case game.Slow:
		return LabelSlow
	case game.Fast:
		return LabelFast
	case game.Unlimited:
		return LabelUnlimited
	default:
		return ""
	}
}

func getStyle(msg string) ui.Style {
	switch msg {
	case LabelPaused:
		return pausedStyle
	case LabelSlow:
		return slowStyle
	case LabelFast:
		return fastStyle
	case LabelUnlimited:
		return unlimitedStyle
	default:
		return ui.NewStyle(ui.ColorWhite)
	}
}

// NewStatus creates a new status side panel.
func NewStatus(x, y, width, height int, msg string) *Status {
	l := widgets.NewList()
	l.TextStyle = getStyle(msg)

	// The reason for having three elements in this list is due to a current
	// bug in termui: text style is ignored for the first element.
	// https://github.com/gizak/termui/issues/248
	l.Rows = []string{
		"",
		msg,
		"",
	}

	s := &Status{
		List:    l,
		width:   width,
		height:  height,
		message: msg,
	}
	s.SetPosition(x, y, width, height)

	return s
}

// SetStatus sets the status message and styling.
func (s *Status) SetStatus(msg string) {
	s.message = msg
	s.TextStyle = getStyle(msg)
	s.Rows[1] = s.message
}

// SetPosition sets the position of the panel.
func (s *Status) SetPosition(x, y, width, height int) {
	s.SetRect(x, y, x+width, y+height)
}
