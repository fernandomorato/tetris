package tetris

import (
	"fmt"
	"math/rand"
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
			t.spawnNewPiece()
			t.hasPendingPiece = true
		}
		// update board status
		ok := t.MovePendingPieceDown()
		if !ok {
			t.hasPendingPiece = false
		}
		// print board
		t.printBoard()
	}
}

func (t *Tetris) spawnNewPiece() {
	// starts at -2, 3
	t.pendingPiece = Tetrominoes[rand.Intn(len(Tetrominoes))]
	for i := range t.pendingPiece {
		t.pendingPiece[i][0] -= 2
		t.pendingPiece[i][1] += 3
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

func (t *Tetris) isDrawablePosition(x int, y int) bool {
	return x < t.rows && 0 <= y && y < t.columns
}

func (t *Tetris) isValidPosition(x int, y int) bool {
	return 0 <= x && x < t.rows && 0 <= y && y < t.columns
}

func (t *Tetris) canPlacePiece(p Piece) bool {
	for i := range p {
		if !t.isDrawablePosition(p[i][0], p[i][1]) || (t.isValidPosition(p[i][0], p[i][1]) && t.board[p[i][0]][p[i][1]] == 1) {
			return false
		}
	}
	return true
}

func (t *Tetris) RotatePiece() bool {
	newPiece := t.pendingPiece
	// for i := range newPiece {
	// 	newPiece[i][1]--
	// }
	return t.updatePendingPiece(newPiece)
}

func (t *Tetris) MovePendingPieceRight() bool {
	newPiece := t.pendingPiece
	for i := range newPiece {
		newPiece[i][1]++
	}
	return t.updatePendingPiece(newPiece)
}

func (t *Tetris) MovePendingPieceLeft() bool {
	newPiece := t.pendingPiece
	for i := range newPiece {
		newPiece[i][1]--
	}
	return t.updatePendingPiece(newPiece)
}

func (t *Tetris) MovePendingPieceDown() bool {
	newPiece := t.pendingPiece
	for i := range newPiece {
		newPiece[i][0]++
	}
	return t.updatePendingPiece(newPiece)
}

func (t *Tetris) updatePendingPiece(newPiece Piece) bool{
	ok := false
	t.drawPiece(t.pendingPiece, 0)
	if t.canPlacePiece(newPiece) {
		ok = true
		t.pendingPiece = newPiece
		t.hasPendingPiece = true
	}
	t.drawPiece(t.pendingPiece, 1)
	return ok
}

func (t *Tetris) drawPiece(p Piece, value int) {
	for i := range p {
		if t.isValidPosition(p[i][0], p[i][1]) {
			t.mux.Lock()
			t.board[p[i][0]][p[i][1]] = value
			t.mux.Unlock()
		}
	}
}
