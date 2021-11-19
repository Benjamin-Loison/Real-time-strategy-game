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
)

// Main fucntion for the part of the client that chats both with the server and
// the gui part.
func run_client(config Configuration_t, players *[]map[string]Unit, gmap *Map, chan_client chan string) {
	// Verbose
	logging("CLIENT", "The client id is " + config.Pseudo)

    chan_server := make(chan string, 2)

	conn, err := net.Dial("tcp", config.Server.Hostname + ":" + strconv.Itoa(config.Server.Port))
	if err != nil {
		logging("Connection", fmt.Sprintf("Error durig TCP dial: %v", err))
		logging("Connection", fmt.Sprintf("\tHostname: %s, port: %d", config.Server.Hostname, config.Server.Port))
        chan_client<-"QUIT"
		return
	}
	defer conn.Close()
	logging("CLIENT",
		fmt.Sprintf("Connection established with %s:%d", config.Server.Hostname, config.Server.Port))

    //GET MAP

    netData, err := bufio.NewReader(conn).ReadString('\n')
	if (err != nil) {
        logging("Socket", fmt.Sprintf("error while reading MAP INFO from server: %v", err))
        chan_client<-"QUIT"
        return
    }

    map_info_data := strings.TrimSpace(string(netData))

    var map_info = &ServerMessage{}
    err = json.Unmarshal([]byte(string(map_info_data)), map_info)
    Check(err)
    //fmt.Printf(" wtf %b", gmap.mut == nil)
    *gmap = map_info.GameMap


    //GET UNITS
    netData, err = bufio.NewReader(conn).ReadString('\n')
	if (err != nil) {
        logging("Socket", fmt.Sprintf("error while reading INIT UNITS from server: %v", err))
        chan_client<-"QUIT"
        return
    }

    players_data := strings.TrimSpace(string(netData))

    var players_info = &ServerMessage{}
    err = json.Unmarshal([]byte(string(players_data)), players_info)
    Check(err)
    //fmt.Printf(" wtf %b", gmap.mut == nil)
    *gmap = map_info.GameMap


    chan_client<-"OK"

	for {
		select {
        case s1 := <-chan_client:
            if s1 == "QUIT" {
                chan_server<-"QUIT"
                os.Exit(0)
            }
        case s2 := <-chan_server:
            fmt.Print(s2)
            //err = json.Unmarshal([]byte(string(s2)), gmap)
            //Check(err)
        default:
		}
	}
}
