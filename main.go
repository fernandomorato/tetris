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
		}
		t.MovePendingPiece(int(ch))
	}
}
