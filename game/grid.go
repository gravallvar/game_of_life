package game

import (
	"image"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gravallvar/game_of_life/util"
)

const (
	// margin for keep cells from being rendered outside the grid border.
	cellMargin = 2
	// minimum matrix width/height to simplify cell generation logic.
	minMatrixSize = 3
)

func newMatrix(random bool, width, height int) [][]bool {
	matrix := make([][]bool, width)

	for x := 0; x < width; x++ {
		matrix[x] = make([]bool, height)

		if random {
			for y := 0; y < height; y++ {
				matrix[x][y] = util.GetRandomCellState()
			}
		}
	}

	return matrix
}

// Grid - ui panel and state container for cell matrix.
// Overrides ui.Block to enable rendering cells.
type Grid struct {
	ui.Block
	cells [][]bool
	// Lock for avoiding concurrency issues during resize events.
	mutex  sync.Mutex
	cursor *cursor
	width  int
	height int
}

// NewGrid creates a new grid.
func NewGrid(width, height int, random bool) *Grid {
	g := &Grid{
		Block:  *ui.NewBlock(),
		cells:  newMatrix(random, width-cellMargin, height-cellMargin),
		cursor: &cursor{x: 0, y: 0},
		width:  width,
		height: height,
	}

	g.Title = "Game Of Life"
	g.SetRect(0, 0, width, height)

	return g
}

// SetSize resizes grid ui panel.
// NOTE: this function does not change the underlying cell matrix. Cells
// outside the panel are simply not rendered.
func (g *Grid) SetSize(x, y, width, height int) {
	g.Lock()
	defer g.Unlock()

	g.width = width
	g.height = height
	g.SetRect(x, y, x+width, x+height)
}

// GetWidth returns grid ui panel width.
func (g *Grid) GetWidth() int {
	return g.width
}

// GetHeight returns grid ui panel height.
func (g *Grid) GetHeight() int {
	return g.height
}

// ResetCells clears the current cell matrix and generates a new one.
// NOTE: this function takes potential ui panel resizing into account and
// generates a matrix of appropriate size.
func (g *Grid) ResetCells(random bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	matrixWidth := g.width - cellMargin
	if matrixWidth < minMatrixSize {
		matrixWidth = minMatrixSize
	}

	matrixheight := g.height - cellMargin
	if matrixheight < minMatrixSize {
		matrixheight = minMatrixSize
	}

	g.cells = newMatrix(random, matrixWidth, matrixheight)
}

func (g *Grid) drawCell(b *ui.Buffer, x, y int) {
	if x < g.width-cellMargin && y < g.height-cellMargin {
		b.SetCell(ui.NewCell(' ',
			getCellStyle(g.cells[x][y])),
			image.Pt(g.Inner.Min.X+x, g.Inner.Min.Y+y),
		)
	}
}

func (g *Grid) drawCursor(b *ui.Buffer) {
	x, y := g.cursor.x, g.cursor.y
	b.SetCell(ui.NewCell(' ',
		g.cursor.getStyle(g.cells[x][y])),
		image.Pt(g.Inner.Min.X+x, g.Inner.Min.Y+y),
	)
}

// Draw renders the grid ui panel. Necessary for implementing the termui
// Drawable interface.
func (g *Grid) Draw(b *ui.Buffer) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.Block.Draw(b)

	for x := range g.cells {
		for y := range g.cells[x] {
			g.drawCell(b, x, y)
		}
	}

	// Reset cursor position if its coordinates are outside the cell matrix.
	if g.cursor.x >= len(g.cells) || g.cursor.y >= len(g.cells[0]) {
		g.cursor.x = 0
		g.cursor.y = 0
	}
	g.drawCursor(b)
}

// NextGeneration generates the next generation of cells for the current matrix.
func (g *Grid) NextGeneration() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	nextCells := newMatrix(false, len(g.cells), len(g.cells[0]))

	for x := range g.cells {
		for y := range g.cells[x] {
			aliveAdjacents := g.getAliveAdjacents(x, y)
			switch {
			// Spawn new cell
			case !g.cells[x][y] && aliveAdjacents == 3:
				nextCells[x][y] = true
			// Cell is kept alive
			case g.cells[x][y] && (aliveAdjacents == 2 || aliveAdjacents == 3):
				nextCells[x][y] = true
			// Cell dies or remains dead.
			default:
				nextCells[x][y] = false
			}
		}
	}

	g.cells = nextCells
}

// getAliveAdjacents returns the number of living cells adjacent to the given
// cell in the matrix.
func (g *Grid) getAliveAdjacents(x, y int) int {
	adjacents := make([]bool, 0, 8)

	// Left column
	if x > 0 {
		switch y {
		case 0:
			adjacents = append(adjacents, g.cells[x-1][y:y+2]...)
		case len(g.cells[x]) - 1:
			adjacents = append(adjacents, g.cells[x-1][y-1:y+1]...)
		default:
			adjacents = append(adjacents, g.cells[x-1][y-1:y+2]...)
		}
	}

	// Middle column
	switch y {
	case 0:
		adjacents = append(adjacents, g.cells[x][y+1])
	case len(g.cells[x]) - 1:
		adjacents = append(adjacents, g.cells[x][y-1])
	default:
		adjacents = append(adjacents, g.cells[x][y-1], g.cells[x][y+1])
	}

	// Right column
	if x < len(g.cells)-1 {
		switch y {
		case 0:
			adjacents = append(adjacents, g.cells[x+1][y:y+2]...)
		case len(g.cells[x]) - 1:
			adjacents = append(adjacents, g.cells[x+1][y-1:y+1]...)
		default:
			adjacents = append(adjacents, g.cells[x+1][y-1:y+2]...)
		}
	}

	aliveAdjacents := 0
	for _, cell := range adjacents {
		if cell {
			aliveAdjacents++
		}
	}

	return aliveAdjacents
}

// MoveCursor moves cursor.
func (g *Grid) MoveCursor(d Direction) {
	c := g.cursor

	switch d {
	case Up:
		if c.y > 0 {
			g.cursor.setCoordinates(c.x, c.y-1)
		}
	case Down:
		if c.y < len(g.cells[c.x])-1 {
			g.cursor.setCoordinates(c.x, c.y+1)
		}
	case Left:
		if c.x > 0 {
			g.cursor.setCoordinates(c.x-1, c.y)
		}
	case Right:
		if c.x < len(g.cells)-1 {
			g.cursor.setCoordinates(c.x+1, c.y)
		}
	}
}

// FlipCell changes cell state for cell marked by cursor.
func (g *Grid) FlipCell() {
	x, y := g.cursor.x, g.cursor.y
	g.cells[x][y] = !g.cells[x][y]
}
