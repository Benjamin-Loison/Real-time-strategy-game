package main

import (
	"bufio"
	//"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"encoding/json"

    "rts/utils"
)

var (
    serv_conn net.Conn
)

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

// Main fucntion for the part of the client that chats both with the server and
// the gui part.
func run_client(config Configuration_t, players *[]utils.Player, gmap *utils.Map, chan_client chan string) {
	// Verbose
	utils.Logging("CLIENT", "The client id is " + config.Pseudo)

	chan_server := make(chan string, 2)

    var err error

	serv_conn, err = net.Dial("tcp", config.Server.Hostname + ":" + strconv.Itoa(config.Server.Port))
	if err != nil {
		utils.Logging("Connection", fmt.Sprintf("Error durig TCP dial: %v", err))
		utils.Logging("Connection", fmt.Sprintf("\tHostname: %s, port: %d", config.Server.Hostname, config.Server.Port))
		chan_client<-"QUIT"
		return
	}
	defer serv_conn.Close()
	utils.Logging("CLIENT",
		fmt.Sprintf("Connection established with %s:%d", config.Server.Hostname, config.Server.Port))

	//GET MAP
	buffer := bufio.NewReader(serv_conn)

	utils.Logging("client", "requesting map")
	netData, err := buffer.ReadString('\n')
	utils.Logging("client", "obtained map")
	if (err != nil) {
		utils.Logging("Socket", fmt.Sprintf("error while reading MAP INFO from server: %v", err))
		chan_client<-"QUIT"
		return
	}

	map_info_data := strings.TrimSpace(string(netData))

	var map_info = &utils.ServerMessage{}
	err = json.Unmarshal([]byte(string(map_info_data)), map_info)
	utils.Check(err)
	//fmt.Printf(" wtf %b", gmap.mut == nil)
	*gmap = map_info.GameMap
	client_id = map_info.Id


	//GET UNITS
	utils.Logging("client", "requesting units")
	netData, err = buffer.ReadString('\n')
	utils.Logging("client", "obtained units")
	if (err != nil) {
		utils.Logging("Socket", fmt.Sprintf("error while reading INIT UNITS from server: %v", err))
		chan_client<-"QUIT"
		return
	}

	players_data := strings.TrimSpace(string(netData))

	//print(players_data)

	var players_info = &utils.ServerMessage{}
	err = json.Unmarshal([]byte(string(players_data)), players_info)
	utils.Check(err)

	*players = players_info.Players

    // now we wait for the go
	utils.Logging("client", "waiting for go")
	netData, err = buffer.ReadString('\n')
	utils.Logging("client", "go received")
	if (err != nil) {
		utils.Logging("Socket", fmt.Sprintf("error while reading MAP INFO from server: %v", err))
		chan_client<-"QUIT"
		return
	}

	chan_client<-"OK"

	go listenServer(serv_conn, chan_server)
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
			}
			//err = json.Unmarshal([]byte(string(s2)), gmap)
			//Check(err)
		default:
		}
	}
}
