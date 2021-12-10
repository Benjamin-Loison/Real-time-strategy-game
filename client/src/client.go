package main

import (
	"bufio"
	"strconv"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	"encoding/json"
	"rts/utils"
    "rts/events"
)


/*           +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+
             | Main function:                                                |
             | --------------                                                |
             |                                                               |
             | This function deals with the communication between the client |
             | and the server (to which it has to connect).                  |
             |                                                               |
             | An auxiliary function `listenServer` listens to the connection|
             | and transmits to this function the incomming messages using a |
             | channel `chan_from_server`.                                        |
             | A channel `chan_link_gui` is used in order to communicate with|
             | the gui part of the client.                                   |
             +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+*/
func run_client(config Configuration_t,
			players *[]utils.Player,
			gmap *utils.Map,
			chan_link_gui chan string,
			chan_gui_link chan string) {
	// Verbose
	utils.Logging("CLIENT", "The client id is " + config.Pseudo)

	chan_from_server := make(chan string, 2)

    var err error// defines a single variable to check the functions errors.
	serv_conn := connectToServer(config.Server.Hostname, config.Server.Port)

	defer serv_conn.Close()

	// Fetch map, units and players information
	GetMapInfo(serv_conn, gmap, players)

	// Launch a goroutine that listens to the server
	go listenServer(serv_conn, chan_from_server)

    // Wait for the server to launch the game
	utils.Logging("client", "waiting for go")
	go_rec := false
	for !go_rec {
		select {
			case s := <- chan_from_server:
				if s == "GO" {
					go_rec = true
				}
		}
	}
	utils.Logging("client", "go received")

	if (err != nil) {
		utils.Logging("Socket",
			fmt.Sprintf("error while reading GO from server: %v", err))
		chan_link_gui<-"QUIT"
		return
	}

	// Main loop
	writer := bufio.NewWriter(serv_conn)
	keepGoing := true
	go func() {
		for keepGoing {
			select {
			case s1 := <-chan_gui_link:
				if s1 == "QUIT" {
					chan_from_server<-"QUIT"
					keepGoing = false
					os.Exit(0)
				} else {
					_,err := writer.Write([]byte(s1 + "\n"))
					writer.Flush()
					utils.Check(err)
				}
			}
		}
	} ()
	go func() {
		for keepGoing {
			select {
				case s2 := <-chan_from_server:
					if s2 == "QUIT" {
						chan_link_gui<-"QUIT"
						os.Exit(0)
					} else if strings.HasPrefix(s2, "CHAT:") {
						chan_link_gui<- s2
						utils.Logging("client", fmt.Sprintf("I recieved '%s'", s2))
					}   else {
                        var event = &events.Event{}
                        err := json.Unmarshal([]byte(s2), event)
                        utils.Check(err)
                        switch event.EventType {
                        case events.ServerUpdate:
                            var update = &events.ServerUpdate_e{}
                            err := json.Unmarshal([]byte(event.Data), update)
                            utils.Check(err)
                            playersRWLock.Lock()
                            for _, u := range update.Units {
                                key := strconv.Itoa(int(u.Id))
                                _,ok := (*players)[0].Units[key]
                                if ok {
                                    (*players)[0].Units[key] = u
                                    continue
                                }
                                _,ok = (*players)[1].Units[key]
                                if ok {
                                    (*players)[1].Units[key] = u
                                }
                            }


                            playersRWLock.Unlock()

						case events.BuildEvent :
							var build = &events.BuildBuilding_e{}
							err := json.Unmarshal([]byte(event.Data), build)
							utils.Check(err)
							gmap.Grid[build.Position_x][build.Position_y].Tile_Type = utils.Rock
						

                        default:
                            utils.Logging("client", "Unknown event\n")
                        }
                    }
				}
			}
	}()
	for keepGoing {
		time.Sleep(60 * time.Second)
	}
	os.Exit(0)
}



/*           +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+
             | Auxiliary function: listens to the server and dump the       |
             | incomming packets as string into the channel.                |
             | The data is expected to end with '\n', character that is not |
             | forwarded into the chan.                                     |
             +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+ */
func listenServer(conn net.Conn, channel chan string) {
	reader := bufio.NewReader(conn)
	for {
		netData, err := reader.ReadString('\n')
		utils.Check(err)
		netData = strings.TrimSpace(string(netData))
		channel <- netData
	}
}



/*          +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+
            | Auxiliary function that initiates a connection with the server |
            +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+*/
func connectToServer(hostname string, port int) (net.Conn) {
	serv_conn, err := net.Dial("tcp", hostname + ":" + strconv.Itoa(port))
	utils.Check(err)
	utils.Logging("CLIENT",
		fmt.Sprintf("Connection established with %s:%d",
			hostname, port))
	return serv_conn
}

// The function ask the server for the map data, and stores the results in
// the given pointers
func GetMapInfo(serv_conn net.Conn, gmap *utils.Map, players *[]utils.Player) () {
	// Fetch the map from the server
	buffer := bufio.NewReader(serv_conn)

	utils.Logging("client", "requesting map")
	netData, err := buffer.ReadString('\n')
	utils.Logging("client", "obtained map")
	utils.Check(err)

	map_info_data := strings.TrimSpace(string(netData))

	var map_info = &utils.ServerMessage{}
	err = json.Unmarshal([]byte(string(map_info_data)), map_info)
	utils.Check(err)
	*gmap = map_info.GameMap
	client_id = map_info.Id


	// Fetch the unit from the server
	utils.Logging("client", "requesting units")
	netData, err = buffer.ReadString('\n')
	utils.Logging("client", "obtained units")
	utils.Check(err)

	players_data := strings.TrimSpace(string(netData))

	var players_info = &utils.ServerMessage{}
	err = json.Unmarshal([]byte(string(players_data)), players_info)
	utils.Check(err)

	*players = players_info.Players
}

