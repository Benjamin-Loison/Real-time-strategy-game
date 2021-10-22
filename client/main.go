package main

import (
    "os"
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/png"
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
    sprite = loadImageFromFile("media/sprites/Dino_blue.png")
	treeSprite = loadImageFromFile("media/sprites/baum.png")
	dino = Unit{0, 0, 0, sprite.SubImage(image.Rect(0, 0, 24, 24)).(*ebiten.Image)}
	tree = Unit{200, 200, 0, treeSprite}
	camera = Unit{0, 0, 0, sprite}	// TODO should be change to another more adapted type
	zoomFactor = 1.0
)


type Unit struct {
    x,y float64 //position
    r   int //size of collision circle
    sprite *ebiten.Image
}

type Game struct {
	keys []ebiten.Key
}

func (g *Game) Update() error {
	//////////// Handling Keyboard events ////////////
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	for _, p := range g.keys {
		switch s := p.String(); s {
		case "S":
			dino.y += 5
		case "Z":
			dino.y -= 5
		case "Q":
			dino.x -= 5
		case "D":
			dino.x += 5

		case "ArrowUp":
			camera.y -= 5
		case "ArrowDown":
			camera.y += 5
		case "ArrowLeft":
			camera.x -= 5
		case "ArrowRight":
			camera.x += 5

		case "I":
			zoomFactor += 0.5
		case "K":
			zoomFactor -= 0.5
		}
	}
	return nil
}

func loadImageFromFile(path string) *ebiten.Image {
    f, err := os.Open(path)
    if err != nil {
		log.Fatal(err)
    }
    defer f.Close()

    img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
    return ebiten.NewImageFromImage(img)
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

func (u Unit) drawUnit(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.ColorM.Reset()
	op.GeoM.Translate(u.x, u.y)
	op.GeoM.Translate(camera.x, camera.y)
	op.GeoM.Scale(zoomFactor, zoomFactor)
	screen.DrawImage(u.sprite, op)
}
func getSelctionRect() (int, int, int, int) {
	return startSelection[0], startSelection[1], endSelection[0], endSelection[1]
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
    ////////////
	/*
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Reset()
    op.ColorM.Reset()
    op.GeoM.Scale(6.0,6.0)
	op.GeoM.Translate(dino.x, dino.y)
	op.GeoM.Translate(camera.x, camera.y)
    screen.DrawImage(dino.sprite.SubImage(image.Rect(0, 0, 24, 24)).(*ebiten.Image), op )
	*/
	dino.drawUnit(screen)
	tree.drawUnit(screen)
    ////////////
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
