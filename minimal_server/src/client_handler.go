package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"rts/events"
	"rts/utils"
)

func load_map (writer *bufio.Writer, id int) {
	init_json := utils.ServerMessage { utils.MapInfo, gmap, nil , id}
	init_marshall, err := json.Marshal(init_json)
	utils.Check(err)
	init_message := []byte(string(init_marshall)+"\n")

	_, err = writer.Write(init_message)
	utils.Check(err)
	writer.Flush()

	units_json := utils.ServerMessage { utils.MapInfo,utils.Map{}, Players, id}
	units_marshall, err := json.Marshal(units_json)
	utils.Check(err)
	units_message := []byte(string(units_marshall)+"\n")

	_, err = writer.Write(units_message)
	utils.Check(err)
	writer.Flush()
}

func wait_for_start(writer *bufio.Writer, main_chan chan string, id int) {
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
	utils.Logging("CLIENT_HANDLER", fmt.Sprintf("GO (%d)", id))
	_,err := writer.Write([]byte("GO\n"))
	utils.Check(err)
	writer.Flush()
}

func client_handler(conn net.Conn, map_path string, main_chan chan string, id int) {
	// Close the connection if the handler is exited
	defer conn.Close()

	writer := bufio.NewWriter(conn)

	// Send the map to the current client and wait for the other player
	load_map(writer, id)
	wait_for_start(writer, main_chan, id)



	// Main Event loop
	// starting the client listener
	listener_chan := make(chan string)
	go listenClient(conn, listener_chan)

	keepGoing := true
	utils.Logging("CLIENT_HANDLER", fmt.Sprintf("Entering the main event loop (%d)", id))
	for keepGoing {
		select {
			case x, _ :=<-main_chan :
				if x == "QUIT" {
					writer.Write([]byte("QUIT\n"))
					writer.Flush()
					utils.Logging("CLIENT_HANDLER",
						fmt.Sprintf("(%d) stop sent: %s",
							id,
							x))
					keepGoing = false
				} else if strings.HasPrefix(x, "CHAT:") {
					utils.Logging("CLIENT_HANDLER",
						fmt.Sprintf("(%d) chat string will be send: %s",
							id,
							x))
					writer.Write([]byte(fmt.Sprintf("%s\n", x)))
					writer.Flush()
					utils.Logging("CLIENT_HANDLER",
						fmt.Sprintf("(%d) chat string sent: %s",
							id,
							x))
				} else {
					utils.Logging("CLIENT_HANDLER",
						fmt.Sprintf("(%d) did not send %s",
							id,
							x))
				}
			case x, _ :=<-listener_chan:
				utils.Logging("CLIENT_HANDLER","Received info from listener")
				if x == "QUIT" {
					utils.Logging("CLIENT_HANDLER", "Error when listening to the client")
					main_chan<-"CLIENT_ERROR"
					keepGoing = false
				} else {
					var client_event = &events.Event{}
					err := json.Unmarshal([]byte(x), client_event)
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

