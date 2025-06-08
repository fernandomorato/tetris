package tetris

import (
	"fmt"
	"math"
)

type Piece [4]Position

var (
	//	. . . .
	//	. . . .
	// [][][][]
	//	. . . .
	I = Piece{{0, 0}, {0, 1}, {0, 2}, {0, 3}}
	// [][][]
	// [] . .
	//  . . .
	L = Piece{{0, 0}, {0, 1}, {0, 2}, {1, 0}}
	// [][][]
	//  .[] .
	//  . . .
	T = Piece{{0, 0}, {0, 1}, {0, 2}, {1, 1}}
	// [][][]
	//  . .[]
	//  . . .
	J = Piece{{0, 0}, {0, 1}, {0, 2}, {1, 2}}
	// [][] .
	//  .[][]
	//  . . .
	Z = Piece{{0, 0}, {0, 1}, {1, 1}, {1, 2}}
	// [][]
	// [][]
	O = Piece{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	//  .[][]
	// [][] .
	//  . . .
	S = Piece{{0, 1}, {0, 2}, {1, 0}, {1, 1}}
)

var Tetrominoes = [7]Piece{I, L, T, J, Z, O, S}

func (p *Piece) rotate() {
	minX := math.MaxInt
	minY := math.MaxInt
	for i := range p {
		minX = min(minX, p[i].x)
		minY = min(minY, p[i].y)
	}
	for i := range p {
		p[i].x -= minX
		p[i].y -= minY
	}
	var cx, cy float64
	switch *p {
	case T, S, Z, L, J:
		cx, cy = 0.5, 1.5
	case I:
		cx, cy = -0.5, 1.5
	case O:
		cx, cy = 0.5, 0.5
	}
	fmt.Println(cx, cy)

	// Rotation: (x, y) -> (y, -x), centered on (cx, cy)

	for i := range p {
		fmt.Printf("%d with %v\n", i, p[i])
		nx := float64(p[i].x) - cx
		ny := float64(p[i].y) - cy
		fmt.Println(nx, ny)
		nx, ny = ny, -nx
		fmt.Println(nx, ny)
		nx += cx
		ny += cy
		fmt.Println(nx, ny)
		p[i] = Position{
			x: int(nx),
			y: int(ny),
		}
	}
	for i := range p {
		p[i].x += minX
		p[i].y += minY
	}
}

func (p *Piece) drop() {
	for i := range p {
		p[i].x++
	}
}

func (p *Piece) left() {
	for i := range p {
		p[i].y--
	}
}

func (p *Piece) right() {
	for i := range p {
		p[i].y++
	}
}
