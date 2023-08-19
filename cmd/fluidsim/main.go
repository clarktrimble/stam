package main

import (
	"github.com/clarktrimble/stam/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	gridSize = 80
	scale    = 8
)

func main() {

	screenSize := gridSize * scale

	ebiten.SetWindowSize(screenSize, screenSize)
	ebiten.SetWindowTitle("Fluid Fumble")

	game, err := game.New(gridSize, scale)
	if err != nil {
		panic(err)
	}

	err = ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}
