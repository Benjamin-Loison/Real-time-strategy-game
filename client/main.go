package main

import (
	"os"
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"math"
	_ "image/png"
	"image/color"
	"image/draw"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	ZOOM_STEP = 0.01
)

var (
	Red	   = color.RGBA{255,0,0,255}
	inSelection = false
	startSelection = [2]int{0, 0}
	endSelection = [2]int{0, 0}
	sprite = loadImageFromFile("media/sprites/Dino_blue.png")
	treeSprite = loadImageFromFile("media/sprites/baum.png")
	dino = Entity{0, 0, 12, sprite.SubImage(image.Rect(0, 0, 24, 24)).(*ebiten.Image), 6, false}
	tree = Entity{640, 360, 32, treeSprite, 1.0, false}
	camera = Entity{0, 0, 0,  nil , 1, false}	// TODO should be change to another more adapted type
	onScreenMap = Map{0, 0, 0, 0, make([]Entity, 0), make([]Entity, 0), make([]Entity, 0)}
	zoomFactor = 1.0

	client_chan chan string
)


type location_type int64
const (
	Floor location_type = 1
	ElfBuilding = 2
	OrcBuilding = 3
	HumanBuilding = 4
	)

type Location struct {
	loc_type location_type
	int_arg int
	str_arg string
	// TODO options map[string]string
}

type Entity struct {
	x,y float64 //position
	r   int //size of collision circle
	sprite *ebiten.Image
	sprite_base_scale float64
	selected bool
}

type Map struct {
	origin_x int
	origin_y int
	w, h int
	buildings []Entity
	floor []Entity
	entities []Entity
}


type Game struct {
	keys []ebiten.Key
	friendlyEntities []*Entity
	envEntities []*Entity

	envLayer *ebiten.Image
	entityLayer *ebiten.Image
	debugLayer *ebiten.Image
	guiLayer *ebiten.Image
}

func update_map(init_x , init_y, w, h int, location_list []Location) {
	
	onScreenMap.origin_x = init_x
	onScreenMap.origin_y = init_y
	onScreenMap.w = w
	onScreenMap.h = h
	onScreenMap.buildings = make([]Entity, 0)
	onScreenMap.floor = make([]Entity, 0)

	for _, l := range location_list {
		switch(l.loc_type) {
			case ElfBuilding:
				onScreenMap.buildings = append(onScreenMap.buildings, tree)
			case OrcBuilding:
				onScreenMap.buildings = append(onScreenMap.buildings, tree)
			case HumanBuilding:
				onScreenMap.buildings = append(onScreenMap.buildings, tree)
			case Floor:
				continue
			default:
				continue
		}
	}

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

func (e Entity) getScreenTranslation() (*ebiten.DrawImageOptions) {

	op := &ebiten.DrawImageOptions{}
	iw,ih := e.sprite.Size()
	
	op.GeoM.Reset()
	op.ColorM.Reset()
	
	op.GeoM.Translate( - float64(iw)/2 , - float64(ih)/2 )
	op.GeoM.Translate(e.x*zoomFactor, e.y*zoomFactor)
	op.GeoM.Translate(-camera.x*zoomFactor, -camera.y*zoomFactor)
	op.GeoM.Translate(-screenWidth*zoomFactor/2, -screenHeight*zoomFactor/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	
	return op
}
func (e Entity) courgette() (*ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	//iw,ih := e.sprite.Size()
	
	op.GeoM.Reset()
	op.ColorM.Reset()

	op.GeoM.Scale(zoomFactor*e.sprite_base_scale, zoomFactor*e.sprite_base_scale)

	return op
}


func (e Entity) drawHitbox(screen *ebiten.Image) {
	op := e.getScreenTransform()
	iw, ih := e.sprite.Size()
	//x1, _ := op.GeoM.Apply(0, 0)
	x2, y2 := op.GeoM.Apply(float64(iw)/2, float64(ih)/2)
	//r1, _ := op.GeoM.Apply(float64(e.r), float64(e.r))
	//drawWireRect(screen, x1, y1, x2-x1, y2-y1, Red)
	//hitbox := createCircle(int(r1 - x1))
	rr :=zoomFactor*e.sprite_base_scale*float64(e.r)
	hitbox := createCircle(int(rr))
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x2-rr, y2-rr)
	screen.DrawImage(hitbox, op)
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
	screen.Clear()
	g.envLayer.Clear()
	g.entityLayer.Clear()
	g.debugLayer.Clear()
	g.guiLayer.Clear()

	screen.Fill(color.White)	// while there is no background

	// drawing the elements
	for _, e := range onScreenMap.floor {
		e.drawEntity(g.envLayer)
		e.drawHitbox(g.debugLayer)
	}
	for _, e := range onScreenMap.buildings {
		e.drawEntity(g.entityLayer)
		e.drawHitbox(g.debugLayer)
	}

	////////////
	drawSelectionRect(g.guiLayer)

	screen.DrawImage(g.envLayer, nil)
	screen.DrawImage(g.entityLayer, nil)
	screen.DrawImage(g.guiLayer, nil)
	screen.DrawImage(g.debugLayer, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if math.Abs(xx*xx+yy*yy - rr*rr) <= rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

func createCircle(r int) (*ebiten.Image) {
	dst := image.NewRGBA(image.Rect(0, 0, 2*r, 2*r))
	redctangle := image.NewRGBA(image.Rect(0, 0, 2*r, 2*r))
	draw.Draw(redctangle, redctangle.Bounds(), &image.Uniform{Red}, image.ZP, draw.Src)
	draw.DrawMask(dst, redctangle.Bounds(), redctangle, image.ZP, &circle{image.Point{r, r}, r}, image.ZP, draw.Over)
	return ebiten.NewImageFromImage(dst)
}


// Function to initialize the game
func Init(g *Game) {
	g.friendlyEntities = append(g.friendlyEntities, &dino)
	g.envEntities = append(g.envEntities, &tree)

	g.envLayer	= ebiten.NewImage(screenWidth, screenHeight)
	g.entityLayer = ebiten.NewImage(screenWidth, screenHeight)
	g.debugLayer  = ebiten.NewImage(screenWidth, screenHeight)
	g.guiLayer	= ebiten.NewImage(screenWidth, screenHeight)

	onScreenMap.buildings = append(onScreenMap.buildings, tree)
}

func get_player_location() (int, int) {
	return int(camera.x), int(camera.y)
}
func set_player_location(player_id string, x int, y int) {}


func main(){
	// starting the client
	go startClient(&client_chan)

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
