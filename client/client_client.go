package main
import (
	"fmt"
	"net"
)
func to_server(conn *net.Conn, query_type string, query_str string) {
	// Generate a query id:
	var id = ""
	var ok = false
	for (!ok) {
		id = random_id(10)
		_, ok = client_queries[id]
	}

	// Add the query to our pool
	client_queries[id] = query_type + ":" + query_str

	fmt.Println("Sending: ``Q" + id + "." + query_type + ":" + query_str + "''")
	// Send the query
	_, _ = (*conn).Write([]byte("A" + id + "." + query_type + ":" + query_str))
}

