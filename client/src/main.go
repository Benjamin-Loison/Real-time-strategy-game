package main

import (
	_ "image/png"
	"os"
)

var (
    config Configuration_t
)

type ServerMap Map

func main(){
	// Load the configuration
	config = loadConfig("conf/conf.json")

    game_map := LoadMap("../maps/map.json")

    chan_client := make(chan string, 2)

	// starting a co-process to deal with the server
    go run_client(config,&game_map,chan_client)

    //wait to have received all data before starting gui
    for {
        select {
        case x, _ := <-chan_client:
            if x == "OK" {
                break
            }else if x == "PANIC" {
                logging("client","Error while retrieving map data")
                os.Exit(-1)
            }
        default:
        }
    }

    //running gui
    RunGui(&game_map, config)

    chan_client <- "QUIT"
}
