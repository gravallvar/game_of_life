package game

import (
	ui "github.com/gizak/termui/v3"
)

var deadStyle = ui.NewStyle(ui.ColorClear)
var aliveStyle = ui.NewStyle(ui.ColorClear, ui.ColorGreen)

// getCellStyle returns the appropriate style for a given cell state.
func getCellStyle(alive bool) ui.Style {
	if alive {
		return aliveStyle
	}

	return deadStyle
}
