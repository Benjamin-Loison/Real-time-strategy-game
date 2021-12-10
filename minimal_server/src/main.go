package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"rts/events"
	"rts/factory"
	"rts/utils"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	serverSpeed = 10
    ffstep = utils.TileSize
)

var (
	Players []utils.Player
    PlayersRWLock = sync.RWMutex{}
	channels = make(map[int]chan string)
	updater_chan = make(chan events.Event)
	gmap utils.Map
	serverTime uint32
    genSeed = 0
)

func broadcast(channels map[int]chan string, msg string) {
	utils.Logging("broadcast", fmt.Sprintf("(0) (%s)", msg))
	channels[0] <- msg // Sending to player 0
	utils.Logging("broadcast", fmt.Sprintf("(1) (%s)", msg))
	channels[1] <- msg // Sending to player 1
	utils.Logging("broadcast", fmt.Sprintf("done(%s).", msg))
}

func register(id int) chan string {
	channels[id] = make(chan string)
	return channels[id]
}

func updater(channels map[int]chan string, stopper_chan chan string){
	for{
		select{
		case e := <-updater_chan:
			switch e.EventType {
				case events.ChatEvent:
					utils.Logging("Updater (CHAT)", fmt.Sprintf("%s", e.Data))
					broadcast(channels, fmt.Sprintf("CHAT:%s", e.Data))
					break
                case events.MoveUnits:
                    event := &events.MoveUnits_e{}
                    err := json.Unmarshal([]byte(e.Data),event)
                    utils.Check(err)
                    PlayersRWLock.Lock()
                    flowField := utils.PathFinding(gmap,event.Dest, ffstep)
                    if flowField != nil {
                        for _ , us := range event.Units {
                            for _ , p := range Players {
                                u, ok := p.Units[us]
                                if ok {
                                    u.FlowTarget = event.Dest
                                    u.FlowStep = ffstep
                                    u.FlowField = &flowField
                                }
                            }
                        }
                    }
                    PlayersRWLock.Unlock()

				default:
					utils.Logging("Updater", "Received an unhandeled event")
					break
			}
		case s := <-stopper_chan:
			if s == "QUIT" {
				utils.Logging("Updater", "Quitting")
				os.Exit(0)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	// Load the config file
	conf := loadConfig("conf/conf.json")
	gmap = utils.LoadMap(conf.MapPath)

	// initializing
	Players = append(Players, utils.Player{Units: map[string]factory.Unit{},Seed: 0} )
	Players = append(Players, utils.Player{Units: map[string]factory.Unit{},Seed: 0} )

	utils.InitializePlayer(&gmap, factory.Player1,&Players[0].Units, &genSeed)
	utils.InitializePlayer(&gmap, factory.Player2,&Players[1].Units, &genSeed)

	// Listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",conf.Hostname, conf.Port))
	utils.Check(err)

	nb_clients := 0

	for nb_clients < 2 {
		conn, err := listener.Accept()

		utils.Check(err)
		go client_handler(conn, conf.MapPath, register(nb_clients), nb_clients)
		nb_clients += 1
	}
	// Close the listener as soon as the two clients are connected
	listener.Close()

	// Launch the updater function and start the game
	go updater(channels, register(-1))
	utils.Logging("Server", "Send start!")
	broadcast(channels, "START")
	go gameLoop(register(-2))

	// wait for end game
	for nb_clients > 0{
		for _, c := range channels {
			select {
			case s := <-c :
				switch(s) {
				case "FINISHED" :
					nb_clients--
					utils.Logging("Server", fmt.Sprintf("%d remaining clients.", nb_clients))
					break
				case "STOP" :
					nb_clients = 0
					break
				case "CLIENT_ERROR" :
					utils.Logging("Server", "Received client error")
					break
				default:
					break
				}
			default:
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	broadcast(channels, "QUIT")
	os.Exit(0)
}

func gameLoop(quit chan string){
	gameOver := false
	for !gameOver {
		select {
		case s :=<-quit :
			if s=="QUIT"{
				return
			}
			break
		default: // not quitting, updating game state
			serverTime += 1
			//utils.Logging("GameLoop",fmt.Sprintf("Server time = %d",serverTime))
            PlayersRWLock.Lock()
            var toBeUpdated []string
            for _,p := range Players {
                for k, u := range p.Units {
                    if u.FlowField !=nil {
                        toBeUpdated = append(toBeUpdated,k)
                        x := u.X / ffstep
                        y := u.Y / ffstep
                        dir := rl.Vector2Normalize((*u.FlowField)[x][y])
                        u.X += int32(dir.X*u.Speed*float32(utils.TileSize))
                        u.Y += int32(dir.X*u.Speed*float32(utils.TileSize))
                        newCoord := rl.Vector2{X: float32(u.X) , Y : float32(u.Y)}
                        if rl.Vector2Distance(newCoord, u.FlowTarget) <= 1.5*float32(utils.TileSize) {
                            u.FlowField = nil
                        }
                    }
                }
            }
            var updatedUnits []factory.Unit
            for _, k := range toBeUpdated {
                val, ok := Players[0].Units[k]
                if ok {
                    updatedUnits = append(updatedUnits, val)
                    continue
                }
                val, ok = Players[1].Units[k]
                if ok {
                    updatedUnits = append(updatedUnits, val)
                    continue
                }
            }
            updateEvent := events.ServerUpdate_e{Units: updatedUnits}
            dataupdate, err := json.Marshal(updateEvent)
            utils.Check(err)
            event := events.Event{EventType: events.ServerUpdate, Data: string(dataupdate)+"\n"}
            dataevent, errr := json.Marshal(event)
            utils.Check(errr)

            PlayersRWLock.Unlock()

			broadcast(channels, string(dataevent))

            break
		}
		time.Sleep(serverSpeed * time.Millisecond)
	}
	quit<-"STOP"
}
