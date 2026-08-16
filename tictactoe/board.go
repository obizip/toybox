package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SquareType int

const (
	SquareEmpty SquareType = iota
	SquareX
	SquareO
)

type Board struct {
	x, y, length int
	squares      [9]SquareType
}

func NewBoard() *Board {
	length := min(screenWidth, screenHeight-infoHeight) - 10
	x := screenWidth/2 - length/2
	y := (screenHeight+infoHeight)/2 - length/2
	var board [9]SquareType

	return &Board{
		x, y, length,
		board,
	}
}

func (b *Board) IsLinedUp() bool {
	ss := b.squares
	for i := range 3 {
		if ss[i*3+0] != SquareEmpty && ss[i*3+0] == ss[i*3+1] && ss[i*3+1] == ss[i*3+2] {
			return true
		}
		if ss[0+i] != SquareEmpty && ss[0+i] == ss[3+i] && ss[3+i] == ss[6+i] {
			return true
		}
	}
	if ss[0] != SquareEmpty && ss[0] == ss[4] && ss[4] == ss[8] {
		return true
	}
	if ss[2] != SquareEmpty && ss[2] == ss[4] && ss[4] == ss[6] {
		return true
	}

	return false
}

func (b *Board) IsFull() bool {
	for _, square := range b.squares {
		if square == SquareEmpty {
			return false
		}
	}
	return true
}

func (b *Board) Draw(screen *ebiten.Image, info string) {
	ebitenutil.DebugPrintAt(screen, info, int(b.x), int(b.y-infoHeight))

	// Render outer rectangle
	vector.StrokeRect(
		screen,
		float32(b.x),
		float32(b.y),
		float32(b.length),
		float32(b.length),
		2.,
		color.White,
		true,
	)

	// Render squares
	squareLength := b.length / 3
	for i := range 2 {
		currentLength := squareLength * (i + 1)

		vector.StrokeLine(
			screen,
			float32(b.x+currentLength),
			float32(b.y),
			float32(b.x+currentLength),
			float32(b.y+b.length),
			2,
			color.White,
			true,
		)

		vector.StrokeLine(
			screen,
			float32(b.x),
			float32(b.y+currentLength),
			float32(b.x+b.length),
			float32(b.y+currentLength),
			2,
			color.White,
			true,
		)
	}

	for i, square := range b.squares {
		b.renderSquare(screen, i, square)
	}
}

func (b *Board) renderSquare(screen *ebiten.Image, index int, square SquareType) {
	i, j := index/3, index%3
	squareLength := b.length / 3

	centerY := b.y + squareLength*(i+1) - squareLength/2
	centerX := b.x + squareLength*(j+1) - squareLength/2

	margin := squareLength / 10

	switch square {
	case SquareEmpty:
		return
	case SquareO:
		vector.StrokeCircle(screen,
			float32(centerX),
			float32(centerY),
			float32(squareLength/2-margin),
			2,
			color.White,
			true,
		)
	case SquareX:
		vector.StrokeLine(
			screen,
			float32(centerX+margin)-float32(squareLength)/2,
			float32(centerY+margin)-float32(squareLength)/2,
			float32(centerX-margin)+float32(squareLength)/2,
			float32(centerY-margin)+float32(squareLength)/2,
			2,
			color.White,
			true,
		)
		vector.StrokeLine(
			screen,
			float32(centerX-margin)+float32(squareLength)/2,
			float32(centerY+margin)-float32(squareLength)/2,
			float32(centerX+margin)-float32(squareLength)/2,
			float32(centerY-margin)+float32(squareLength)/2,
			2,
			color.White,
			true,
		)
		return
	}
}

func (b *Board) HandleClick(x, y int, isX bool) {
	if x <= b.x || x >= b.x + b.length {
		return
	}
	if y <= b.y || y >= b.y + b.length {
		return
	}

	squareLength := b.length / 3
	for i := range 3 {
		for j := range 3 {
			squareY := b.y + squareLength*i
			squareX := b.x + squareLength*j

			if !(squareX <= x && x <= squareX+squareLength &&
				squareY <= y && y <= squareY+squareLength) {
				continue
			}

			if b.squares[i*3+j] != SquareEmpty {
				continue
			}

			if isX {
				b.squares[i*3+j] = SquareX
			} else {
				b.squares[i*3+j] = SquareO
			}
			return
		}
	}
}
