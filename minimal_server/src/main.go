package main

import(
	"net"
	"encoding/json"
	)


func client_handler(conn net.Conn, map_path string) {
	init_json := ServerMessage { MapInfo, LoadMap(map_path), nil }
    init_marshall, err := json.MarshalIndent(init_json,"","   ")
	Check(err)
	init_message := []byte(init_marshall)

	conn.Write(init_message)

	return
}

func main() {
	// Load the config file
	conf := loadConfig("conf/conf.json")
	
	// Listen
	listener, err := net.Listen("tcp", conf.Hostname + ":" + string(conf.Port))
	Check(err)
	defer listener.Close()

	for {
		// Wait for a new connection
		conn, err := listener.Accept()
		Check(err)
		go client_handler(conn, conf.MapPath)
	}
}

