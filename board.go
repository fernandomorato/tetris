package main

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
	mux             *sync.Mutex
}

type Piece struct {
	x int
	y int
}

func (t *Tetris) init() {
	t.rows = 20
	t.columns = 10
	t.boardLeftOffset = 20
	t.board = make([][]int, t.rows)
	for i := range t.rows {
		t.board[i] = make([]int, t.columns)
	}
	t.ticker = time.NewTicker(100 * time.Millisecond)
	t.mux = &sync.Mutex{}
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
		t.movePendingPieceDown()
		// print board
		t.printBoard()
	}
}

func (t *Tetris) printBoard() {
	for i := range t.rows {
		t.printLine(i)
	}
	t.printLastLine()
}

func (t *Tetris) printLine(i int) {
	t.printOffset()
	fmt.Print("<!")
	for j := range t.columns {
		cell := " ."
		if t.board[i][j] == 1 {
			cell = "[]"
		}
		fmt.Printf("%s", cell)
	}
	fmt.Println("!>")
}

func (t *Tetris) printLastLine() {
	t.printOffset()
	fmt.Print("  ")
	for range t.columns {
		fmt.Print("\\/")
	}
	fmt.Println("  ")
}

func (t *Tetris) printOffset() {
	for range t.boardLeftOffset {
		fmt.Print("  ")
	}
}

func (t *Tetris) isValidPosition(x int, y int) bool {
	return 0 <= x && x < t.rows && 0 <= y && y < t.columns
}

func (t *Tetris) canDropPendingPiece() bool {
	if t.isValidPosition(t.pendingPiece.x+1, t.pendingPiece.y) {
		return t.board[t.pendingPiece.x+1][t.pendingPiece.y] != 1
	}
	return false
}

func (t *Tetris) canPlacePiece(p Piece) bool {
	return t.isValidPosition(p.x, p.y) && t.board[p.x][p.y] != 1
}

func (t *Tetris) movePendingPieceRight() {
	newPiece := t.pendingPiece
	newPiece.y++
	if t.canPlacePiece(newPiece) {
		t.hasPendingPiece = true
		t.setBoardPosition(t.pendingPiece, 0)
		t.setBoardPosition(newPiece, 1)
		t.pendingPiece = newPiece
			t.hasPendingPiece = true
	}
}

func (t *Tetris) movePendingPieceLeft() {
	newPiece := t.pendingPiece
	newPiece.y--
	if t.canPlacePiece(newPiece) {
		t.hasPendingPiece = true
		t.setBoardPosition(t.pendingPiece, 0)
		t.setBoardPosition(newPiece, 1)
		t.pendingPiece = newPiece
	}
}

func (t *Tetris) movePendingPieceDown() {
	newPiece := t.pendingPiece
	newPiece.x++
	if t.canPlacePiece(newPiece) {
		t.hasPendingPiece = true
		t.setBoardPosition(t.pendingPiece, 0)
		t.setBoardPosition(newPiece, 1)
		t.pendingPiece = newPiece
	} else {
		t.hasPendingPiece = false
	}
}

func (t *Tetris) setBoardPosition(p Piece, value int) {
	if t.isValidPosition(p.x, p.y) {
		t.mux.Lock()
		t.board[p.x][p.y] = value
		t.mux.Unlock()
	}
}
