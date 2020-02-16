package game

import (
	ui "github.com/gizak/termui/v3"
)

// Direction - enum for cursor direction
type Direction int

// Cursor directions.
const (
	Up    Direction = iota
	Down  Direction = iota
	Left  Direction = iota
	Right Direction = iota
)

type cursor struct {
	x int
	y int
}

func (c *cursor) setCoordinates(x, y int) {
	c.x = x
	c.y = y
}

var deadMarkedStyle = ui.NewStyle(ui.ColorClear, ui.ColorCyan)
var aliveMarkedStyle = ui.NewStyle(ui.ColorClear, ui.ColorYellow)

// getStyle returns the appropriate style for a cursor. Different styles are
// returned depending on the state of the marked cell.
func (c *cursor) getStyle(aliveMarked bool) ui.Style {
	if aliveMarked {
		return aliveMarkedStyle
	}

	return deadMarkedStyle
}
