package game

import (
	"image"

	ui "github.com/gizak/termui/v3"
	"github.com/gravallvar/game_of_life/util"
)

const gridWidth, gridHeight = 64, 22

type Grid struct {
	ui.Block
	currentCells [][]cell
	nextCells    [][]cell
	cursor       *cursor
}

func NewGrid(random bool) *Grid {
	return &Grid{
		Block:        *ui.NewBlock(),
		currentCells: newCellMatrix(random),
		nextCells:    newCellMatrix(false),
		cursor:       &cursor{x: 0, y: 0},
	}
}

func newCellMatrix(random bool) [][]cell {
	cells := make([][]cell, gridWidth)

	for x := 0; x < gridWidth; x++ {
		cells[x] = make([]cell, gridHeight)
		for y := 0; y < gridHeight; y++ {
			cell := cell{
				alive: false,
				x:     x,
				y:     y,
			}
			if random {
				cell.alive = util.GetRandomCellState()
			}

			cells[x][y] = cell
		}
	}

	return cells
}

func (g *Grid) ResetCells(random bool) {
	g.currentCells = newCellMatrix(random)
	g.nextCells = newCellMatrix(false)
}

func (g *Grid) setCell(b *ui.Buffer, c cell) {
	b.SetCell(ui.NewCell(' ', c.getStyle()), image.Pt(g.Inner.Min.X+c.x, g.Inner.Min.Y+c.y))
}

func (g *Grid) Draw(b *ui.Buffer) {
	g.Block.Draw(b)

	for x := range g.currentCells {
		for y := range g.currentCells[x] {
			g.setCell(b, g.currentCells[x][y])
		}
	}

	b.SetCell(ui.NewCell(' ', g.cursor.getStyle()), image.Pt(g.Inner.Min.X+g.cursor.x, g.Inner.Min.Y+g.cursor.y))
}

func (g *Grid) NextGeneration() {
	for x := range g.currentCells {
		for y := range g.currentCells[x] {
			aliveAdjacents := g.getAliveAdjacents(x, y)
			switch {
			case !g.currentCells[x][y].alive && aliveAdjacents == 3:
				g.nextCells[x][y].alive = true
			case g.currentCells[x][y].alive && (aliveAdjacents == 2 || aliveAdjacents == 3):
				g.nextCells[x][y].alive = true
			default:
				g.nextCells[x][y].alive = false
			}
		}
	}

	g.currentCells = g.nextCells
	g.nextCells = newCellMatrix(false)
}

func (g *Grid) getAliveAdjacents(x, y int) int {
	adjacents := make([]cell, 0, 8)

	// Left column
	if x > 0 {
		switch {
		case y == 0:
			adjacents = append(adjacents, g.currentCells[x-1][y:y+2]...)
		case y == len(g.currentCells[y])-1:
			adjacents = append(adjacents, g.currentCells[x-1][y-1:y+1]...)
		default:
			adjacents = append(adjacents, g.currentCells[x-1][y-1:y+2]...)
		}
	}

	// Middle column
	switch y {
	case 0:
		adjacents = append(adjacents, g.currentCells[x][y+1])
	case len(g.currentCells[y]) - 1:
		adjacents = append(adjacents, g.currentCells[x][y-1])
	default:
		adjacents = append(adjacents, g.currentCells[x][y-1], g.currentCells[x][y+1])
	}

	// Right column
	if x < len(g.currentCells)-1 {
		switch {
		case y == 0:
			adjacents = append(adjacents, g.currentCells[x+1][y:y+2]...)
		case y == len(g.currentCells[y])-1:
			adjacents = append(adjacents, g.currentCells[x+1][y-1:y+1]...)
		default:
			adjacents = append(adjacents, g.currentCells[x+1][y-1:y+2]...)
		}
	}

	aliveAdjacents := 0
	for _, c := range adjacents {
		if c.alive {
			aliveAdjacents++
		}
	}

	return aliveAdjacents
}

func (g *Grid) MoveCursor(d Direction) {
	c := g.cursor
	switch d {
	case Up:
		g.cursor.setPoint(c.x, c.y-1)
	case Down:
		g.cursor.setPoint(c.x, c.y+1)
	case Left:
		g.cursor.setPoint(c.x-1, c.y)
	case Right:
		g.cursor.setPoint(c.x+1, c.y)
	}
}

func (g *Grid) FlipCell() {
	markedCell := g.currentCells[g.cursor.x][g.cursor.y]
	g.currentCells[g.cursor.x][g.cursor.y].alive = !markedCell.alive
}
