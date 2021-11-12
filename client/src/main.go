package main

import (
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
)
var client_chan chan string

func main(){
	// starting the client
	go startClient(&client_chan)

	// initializing the game
	g := &Game{}
	Init(g)
	//initializing ebiten
	ebiten.SetWindowSize(screenWidth,screenHeight)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("EHO: Elves, humans and orks")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
