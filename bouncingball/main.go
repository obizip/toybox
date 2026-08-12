package main

import (
	"log"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Ball struct {
	cx float32
	cy float32
	r float32
	dx float32
	dy float32
}

func (b *Ball) draw(screen *ebiten.Image) {
	vector.FillCircle(screen, b.cx, b.cy, b.r, color.White, true)
}

func (b *Ball) update() {
	b.cx = b.cx + b.dx
	b.cy = b.cy + b.dy
}

type Game struct{
	ball *Ball
}

func NewGame() *Game {
	return &Game{
		ball: &Ball{60.0, 60.0, 30.0, 5.0, 5.0},
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ball.draw(screen)
}

func (g *Game) Update() error {
	g.ball.update()

	ball := g.ball
	if screenWidth <= ball.cx + ball.r || ball.cx - ball.r <= 0 {
		ball.dx = -ball.dx
	}
	if screenHeight <= ball.cy + ball.r || ball.cy - ball.r <= 0 {
		ball.dy = -ball.dy
	}
	
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
