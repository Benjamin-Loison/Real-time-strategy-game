package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"rts/events"
	"rts/utils"
	"rts/factory"
)

const (
	serverSpeed = 10
)

var (
	Players []utils.Player
	channels = make(map[int]chan string)
	updater_chan = make(chan events.Event)
	gmap utils.Map
	serverTime uint32
)

func broadcast(channels map[int]chan string, msg string) {
	utils.Logging("broadcast", msg)
	channels[0] <- msg // Sending to player 0
	channels[1] <- msg // Sending to player 1
}

func register(id int) chan string {
	channels[id] = make(chan string)
	return channels[id]
}

func updater(channels map[int]chan string, stopper_chan chan string){
	for{
		select{
		case e := <-updater_chan:
			switch e.EventType {
				case events.ChatEvent:
					utils.Logging("Updater (CHAT)", fmt.Sprintf("%s", e.Data))
					broadcast(channels, fmt.Sprintf("CHAT:%s", e.Data))
					break
				default:
					utils.Logging("Updater", "Received an unhandeled event")
					break
			}
		case s := <-stopper_chan:
			if s == "QUIT" {
				utils.Logging("Updater", "Quitting")
				os.Exit(0)
			}
		default:
			break
		}
	}
}

func main() {
	// Load the config file
	conf := loadConfig("conf/conf.json")
	gmap = utils.LoadMap(conf.MapPath)

	// initializing
	Players = append(Players, utils.Player{Units: map[string]factory.Unit{},Seed: 0} )
	Players = append(Players, utils.Player{Units: map[string]factory.Unit{},Seed: 0} )

	utils.InitializePlayer(&gmap, factory.Player1,&Players[0].Units, &Players[0].Seed)
	utils.InitializePlayer(&gmap, factory.Player2,&Players[1].Units, &Players[1].Seed)

	// Listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",conf.Hostname, conf.Port))
	utils.Check(err)

	nb_clients := 0

	for nb_clients < 2 {
		conn, err := listener.Accept()

		utils.Check(err)
		go client_handler(conn, conf.MapPath, register(nb_clients), nb_clients)
		nb_clients += 1
	}
	// Close the listener as soon as the two clients are connected
	listener.Close()

	// Launch the updater function and start the game
	go updater(channels, register(-1))
	utils.Logging("Server", "Send start!")
	broadcast(channels, "START")
	go gameLoop(register(-2))

	//wait for end game
	for nb_clients > 0{
		for _, c := range channels {
			select {
			case s := <-c :
				switch(s) {
				case "FINISHED" :
					nb_clients--
					utils.Logging("Server", fmt.Sprintf("%d remaining clients.", nb_clients))
					break
				case "STOP" :
					nb_clients = 0
					break
				case "CLIENT_ERROR" :
					utils.Logging("Server", "Received client error")
					break
				default:
					break
				}
			default:
				break
			}
		}
	}
	broadcast(channels, "QUIT")
	os.Exit(0)
}

func gameLoop(quit chan string){
	gameOver := false
	for !gameOver {
		select {
		case s :=<-quit :
			if s=="QUIT"{
				return
			}
			break
		default: // not quitting, updating game state
			serverTime += 1
			//utils.Logging("GameLoop",fmt.Sprintf("Server time = %d",serverTime))
			break
		}
		time.Sleep(serverSpeed * time.Millisecond)
	}
	quit<-"STOP"
}
