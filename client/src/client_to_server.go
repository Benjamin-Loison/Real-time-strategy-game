package main
import (
	"fmt"
	"net"
)

// The following function sends an answer to the server
func answer_to_server(conn *net.Conn, query_id string, query_ans string) {
	str := "A" + query_id + "." + query_ans + "\n"
	logging("Sending the answer: ", str)
	_, _ = (*conn).Write([]byte(str))
}

func status_to_server(conn *net.Conn, query_id string, query_status string) {
	str := "S" + query_id + "." + query_status + "\n"
	logging("Sending the status: ", str)
	_, _ = (*conn).Write([]byte(str))
}


func query_to_server(conn *net.Conn, query_type string, query_str string) {
	// Generate a query id:
	var id = ""
	id = random_id(10)

	// Add the query to our pool
	client_queries[id] = query_type + ":" + query_str

	query := fmt.Sprintf("Q%s.%s:%s", id, query_type, query_str)
	logging("to_server", fmt.Sprintf("Sending query %s", query))
	// Send the query
	_, _ = (*conn).Write([]byte(query + "\n"))
}

