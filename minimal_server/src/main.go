package main

import(
	"net"
	)


func client_handler(conn net.Conn) {

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
		go client_handler(conn)
	}
}

