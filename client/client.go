package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"bufio"
	"os"
	"time"
)


func logging(src string, message string) {
	fmt.Println(time.Now().String() + "[" + src + "] " + message)
}

// Two global maps store the queries that are to be treated (either by us, or
// by the server)
var (
	server_queries map[string]string
	client_queries map[string]string
	client_id string
	host string
	port int
	gui_chan chan string
)

func handle_server(c net.Conn, channel chan string) {
	logging("server handler", "The server handler has started.")
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if (err != nil) {
			logging("CLIENT", "error while reading from server")
			continue
		}
		channel <- strings.TrimSpace(string(netData))
	}
}

func handle_local(channel chan string) {
	logging("stdin handler", "The stdin handler has started.")
	var reader = bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		channel <- message
	}
}

func startClient(gui_chan_ptr *chan string) {
	/* Useful variables:
		running
			controls the main loop
		clientID
			random ID for the client
	*/
	running := true
	client_id = random_id(10)
	server_queries = make(map[string]string)
	client_queries = make(map[string]string)
	chan_stdin := make(chan string)
	chan_server := make(chan string)
	*gui_chan_ptr = make(chan string)
	host = "127.0.0.1"
	port = 10000

	// Verbose
	logging("CLIENT", "The client id is " + client_id)

	conn, err := net.Dial("tcp", host + ":" + strconv.Itoa(port))
	if err != nil {
		logging("CLIENT", fmt.Sprintf("Error durig TCP dial: %v", err))
		return
	}
	logging("CLIENT",
		fmt.Sprintf("Connection established with %s:%d", host, port))

	go handle_server(conn, chan_server)
	go handle_local(chan_stdin)

	logging("CLIENT", "Main loop is starting.")
	for running {
		select {
			case s1 := <-chan_stdin :
				// Recieving s1 from the terminal
				switch strings.Trim(s1, "\n") {
					case "!!QUIT":
						running = false
						gui_chan <- "QUIT"
					case "!!STATUS":
						fmt.Println("[CLIENT]: Status: ...")
					default:
						logging("CLIENT",
							fmt.Sprintf("Command %s unknown, querying info", s1))
						query_to_server(&conn, "info", "")
				}
			case s2 := <-chan_server:
				// Recieving s2 from the server
				from_server(&conn, s2)
		}
	}

	conn.Close()
}

