package main

import(
	"net"
	"encoding/json"
    "fmt"
	)


func client_handler(conn net.Conn, map_path string) {
	// Close the connection if the handler is exited
	defer conn.Close()

	// Load the map and send it to the client
	init_json := ServerMessage { MapInfo, LoadMap(map_path), nil }
    init_marshall, err := json.MarshalIndent(init_json,"","   ")
	Check(err)
	init_message := []byte(init_marshall)

	conn.Write(init_message)

	// Exit
	return
}

func main() {
	// Load the config file
	conf := loadConfig("conf/conf.json")

	// Listen
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",conf.Hostname, conf.Port))
	Check(err)
	defer listener.Close()

	for {
		// Wait for a new connection
		conn, err := listener.Accept()
		Check(err)
		go client_handler(conn, conf.MapPath)
	}
}

