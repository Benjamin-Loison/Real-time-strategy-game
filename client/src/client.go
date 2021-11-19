package main

import (
	"bufio"
	//"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Main fucntion for the part of the client that chats both with the server and
// the gui part.
func run_client(config Configuration_t, gmap *Map, chan_client chan string) {
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

	// Listen the standart input and the server in co-processes
	go handle_server(conn, chan_server)
	//go handle_local(chan_stdin)

	logging("CLIENT", "Main loop is starting.")
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

// This function reads from the server and sends back to the main function the
// recieved messages
func handle_server(c net.Conn, channel chan string) {
	time.Sleep(time.Second)
	logging("server handler", "Starting routing traffic.")
	for {
		// Exit
		select {
			case x, _ := <-channel :
				if x == "QUIT" {
	                logging("server handler", "Stopping routing traffic.")
					return
				}
			default :
				// nop
		}
		// Read
		netData, err := bufio.NewReader(c).ReadString('\n')
		if (err != nil) {
			logging("Socket", fmt.Sprintf("error while reading from server: %v", err))
			time.Sleep(1 * time.Second)
			continue
		}
        fmt.Print(strings.TrimSpace(string(netData)))
		channel <- strings.TrimSpace(string(netData))
	}
}
