package utils

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"strconv"
	"time"
    "strings"
    "rts/factory"

	"github.com/gen2brain/raylib-go/raylib"
)

func Logging(src string, msg string) {
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


const (

    TileSize  int32 = 32
	fontSize int32 = TileSize/4
	unit_size float32 = 0.4+float32(fontSize)
)

type Tile struct {
    Tile_Type TileType `json:"Tile_Type"`
    Startpoint factory.Owner  `json:"Startpoint"`
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
            m[i][j] = Tile{Tile_Type: None, Startpoint: factory.NoOne}
        }
    }
    return Map{width, height, m}
}

func DrawTile(x,y int32,tileType TileType, startpoint factory.Owner){
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


type ServerMessageType int32

const (
	MapInfo ServerMessageType = 0
	Startingfactory = 1
	Update = 2
)


func getId(seed *int) int {
	ret := *seed
	(*seed) ++
	return ret
}

type Player struct {
	Units map[string]factory.Unit `json:"Units"`
	Seed int `json:"Seed"`
}

type ServerMessage struct {
	MessageType ServerMessageType `json:"MessageType"`
	GameMap Map `json:"GameMap"`
	Players []Player `json:"Players"`
	Id int `json:"Id"`
}

//Raylib represents keys as int32 values
func KeyOfString(s string)(int32) {
	switch l := len(s); l {
		case 0:
			panic("An empty string cannot represent a key in the configuration file.")
		case 1:
			return int32([]rune(strings.ToUpper(s))[0])
		default:
			switch s {
				case "SPACE":
					return rl.KeySpace
				case "RIGHT":
					return rl.KeyRight
				case "LEFT":
					return rl.KeyLeft
				case "DOWN":
					return rl.KeyDown
				case "UP":
					return rl.KeyUp
				default:
					panic("Not implemented: recognition of non-ascii characters and description")
			}
	}
}


func InitializePlayer(gmap *Map, own factory.Owner, units *map[string]factory.Unit,id *int){
    for i := int32(0) ; i < gmap.Width ; i++ {
        for j := int32(0) ; j < gmap.Height ; j++ {
            if gmap.Grid[i][j].Startpoint == own {
                id_unit := getId(id)
                (*units)[strconv.Itoa(id_unit)] = factory.MakeHumanPeon(TileSize*i+TileSize/2.0, TileSize*j+TileSize/2.0, int32(id_unit), own)
            }
        }
    }
}

func DrawUnit(u factory.Unit, owned bool){
	if owned {
        rl.DrawCircle(u.X,u.Y,unit_size,rl.Blue)
	}else{
        rl.DrawCircle(u.X,u.Y,unit_size,rl.Red)
	}
	rl.DrawText(u.Name,u.X-fontSize/2,u.Y-fontSize/2,fontSize,rl.White)
}


