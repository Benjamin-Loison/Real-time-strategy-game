package main

import (
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)

var (
	// Stores teh server ID string
	serverID string

	// server_queries (resp. client_queries) map[string] string store the active queries,
	// i.e. queries which have not been set valid by the return of an "ok"
	// status (resp. while no correcr answer has been sent).
	server_queries map[string]string
	client_queries map[string]string

	// The random client ID
	client_id string

	// The ip address and port of the server
	host string
	port int

	// The channel that interacts with the gui part of the code (not required at
	// the moment)
	gui_chan chan string
	map_chan chan ServerMap
)

type ServerMap struct {
	x_init, y_init, w, h int
	loc []Location

}

func main(){
	// Load the configuration
	config := loadConfig("conf/conf.json")

	// starting a co-process to deal with the server
	go run_client(&client_chan, &map_chan, config)

	// initializing the game
	g := &Game{}
	Init(g,config.Keys)
	//initializing ebiten
	ebiten.SetWindowSize(screenWidth,screenHeight)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("EHO: Elves, humans and orks")
	if err := ebiten.RunGame(g); err != nil {
		logging("GUI", "Unexpected failure of the graphical engine.")
		// Question: est-ce qu'on garde log? Si oui, est-ce que c'est aps
		// overkill (log systèmes ? ou  stderr, quel context ?)
		log.Fatal(err)
		return
	}
}
