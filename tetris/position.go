package tetris

type Position struct {
	x int
	y int
}

func (p *Position) rotateClockwise(c Position) {
	// 90 degree clockwise rotation of p(x, y) around c(x, y)
	p.x -= c.x
	p.y -= c.y
	p.x, p.y = p.y, -p.x
	p.x += c.x
	p.y += c.y
}
