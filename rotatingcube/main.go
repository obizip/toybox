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
	origin = &Vector3{
		screenWidth / 2,
		screenHeight / 2,
		0,
	}
)

type Cube struct {
	x, y                           float64
	a1, b1, c1, d1, a2, b2, c2, d2 Vector3
}

func NewCube(x, y float64, width, height, length float64) Cube {
	return Cube{
		x, y,
		Vector3{-width / 2, height / 2, -length / 2},
		Vector3{-width / 2, height / 2, length / 2},
		Vector3{width / 2, height / 2, length / 2},
		Vector3{width / 2, height / 2, -length / 2},

		Vector3{-width / 2, -height / 2, -length / 2},
		Vector3{-width / 2, -height / 2, length / 2},
		Vector3{width / 2, -height / 2, length / 2},
		Vector3{width / 2, -height / 2, -length / 2},
	}
}

func (c Cube) Rotate(axis Vector3, angle float64) Cube {
	return Cube{
		c.x, c.y,
		c.a1.Rotate(axis, angle),
		c.b1.Rotate(axis, angle),
		c.c1.Rotate(axis, angle),
		c.d1.Rotate(axis, angle),
		c.a2.Rotate(axis, angle),
		c.b2.Rotate(axis, angle),
		c.c2.Rotate(axis, angle),
		c.d2.Rotate(axis, angle),
	}
}

func DrawLine(screen *ebiten.Image, x, y float64, v, w Vector3) {
	vector.StrokeLine(
		screen,
		float32(v.x+x),
		float32(v.y+y),
		float32(w.x+x),
		float32(w.y+y),
		2,
		color.White,
		true,
	)
}

func (c *Cube) Draw(screen *ebiten.Image) {
	DrawLine(screen, c.x, c.y, c.a1, c.b1)
	DrawLine(screen, c.x, c.y, c.a1, c.d1)
	DrawLine(screen, c.x, c.y, c.a1, c.a2)

	DrawLine(screen, c.x, c.y, c.c1, c.b1)
	DrawLine(screen, c.x, c.y, c.c1, c.d1)
	DrawLine(screen, c.x, c.y, c.c1, c.c2)

	DrawLine(screen, c.x, c.y, c.b2, c.a2)
	DrawLine(screen, c.x, c.y, c.b2, c.c2)
	DrawLine(screen, c.x, c.y, c.b2, c.b1)

	DrawLine(screen, c.x, c.y, c.d2, c.a2)
	DrawLine(screen, c.x, c.y, c.d2, c.c2)
	DrawLine(screen, c.x, c.y, c.d2, c.d1)
}

type Game struct {
	cube Cube
}

func NewGame() *Game {
	return &Game{cube: NewCube(origin.x, origin.y, 80, 80, 80)}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.cube.Draw(screen)
}

func (g *Game) Update() error {
	t := float64(ebiten.Tick()%3600) / 3600 * 2.0 * math.Pi
	axis := Vector3{
		x: math.Sin(t * 3.0),
		y: math.Cos(t * 2.0),
		z: math.Sin(t * 1.0),
	}
	angle := 2.
	g.cube = g.cube.Rotate(axis, angle)

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Rotating Cube")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
