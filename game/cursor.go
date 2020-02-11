package game

import (
	ui "github.com/gizak/termui/v3"
)

type Direction int

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

var cursorStyle = ui.NewStyle(ui.ColorClear, ui.ColorRed)

func (c *cursor) getStyle() ui.Style {
	return cursorStyle
}

func (c *cursor) setPoint(x, y int) {
	c.x = x
	c.y = y
}
