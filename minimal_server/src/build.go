
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
	BuildDuration int `json:"BuildDuration"`
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

func isIn(slice []Building, element string) bool {
	for i := 0 ; i < len(slice) ; i ++ {
		if slice[i].Name == element {
			return true
		}
	}
	return false
}

func CheckRights(tree TechnologicalTree_t, currentBuildings []Building, b Building) bool {
	if tree.Name == b.Name || isIn(currentBuildings, b.Name) {
		return true
	}
	res := false
	if isIn(currentBuildings, tree.Name) {
		for i := 0 ; i < len(tree.Children) ; i ++ {
			res = res || CheckRights(tree.Children[i], currentBuildings, b)
		}
	}
	return res
}

func Build(b Building, t TechnologicalTree_t) {
	fmt.Println(b.Name)
	switch (b.Name) {
		case "TownHall":
			utils.Logging("Build", "TownHall added.")
			gmap.Grid[b.Position_x][b.Position_y].Tile_Type = utils.TownHall
			break
		case "House":
			utils.Logging("Build", "House added.")
			gmap.Grid[b.Position_x][b.Position_y].Tile_Type = utils.House
			break
		default:
			utils.Logging("Build", "Unknown type: " + b.Name)
			break
	}
}

func getDuration(name string, tree TechnologicalTree_t) int {
	if tree.Name == name {
		return tree.BuildDuration
	}
	for i := 0 ; i < len(tree.Children) ; i ++ {
		// Checking child i
		res := getDuration(name, tree.Children[i])
		if res >= 0 {
			return res
		}
	}
	return -1
}

