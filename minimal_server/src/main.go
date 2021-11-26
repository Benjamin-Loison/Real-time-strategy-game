package main

import(
	"bufio"
	"net"
	"encoding/json"
	"fmt"
)

var (
	Players []Player
	main_chan = make(chan string,2)
	gmap Map
)

func getId(seed *int) int {
	ret := *seed
	(*seed) ++
	return ret
}


func broadcast(channels map[int]chan string, msg string) {
	for _, val := range channels {
		val <- msg
	}
}

func client_handler(conn net.Conn, map_path string, main_chan chan string, id int) {
	// Close the connection if the handler is exited
	defer conn.Close()

	// Load the map and send it to the client
	init_json := ServerMessage { MapInfo, gmap, nil , id}
	init_marshall, err := json.Marshal(init_json)
	Check(err)
	init_message := []byte(string(init_marshall)+"\n")

	buffer := bufio.NewWriter(conn)
	_, err = buffer.Write(init_message)
	Check(err)
	buffer.Flush()

	units_json := ServerMessage { MapInfo,Map{}, Players, id}
	units_marshall, err := json.Marshal(units_json)
	Check(err)
	units_message := []byte(string(units_marshall)+"\n")

	_, err = buffer.Write(units_message)
	Check(err)
	buffer.Flush()

	// Wait for the other player
	logging("CLIENT_HANDLER", fmt.Sprintf("Waiting for other player (%d)", id))
	keepGoing := true
	for keepGoing {
		select {
			case x, _ :=<- main_chan :
				if x == "START" {
					keepGoing = false
				}
		}
	}

	logging("CLIENT_HANDLER", fmt.Sprintf("Entering the main event loop (%d)", id))
	// Main Event loop
	keepGoing = true
	for keepGoing {
		select {
			case x, _ :=<- main_chan :
				if x == "QUIT" {
					conn.Write([]byte("QUIT\n"))
					keepGoing = false
				}
		}
	}

	// Exit
	main_chan<-"FINISHED"
	logging("CLIENT_HANDLER", fmt.Sprintf("%d quits.", id))
	return
}

func main() {
	// Load the config file
	conf := loadConfig("conf/conf.json")
	gmap = LoadMap(conf.MapPath)

	// initializing
	Players = append(Players, Player{Units: map[string]Unit{},Seed: 0} )
	Players = append(Players, Player{Units: map[string]Unit{},Seed: 0} )

	channels := make(map[int]chan string)

	initializePlayer(&gmap, Player1,&Players[0].Units, &Players[0].Seed)
	initializePlayer(&gmap, Player2,&Players[1].Units, &Players[1].Seed)

	// Listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",conf.Hostname, conf.Port))
	Check(err)
	defer listener.Close()

	ok := 0

	for {
		// Wait for a new connection
		if ok == 2 {
			break
		}
		conn, err := listener.Accept()

		Check(err)
		channels[ok] = make(chan string)
		go client_handler(conn, conf.MapPath,channels[ok], ok)
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
			default:
			}
		}
	}
}

