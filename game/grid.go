package game

import (
	"image"

	ui "github.com/gizak/termui/v3"
	"github.com/gravallvar/game_of_life/util"
)

const gridWidth, gridHeight = 64, 22

type Grid struct {
	ui.Block
	currentCells [][]bool
	nextCells    [][]bool
}

func newCellMatrix(random bool) [][]bool {
	cells := make([][]bool, gridWidth)

	for x := 0; x < gridWidth; x++ {
		cells[x] = make([]bool, gridHeight)
		for y := 0; y < gridHeight; y++ {
			cell := false
			if random {
				cell = util.GetRandomCellState()
			}

			cells[x][y] = cell
		}
	}

	return cells
}

func NewGrid(random bool) *Grid {
	return &Grid{
		Block:        *ui.NewBlock(),
		currentCells: newCellMatrix(random),
		nextCells:    newCellMatrix(false),
	}
}

func (g *Grid) setCell(b *ui.Buffer, c bool, x, y int) {
	b.SetCell(ui.NewCell(' ', getCellStyle(c)), image.Pt(g.Inner.Min.X+x, g.Inner.Min.Y+y))
}

func (g *Grid) Draw(b *ui.Buffer) {
	g.Block.Draw(b)

	for x := range g.currentCells {
		for y := range g.currentCells[x] {
			g.setCell(b, g.currentCells[x][y], x, y)
		}
	}
}

func (g *Grid) NextGeneration() {
	for x := range g.currentCells {
		for y := range g.currentCells[x] {
			aliveAdjacents := g.getAliveAdjacents(x, y)
			switch {
			case !g.currentCells[x][y] && aliveAdjacents == 3:
				g.nextCells[x][y] = true
			case g.currentCells[x][y] && (aliveAdjacents == 2 || aliveAdjacents == 3):
				g.nextCells[x][y] = true
			default:
				g.nextCells[x][y] = false
			}
		}
	}

	g.currentCells = g.nextCells
	g.nextCells = newCellMatrix(false)
}

func (g *Grid) getAliveAdjacents(x, y int) int {
	adjacents := make([]bool, 0, 8)

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
		if c {
			aliveAdjacents++
		}
	}

	return aliveAdjacents
}
