package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	// 中心を原点とする
	origin = &Vector{float32(screenWidth / 2), float32(screenHeight / 2)}
)

type Vector struct {
	x, y float32
}

func (v *Vector) Draw(screen *ebiten.Image) {
	maxRadius := min(screenWidth, screenHeight) / 2
	mulX := float32(maxRadius / 10)
	mulY := float32(maxRadius / 10)

	vector.StrokeLine(
		screen,
		origin.x,
		origin.y,
		origin.x+v.x*mulX,
		origin.y-v.y*mulY,
		2,
		color.White,
		true,
	)
}

func (v *Vector) Rotate(degree int) {
	radian := float64(degree) * math.Pi / 180
	newX := math.Cos(radian)*float64(v.x) - math.Sin(radian)*float64(v.y)
	newY := math.Sin(radian)*float64(v.x) + math.Cos(radian)*float64(v.y)

	v.x = float32(newX)
	v.y = float32(newY)
}

type Game struct{
	v *Vector
}

func NewGame() *Game {
	return &Game{
		v: &Vector{0, 8},
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.v.Draw(screen)
}

func (g *Game) Update() error {
	g.v.Rotate(1)

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Rotating Line")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
