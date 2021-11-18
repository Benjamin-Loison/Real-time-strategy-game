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


var serverID string


func logging(src string, message string) {
	fmt.Println(time.Now().Format(time.ANSIC) + "[" + src + "] " + message)
}

var (
	// Controls the main event loop
	running bool

	// server_queries (resp. client_queries) map[string] string store the active queries,
	// i.e. queries which have not been set valid by the return of an "ok"
	// status (resp. while no correcr answer has been sent).
	server_queries map[string]string
	client_queries map[string]string

	// The random client ID
	client_id string

	// The ip address and port of the server
	host string
	port int

	// The channel that interacts with the gui part of the code (not required at
	// the moment)
	gui_chan chan string
)

// This function reads from the server and sends back to the main function the
// recieved messages
func handle_server(c net.Conn, channel chan string) {
	logging("server handler", "Starting routing traffic.")
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if (err != nil) {
			logging("Socket", fmt.Sprintf("error while reading from server: %v", err))
			time.Sleep(1 * time.Second)
			continue
		}
		channel <- strings.TrimSpace(string(netData))
	}
	logging("server handler", "Stopping routing traffic.")
}

func handle_local(channel chan string) {
	logging("stdin", "Starting routing traffic.")
	var reader = bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		logging("stdin", "⇐ " + message)
		channel <- message
	}
	logging("stdin", "Stopping routing traffic.")
}

func run_client(gui_chan_ptr *chan string, config Configuration, running *bool) {
	/* Useful variables:
		running
		*	controls the main loop
		clientID
			random ID for the client
	*/
	client_id = random_id(10)
	server_queries = make(map[string]string)
	client_queries = make(map[string]string)
	chan_stdin := make(chan string)
	chan_server := make(chan string)
	*gui_chan_ptr = make(chan string)

	// Verbose
	logging("CLIENT", "The client id is " + client_id)

	conn, err := net.Dial("tcp", config.hostname + ":" + strconv.Itoa(config.port))
	if err != nil {
		logging("Connection", fmt.Sprintf("Error durig TCP dial: %v", err))
		logging("Connection", fmt.Sprintf("\tHostname: %s, port: %d", config.hostname, config.port))
		return
	}
	defer conn.Close()
	logging("CLIENT",
		fmt.Sprintf("Connection established with %s:%d", host, port))

	// Listen the standart input and the server in co-processes
	go handle_server(conn, chan_server)
	go handle_local(chan_stdin)

	// Wait for the co-processes to dump their initial messages
	time.Sleep(5 * time.Second)

	// Initiate the game: ask the server for its identification number
	query_to_server(&conn, "info", "")// Warning: static parameters.
	time.Sleep(time.Second)
	query_to_server(&conn, "map", "0,0,100,100")// Warning: static parameters.
	//query_to_server(&conn, "location", "0,0,100,100")// Warning: static parameters.

	logging("CLIENT", "Main loop is starting.")
	for *running {
		select {
			case s1 := <-chan_stdin :
				// Recieving s1 from the terminal
				switch strings.Trim(s1, "\n") {
					case "!!QUIT":
						gui_chan <- "QUIT"
						*running = false
						return
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

}

