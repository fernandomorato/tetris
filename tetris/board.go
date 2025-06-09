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

var Tetrominoes = [7]Piece{I, L, T, J, Z, O, S}

type Tetris struct {
	rows            int
	columns         int
	boardLeftOffset int
	board           [][]int
	currentPiece    Piece
	nextPiece       Piece
	hasPendingPiece bool
	ticker          *time.Ticker
	mux             *sync.Mutex
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
	t.getNextPiece()
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
			t.mux.Lock()
			t.MovePiece(int(ch))
			t.mux.Unlock()
		}
	}
}

func (t *Tetris) getNextPiece() {
	var newPiece Piece
	index := rand.Intn(len(Tetrominoes))
	newPiece = append(newPiece, Tetrominoes[index]...)
	t.nextPiece = newPiece
}

func (t *Tetris) refresh() {
	for ; ; <-t.ticker.C {
		// if has no pending piece:
		// 1 - check if there are completed lines and update the board accordingly
		// 2 - spawns new piece
		t.mux.Lock()
		if !t.hasPendingPiece {
			t.processCompletedLines()
			ok := t.spawnNewPiece()
			if !ok {
				fmt.Println("game over")
				os.Exit(0)
			}
			t.hasPendingPiece = ok
		} else {
			t.hasPendingPiece = t.MovePiece(KeyJ)
		}
		t.mux.Unlock()
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
	for j := range t.columns {
		cell := " ."
		if t.board[i][j] == 1 {
			cell = "[]"
		}
		line += fmt.Sprintf("%s", cell)
	}
	line += "!>"
	if i >= 7 && i <= 11 {
		line += t.printNextPiece(i)
	}
	line += "\n"
	return line
}

func (t *Tetris) printNextPiece(i int) string {
	line := "  "
	if i == 7 {
		line += "   NEXT"
	} else {
		for j := range 5 {
			foundPosition := false
			for _, pos := range t.nextPiece {
				if i-9 == pos.x && j == pos.y - 3 {
					foundPosition = true
				}
			}
			if foundPosition {
				line += "[]"
			} else {
				line += " ."
			}
		}
	}
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

func (t *Tetris) canDrawPiece(p Piece) bool {
	for i := range p {
		cond1 := !t.isValidPosition(p[i]) && (t.columns <= p[i].y || p[i].y < 0 || p[i].x >= t.rows)
		cond2 := t.isValidPosition(p[i]) && t.board[p[i].x][p[i].y] == 1
		if cond1 || cond2 {
			return false
		}
	}
	return true
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
	if !t.hasPendingPiece {
		return false
	}
	var newPiece Piece
	newPiece = append(newPiece, t.currentPiece...)
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

func (t *Tetris) spawnNewPiece() bool {
	// Spawns a random piece. The spawn position of each Tetromino is fixed
	var newPiece Piece
	newPiece = append(newPiece, t.nextPiece...)
	ok := t.canPlacePiece(newPiece)
	t.currentPiece = newPiece
	t.getNextPiece()
	t.drawPiece(1)
	t.printBoard()
	return ok
}

func (t *Tetris) updatePiece(newPiece Piece) bool {
	ok := false
	t.drawPiece(0)
	if ok = t.canDrawPiece(newPiece); ok {
		t.currentPiece = newPiece
	}
	t.drawPiece(1)
	if ok {
		t.printBoard()
	}
	return ok
}

func (t *Tetris) drawPiece(value int) {
	for _, pos := range t.currentPiece {
		if t.isValidPosition(pos) {
			t.board[pos.x][pos.y] = value
		}
	}
}

func (t *Tetris) processCompletedLines() {
	// set completed lines to zero
	for i := range t.rows {
		filled := 0
		for j := range t.columns {
			filled += t.board[i][j]
		}
		if filled == t.columns {
			for j := range t.columns {
				t.board[i][j] = 0
			}
			for ni := i - 1; ni >= 0; ni-- {
				for j := range t.columns {
					t.board[ni][j], t.board[ni+1][j] = t.board[ni+1][j], t.board[ni][j]
				}
			}
		}
	}
}
