package main

import(
	"net"
	"encoding/json"
    "fmt"
)

var (
    Players []Player
    IdSeeds []int
    main_chan = make(chan string,2)
    gmap Map
)

func getId(seed *int) int {
    ret := *seed
    (*seed) ++
    return ret
}

func client_handler(conn net.Conn, map_path string, main_chan chan string, id int) {
	// Close the connection if the handler is exited
	defer conn.Close()

	// Load the map and send it to the client
	init_json := ServerMessage { MapInfo, gmap, nil , id}
    init_marshall, err := json.Marshal(init_json)
	Check(err)
	init_message := []byte(string(init_marshall)+"\n")

	conn.Write(init_message)

	units_json := ServerMessage { MapInfo,Map{}, Players, id}
    units_marshall, err := json.Marshal(units_json)
	Check(err)
	units_message := []byte(string(units_marshall)+"\n")

	conn.Write(units_message)
	// Exit
    main_chan<-"FINISHED"
	return
}


func main() {
	// Load the config file
	conf := loadConfig("conf/conf.json")
    gmap = LoadMap(conf.MapPath)

    // initializing
    Players = append(Players, Player{} )
    IdSeeds = append(IdSeeds, 0)
    Players = append(Players, Player{} )
    IdSeeds = append(IdSeeds, 0)

    initializePlayer(&gmap, Player1,&Players[0].Units,&IdSeeds[0])
    initializePlayer(&gmap, Player2,&Players[1].Units,&IdSeeds[1])

	// Listen
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",conf.Hostname, conf.Port))
	Check(err)
	defer listener.Close()

    ok := 0

	for {
		// Wait for a new connection
        if ok == 2 {
            break
        }
		conn, err := listener.Accept()

        Check(err)
		go client_handler(conn, conf.MapPath,main_chan, ok)
		ok += 1
	}
    //start game
    for {
        if ok == 0 {
            break
        }
        select {
        case s := <-main_chan :
            if s == "FINISHED" {
                ok--
            }
        default:
        }
    }
}

