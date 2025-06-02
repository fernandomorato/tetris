package tetris

type Piece [4]Position

var Tetrominoes = [7]Piece{
	{ // [][][][]
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
	},
	{ //     []
		// [][][]
		{0, 2},
		{1, 0},
		{1, 1},
		{1, 2},
	},
	{ //   []
		// [][][]
		{0, 1},
		{1, 0},
		{1, 1},
		{1, 2},
	},
	{ // []
		// [][][]
		{0, 0},
		{1, 0},
		{1, 1},
		{1, 2},
	},
	{ // [][]
		//   [][]
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 2},
	},
	{ // [][]
		// [][]
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	},
	{ //   [][]
		// [][]
		{0, 1},
		{0, 2},
		{1, 0},
		{1, 1},
	},
}

func (p *Piece) rotation() Piece {
	cx := 1.5
	cy := 1.5
	if *p == Tetrominoes[5] {
		return *p
	} else if *p == Tetrominoes[0] {
		cx = 1.0
		cy = 1.0
	}
	// Rotation: (x, y) -> (y, -x), centered on (cx, cy)
	newPiece := Piece{}
	for i := range p {
		newPiece[i].x = int(float64(p[i].y) - cy + cx)
		newPiece[i].y = int(-(float64(p[i].x) - cx) + cy)
	}
	return newPiece
}
