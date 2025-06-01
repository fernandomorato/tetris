package tetris

import (
	"fmt"
	"sync"
	"time"
)

type Tetris struct {
	rows            int
	columns         int
	boardLeftOffset int
	board           [][]int
	pendingPiece    Piece
	hasPendingPiece bool
	ticker          *time.Ticker
	mux             *sync.RWMutex
}

func (t *Tetris) Init() {
	t.rows = 20
	t.columns = 10
	t.boardLeftOffset = 20
	t.board = make([][]int, t.rows)
	for i := range t.rows {
		t.board[i] = make([]int, t.columns)
	}
	t.ticker = time.NewTicker(100 * time.Millisecond)
	t.mux = &sync.RWMutex{}
	go t.refresh()
}

func (t *Tetris) refresh() {
	for ; ; <-t.ticker.C {
		// if has no pending piece, create one at -1, 4
		if !t.hasPendingPiece {
			t.pendingPiece = Piece{
				x: -1,
				y: 4,
			}
			t.hasPendingPiece = true
		}
		// update board status
		t.MovePendingPieceDown()
		// print board
		t.printBoard()
	}
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

func (t *Tetris) isValidPosition(x int, y int) bool {
	return 0 <= x && x < t.rows && 0 <= y && y < t.columns
}

func (t *Tetris) canPlacePiece(p Piece) bool {
	return t.isValidPosition(p.x, p.y) && t.board[p.x][p.y] != 1
}

func (t *Tetris) MovePendingPieceRight() {
	newPiece := t.pendingPiece
	newPiece.y++
	if t.canPlacePiece(newPiece) {
		t.updatePendingPiece(newPiece)
	}
}

func (t *Tetris) MovePendingPieceLeft() {
	newPiece := t.pendingPiece
	newPiece.y--
	if t.canPlacePiece(newPiece) {
		t.updatePendingPiece(newPiece)
	}
}

func (t *Tetris) MovePendingPieceDown() {
	newPiece := t.pendingPiece
	newPiece.x++
	if t.canPlacePiece(newPiece) {
		t.updatePendingPiece(newPiece)
	} else {
		t.hasPendingPiece = false
	}
}

func (t *Tetris) updatePendingPiece(newPiece Piece) {
	t.hasPendingPiece = true
	t.setBoardPosition(t.pendingPiece, 0)
	t.setBoardPosition(newPiece, 1)
	t.pendingPiece = newPiece
}

func (t *Tetris) setBoardPosition(p Piece, value int) {
	if t.isValidPosition(p.x, p.y) {
		t.mux.Lock()
		t.board[p.x][p.y] = value
		t.mux.Unlock()
	}
}
