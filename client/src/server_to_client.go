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



// The following decide what is to be done when the server sends a query to the client
func manage_server_query(conn *net.Conn, query_id string, query_type string, query_str string) {
	// Verbose
	logging("Server → Client",
		fmt.Sprintf("Query <id: %s>:\n\ttype: %s\n\toptions: %s",
			query_id, query_type, query_str))

	// Add the query to the map
	server_queries[query_id] = query_type + ":" + query_str

	// Manage the query depending on its type:
	switch(query_type) {
	case "info":
		// Send back the client's ID
		answer_to_server(conn, query_id, client_id)
	case "map":
		logging("Server → Client", "Ignoring: the server should not ask for the map.")
	case "location":
		switch (query_str) {
			case "get":
				x, y := get_player_location()
				answer_to_server(conn, query_id, strconv.Itoa(x) + "," + strconv.Itoa(y))
			default:
				if(strings.HasPrefix(query_str, "set,")) {
					splitted := strings.Split(query_str, ",")
					player_id := splitted[1]
					x_coord, _ := strconv.Atoi(splitted[2])
					y_coord, _ := strconv.Atoi(splitted[3])
					set_player_location(player_id, x_coord, y_coord)
				} else {
					logging("Server → Client", "Ill-formed query.")
					delete(server_queries, query_id)
				}
		}
	default:
		logging("Server → Client", "Unknown query type.")
		delete(server_queries, query_id)
	}
}

func manage_server_answer(conn *net.Conn, answer_id string, answer_str string) {
	// Verbose
	logging("Server → Client",
		fmt.Sprintf("Recieving answer <id: %s>: answer= %s",
			answer_id, answer_str))
	// Checks that the query exists
	query_str, ok := client_queries[answer_id]
	if(!ok) {
		logging("Server → Client", "Unknown query id")
		status_to_server(conn, answer_id, "nok")
		return
	} else {
		splitted := strings.Split(strings.Trim(query_str, "\n"), ":")
		switch(splitted[0]) {
			case "info":
				serverID = answer_str
				fmt.Println("[SERVER] ID: " + serverID)
				status_to_server(conn, answer_id, "ok")
			case "map":
				//logging("Server → Client", "The query was " + query_str)
				logging("Server → Client", fmt.Sprintf("[map] <%s> %s", answer_id, answer_str))
				// Split the whole answer
				split_answer := strings.Split(answer_str, ",,")

				w, _ := strconv.Atoi(strings.Split(splitted[1], ",")[2])
				h, _ := strconv.Atoi(strings.Split(splitted[1], ",")[3])

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
					location_list = append(location_list,
						Location {elem_type, elem_int, elem_str})
				}
				
				// Update the map
				map_chan <- ServerMap{init_x, init_y, w, h, location_list}
				status_to_server(conn, answer_id, "ok")
			default:
				logging("Server → Client", "Unknown query type")
				status_to_server(conn, answer_id, "nok")
				return
		}
	}
	delete(client_queries, answer_id)
}

func from_server(conn *net.Conn, str string) {
	logging("Server → Client", "Recieving " + str)
	switch(str[0]) {
	case 'Q':
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
			logging("Server → Client", "Ignoring: Ill-formed query: " + str)
		}
	case 'A':
		// Parsing the answer
		var rg = regexp.MustCompile(`A[ 0-9a-zA-Z]+\.[a-zA-Z0-9 ]*`)
		var idrg = regexp.MustCompile(`[0-9a-zA-Z ]+`)
		if (rg.MatchString(str)) {// the query is valid
			var a_id = (idrg.FindString(str))[1:]
			var a_str = str[len(a_id) + 2:]
			// Managing the parsed answer
			manage_server_answer(conn, a_id, a_str)
		} else {
			logging("Server → Client", "Ignoring: Ill-formed answer: " + str)
		}
	case 'S':
		// Parsing the status
		var rg = regexp.MustCompile(`S[0-9a-zA-Z]+\.([o|O][k|K]|[n|N][o|O][k|K])`)
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
			logging("Server → Client", "Ignoring: Ill-formed status: " + str)
		}
	default:
		logging("Server → Client", "Ignoring: Ill-formed message: " + str)
	}
}

