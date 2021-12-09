package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
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
	for _, val := range channels {
		val <- msg
	}
}

func client_handler(conn net.Conn, map_path string, main_chan chan string, id int) {
	// Close the connection if the handler is exited
	defer conn.Close()

	// Load the map and send it to the client
	init_json := utils.ServerMessage { utils.MapInfo, gmap, nil , id}
	init_marshall, err := json.Marshal(init_json)
	utils.Check(err)
	init_message := []byte(string(init_marshall)+"\n")

	buffer := bufio.NewWriter(conn)
	_, err = buffer.Write(init_message)
	utils.Check(err)
	buffer.Flush()

	units_json := utils.ServerMessage { utils.MapInfo,utils.Map{}, Players, id}
	units_marshall, err := json.Marshal(units_json)
	utils.Check(err)
	units_message := []byte(string(units_marshall)+"\n")

	_, err = buffer.Write(units_message)
	utils.Check(err)
	buffer.Flush()

	// Wait for the other player
	utils.Logging("CLIENT_HANDLER", fmt.Sprintf("Waiting for other player (%d)", id))
	keepGoing := true
	for keepGoing {
		select {
			case x, _ :=<-main_chan :
				if x == "START" {
					keepGoing = false
				}
		}
	}
	/////////////////// Starting so sending go to the client
	_,err = buffer.Write([]byte("GO\n"))
	utils.Check(err)
	buffer.Flush()

	utils.Logging("CLIENT_HANDLER", fmt.Sprintf("Entering the main event loop (%d)", id))
	// Main Event loop
	// starting the client listener
	listener_chan := make(chan string)
	go listenClient(conn,listener_chan )

	keepGoing = true
	for keepGoing {
		select {
			case x, _ :=<-main_chan :
				if x == "QUIT" {
					conn.Write([]byte("QUIT\n"))
					keepGoing = false
				} else if strings.HasPrefix(x, "CHAT:") {
					conn.Write([]byte(fmt.Sprintf("%s\n", x)))
				}
			case x, _ :=<-listener_chan:
				utils.Logging("CLIENT_HANDLER","Received info from listener")
				if x == "QUIT" {
					utils.Logging("CLIENT_HANDLER", "Error when listening to the client")
					main_chan<-"CLIENT_ERROR"
					keepGoing = false
				}else {
					var client_event = &events.Event{}
					err = json.Unmarshal([]byte(x), client_event)
					if err != nil {
						utils.Logging("CLIENT_HANDLER", fmt.Sprintf("Error when receiving event from client (%d)", id))
					}
					// should now send to the updater
					utils.Logging("CLIENT_HANDLER","Sending info to updater")
					updater_chan<-*client_event
					utils.Logging("CLIENT_HANDLER","info for updater sent")
				}
		}
	}

	// Exit
	main_chan<-"FINISHED"
	utils.Logging("CLIENT_HANDLER", fmt.Sprintf("%d quits.", id))
	return
}

func listenClient(conn net.Conn, channel chan string) {
	// TODO utiliser la fonction register, pour assurer la mort du listenClient
	reader := bufio.NewReader(conn)
	for {
		utils.Logging("Listener","listening to client")
		netData, err := reader.ReadString('\n')
		netData = strings.TrimSpace(string(netData))
		utils.Logging("Listener","received from client: " + netData)
		if err == nil {
			channel <- netData
		} else {
			channel <- "QUIT"
			return
		}
	}
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
					utils.Logging("Updater", "Received an event")
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
	defer listener.Close()

	nb_clients := 0

	for nb_clients < 2 {
		conn, err := listener.Accept()

		utils.Check(err)
		go client_handler(conn, conf.MapPath,register(nb_clients), nb_clients)
		nb_clients += 1
	}

	go updater(channels, register(-1))
	broadcast(channels, "START")

	//start game
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
