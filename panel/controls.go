package panel

import (
	"github.com/gizak/termui/v3/widgets"
)

// Controls - side panel displaying controls.
// Wraps widgets.List to allow easy rendering.
type Controls struct {
	*widgets.List
	width  int
	height int
}

// NewControls creates a new controls side panel.
func NewControls(x, y, width, height int) *Controls {
	l := widgets.NewList()
	l.Title = "Controls"
	l.Rows = []string{
		"[arrow keys] Move cursor",
		"[1,2,3]      Set speed",
		"[space]      Pause",
		"[c]          Clear grid",
		"[r]          Random grid",
		"[q,ctrl+c]   Quit",
		"",
		"----- While paused -----",
		"[n]   Next generation",
		"[tab] Flip cell",
	}

	c := &Controls{
		List:   l,
		width:  width,
		height: height,
	}
	c.SetPosition(x, y, width, height)

	return c
}

// SetPosition sets the position of the panel.
func (c *Controls) SetPosition(x, y, width, height int) {
	c.SetRect(x, y, x+width, y+height)
}
