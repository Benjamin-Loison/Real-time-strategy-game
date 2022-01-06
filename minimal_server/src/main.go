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
	buildingsToBuild = make([]Building, 0)
	buildingsToBuild_txt = make([]string, 0)
	technologicalTree = LoadTechnologicalTree()
	passedBuilts = make([]Building, 0)
)

// broadcast send msg to all the channels to which clients listen
func broadcast(channels map[int]chan string, msg string) {
	channels[0] <- msg // Sending to player 0
	channels[1] <- msg // Sending to player 1
	//utils.Logging("broadcast", fmt.Sprintf("done(%s).", msg))
}

// Updates the channels variable to have a new channel identified by [id]
func register(id int) chan string {
	channels[id] = make(chan string)
	return channels[id]
}

// The updater function's goal is to handle the events that are
// sent on the updater_chan channel and to broadcast them
// to the client when needed
func updater(channels map[int]chan string, stopper_chan chan string){
	for{
		select{
		case e := <-updater_chan:
			switch e.EventType {
				case events.ChatEvent:
					broadcast(channels, fmt.Sprintf("CHAT:%s", e.Data))
					break
				case events.AttackUnit:
					
                case events.MoveUnits:
					utils.Logging("Updater (UPDATE)", fmt.Sprintf("%s", e.Data))
                    event := &events.MoveUnits_e{}
                    err := json.Unmarshal([]byte(e.Data),event)
                    utils.Check(err)
                    PlayersRWLock.Lock()
                    flowField := &[][]rl.Vector2{}
                    *flowField = utils.PathFinding(gmap,event.Dest, ffstep)
                    fmt.Println(event)
                    if *flowField != nil && len(*flowField) > 0 {
                        for _ , us := range event.Units {
                            for _ , p := range Players {
                                u, ok := p.Units[us]
                                if ok {
					                //utils.Logging("Updater", "TARGET FOUND &@~~ø~đ")
                                    u.FlowTarget = event.Dest
                                    u.FlowStep = ffstep
                                    u.FlowField = flowField
                                    p.Units[us] = u
                                }
                            }
                        }
                    }
                    //fmt.Println(flowField)
                    //fmt.Println(Players)
                    PlayersRWLock.Unlock()
				case events.BuildEvent:
					// Load the event data
					event := &events.BuildBuilding_e{}
					err := json.Unmarshal([]byte(e.Data),event)
					utils.Check(err)

					// Remarshal the event and broadcast
					txt, err := json.Marshal(e)
					utils.Check(err)
					buildingsToBuild_txt = append(buildingsToBuild_txt, string(txt))

					// Check that the build is authorized wrt the techTree
					// Add the building to the ToBuild slice
					buildingsToBuild = append(buildingsToBuild,
						Building{
							Name: event.BuildingName,
							Position_x: int32(event.Position_x),
							Position_y: int32(event.Position_y),
							BuildDuration: time.Duration(1000000000 * getDuration(event.BuildingName, technologicalTree)),
							BuildStartingTime: time.Now()})
					fmt.Println(buildingsToBuild)


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
	time.Sleep(time.Second)

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

// Game loop function where most the game logic is handled
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
            PlayersRWLock.Lock()
			// Moving units to their destination
            var toBeUpdated []string
            for _,p := range Players {
                for k, u := range p.Units {
					x := u.X / ffstep
					y := u.Y / ffstep
					pos := rl.Vector2{X: float32(u.X), Y: float32(u.Y)}
					dir := rl.Vector2Zero()
					if u.FlowField !=nil && len(*u.FlowField) > 0 {
						dir = (*u.FlowField)[x][y]
					}
					sepDir := rl.Vector2Zero()
					num_others := 0
					for k_other, u_other := range p.Units {
						pos_other := rl.Vector2{X: float32(u_other.X), Y: float32(u_other.Y)}
						dist := rl.Vector2Distance(pos, pos_other)
						if k_other != k && dist < 5.0*utils.Unit_size {
							num_others++
							normalized_vect := rl.Vector2Normalize(rl.Vector2Subtract(pos, pos_other))
							sepDir = rl.Vector2Add(sepDir, rl.Vector2Scale(normalized_vect, 1.0/dist))
						}
					}
					if num_others != 0 {
						sepDir = rl.Vector2Scale(sepDir, 100.0/float32(num_others))
					}
					dir = rl.Vector2Add(dir, sepDir)
					//u.X += int32(dir.X*u.Speed*float32(utils.TileSize))
					//u.Y += int32(dir.X*u.Speed*float32(utils.TileSize))
					new_X := u.X + int32(dir.X*2)
					new_Y := u.Y + int32(dir.Y*2)
					newCoord := rl.Vector2{X: float32(new_X) , Y : float32(new_Y)}
					
					canMove := true
					// Checking if there is a unit on the path
					for _, q := range Players {
						for _, v := range q.Units {
							pos_v := rl.Vector2{X: float32(v.X), Y: float32(v.Y)}
							if v != u && (rl.Vector2Distance(newCoord, pos_v) < 2.5*utils.Unit_size) {
								canMove = false
								// u will try to move in the same direction as v, if v is heading to the same target
								// TODO: this will break not going through walls
								/*if v.FlowTarget == u.FlowTarget && v.FlowField != nil && len(*v.FlowField)  > 0 &&
									((*v.FlowField)[x][y].X != 0 || (*v.FlowField)[x][y].Y != 0){
									delta_x := (*v.FlowField)[x][y].X - dir.X
									delta_y := (*v.FlowField)[x][y].Y - dir.Y
									(*u.FlowField)[x][y] = (*v.FlowField)[x][y]
									(*u.FlowField)[x][y].X += delta_x
									(*u.FlowField)[x][y].Y += delta_y

								}*/
								break
							}
						}
						if !(canMove) {
							break
						}
					}

					if !(canMove) {
						// TODO: change the way conflicts are handled

						continue
					}

					u.X = new_X
					u.Y = new_Y

					// Checking if target was reached
					if rl.Vector2Distance(newCoord, u.FlowTarget) <= 1.5*float32(utils.TileSize) {
						u.FlowField = nil
					}
					toBeUpdated = append(toBeUpdated,k)
					p.Units[k] = u
                }
            }
			// Sending update to client if necessary
			//utils.Logging("GameLoop",fmt.Sprintf("calculating update update %d", len(toBeUpdated)))
            //fmt.Println(Players)
            if len(toBeUpdated)>0 {
			    //utils.Logging("GameLoop","There is an update")
				var updatedUnits []factory.Unit
                for _, k := range toBeUpdated {
					// Handling the units of player 0
                    val, ok := Players[0].Units[k]
                    if ok {
                        updatedUnits = append(updatedUnits, val)
                        continue
                    }
					// Handling the units of player 1
                    val, ok = Players[1].Units[k]
                    if ok {
                        updatedUnits = append(updatedUnits, val)
                    }
                }
                updateEvent := events.ServerUpdate_e{Units: updatedUnits}
                dataupdate, err := json.Marshal(updateEvent)
                utils.Check(err)
                event := events.Event{EventType: events.ServerUpdate, Data: string(dataupdate)}
                dataevent, errr := json.Marshal(event)
                utils.Check(errr)
				//utils.Logging("gameLoop", fmt.Sprintf("Sending update %s", string(dataevent)))


                broadcast(channels, string(dataevent))
            }
            PlayersRWLock.Unlock()

			if len(buildingsToBuild) > 0 {
				new_toBuild := make([]Building, 0)
				new_toBuild_txt := make([]string, 0)
				for i, e := range buildingsToBuild {
					if !CheckRights(technologicalTree, passedBuilts, e) {
						broadcast(channels, "CHAT:SERVER: A building was not built.")
					} else {
						if time.Since(e.BuildStartingTime) > e.BuildDuration {
							// Build
							utils.Logging("Build", fmt.Sprintf("%d has passed", e.BuildDuration/1000000000))
							Build(e, technologicalTree)
							broadcast(channels, buildingsToBuild_txt[i])
							passedBuilts = append(passedBuilts, e)
						} else {
							new_toBuild = append(new_toBuild, e)
							new_toBuild_txt = append(new_toBuild_txt, buildingsToBuild_txt[i])
						}
					}
				}
				buildingsToBuild = new_toBuild
				buildingsToBuild_txt = new_toBuild_txt
			}
		}
		time.Sleep(serverSpeed * time.Millisecond)
	}
	quit<-"STOP"
}
