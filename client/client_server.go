package main
import (
	"fmt"
	"net"
	"regexp"
	)

/*
	Thif file is used to manage all message that come from the server.

	There are at least three functions below to:
		- understand a query from the server
		- send an answer to the server
		- fetch the status of the server on a particular query
*/



// The following function sends an answer to the server
func answer_server(conn *net.Conn, query_id string, query_ans string) {
	_, _ = (*conn).Write([]byte("A" + query_id + "." + query_ans))
}

// The following decide what is to be done whenthe server sends us a query
func manage_server_query(conn *net.Conn, query_id string, query_type string, query_str string) {
	// Verbose
	fmt.Println("New query:\n\tid: " + query_id + "\n\t\ttype: " + query_type +
		"\n\t\tquery: " + query_str + "\n")

	// Add the query to the map
	server_queries[query_id] = query_type + ":" + query_str

	// Manage the query depending on its type:
	switch(query_type) {
	case "info":
		// Answers.
		answer_server(conn, query_id, client_id)
	default:
		fmt.Println("Unknown query type %s.\n", query_type)
		delete(server_queries, query_id)
	}
}

func manage_server_answer(answer_id string, answer_str string) {
	// Checks that the query exists
	var ok bool
	_, ok = client_queries[answer_id]
	if(!ok) {
		fmt.Println("New answer to unknown query!\n\tquery id: %s\n\tanswer:%s\n", answer_id, answer_str)
	} else {
		fmt.Println("Answer:\n\tquery id: %s\n\tanswer:%s\n", answer_id, answer_str)
	}
	delete(client_queries, answer_id)
}

func from_server(conn *net.Conn, str string) {
	switch(str[0]) {
	case 'Q':
		// Verbosity
		fmt.Println("Server sends new query: '%s'", str)
		// Parsing the query
		var rg = regexp.MustCompile(`Q[ 0-9a-zA-Z]+\.[a-zA-Z 0-9]+:[a-zA-Z0-9 ]*`)
		var idrg = regexp.MustCompile(`[0-9a-zA-Z ]+`)
		var typerg = regexp.MustCompile(`[0-9a-zA-Z ]+`)
		if (rg.MatchString(str)) {// the query is valid
			var q_id = (idrg.FindString(str))[1:]
			var q_type = typerg.FindString(str[len(q_id)+1:])
			var q_str = str[len(q_id) + len(q_type)+3:]
			// Managing the parsed query
			manage_server_query(conn, q_id, q_type, q_str)
		} else {
			fmt.Printf("The query is ill-formed!\n")
		}
	case 'A':
		// Verbosity
		fmt.Println("Server sends answer: '%s'", str)
		// Parsing the answer
		var rg = regexp.MustCompile(`A[ 0-9a-zA-Z]+\.[a-zA-Z0-9 ]*`)
		var idrg = regexp.MustCompile(`[0-9a-zA-Z ]+`)
		if (rg.MatchString(str)) {// the query is valid
			var a_id = (idrg.FindString(str))[1:]
			var a_str = str[len(a_id) + 2:]
			// Managing the parsed answer
			manage_server_answer(a_id, a_str)
		} else {
			fmt.Printf("The answer is ill-formed!\n")
		}
	case 'S':
		// Verbosity
		fmt.Println("Server sends status: '%s'", str)
		// Parsing the status
		var rg = regexp.MustCompile(`S[0-9a-zA-Z]+\.(ok|nok)`)
		var idrg = regexp.MustCompile(`[0-9a-zA-Z ]+`)
		if (rg.MatchString(str)) {// the query is valid
			var s_id = (idrg.FindString(str))[1:]
			var s_str = str[len(s_id) + 2:]
			// Deleting the query from our map
			delete(server_queries, s_id)
			if(s_str == "nok") {
				panic("The server is angry at us!!!")
			}
		} else {
			fmt.Printf("The status is ill-formed!\n")
		}
	default:
		fmt.Println("The server sent '%s':\n\tbad query.\n", str)
	}
}

