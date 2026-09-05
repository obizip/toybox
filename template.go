package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct{}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bouncing Ball")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
