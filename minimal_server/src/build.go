
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
	"fmt"
	"rts/utils"

	//rl "github.com/gen2brain/raylib-go/raylib"
)

type Building struct {
	Name string
	Position_x int32
	Position_y int32
	BuildDuration time.Duration
	BuildStartingTime time.Time
}

type TechnologicalTree_t struct {
	Name string `json:"Name"`
	Children []TechnologicalTree_t `json:"Children"`
}

func LoadTechnologicalTree() TechnologicalTree_t {
	// Read the configuration file
	file, err := ioutil.ReadFile("conf/topological_tree.conf")
	if err != nil {
		utils.Logging("TechoTree", fmt.Sprintf("Cannot open the config file: %v", err))
		os.Exit(-1)
	}

	// Convert the json content into a Configuration structure
	var technoTree = &TechnologicalTree_t{}
	err = json.Unmarshal(file, &technoTree)
	if err != nil {
		utils.Logging("TechnoTree", fmt.Sprintf("Cannot parse the config file: %v", err))
		os.Exit(1)
	}
	
	fmt.Println(*technoTree)

	return *technoTree
}

func Build(b Building, t TechnologicalTree_t) {
	gmap.Grid[b.Position_x][b.Position_y].Tile_Type = utils.Rock
}

