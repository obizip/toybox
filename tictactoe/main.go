package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState = int

const (
	screenWidth  = 640
	screenHeight = 480
	infoHeight   = 30

	GamePlay GameState = iota
	GameFinished
)

type Game struct {
	isX   bool
	board *Board
	state GameState
	info  string
}

func NewGame() *Game {
	return &Game{
		isX:   true,
		board: NewBoard(),
		state: GamePlay,
		info: "",
	}
}

func (g *Game) Init() {
	game := NewGame()
	g.isX = game.isX
	g.board = game.board
	g.state = game.state
	g.info = game.info
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.Draw(screen, g.info)
}

func (g *Game) Update() error {
	switch g.state {
	case GamePlay:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			cursorX, cursorY := ebiten.CursorPosition()
			g.board.HandleClick(cursorX, cursorY, g.isX)
			if g.board.IsLinedUp() {
				g.state = GameFinished
				g.info = "Winner: " + g.getPlayerName()
				return nil
			}
			if g.board.IsFull() {
				g.state = GameFinished
				g.info = "Draw"
				return nil
			}

			g.isX = !g.isX
		}
		g.info = "Player: " + g.getPlayerName()

	case GameFinished:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			g.info = "New"
			g.Init()
		}
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) getPlayerName() string {
	if g.isX {
		return "X"
	} else {
		return "O"
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tic-Tac-Toe")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
