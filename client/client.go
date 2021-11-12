package main

import (
	"fmt"
	"net"
)

// Two global maps store the queries that are to be treated (either by us, or
// by the server)
var server_queries map[string]string
var client_queries map[string]string
var client_id string


func main() {
	// Define a client Id:
	client_id = random_id(10)
	fmt.Println("Client id: " + client_id)

	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
    //fmt.Printf("coming1")
	keep_going := true

	in_stdin := make(chan string)
	in_server := make(chan string)

	server_queries = make(map[string]string)
	client_queries = make(map[string]string)


	for keep_going {
		select {
		case s1 := <-in_stdin :
			fmt.Println (s1)
			to_server(&conn, "info", "myinformation")
		case s2 := <-in_server:
			from_server(&conn, s2)
		}
	}

	conn.Close()
}

