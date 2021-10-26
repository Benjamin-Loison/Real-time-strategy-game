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
    ZOOM_STEP = 0.01
)

var (
	Red       = color.RGBA{255,0,0,255}
	inSelection = false
	startSelection = [2]int{0, 0}
	endSelection = [2]int{0, 0}
    sprite = loadImageFromFile("media/sprites/Dino_blue.png")
	treeSprite = loadImageFromFile("media/sprites/baum.png")
	dino = Entity{0, 0, 12, sprite.SubImage(image.Rect(0, 0, 24, 24)).(*ebiten.Image), 6, false}
	tree = Entity{640, 360, 32, treeSprite, 1.0, false}
	camera = Entity{0, 0, 0,  nil , 1, false}	// TODO should be change to another more adapted type
	zoomFactor = 1.0
)


type Entity struct {
    x,y float64 //position
    r   int //size of collision circle
    sprite *ebiten.Image
    sprite_base_scale float64
    selected bool
}

type Game struct {
	keys []ebiten.Key
    friendlyEntities []*Entity
    envEntities []*Entity
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
			zoomFactor += ZOOM_STEP
		case "K":
			zoomFactor -= ZOOM_STEP
		}
	}
    ////////////////////////////////////////////////
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


func drawWireRect(screen *ebiten.Image, x, y, w, h float64, c color.Color) {
	ebitenutil.DrawLine(screen, x   , y  , x+w , y   , c)
	ebitenutil.DrawLine(screen, x   , y  , x   , y+h , c)
	ebitenutil.DrawLine(screen, x+w , y  , x+w , y+h , c)
	ebitenutil.DrawLine(screen, x   , y+h, x+w , y+h , c)

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
        x  := float64(startSelection[0])
        y  := float64(startSelection[1])
        dx := float64(endSelection[0])
        dy := float64(endSelection[1])

		drawWireRect(screen, x, y, dx, dy, Red)
	}
}


func (e Entity) getScreenTransform() (*ebiten.DrawImageOptions) {

	op := &ebiten.DrawImageOptions{}
    iw,ih := e.sprite.Size()
	
	op.GeoM.Reset()
	op.ColorM.Reset()
	
	op.GeoM.Translate( - float64(iw)/2 , - float64(ih)/2 )
	op.GeoM.Scale(zoomFactor*e.sprite_base_scale, zoomFactor*e.sprite_base_scale)
	op.GeoM.Translate(e.x*zoomFactor, e.y*zoomFactor)
	op.GeoM.Translate(-camera.x*zoomFactor, -camera.y*zoomFactor)
	op.GeoM.Translate(-screenWidth*zoomFactor/2, -screenHeight*zoomFactor/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	
	return op
}


func (e Entity) drawHitbox(screen *ebiten.Image) {
	op := e.getScreenTransform()
	iw, ih := e.sprite.Size()
	x1, y1 := op.GeoM.Apply(0, 0)
	x2, y2 := op.GeoM.Apply(float64(iw), float64(ih))
	drawWireRect(screen, x1, y1, x2-x1, y2-y1, Red)
}


func (e Entity) drawEntity(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.ColorM.Reset()
	op = e.getScreenTransform()
	screen.DrawImage(e.sprite, op)
}
func getSelctionRect() (int, int, int, int) {
	return startSelection[0], startSelection[1], endSelection[0], endSelection[1]
}

func (g *Game) Draw(screen *ebiten.Image) {
    // clearing the screen
	screen.Fill(color.White)

    // drawing the elements
    for _, e := range g.envEntities {
        e.drawEntity(screen)
		e.drawHitbox(screen)
    }
    for _, e := range g.friendlyEntities {
        e.drawEntity(screen)
		e.drawHitbox(screen)
    }

    ////////////
	drawSelectionRect(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Function to initialize the game
func Init(g *Game) {
    g.friendlyEntities = append(g.friendlyEntities, &dino)
    g.envEntities = append(g.envEntities, &tree)
}

func main(){
    // initializing the game
    g := &Game{}
    Init(g)
    //initializing ebiten
	ebiten.SetWindowSize(screenWidth,screenHeight)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("EHO: Elves, humans and orks")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
