package main

import (
    "net"
    "fmt"
    "strconv"
    "bufio"
    "time"
    "strings"
)

// Main fucntion for the part of the client that chats both with the server and
// the gui part.
func run_client(config Configuration_t) {
	// Verbose
	logging("CLIENT", "The client id is " + config.Pseudo)

    chan_server := make(chan string, 2)

	conn, err := net.Dial("tcp", config.Server.Hostname + ":" + strconv.Itoa(config.Server.Port))
	if err != nil {
		logging("Connection", fmt.Sprintf("Error durig TCP dial: %v", err))
		logging("Connection", fmt.Sprintf("\tHostname: %s, port: %d", config.Server.Hostname, config.Server.Port))
		return
	}
	defer conn.Close()
	logging("CLIENT",
		fmt.Sprintf("Connection established with %s:%d", config.Server.Hostname, config.Server.Port))

	// Listen the standart input and the server in co-processes
	go handle_server(conn, chan_server)
	//go handle_local(chan_stdin)

	// Wait for the co-processes to dump their initial messages
	time.Sleep(time.Second)

	// Initiate the game: ask the server for its identification number
	//query_to_server(&conn, "info", "")// Warning: static parameters.
	//query_to_server(&conn, "map", "0,0,100,100")// Warning: static parameters.
	//query_to_server(&conn, "location", "0,0,100,100")// Warning: static parameters.

	logging("CLIENT", "Main loop is starting.")
	for {
		select {
			case s2 := <-chan_server:
				// Recieving s2 from the server
                print(s2)
            default:
		}
	}

}

// This function reads from the server and sends back to the main function the
// recieved messages
func handle_server(c net.Conn, channel chan string) {
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
		channel <- strings.TrimSpace(string(netData))
	}
}
