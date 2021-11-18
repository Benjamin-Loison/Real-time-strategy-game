package main

import (
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)

func main(){
	// Load the configuration
	config := loadConfig("conf/conf.json")

	// starting the client
	running = true
	go run_client(&client_chan, config, &running)

	// initializing the game
	g := &Game{}
	Init(g)
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
