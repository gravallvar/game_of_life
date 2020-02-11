package main

import (
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/gravallvar/game_of_life/game"
)

var paused bool = false

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "Controls"
	l.Rows = []string{
		"[p] Pause",
		"[c] Clear grid",
		"[r] Random grid",
		"[q, ctrl+c] Quit",
		"[arrow keys] Move cursor",
		"--- While paused ---",
		"[n] Next generation",
		"[space] Flip cell",
	}
	l.SetRect(68, 0, 100, 10)

	grid := game.NewGrid(true)
	grid.Title = "Game of life"
	grid.SetRect(0, 0, 66, 24)

	go gameLoop(grid, l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "p":
			paused = !paused
		case "r":
			grid.ResetCells(true)
		case "c":
			grid.ResetCells(false)
		case "n":
			if paused {
				grid.NextGeneration()
			}
		case "<Up>":
			grid.MoveCursor(game.Up)
		case "<Down>":
			grid.MoveCursor(game.Down)
		case "<Left>":
			grid.MoveCursor(game.Left)
		case "<Right>":
			grid.MoveCursor(game.Right)
		case "<Space>":
			if paused {
				grid.FlipCell()
			}
		}
	}
}

func gameLoop(grid *game.Grid, list *widgets.List) {
	for {
		if !paused {
			grid.NextGeneration()
		}
		ui.Render(grid, list)
		time.Sleep(80 * time.Millisecond)
	}
}
