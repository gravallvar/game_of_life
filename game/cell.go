package game

import (
	ui "github.com/gizak/termui/v3"
)

var deadStyle = ui.NewStyle(ui.ColorClear)
var aliveStyle = ui.NewStyle(ui.ColorClear, ui.ColorCyan)

type cell struct {
	alive bool
	x     int
	y     int
}

func (c cell) getStyle() ui.Style {
	if c.alive {
		return aliveStyle
	}

	return deadStyle
}
