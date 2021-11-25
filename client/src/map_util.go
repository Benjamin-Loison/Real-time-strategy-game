package main

import (
	"image/color"
    "io/ioutil"
    "fmt"
    "os"
    "time"
    "sync"
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

    TileSize  int32 = 32
    fontSize int32 = TileSize/4
    unit_size float32 = 0.4+float32(fontSize)
)

func logging(src string, message string) {
	fmt.Println(time.Now().Format(time.ANSIC) + "[" + src + "] " + message)
}

type Tile struct {
    Tile_Type TileType `json:"Tile_Type"`
    Startpoint Owner  `json:"Startpoint"`
}

type Map struct {
    Width int32 `json:"Width"`
    Height int32 `json:"Height"`
    Grid [][]Tile `json:"Grid"`
    mut *sync.Mutex
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
    return Map{Width :width, Height: height, Grid: m, mut: &sync.Mutex{}}
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
}

func DrawUnit(u Unit, owned bool){
    if owned {
        rl.DrawCircle(u.X,u.Y,unit_size,rl.Blue)
    }else{
        rl.DrawCircle(u.X,u.Y,unit_size,rl.Red)
    }
    rl.DrawText(u.Name,u.X-fontSize/2,u.Y-fontSize/2,fontSize,rl.White)
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
    //var m sync.Mutex
    result_map.mut = &sync.Mutex{}
    return *result_map
}
