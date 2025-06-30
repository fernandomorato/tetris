package main

import (
	"log"

	"github.com/fernandomorato/tetris/tetris"
	"github.com/gdamore/tcell/v2"
)

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	t := tetris.Tetris{}
	t.Init(s)
}
