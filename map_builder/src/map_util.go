package main

import (
	"image/color"
    "io/ioutil"
    "fmt"
    "os"
    "encoding/json"
	"github.com/gen2brain/raylib-go/raylib"
)

type TileType int64

var (
    GroundColor = color.RGBA{20,99,6,255}
)

const (
    Rock TileType = 0
	Tree          = 1
	Gold          = 2
    None          = 3
)

type Owner int64

const (
    Player1 Owner = 1
    Player2       = 2
    NoOne         = 0

    tileSize  int32 = 32
    fontSize int32 = 20
)

type Tile struct {
    Tile_Type TileType `json:"Tile_Type"`
    Startpoint Owner  `json:"Startpoint"`
}

type Map struct {
    Width int32 `json:"Width"`
    Height int32 `json:"Height"`
    Grid [][]Tile `json:"Grid"`
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func makeMap(width, height int32) Map {
    m := make([][]Tile, width)
    for i := range m {
        m[i] = make([]Tile, height)
        for j := range m[i] {
            m[i][j] = Tile{Tile_Type: None, Startpoint: NoOne}
        }
    }
    return Map{width, height, m}
}

func drawTile(x,y int32,tileType TileType, startpoint Owner){
    switch tileType {
    case Rock :
        rl.DrawCircle(tileSize*x+tileSize/2.0,tileSize*y+tileSize/2.0,0.5*float32(tileSize),rl.DarkGray)
    case Gold :
        rl.DrawCircle(tileSize*x+tileSize/2.0,tileSize*y+tileSize/2.0,0.5*float32(tileSize),rl.Gold)
    case Tree :
        rl.DrawCircle(tileSize*x+tileSize/2.0,tileSize*y+tileSize/2.0,0.5*float32(tileSize),rl.Green)
    default:
    }
    switch startpoint {
    case Player1 :
        rl.DrawText( "P1" ,tileSize*x+tileSize/5.0,tileSize*y+tileSize/5.0, fontSize,rl.Blue)
    case Player2 :
        rl.DrawText( "P2" ,tileSize*x+tileSize/5.0,tileSize*y+tileSize/5.0, fontSize,rl.Red)
    default:
    }
}

func drawMap(gmap Map) {
    rl.DrawRectangle(0,0,gmap.Width*tileSize,gmap.Height*tileSize, GroundColor)

    for i:= int32(0); i < gmap.Width ; i++ {
        for j:= int32(0); j < gmap.Height ; j++ {
            drawTile(i,j,gmap.Grid[i][j].Tile_Type, gmap.Grid[i][j].Startpoint)
        }
    }
}

func saveMap(gmap Map){
    map_json, err := json.MarshalIndent(gmap,"","   ")

    check(err)

    data := []byte(string(map_json))
    err = os.WriteFile("./map.json", data, 0644)

    check(err)
}

func printMap(gmap Map){
    map_json, err := json.MarshalIndent(gmap,"","   ")

    check(err)
    fmt.Printf(string(map_json))
}

func loadMap(path string) Map {
    map_data, err := ioutil.ReadFile(path)
    check(err)
    var result_map = &Map{}
    err = json.Unmarshal([]byte(string(map_data)), result_map)
    check(err)
    return *result_map
}
