package main
import (
	"fmt"
	"net"
)

func query_to_server(conn *net.Conn, query_type string, query_str string) {
	// Generate a query id:
	var id = ""
	id = random_id(10)

	// Add the query to our pool
	client_queries[id] = query_type + ":" + query_str

	query := fmt.Sprintf("Q%s.%s:%s", id, query_type, query_str)
	logging("to_server", fmt.Sprintf("Sending query %s", query))
	// Send the query
	_, _ = (*conn).Write([]byte(query))
}

