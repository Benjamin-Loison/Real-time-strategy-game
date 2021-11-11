package main
import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	)

/*
	This file is used to manage all message that come from the server.

	There are at least three functions below to:
		- understand a query from the server
		- send an answer to the server
		- fetch the status of the server on a particular query
*/



// The following function sends an answer to the server
func answer_server(conn *net.Conn, query_id string, query_ans string) {
	_, _ = (*conn).Write([]byte("A" + query_id + "." + query_ans))
}

// The following decide what is to be done when the server sends a query to the client
func manage_server_query(conn *net.Conn, query_id string, query_type string, query_str string) {
	// Verbose
	fmt.Println("New query:\n\tid: " + query_id + "\n\t\ttype: " + query_type +
		"\n\t\tquery: " + query_str + "\n")

	// Add the query to the map
	server_queries[query_id] = query_type + ":" + query_str

	// Manage the query depending on its type:
	switch(query_type) {
	case "info":
		// Send back the client's ID
		answer_server(conn, query_id, client_id)
	case "map":
		fmt.Println("The server should not ask for the map.")
	case "location":
		switch (query_str) {
			case "get":
				x, y := get_player_location()
				answer_server(conn, query_id, strconv.Itoa(x) + "," + strconv.Itoa(y))
			default:
				if(strings.HasPrefix(query_str, "set,")) {
					splitted := strings.Split(query_str, ",")
					player_id := splitted[1]
					x_coord, _ := strconv.Atoi(splitted[2])
					y_coord, _ := strconv.Atoi(splitted[3])
					set_player_location(player_id, x_coord, y_coord)
				} else {
					fmt.Println("[CLIENT] Ill-formed query from server")
				}
		}
	default:
		fmt.Println("Unknown query type %s.\n", query_type)
		delete(server_queries, query_id)
	}
}

func location_type_from_str(s string) location_type {
	switch (s) {
		case "Floor":
			return Floor
		case "ElfBuilding":
			return ElfBuilding
		case "OrcBuilding":
			return OrcBuilding
		case "HumanBuilding":
			return HumanBuilding
		default:
			return Floor
	}
}

func manage_server_answer(answer_id string, answer_str string) {
	// Checks that the query exists
	var ok bool
	query_str, ok := client_queries[answer_id]
	if(!ok) {
		fmt.Println("New answer to unknown query!\n\tquery id: %s\n\tanswer:%s\n", answer_id, answer_str)
	} else {
		splitted := strings.Split(strings.Trim(query_str, "\n"), ":")
		switch(splitted[0]) {
			case "info":
				fmt.Println("[SERVER] ID: " + splitted[1])
			case "map":
				// Split the whole answer
				split_answer := strings.Split(answer_str, ",,")

				w := strings.Split(splitted[1], ",")[3]
				h := strings.Split(splitted[1], ",")[4]

				// Fetch the initiam location
				initial_position := strings.Split(split_answer[0], ",")
				init_x, _ := strconv.Atoi(initial_position[0])
				init_y, _ := strconv.Atoi(initial_position[1])

				// Fetch the descrition of the map
				description_list := split_answer[1:]
				location_list := make([]Location, 0)
				for _, elem := range description_list {
					elem_split := strings.Split(elem, ",")
					elem_type := location_type_from_str(elem_split[0])
					elem_int, _ := strconv.Atoi(elem_split[1])
					elem_str := elem_split[2]
					location_list := append(location_list,
						Location {elem_type, elem_int, elem_str})
				}
				
				// Update the map
				update_map(init_x, init_y, w, h, location_list)
			default:
				fmt.Println("Got an answer for a query of unknown type.")
		}
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

