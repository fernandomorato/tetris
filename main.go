package main

import (
	"os"

	"github.com/eiannone/keyboard"
	"github.com/fernandomorato/tetris/tetris"
)

func main() {
	t := tetris.Tetris{}
	t.Init()
	for {
		ch, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyCtrlC {
			os.Exit(0)
		} else if ch == 104 { // h - left
			t.MovePendingPieceLeft()
		} else if ch == 108 { // l - right
			t.MovePendingPieceRight()
		} else if ch == 106 { // j - down
			t.MovePendingPieceDown()
		} else if ch == 107 { // k
			t.RotatePiece()
		}
	}
}
