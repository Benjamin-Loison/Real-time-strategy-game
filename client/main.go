package main

import (
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	//"image"
	"image/color"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

var (
	Red       = color.RGBA{255,0,0,255}
	inSelection = false
	startSelection = [2]int{0, 0}
	endSelection = [2]int{0, 0}
)


type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func drawSelectionRect(screen *ebiten.Image) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !inSelection {
			inSelection = true
			startSelection[0],startSelection[1]  = ebiten.CursorPosition()
			endSelection[0],endSelection[1]  = ebiten.CursorPosition()
			endSelection[0] -= startSelection[0]
			endSelection[1] -= startSelection[1]
		} else {
			endSelection[0],endSelection[1]  = ebiten.CursorPosition()
			endSelection[0] -= startSelection[0]
			endSelection[1] -= startSelection[1]
		}
	} else {
	  inSelection = false
	}

	if inSelection {
		ebitenutil.DrawRect(screen, float64(startSelection[0]), float64(startSelection[1]), float64(endSelection[0]), float64(endSelection[1]), Red)
	}
}

func getSelctionRect() (int, int, int, int) {
	return startSelection[0], startSelection[1], endSelection[0], endSelection[1]
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	ebitenutil.DrawRect(screen, screenWidth/2 ,screenHeight/2, screenWidth/4, screenHeight/4, Red)
	drawSelectionRect(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main(){
	ebiten.SetWindowSize(1280,720)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("EHO: Elves, humans and orks")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
