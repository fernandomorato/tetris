package tetris

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const KeyUp = 107
const KeyDown = 106
const KeyLeft = 104
const KeyRight = 108

type Tetris struct {
	rows                 int
	columns              int
	boardLeftOffset      int
	board                [][]int
	pendingPiece         Piece
	pendingPiecePosition Position
	hasPendingPiece      bool
	ticker               *time.Ticker
	mux                  *sync.RWMutex
}

type Position struct {
	x int
	y int
}

func (t *Tetris) Init() {
	t.rows = 20
	t.columns = 10
	t.boardLeftOffset = 20
	t.board = make([][]int, t.rows)
	for i := range t.rows {
		t.board[i] = make([]int, t.columns)
	}
	t.ticker = time.NewTicker(200 * time.Millisecond)
	t.mux = &sync.RWMutex{}
	go t.refresh()
}

func (t *Tetris) refresh() {
	for ; ; <-t.ticker.C {
		// if has no pending piece:
		// 1 - check if there are completed lines and update the board accordingly
		// 2 - spawns new piece
		if !t.hasPendingPiece {
			t.processCompletedLines()
			t.spawnNewPiece()
		}
		// update board status
		t.hasPendingPiece = t.MovePendingPiece(KeyDown)
		t.printBoard()
	}
}

func (t *Tetris) spawnNewPiece() {
	// new pieces starts at position -2, 3
	t.pendingPiece = Tetrominoes[rand.Intn(len(Tetrominoes))]
	t.pendingPiecePosition = Position{x: -2, y: 3}
	t.hasPendingPiece = true
}

func (t *Tetris) printBoard() {
	board := ""
	for i := range t.rows {
		board += t.printLine(i)
	}
	board += t.printLastLine()
	fmt.Println(board)
}

func (t *Tetris) printLine(i int) string {
	line := ""
	line += t.printOffset()
	line += "<!"
	t.mux.RLock()
	for j := range t.columns {
		cell := " ."
		if t.board[i][j] == 1 {
			cell = "[]"
		}
		line += fmt.Sprintf("%s", cell)
	}
	t.mux.RUnlock()
	line += "!>\n"
	return line
}

func (t *Tetris) printLastLine() string {
	line := ""
	line += t.printOffset()
	line += "  "
	for range t.columns {
		line += "\\/"
	}
	line += "  \n"
	return line
}

func (t *Tetris) printOffset() string {
	offset := ""
	for range t.boardLeftOffset {
		offset += fmt.Sprint("  ")
	}
	return offset
}

func (t *Tetris) isDrawablePosition(p Position) bool {
	return p.x < t.rows && 0 <= p.y && p.y < t.columns
}

func (t *Tetris) isValidPosition(p Position) bool {
	return 0 <= p.x && p.x < t.rows && 0 <= p.y && p.y < t.columns
}

func (t *Tetris) canPlacePiece(p Piece, pos Position) bool {
	for i := range p {
		newPosition := Position{p[i].x + pos.x, p[i].y + pos.y}
		if !t.isDrawablePosition(newPosition) || (t.isValidPosition(newPosition) && t.board[newPosition.x][newPosition.y] == 1) {
			return false
		}
	}
	return true
}

func (t *Tetris) MovePendingPiece(key int) bool {
	newPiece := t.pendingPiece
	newPosition := t.pendingPiecePosition
	switch key {
	case KeyUp:
		newPiece = newPiece.rotation()
	case KeyDown:
		newPosition.x++
	case KeyLeft:
		newPosition.y--
	case KeyRight:
		newPosition.y++
	default:
		return false
	}
	return t.updatePendingPiece(newPiece, newPosition)
}

func (t *Tetris) updatePendingPiece(newPiece Piece, newPosition Position) bool {
	ok := false
	t.drawPendingPiece(0)
	if t.canPlacePiece(newPiece, newPosition) {
		ok = true
		t.pendingPiece = newPiece
		t.pendingPiecePosition = newPosition
		t.hasPendingPiece = true
	}
	t.drawPendingPiece(1)
	if ok {
		t.printBoard()
	}
	return ok
}

func (t *Tetris) drawPendingPiece(value int) {
	p := t.pendingPiece
	pos := t.pendingPiecePosition
	for i := range p {
		if t.isValidPosition(Position{x: pos.x + p[i].x, y: pos.y + p[i].y}) {
			t.mux.Lock()
			t.board[pos.x+p[i].x][pos.y+p[i].y] = value
			t.mux.Unlock()
		}
	}
}

func (t *Tetris) processCompletedLines() {
	// set completed lines to zero
	for i := t.rows - 1; i >= 0; i-- {
		filled := 0
		for j := range t.columns {
			filled += t.board[i][j]
		}
		if filled == t.columns {
			t.mux.Lock()
			for j := range t.columns {
				t.board[i][j] = 0
			}
			t.mux.Unlock()
		}
	}

	// drop lines
	for i := t.rows - 2; i >= 0; i-- {
		for ni := i + 1; ni < t.rows; ni++ {
			isBelowEmpty := true
			for j := 0; j < t.columns && isBelowEmpty; j++ {
				t.mux.RLock()
				if t.board[ni][j] == 1 {
					isBelowEmpty = false
				}
				t.mux.RUnlock()
			}
			if !isBelowEmpty {
				break
			}
			t.mux.Lock()
			for j := range t.columns {
				t.board[ni][j] = t.board[ni-1][j]
				t.board[ni-1][j] = 0
			}
			t.mux.Unlock()
		}
	}
}
