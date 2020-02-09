package main

import (
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gravallvar/game_of_life/game"
)

var pause bool = false

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	grid := game.NewGrid(true)
	grid.Title = "Game of life"
	grid.SetRect(0, 0, 66, 24)

	go gameLoop(grid)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "p":
			pause = !pause
		}
	}
}

func gameLoop(grid *game.Grid) {
	for {
		if !pause {
			grid.NextGeneration()
		}
		ui.Render(grid)
		time.Sleep(80 * time.Millisecond)
	}
}
