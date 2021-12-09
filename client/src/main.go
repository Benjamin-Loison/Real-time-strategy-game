package main

import (
	_ "image/png"
	"os"
    "rts/utils"
)

var (
	config Configuration_t
	config_menus MenuConfiguration_t
	players []utils.Player
	client_id int
)

type ServerMap utils.Map

func main(){
	// Load the configuration
	config = loadConfig("conf/conf.json")
	config_menus = loadTextMenus("conf/menus.json")
	//TODO: check theat there is no conflict between the two configuration files

	game_map := utils.Map{}

	chan_client := make(chan string, 3)

	// starting a co-process to deal with the server
	go run_client(config,&players,&game_map,chan_client)

	//wait to have received all data before starting gui

	ok := false

	for {
		if ok {
			break
		}
		select {
		case x, _ := <-chan_client:
			if x == "OK" {
				ok = true
			}else if x == "QUIT" {
				utils.Logging("client","Error while retrieving map data")
				os.Exit(-1)
			}
		default:
		}
	}

	//running gui
	RunGui(&game_map, &players, config, config_menus,chan_client)

	chan_client <- "QUIT"
}
