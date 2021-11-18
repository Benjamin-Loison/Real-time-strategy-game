package main

import (
	_ "image/png"
)

var (
    config Configuration_t
)

type ServerMap Map

func main(){
	// Load the configuration
	config = loadConfig("conf/conf.json")

    game_map := LoadMap("../maps/map.json")

	// starting a co-process to deal with the server
    go run_client(config)

    //running gui
    RunGui(&game_map, config)
}
