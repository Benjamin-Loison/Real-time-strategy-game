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
)


type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  screen.Fill(color.White)
  ebitenutil.DrawRect(screen, 2*screenWidth/3 ,2*screenHeight/3, screenWidth, screenHeight, Red)
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
