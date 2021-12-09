package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
    "os"
	"strings"

	"rts/events"
	"rts/utils"
)

var (
	Players []utils.Player
	channels = make(map[int]chan string)
    updater_chan = make(chan events.Event)
	gmap utils.Map
)

func broadcast(channels map[int]chan string, msg string) {
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
    // ajouter soi même dans la liste des channels pour bien être fermé
    // -> to do so, créer une fonction register.
    reader := bufio.NewReader(conn)
    for {
        utils.Logging("Listener","listening to client")
        netData, err := reader.ReadString('\n')
        utils.Logging("Listener","received from client")
        netData = strings.TrimSpace(string(netData))
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

func updater(stopper_chan chan string){
    for{
        select{
        case e := <-updater_chan:
            switch e {
                default:
                    utils.Logging("Updater", "Received an event")
            }
        case s := <-stopper_chan:
            if s == "QUIT" {
                utils.Logging("Updater", "Quitting")
                os.Exit(0)
            }
        }
    }
}

func main() {
	// Load the config file
	conf := loadConfig("conf/conf.json")
	gmap = utils.LoadMap(conf.MapPath)

	// initializing
	Players = append(Players, utils.Player{Units: map[string]utils.Unit{},Seed: 0} )
	Players = append(Players, utils.Player{Units: map[string]utils.Unit{},Seed: 0} )

	utils.InitializePlayer(&gmap, utils.Player1,&Players[0].Units, &Players[0].Seed)
	utils.InitializePlayer(&gmap, utils.Player2,&Players[1].Units, &Players[1].Seed)

    go updater(register(-1))

	// Listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",conf.Hostname, conf.Port))
	utils.Check(err)
	defer listener.Close()

	ok := 0

	for {
		// Wait for a new connection
		if ok == 2 {
			break
		}
		conn, err := listener.Accept()

		utils.Check(err)
		go client_handler(conn, conf.MapPath,register(ok), ok)
		ok += 1
	}

	broadcast(channels, "START")
	//
	//broadcast(channels, "QUIT")

	//start game
	for {
		if ok == 0 {
			break
		}
		for _, c := range channels {
			select {
			case s := <-c :
				if s == "FINISHED" {
					ok--
				}
                if s == "CLIENT_ERROR" {
					utils.Logging("Server", "Received client error")
				}
			default:
			}
		}
	}
}

