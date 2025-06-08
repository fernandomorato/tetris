package tetris

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	KeyW          = 119
	KeyA          = 97
	KeyS          = 115
	KeyD          = 100
	KeyK          = 107
	KeyJ          = 106
	KeyH          = 104
	KeyL          = 108
	RefreshRateMs = 250
)

type Tetris struct {
	rows            int
	columns         int
	boardLeftOffset int
	board           [][]int
	piece           Piece
	hasPendingPiece bool
	ticker          *time.Ticker
	mux             *sync.Mutex
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
	t.ticker = time.NewTicker(RefreshRateMs * time.Millisecond)
	t.mux = &sync.Mutex{}
	go t.refresh()
	for {
		ch, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		switch key {
		case keyboard.KeyCtrlC:
			os.Exit(0)
		default:
			t.MovePiece(int(ch))
		}
	}
}

func (t *Tetris) refresh() {
	for ; ; <-t.ticker.C {
		// if has no pending piece:
		// 1 - check if there are completed lines and update the board accordingly
		// 2 - spawns new piece
		if !t.hasPendingPiece {
			t.processCompletedLines()
			t.spawnNewPiece()
		} else {
			t.hasPendingPiece = t.MovePiece(KeyJ)
		}
		t.printBoard()
	}
}

func (t *Tetris) spawnNewPiece() {
	newPiece := Tetrominoes[rand.Intn(len(Tetrominoes))]
	offset := 4
	if newPiece == I {
		offset--
	}
	for i := range newPiece {
		newPiece[i].y += offset
	}
	t.piece = newPiece
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
	for j := range t.columns {
		cell := " ."
		if t.board[i][j] == 1 {
			cell = "[]"
		}
		line += fmt.Sprintf("%s", cell)
	}
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

func (t *Tetris) isValidPosition(p Position) bool {
	return 0 <= p.x && p.x < t.rows && 0 <= p.y && p.y < t.columns
}

func (t *Tetris) canPlacePiece(p Piece) bool {
	for i := range p {
		if !t.isValidPosition(p[i]) || t.board[p[i].x][p[i].y] == 1 {
			return false
		}
	}
	return true
}

func (t *Tetris) MovePiece(key int) bool {
	newPiece := t.piece
	t.mux.Lock()
	defer t.mux.Unlock()
	switch key {
	case KeyK, KeyW:
		newPiece.rotate()
	case KeyJ, KeyS:
		newPiece.drop()
	case KeyH, KeyA:
		newPiece.left()
	case KeyL, KeyD:
		newPiece.right()
	}
	return t.updatePiece(newPiece)
}

func (t *Tetris) updatePiece(newPiece Piece) bool {
	ok := false
	t.drawPiece(0)
	if ok = t.canPlacePiece(newPiece); ok {
		t.piece = newPiece
	}
	t.drawPiece(1)
	if ok {
		t.printBoard()
	}
	return ok
}

func (t *Tetris) drawPiece(value int) {
	p := t.piece
	for i := range p {
		if t.isValidPosition(p[i]) {
			t.board[p[i].x][p[i].y] = value
		}
	}
}

func (t *Tetris) processCompletedLines() {
	// set completed lines to zero
	t.mux.Lock()
	defer t.mux.Unlock()
	for i := t.rows - 1; i >= 0; i-- {
		filled := 0
		for j := range t.columns {
			filled += t.board[i][j]
		}
		if filled == t.columns {
			for j := range t.columns {
				t.board[i][j] = 0
			}
		}
	}

	// drop lines
	for i := t.rows - 2; i >= 0; i-- {
		for ni := i + 1; ni < t.rows; ni++ {
			isBelowEmpty := true
			for j := 0; j < t.columns && isBelowEmpty; j++ {
				if t.board[ni][j] == 1 {
					isBelowEmpty = false
				}
			}
			if !isBelowEmpty {
				break
			}
			for j := range t.columns {
				t.board[ni][j] = t.board[ni-1][j]
				t.board[ni-1][j] = 0
			}
		}
	}
}
