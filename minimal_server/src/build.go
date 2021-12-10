
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
	"fmt"
	"rts/events"
	"rts/utils"

)

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

func Build(channels map[int]chan string,
			e *events.BuildBuilding_e,
			technoTree TechnologicalTree_t) {
	// Check wether or not the building exists in the
	// technological tree
	//TODO!



	// Check wether or not the player is allowed to build it
	// TODO!



	// Launch the build
	time.Sleep(time.Second)


	// Broadcast a message to keep the players posted
	broadcast(channels, fmt.Sprintf("CHAT:building %s", e.BuildingName))
}
