package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"encoding/json"

    "rts/utils"
)

/*           +~~~~~~~~~~~~~~~~~~+
             | Global variables |
             +~~~~~~~~~~~~~~~~~~+ */

var (
    serv_conn net.Conn
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
             | channel `chan_server`.                                        |
             | A channel `chan_client` is used in order to communicate with  |
             | the gui part of the client.                                   |
             +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+*/
func run_client(config Configuration_t, players *[]utils.Player, gmap *utils.Map, chan_client chan string) {
	// Verbose
	utils.Logging("CLIENT", "The client id is " + config.Pseudo)

	chan_server := make(chan string, 2)

    var err error// defines a single variable to check the functions errors.

	// Connection to the server
	serv_conn, err = net.Dial("tcp", config.Server.Hostname + ":" + strconv.Itoa(config.Server.Port))
	if err != nil {
		utils.Logging("Connection",
			fmt.Sprintf("Error durig TCP dial (%s:%d)): %v",
				config.Server.Hostname,
				config.Server.Port,
				err))
		chan_client<-"QUIT"
		return
	}
	defer serv_conn.Close()
	utils.Logging("CLIENT",
		fmt.Sprintf("Connection established with %s:%d",
			config.Server.Hostname, config.Server.Port))


	// Fetch the map from the server
	buffer := bufio.NewReader(serv_conn)

	utils.Logging("client", "requesting map")
	netData, err := buffer.ReadString('\n')
	utils.Logging("client", "obtained map")
	if (err != nil) {
		utils.Logging("Socket",
			fmt.Sprintf("error while reading MAP INFO from server: %v", err))
		chan_client<-"QUIT"
		return
	}

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
	if (err != nil) {
		utils.Logging("Socket", fmt.Sprintf("error while reading INIT UNITS from server: %v", err))
		chan_client<-"QUIT"
		return
	}

	players_data := strings.TrimSpace(string(netData))

	var players_info = &utils.ServerMessage{}
	err = json.Unmarshal([]byte(string(players_data)), players_info)
	utils.Check(err)

	*players = players_info.Players

    // Wait for the server to launch the game
	utils.Logging("client", "waiting for go")
	netData, err = buffer.ReadString('\n')
	utils.Logging("client", "go received")
	if (err != nil) {
		utils.Logging("Socket",
			fmt.Sprintf("error while reading GO from server: %v", err))
		chan_client<-"QUIT"
		return
	}

	chan_client<-"OK"

	// Launch a goroutine that listens to the server
	go listenServer(serv_conn, chan_server)

	// Main loop
	for {
		select {
		case s1 := <-chan_client:
			if s1 == "QUIT" {
				chan_server<-"QUIT"
				os.Exit(0)
			}
		case s2 := <-chan_server:
			if s2 == "QUIT" {
                chan_client<-"QUIT"
				os.Exit(0)
			} else if strings.HasPrefix(s2, "CHAT:") {
				chan_client<- s2
				utils.Logging("client", fmt.Sprintf("I recieved '%s'", s2))
			}
			//err = json.Unmarshal([]byte(string(s2)), gmap)
			//Check(err)
		default:
		}
	}
}



func listenServer(conn net.Conn, channel chan string) {
	reader := bufio.NewReader(conn)
	for {
		netData, err := reader.ReadString('\n')
		netData = strings.TrimSpace(string(netData))
		if err == nil {
			channel <- netData
		} else {
			channel <- "QUIT"
			return
		}
	}
}

