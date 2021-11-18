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

// Main fucntion for the part of the client that chats both with the server and
// the gui part.
func run_client(gui_chan_ptr *chan string, map_chan *chan ServerMap, config Configuration) {
	// Initialisation of useful variables (descsribed at their declaration in
	// main.go)
	client_id = random_id(10)
	server_queries = make(map[string]string)
	client_queries = make(map[string]string)

	// Channel to speak with the gui part of the app
	*gui_chan_ptr = make(chan string, 2)
	*map_chan = make(chan ServerMap, 2)

	// Channels to speak with the server and the terminal
	chan_stdin := make(chan string, 2)
	chan_server := make(chan string, 2)

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
	time.Sleep(time.Second)

	// Initiate the game: ask the server for its identification number
	query_to_server(&conn, "info", "")// Warning: static parameters.
	query_to_server(&conn, "map", "0,0,100,100")// Warning: static parameters.
	//query_to_server(&conn, "location", "0,0,100,100")// Warning: static parameters.


	// TODO
	testes := ServerMap { 0, 0, 1, 1, []Location{Location{Floor, 0, "a"}} }
	*map_chan <- testes
	// TODO

	logging("CLIENT", "Main loop is starting.")
	for {
		select {
			case s1 := <-chan_stdin :
				// Recieving s1 from the terminal
				switch strings.Trim(s1, "\n") {
					case "QUIT":
						chan_server <- "QUIT"
						chan_stdin <- "QUIT"
						*gui_chan_ptr <- "QUIT"
						logging("CLIENT", "Ciao!")
						return
					case "!!STATUS":
						fmt.Println("[CLIENT]: Status: ...")
					default:
						logging("CLIENT",
							fmt.Sprintf("Command %s unknown. Ignoring", s1), 31)
				}
			case s2 := <-chan_server:
				// Recieving s2 from the server
				from_server(&conn, s2)
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
	logging("server handler", "Stopping routing traffic.")
}

// This function reads from the standart input (terminal) and sends back to the
// main function the recieved messages
func handle_local(channel chan string) {
	logging("stdin", "Starting routing traffic.")
	var reader = bufio.NewReader(os.Stdin)
	for {
		select {
			case x, _ := <-channel :
				if x == "QUIT" {
					return
				}
			default :
				// nop
		}
		message, _ := reader.ReadString('\n')
		channel <- message
	}
	logging("stdin", "Stopping routing traffic.")
}

// Logging fucntion used instead of the "log" package
func logging(src string, message string, color ...int) {
	if len(color) > 0 {
		//fmt.Printf("\033[%dm", color)
	}
	fmt.Println(time.Now().Format(time.ANSIC) + "[" + src + "] " + message)
	if len(color) > 0 {
		//fmt.Printf("\033[0m")
	}
}

