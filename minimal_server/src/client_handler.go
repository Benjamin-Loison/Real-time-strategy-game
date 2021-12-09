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

func client_handler(conn net.Conn, map_path string, main_chan chan string, id int) {
	// Close the connection if the handler is exited
	defer conn.Close()

	writer := bufio.NewWriter(conn)

	// Send the map to the current client and wait for the other player
	load_map(writer, id)
	wait_for_start(writer, main_chan, id)

	// starting the client listener
	listener_chan := make(chan string)
	go listenClient(id, conn, listener_chan)

	// Main Event loop
	utils.Logging("CLIENT_HANDLER", fmt.Sprintf("Entering the main event loop (%d)", id))
	keepGoing := true
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
						fmt.Sprintf("(%d) did not understand %s",
							id,
							x))
				}
				break

			case x, _ :=<-listener_chan:
				if x == "QUIT" {
					utils.Logging("CLIENT_HANDLER",
						fmt.Sprintf("(%d) Error when listening to the client",
							id))
					main_chan<-"CLIENT_ERROR"
					keepGoing = false
				} else {
					var client_event = &events.Event{}
					err := json.Unmarshal([]byte(x), client_event)
					if err != nil {
						utils.Logging("CLIENT_HANDLER",
							fmt.Sprintf("(%d) Error when receiving event from client.", id))
					}
					// should now send to the updater
					utils.Logging("CLIENT_HANDLER",
						fmt.Sprintf("(%d) Sending info to updater", id))
					updater_chan<-*client_event
					utils.Logging("CLIENT_HANDLER",
						fmt.Sprintf("(%d) Info sent to updater", id))
				}
				break
			default:
				break
		}
	}

	// Exit
	main_chan<-"FINISHED"
	utils.Logging("CLIENT_HANDLER", fmt.Sprintf("%d quits.", id))
	return
}

func listenClient(id int, conn net.Conn, channel chan string) {
	// TODO utiliser la fonction register, pour assurer la mort du listenClient
	reader := bufio.NewReader(conn)
	for {
		utils.Logging("Listener", fmt.Sprintf("(%d) listening to client", id))
		netData, err := reader.ReadString('\n')
		netData = strings.TrimSpace(string(netData))
		utils.Logging("Listener",
			fmt.Sprintf("(%d) received from client: %s", id, netData))
		if err == nil {
			channel <- netData
		} else {
			channel <- "QUIT"
			return
		}
	}
}

/*           +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+
             | Auxiliary function:                   |
             | Loads the map, send it to the client. |
             +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+ */
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


/*           +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+
             | Auxiliary function:                                          |
             | Wait for the main function to say that the other players are |
             | here.                                                        |
             +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+ */
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

