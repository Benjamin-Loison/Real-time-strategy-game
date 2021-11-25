package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
)

func logging(src string, msg string) {
	fmt.Println(time.Now().Format(time.ANSIC) + "[" + src + "] " + msg)
}

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

    TileSize  int32 = 32
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

func Check(e error) {
    if e != nil {
        panic(e)
    }
}


func MakeMap(width, height int32) Map {
    m := make([][]Tile, width)
    for i := range m {
        m[i] = make([]Tile, height)
        for j := range m[i] {
            m[i][j] = Tile{Tile_Type: None, Startpoint: NoOne}
        }
    }
    return Map{width, height, m}
}

func DrawTile(x,y int32,tileType TileType, startpoint Owner){
    switch tileType {
    case Rock :
        rl.DrawCircle(TileSize*x+TileSize/2.0,TileSize*y+TileSize/2.0,0.5*float32(TileSize),rl.DarkGray)
    case Gold :
        rl.DrawCircle(TileSize*x+TileSize/2.0,TileSize*y+TileSize/2.0,0.5*float32(TileSize),rl.Gold)
    case Tree :
        rl.DrawCircle(TileSize*x+TileSize/2.0,TileSize*y+TileSize/2.0,0.5*float32(TileSize),rl.Green)
    default:
    }
    switch startpoint {
    case Player1 :
        rl.DrawText( "P1" ,TileSize*x+TileSize/5.0,TileSize*y+TileSize/5.0, fontSize,rl.Blue)
    case Player2 :
        rl.DrawText( "P2" ,TileSize*x+TileSize/5.0,TileSize*y+TileSize/5.0, fontSize,rl.Red)
    default:
    }
}

func DrawMap(gmap Map) {
    rl.DrawRectangle(0,0,gmap.Width*TileSize,gmap.Height*TileSize, GroundColor)

    for i:= int32(0); i < gmap.Width ; i++ {
        for j:= int32(0); j < gmap.Height ; j++ {
            DrawTile(i,j,gmap.Grid[i][j].Tile_Type, gmap.Grid[i][j].Startpoint)
        }
    }
}

func SaveMap(gmap Map){
    map_json, err := json.MarshalIndent(gmap,"","   ")

    Check(err)

    data := []byte(string(map_json))
    err = os.WriteFile("./map.json", data, 0644)

    Check(err)
}

func PrintMap(gmap Map){
    map_json, err := json.MarshalIndent(gmap,"","   ")

    Check(err)
    fmt.Printf(string(map_json))
}

func LoadMap(path string) Map {
    map_data, err := ioutil.ReadFile(path)
    Check(err)
    var result_map = &Map{}
    err = json.Unmarshal([]byte(string(map_data)), result_map)
    Check(err)
    return *result_map
}

func initializePlayer(gmap *Map, own Owner, units *map[string]Unit,id *int){
    for i := int32(0) ; i < gmap.Width ; i++ {
        for j := int32(0) ; j < gmap.Width ; j++ {
            if gmap.Grid[i][j].Startpoint == own {
                id_unit := getId(id)
                (*units)[strconv.Itoa(id_unit)] = Unit{TileSize*i+TileSize/2.0,TileSize*j+TileSize/2.0,"P",int32(id_unit),own}
            }
        }
    }
}
