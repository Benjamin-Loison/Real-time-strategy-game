package main

import (
	"errors"
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
	ZOOM_STEP = 1.015
	cameraBorderThreshold = 100
	zoomMin = 0.2
	zoomMax = 5
)

var (
	Red	   = color.RGBA{255, 0, 0, 255}
	color_selected	 = color.RGBA{0, 0, 255, 255}
	inSelection = false
	startSelection = [2]int{0, 0}
	endSelection = [2]int{0, 0}
	sprite = loadImageFromFile("Dino_blue.png")
	treeSprite = loadImageFromFile("baum.png")
	floorSprite = loadImageFromFile("floor.png")
	dino = Entity{0, 0, 12, sprite.SubImage(image.Rect(0, 0, 24, 24)).(*ebiten.Image), 6, false}
	tree = Entity{640, 360, 32, treeSprite, 1.0, false}
	floor = Entity{640, 360, 32, floorSprite, 1.0, false}
	camera = Entity{0, 0, 0,  nil , 1, false}	// TODO should be change to another more adapted type
	onScreenMap = Map{0, 0, 0, 0, make([]Entity, 0), make([]Entity, 0), make([]*Entity, 0)}
	zoomFactor = 1.0
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
	entities []*Entity
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

// createDummyEntity returns an entity whose position is (x, y) on the map
func createDummyEntity(x, y int) (Entity) {
	return Entity{float64(x), float64(y), 32, treeSprite, 1, false}
}

func createFloor(x, y int) (Entity) {
	return Entity{float64(x), float64(y), 32, floorSprite, 1, false}
}

// update_map updates the onScreenMap variable using informations provided by the client
// init_x, init_y are the coordinates of the upper left corner of the provided area
// w (resp. h) is the width (resp. height) of the provided area
// location_list contains all the points of interest in the provided area
func update_map(init_x , init_y, w, h int, location_list []Location) {
	
	onScreenMap.origin_x = init_x - w/2 - screenWidth
	onScreenMap.origin_y = init_y - h/2 - screenHeight
	onScreenMap.w = w
	onScreenMap.h = h
	onScreenMap.buildings = make([]Entity, 0)
	onScreenMap.floor = make([]Entity, 0)

	for i, l := range location_list {
		switch(l.loc_type) {
			case ElfBuilding:
				onScreenMap.buildings = append(onScreenMap.buildings, createDummyEntity(i%w, i/w))
			case OrcBuilding:
				onScreenMap.buildings = append(onScreenMap.buildings, createDummyEntity(i%w, i/w))
			case HumanBuilding:
				onScreenMap.buildings = append(onScreenMap.buildings, createDummyEntity(i%w, i/w))
			case Floor:
				onScreenMap.floor = append(onScreenMap.floor, createFloor(i%w, i/w))
				logging("Update_map", "I have created a floor, are you pround enough or should I print it ?")
			default:
				continue
		}
	}
	logging("Update", "Map updated.")
}

// Update handle all the operations that should be done at every tick
// for example looking and handling keyboard inputs
func (g *Game) Update() error {
	select {
		case x, _ := <-gui_chan :
			if x == "QUIT" {
				panic(errors.New("Exiting"))
			}
		case x, _ := <- map_chan :
			logging("Update", "map recieved.")
			update_map(x.x_init, x.y_init, x.w, x.h, x.loc)
		default :
			// nop
	}

	//////////// Handling Keyboard events ////////////
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	mov_size := 5 / zoomFactor
	for _, p := range g.keys {
		switch s := p.String(); s {
		// character movement
		case "S":
			dino.y += 5
		case "Z":
			dino.y -= 5
		case "Q":
			dino.x -= 5
		case "D":
			dino.x += 5

		// camera movement
		case "ArrowUp":
			camera.y -= mov_size
			startSelection[1] += int(mov_size * zoomFactor)
		case "ArrowDown":
			camera.y += mov_size
			startSelection[1] -= int(mov_size * zoomFactor)
		case "ArrowLeft":
			camera.x -= mov_size
			startSelection[0] += int(mov_size * zoomFactor)
		case "ArrowRight":
			camera.x += mov_size
			startSelection[0] -= int(mov_size * zoomFactor)

		// zoom changes
		case "I":
			zoomFactor *= ZOOM_STEP
		case "K":
			zoomFactor /= ZOOM_STEP
		}
	}
	////////////////////////////////////////////////
	//////////// Handling camera mouse control /////////////
	x, y := ebiten.CursorPosition()
	if x <= cameraBorderThreshold && x >= 0 {
		camera.x -= mov_size
		startSelection[0] += int(mov_size * zoomFactor)
	} else if x >= screenWidth - cameraBorderThreshold && x <= screenWidth {
		camera.x += mov_size
		startSelection[0] -= int(mov_size * zoomFactor)
	}
	if y <= cameraBorderThreshold && y >= 0 {
		camera.y -= mov_size
		startSelection[1] += int(mov_size * zoomFactor)
	} else if y >= screenHeight - cameraBorderThreshold && y <= screenHeight {
		camera.y += mov_size
		startSelection[1] -= int(mov_size * zoomFactor)
	}
	////////////////////////////////////////////////
	if zoomFactor < zoomMin {
		zoomFactor = zoomMin
	}
	if zoomFactor > zoomMax {
		zoomFactor = zoomMax
	}
	

	return nil
}

func loadImageFromFile(path string) *ebiten.Image {
	path = "media/sprites/" + path
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

// drawWireRect draws a wireframe rectangle
func drawWireRect(screen *ebiten.Image, x, y, w, h float64, c color.Color) {
	ebitenutil.DrawLine(screen, x   , y  , x+w , y   , c)
	ebitenutil.DrawLine(screen, x   , y  , x   , y+h , c)
	ebitenutil.DrawLine(screen, x+w , y  , x+w , y+h , c)
	ebitenutil.DrawLine(screen, x   , y+h, x+w , y+h , c)
}

// drawSelectionRect draws to screen the selection box
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
	} else if inSelection {
		inSelection = false
		x  := float64(startSelection[0])
		y  := float64(startSelection[1])
		dx := float64(endSelection[0])
		dy := float64(endSelection[1])
		selectUnits(x, y, x+dx, y+dy)
	}
	if inSelection {
		x  := float64(startSelection[0])
		y  := float64(startSelection[1])
		dx := float64(endSelection[0])
		dy := float64(endSelection[1])

		drawWireRect(screen, x, y, dx, dy, Red)
	}
}

func getSign(x float64) (float64) {
	if x > 0 {
		return 1
	}
	return -1
}
func abs(x float64) (float64) {
	if x < 0 {
		return -x
	}
	return x
}

/* intersectRectangleCircle returns true iff there is an intersection
 * between a rectangle whose centre is at (x_r, y_r) and of dimension w*h
 * and a circle of radius r and of centre (x_c, y_c) */
func intersectRectangleCircle(x_r, y_r, w, h, x_c, y_c, r float64) (bool){
	dx := abs(x_c - x_r);
	dy := abs(y_c - y_r);
	// dx, dy between the centre of the circle and the centre of the rectangle
	dx = abs(x_c - (x_r+w/2))
	dy = abs(y_c - (y_r+h/2))

    if (dx > (w/2 + r)) { return false; }
    if (dy > (h/2 + r)) { return false; }

    if (dx <= (w/2)) { return true; }
    if (dy <= (h/2)) { return true; }

	dx = dx - w/2
	dy = dy - h/2

	cornerDistance_sq := dx*dx + dy*dy;	// distance squared of the centre of the circle to the closest corner of the rectangle

    return (cornerDistance_sq <= (r*r));
}

/* pointInRectangle returns true iff (x, y) lies in the rectangle
 * defined by a corner at (x_r, y_r) and the *signed* width and height w, h */
 // TODO to fix : it is possible to select from the upper left without intersection
func pointInRectangle(x, y, x_r, y_r, w, h float64) (bool) {
	sign_x := getSign(w)
	sign_y := getSign(h)
	return (sign_x*x_r <= sign_x*x && sign_x*x <= sign_x*(x_r+w) && sign_y*y_r <= sign_y*y && sign_y*y <= sign_y*(y_r+h))
}

func selectUnits(x1, y1, x2, y2 float64) {
	for _, e := range onScreenMap.entities {
		e.selected = false
		op := e.getScreenTranslation()
		e_x, e_y := op.GeoM.Apply(0, 0)
		r := float64(e.r) * zoomFactor*e.sprite_base_scale
		x_r := x1 + (x2 - x1)
		y_r := y1 + (y2 - y1)
		e.selected = intersectRectangleCircle(x_r, y_r, abs(x2-x1), abs(y2-y1), e_x, e_y, r) ||
						pointInRectangle(e_x, e_y, x1, y1, (x2 - x1), (y2 - y1))
	}
}


// getScreenTransform returns the transformation to be applied to draw entity e
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


// getScreenTranslation returns only the translation part of the transformation needed
// to draw e
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

// drawHitbox draws the collision box of e onto screen
func (e Entity) drawHitbox(screen *ebiten.Image) {
	op := e.getScreenTransform()
	iw, ih := e.sprite.Size()
	x2, y2 := op.GeoM.Apply(float64(iw)/2, float64(ih)/2)
	rr :=zoomFactor*e.sprite_base_scale*float64(e.r)
	hitbox := createCircle(int(rr), Red)
	if e.selected {
		hitbox = createCircle(int(rr), color_selected)
	}
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

// getSelctionRect returns the coordinates of the top left corner
// as well as its width and height
func getSelctionRect() (int, int, int, int) {
	return startSelection[0], startSelection[1], endSelection[0], endSelection[1]
}

// Draw is called on every game tick to update the displayed image
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
	for _, e := range onScreenMap.entities {
		e.drawEntity(g.entityLayer)
		e.drawHitbox(g.debugLayer)
	}

	drawSelectionRect(g.guiLayer)
	////////////

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

// createCircle returns an image containing a circle of radius r
func createCircle(r int, c color.Color) (*ebiten.Image) {
	dst := image.NewRGBA(image.Rect(0, 0, 2*r, 2*r))
	redctangle := image.NewRGBA(image.Rect(0, 0, 2*r, 2*r))
	draw.Draw(redctangle, redctangle.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
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
	onScreenMap.floor = append(onScreenMap.floor, floor)
	onScreenMap.entities = append(onScreenMap.entities, &dino)
}

func get_player_location() (int, int) {
	return int(camera.x), int(camera.y)
}
func set_player_location(player_id string, x int, y int) {}

