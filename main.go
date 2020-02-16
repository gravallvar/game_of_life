package main

import (
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gravallvar/game_of_life/game"
	"github.com/gravallvar/game_of_life/panel"
)

const (
	// Dimensions for side panels
	sidePanelsWidth = 27
	controlsHeight  = 12
	statusHeight    = 5
)

var paused bool
var speed = game.Fast

// Determines whether or not ui needs to cleared before being rendered.
var clearUI bool

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Create grid
	termWidth, termHeight := ui.TerminalDimensions()
	gridWidth := termWidth - sidePanelsWidth
	gridHeight := termHeight
	randomCells := true
	grid := game.NewGrid(gridWidth, gridHeight, randomCells)

	// Create side panels
	controls := panel.NewControls(gridWidth, 0, sidePanelsWidth, controlsHeight)
	status := panel.NewStatus(gridWidth, controlsHeight, sidePanelsWidth, statusHeight, panel.GetSpeedLabel(speed))

	// Game logic loop is run in a separate go-routine to handle asynchronous
	// control events.
	go gameLoop(grid, controls, status)

	eventLoop(grid, controls, status)
}

// gameLoop handles game logic and rendering. Every tick a new generation of
// cells are created and rendered unless the game is paused.
func gameLoop(grid *game.Grid, controls *panel.Controls, status *panel.Status) {
	for {
		if !paused {
			grid.NextGeneration()
		}

		// Clear terminal window before rendering if needed.
		if clearUI {
			ui.Clear()
			clearUI = false
		}
		ui.Render(grid, controls, status)

		// Limit framerate if needed.
		if speed != game.Unlimited && !paused {
			time.Sleep(time.Duration(speed) * time.Millisecond)
		}
	}
}

// eventLoop polls for new events and handles them. Events include key presses
// and terminal window resizing.
func eventLoop(grid *game.Grid, controls *panel.Controls, status *panel.Status) {
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		// Speed
		case "1":
			paused = false
			speed = game.Slow
			status.SetStatus(panel.LabelSlow)
		case "2":
			paused = false
			speed = game.Fast
			status.SetStatus(panel.LabelFast)
		case "3":
			paused = false
			speed = game.Unlimited
			status.SetStatus(panel.LabelUnlimited)
		case "<Space>":
			paused = !paused
			if paused {
				status.SetStatus(panel.LabelPaused)
			} else {
				status.SetStatus(panel.GetSpeedLabel(speed))
			}

		// Cursor
		case "<Up>":
			grid.MoveCursor(game.Up)
		case "<Down>":
			grid.MoveCursor(game.Down)
		case "<Left>":
			grid.MoveCursor(game.Left)
		case "<Right>":
			grid.MoveCursor(game.Right)
		case "<Tab>":
			if paused {
				grid.FlipCell()
			}

		// Grid
		case "r":
			grid.ResetCells(true)
		case "c":
			grid.ResetCells(false)
		case "n":
			if paused {
				grid.NextGeneration()
			}

		// Resize
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			grid.SetSize(0, 0, payload.Width-sidePanelsWidth, payload.Height)
			controls.SetPosition(grid.GetWidth(), 0, sidePanelsWidth, controlsHeight)
			status.SetPosition(grid.GetWidth(), controlsHeight, sidePanelsWidth, statusHeight)

			// After resize, terminal needs to be cleared before next render
			// to remove potential cell artefacts.
			clearUI = true

		// Quit
		case "q", "<C-c>":
			return
		}
	}
}
