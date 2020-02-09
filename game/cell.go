package game

import (
	ui "github.com/gizak/termui/v3"
)

var deadStyle = ui.NewStyle(ui.ColorClear)
var aliveStyle = ui.NewStyle(ui.ColorClear, ui.ColorCyan)

func getCellStyle(c bool) ui.Style {
	if c {
		return aliveStyle
	}

	return deadStyle
}
