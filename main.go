package main

import (
	"os"

	"github.com/eiannone/keyboard"
)

func main() {
	t := Tetris{}
	t.init()
	for {
		ch, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyCtrlC {
			os.Exit(0)
		} else if ch == 104 { // h - left
			t.movePendingPieceLeft()
		} else if ch == 108 { // l - right
			t.movePendingPieceRight()
		} else if ch == 106 { // j - down
			t.movePendingPieceDown()
		} else if ch == 107 { // k
			// move up -> rotate
		}
	}
}
