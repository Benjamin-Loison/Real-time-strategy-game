package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"bufio"
	"os"
)

// Two global maps store the queries that are to be treated (either by us, or
// by the server)
var server_queries map[string]string
var client_queries map[string]string
var client_id string

func handle_server(c net.Conn, channel chan string) {
	fmt.Println("[CLIENT] handling the conenction.")
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if (err != nil) {
			fmt.Println("[CLIENT] error while reading from server")
			continue
		}
		channel <- strings.TrimSpace(string(netData))
	}
}

func handle_local(channel chan string) {
	fmt.Println("[CLIENT] Handling stdin")
	var reader = bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		channel <- message
	}
}

func client_main(host string, port int, main_chan chan string) {
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

	// Verbose
	fmt.Println("[CLIENT] Client id: " + client_id)

	conn, err := net.Dial("tcp", host + ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Println("[CLIENT] Connection established with " + host + ":" + strconv.Itoa(port))

	go handle_server(conn, chan_server)
	go handle_local(chan_stdin)

	fmt.Println("[CLIENT] Starting the main client loop.")
	for running {
		select {
			case s1 := <-chan_stdin :
				// Recieving s1 from the terminal
				switch strings.Trim(s1, "\n") {
					case "!!QUIT":
						running = false
						main_chan <- "QUIT"
					case "!!STATUS":
						fmt.Println("[CLIENT]: Status: ...")
					default:
						fmt.Println("Sending " + s1 + "to the server.")
						to_server(&conn, "info", "myinformation")
				}
			case s2 := <-chan_server:
				// Recieving s2 from the server
				from_server(&conn, s2)
		}
	}

	conn.Close()
}

