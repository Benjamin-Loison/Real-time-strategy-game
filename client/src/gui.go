package main

import (
	"fmt"
    //"reflect"
	"os"
	_ "image/png"
	"math"
	"time"
	"encoding/json"
	"rts/events"
	"rts/utils"
	"github.com/gen2brain/raylib-go/raylib"
	"strings"
    "strconv"
)

const (
	screenWidth = 1280
	screenHeight = 720
	zoomFactor = 1.03
	cameraSpeed = 300.0
)

type GameState int32

const (
	StateNone GameState = 1
	StateChat GameState = 2
	StateMenu GameState = 4
	StateWaitClick GameState = 8

    ffstep = utils.TileSize

	helpSize = 13
)

var (
	currentState = StateNone
	currentMenu = -1
	currentAction = -1

	lastInputTime = time.Now()

	selectedUnits = map[string]bool{}
)



/*           +~~~~~~~~~~~~~~~~~~~~~~+
             | Main loop of the gui |
             +~~~~~~~~~~~~~~~~~~~~~~+ */
func RunGui(gmap *utils.Map,
			players *[]utils.Player,
			config Configuration_t, config_menus MenuConfiguration_t,
			chan_link_gui chan string, chan_gui_link chan string) {
	
	ChatText := ""
	has_send := 0
	currentMessages := make([]MessageItem_t, 0)
	lastMessagesUpdate := time.Now()
	need_pause := false

	var currentAction Action_t
    
    var flowField  [][]rl.Vector2 = nil

    var selectionStart rl.Vector2
    var inSelection = false

	rl.SetTraceLog(rl.LogNone)
	rl.InitWindow(screenWidth, screenHeight, "RTS")
	rl.InitAudioDevice()

	xm := rl.LoadMusicStream("music.xm")
	rl.PlayMusicStream(xm)
	pause := false
	defer rl.UnloadMusicStream(xm)
	defer rl.CloseAudioDevice()
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var map_width int32 = gmap.Width
	var map_height int32 = gmap.Height

	map_middle := rl.Vector2{X:float32(utils.TileSize)*float32(map_width)/2.0,Y: float32(utils.TileSize)*float32(map_height)/2.0 }


	camera := rl.NewCamera2D(rl.Vector2{X:screenWidth/2.0,Y:screenHeight/2.0},rl.Vector2{X: float32(utils.TileSize)*float32(map_width)/2.0, Y: float32(utils.TileSize)*float32(map_height)/2.0}, 0, 1.0)

	rl.BeginMode2D(camera)

	skipCharPressed := false

	for !rl.WindowShouldClose() {
		/*           +~~~~~~~~~~~~~~~~~~~~~~~~~~~+
		             | Check that shouldn't quit |
		             +~~~~~~~~~~~~~~~~~~~~~~~~~~~+ */
		select {
		case x, _ := <-chan_link_gui:
			if x == "QUIT" {
				utils.Logging("gui","Forced to quit")
				return
			} else if strings.HasPrefix(x, "CHAT:") {
				currentMessages = NewMessageItem(currentMessages, x[5:], config)
				organizeMessages(currentMessages)
				if has_send > 0 {
					has_send -= 1
				}
			} else {
				utils.Logging("gui", fmt.Sprintf("Je ne comprends pas %s", x))
			}
			break
		default:
			break
		}

		/*           +~~~~~~~~~~~~~~~~~+
		             | Audio |
		             +~~~~~~~~~~~~~~~~~+ */
		rl.UpdateMusicStream(xm) // Update music buffer with new stream data
		// Restart music playing (stop and play)
		if rl.IsKeyPressed(rl.KeySpace) && currentState & StateChat == 0 {
			rl.StopMusicStream(xm)
			rl.PlayMusicStream(xm)
		}

		// Pause/Resume music playing
		if rl.IsKeyPressed(rl.KeyP) && currentState & StateChat == 0 {
			pause = !pause

			if pause {
				rl.PauseMusicStream(xm)
			} else {
				rl.ResumeMusicStream(xm)
			}
		}

		/*           +~~~~~~~~~~~~~~~~~+
		             | Update the game |
		             +~~~~~~~~~~~~~~~~~+ */

		// Check wether the movements have to be detected
		if (currentState & StateNone > 0) {
			offsetThisFrame := cameraSpeed*rl.GetFrameTime()

			if (rl.IsKeyDown(config.Keys.Chat)) {
				utils.Logging("GUI", "Entering chat state.")
				// Prevent an extra T from appearing in the chat message
				rl.GetKeyPressed()
				skipCharPressed = true
				currentState = StateChat
			}
			if (rl.IsKeyDown(config.Keys.Right)){
				//camera.Offset.X -= 2.0
				camera.Target.X += offsetThisFrame
			}
			if (rl.IsKeyDown(config.Keys.Left)){
				//camera.Offset.X += 2.0
				camera.Target.X -= offsetThisFrame
			}
			if (rl.IsKeyDown(config.Keys.Up)){
				//camera.Offset.Y += 2.0
				camera.Target.Y -= offsetThisFrame
			}
			if (rl.IsKeyDown(config.Keys.Down)){
				//camera.Offset.Y -= 2.0
				camera.Target.Y += offsetThisFrame
			}
			if (rl.IsKeyDown(config.Keys.ZoomIn)){
				camera.Zoom *= zoomFactor
			}
			if (rl.IsKeyDown(config.Keys.ZoomOut)){
				camera.Zoom /= zoomFactor
			}
			if (rl.IsKeyDown(config.Keys.Menu)){
				if ((currentState & StateMenu > 0) && time.Since(lastInputTime) > time.Second) {
					utils.Logging("GUI", "Exiting menu mode")
					currentState -= StateMenu
					lastInputTime = time.Now()
				} else if (time.Since(lastInputTime) > time.Second) {
					utils.Logging("GUI", "Entering menu mode")
					currentState += StateMenu
					currentMenu = -1// undefined menu
					lastInputTime = time.Now()
				}
			}
			if (rl.IsKeyDown(config.Keys.ResetCamera)){
				camera.Zoom = 1.0
				camera.Target.X = map_middle.X
				camera.Target.Y = map_middle.Y
			}

			if (rl.IsMouseButtonPressed(rl.MouseRightButton)) && len(selectedUnits)> 0 {
                //fmt.Println(rl.GetScreenToWorld2D(rl.GetMousePosition(),camera))
                //flowField = utils.PathFinding(*gmap,rl.GetScreenToWorld2D(rl.GetMousePosition(),camera),ffstep)

				// We send the order to move units to the server
				dest := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
				units := []string{}

                fmt.Println("SELECTED UNITS ##########")
                fmt.Println(selectedUnits)
                fmt.Println(*players)

                // should maybe use this lock when modifying the enemy player when he is attacked
                playersRWLock.RLock()
			    for _, u := range (*players)[client_id].Units {
                    s := strconv.Itoa(int(u.Id))
                    fmt.Printf("key : %s\n",s)
				    if selectedUnits[ s ] {
					    units = append(units, s)
				    }
			    }
                playersRWLock.RUnlock()

				attack := false
                unitToAttack := ""

                enemy_id := (client_id + 1) % 2

				for k, v := range (*players)[enemy_id].Units {
                    if rl.CheckCollisionPointCircle(rl.Vector2{float32(v.X),float32(v.Y)}, dest, utils.Unit_size) {
						attack = true
                        unitToAttack = k
                        break
					}
                }
                if attack {
                    unitToAttackEntity := (*players)[enemy_id].Units[unitToAttack]
                    for _, v := range units {
                        unit := (*players)[client_id].Units[v]
                        r := unit.AttackRange
                        d := math.Sqrt(float64((unitToAttackEntity.X - unit.X) * (unitToAttackEntity.X - unit.X) + (unitToAttackEntity.Y - unit.Y) * (unitToAttackEntity.Y - unit.Y)))
                        fmt.Printf("%d <? %d", d, r)
                        if d > float64(r) {
                            attack = false
                            break
                        }
                    }
                }
				if attack {
                    fmt.Println("attack")
                    
                    attack := events.AttackUnit_e{Units: units, Unit: unitToAttack}
				    data, err := json.Marshal(attack)
				    utils.Check(err)
				    e := events.Event{EventType: events.AttackUnit, Data: string(data)}
				    e_marsh, err := json.Marshal(e)
				    utils.Check(err)
				    chan_gui_link<-string(e_marsh)

                    isDead := false
                    unit := (*players)[enemy_id].Units[unitToAttack]
                    for _, v := range units {
                        damage := (*players)[client_id].Units[v].AttackAmount
                        fmt.Printf("|%d|", time.Now().Nanosecond())
                        if unit.Health > damage {
                            unit.Health -= damage
                        } else {
                            isDead = true
			                delete((*players)[enemy_id].Units, unitToAttack)
                            break
                        }
                    }
                    if (!isDead) {
                        (*players)[enemy_id].Units[unitToAttack] = unit
                    }
                } else {
				    move := events.MoveUnits_e{Units: units, Dest: dest}
				    data, err := json.Marshal(move)
				    utils.Check(err)
				    e := events.Event{EventType: events.MoveUnits, Data: string(data)}
				    e_marsh, err := json.Marshal(e)
				    utils.Check(err)
				    chan_gui_link<-string(e_marsh)
                }
			}

			if (rl.IsMouseButtonPressed(rl.MouseLeftButton)){
				if inSelection == false {
					inSelection = true
					selectionStart = rl.GetMousePosition()
				}
			}
			if( rl.IsMouseButtonReleased(rl.MouseLeftButton)) {
				inSelection = false
				selectedUnits = map[string]bool{}
				for k, v := range (*players)[client_id].Units {
					if rl.CheckCollisionCircleRec(rl.Vector2{float32(v.X),float32(v.Y)}, utils.Unit_size, getRectangle2Pt(rl.GetScreenToWorld2D( selectionStart, camera),  rl.GetScreenToWorld2D( rl.GetMousePosition(), camera))  ) {
						selectedUnits[k] = true
					}
				}
			}
		}

		// Check wether or not to listen to the keys
		if (currentState & StateChat > 0) {
			// Wait for a message to be entered
			key := rl.GetKeyPressed()
			charPressed := rl.GetCharPressed()
			if (skipCharPressed && charPressed > 0) {
				charPressed = 0
				skipCharPressed = false
			}
			if (key != 0 || charPressed != 0) {
				utils.Logging("GUI", fmt.Sprintf("Key pressed: %d, char pressed: %d", key, charPressed))
			}
			if (key == rl.KeyEnter) {
				// Sending the message if it is not empty ("")
				splitted_message := SplitMessage(ChatText, message_max_len)
				for i := 0 ; i < len(splitted_message) ; i ++ {
					e := events.Event{
						EventType : events.ChatEvent,
						Data: fmt.Sprintf("%s: %s",
							config.Pseudo,
							splitted_message[i])}
					e_marsh, err := json.Marshal(e)
					utils.Check(err)
					chan_gui_link<- string(e_marsh)
					has_send += 1
				}
				// Reset the message and state
				ChatText = ""
				currentState = StateNone
			} else if (key == rl.KeyBackspace) {
				length := len(ChatText)
				if length > 0 {
					ChatText = ChatText[:(length-1)]
				}
			} else if (time.Since(lastInputTime) > 10 * time.Millisecond && charPressed > 0) {
				ChatText += fmt.Sprintf("%c", charPressed)
			}
			if (key > 0 || charPressed > 0) {
				lastInputTime = time.Now()
			}
		}



		// Draw to screenTexture
		//----------------------------------------------------------------------------------
		rl.BeginDrawing();
			rl.ClearBackground(rl.Black);
			rl.BeginMode2D(camera);
				// DRAW MAP
				utils.DrawMap(*gmap)
				// DRAW UNITS
                playersRWLock.Lock()
				for i, player := range *players {
					for k , player_unit := range player.Units {
						_, found := selectedUnits[k]
						utils.DrawUnit(player_unit, i== client_id, found )
					}
				}
                playersRWLock.Unlock()
                // DEBUG
                if flowField != nil {
                    utils.DrawFlowField(flowField,ffstep)
                }

			rl.EndMode2D();
			if( inSelection ){
				selRect := getRectangle2Pt(selectionStart, rl.GetMousePosition() )
				rl.DrawRectangleLines( selRect.ToInt32().X , selRect.ToInt32().Y, selRect.ToInt32().Width, selRect.ToInt32().Height, rl.Magenta)
			}


		// Check wether a menu has to be printed
		if (currentState & StateMenu > 0) {
			if(time.Since(lastInputTime) > 10 * time.Second) {
				currentState -= StateMenu
			}
			if currentMenu < 0 {
				// One must define a menu
				if len(selectedUnits) > 0 {
					currentMenu = 1 // Selection menu
				} else {
					currentMenu = 0 // Default menu
				}
				utils.Logging("GUI",
					fmt.Sprintf("Choosing the current menu: %d", currentMenu))
			}
			menu_options := make(map[int32]MenuElement_t, 0)
			if currentMenu >= 0 {
				// Print the current menu and its elements, and check for its hotkeys:
				current_menu := FindMenuByRef(config_menus.Menus, currentMenu)
				// Menu title
				rl.DrawText(current_menu.Title, 0, 0, 40, rl.Red)
				// Menu options and keys
				offset := 0
				for i := 0 ; i < len(current_menu.Elements) ; i ++ {
					if false {
						offset = offset + 1
					}
					rl.DrawText(
						fmt.Sprintf("%c) %s",
							current_menu.Elements[i].Key,
							current_menu.Elements[i].Name),
						100,
						int32(40 + (20 * (i - offset))),
						15,
						rl.Blue)

					// Adding the menu option to the [menu_options] mapping
					menu_options[current_menu.Elements[i].Key] = current_menu.Elements[i]
				}
			}
			// Check the delay since last interaction
			if(time.Since(lastInputTime) > time.Second) {
				for key, val := range menu_options {
					if(rl.IsKeyDown(key)) {
						switch val.Type {
							case MenuElementSubMenu:
								utils.Logging("GUI",
									fmt.Sprintf("Switching to the menu %d.",
										currentMenu))
								currentMenu = val.Ref
								lastInputTime = time.Now()
								break
							case MenuElementAction:
								lastInputTime = time.Now()
								currentAction = FindActionByRef(config_menus.Actions, val.Ref)
								utils.Logging("GUI",
									fmt.Sprintf("Executing the action %d: %s",
										val.Ref, currentAction.Title))
								if currentAction.Type == ActionQuitGame {
									os.Exit(0)
								} else if currentAction.Type == ActionMappedKeys {
									need_pause = true
									keys := GenKeysSubMenu(config.Keys)
									for i := 0 ; i < len(keys) ; i ++ {
										rl.DrawText(
											keys[i],
											100,
											int32(75 + (20 * i)),
											15,
											rl.Blue)
									}
								} else if currentAction.Type == ActionBuilding {
									mouse_pos,_ := get_mouse_grid_pos(camera, gmap.Width, gmap.Height)
									build := events.BuildBuilding_e {
										Position_x: int32(mouse_pos.X),
										Position_y: int32(mouse_pos.Y),
										BuildingName: currentAction.Title }
									data, err := json.Marshal(build)
									utils.Check(err)
									e := events.Event{EventType: events.BuildEvent, Data: string(data)}
									e_marsh, err := json.Marshal(e)
									utils.Check(err)
									currentMessages = NewMessageItem(currentMessages, "CLIENT:Building a " + currentAction.Title, config)
									organizeMessages(currentMessages)
									chan_gui_link<-string(e_marsh)
								} else {
									utils.Logging("GUI",
										fmt.Sprintf("Unknown action `%s`",
											currentAction.Title))
								}
								break
						}
					}
				}
			}
		}

		// Recompute the location of messages
		if time.Since(lastMessagesUpdate) > time.Second {
			for i := 0 ; i < len(currentMessages) ; i ++ {
				old_len := len(currentMessages)
				currentMessages = FilterMessages(currentMessages)
				if len(currentMessages) < old_len {
					organizeMessages(currentMessages)
				}
			}
		}
		// display chat box
		if (currentState & StateChat > 0) {
			message_position := rl.GetScreenHeight() - (message_font_size + message_padding)
			rl.DrawText(" > " + ChatText + "_", 0,  int32(message_position), int32(message_font_size), message_color_ours)
		}
		// Print chat messages
		for i := 0 ; i < len(currentMessages) ; i ++ {
			if currentMessages[i].Ownership {
				rl.DrawText(currentMessages[i].Message,
					int32(currentMessages[i].Position_x),
					int32(currentMessages[i].Position_y),
					int32(message_font_size),
					message_color_ours)
			} else {
				rl.DrawText(currentMessages[i].Message,
					int32(currentMessages[i].Position_x),
					int32(currentMessages[i].Position_y),
					int32(message_font_size),
					message_color)
			}
		}

		// Print help message
		rl.DrawText("Press ';' then 1 for help",
			int32(rl.GetScreenWidth() / 2),
			int32(rl.GetScreenHeight() - helpSize),
			int32(helpSize),
			rl.Green)

		rl.EndDrawing();
		if need_pause {
			time.Sleep(time.Second)
			need_pause = false
		}
	}
}



/*                +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+
                  | Drawing auxiliary functions |
                  +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+ */
func drawGrid(width int32, height int32) {
	for i := int32(0) ; i <= height ; i++ {
		rl.DrawLine(0, utils.TileSize *i, utils.TileSize*width ,utils.TileSize*i,rl.Red)
	}
	for i := int32(0) ; i <= width ; i++ {
		rl.DrawLine(utils.TileSize*i, 0, utils.TileSize*i ,utils.TileSize*height,rl.Red)
	}
}

func get_mouse_grid_pos(camera rl.Camera2D, width , height int32) (rl.Vector2, bool) {
	mouse_screen_pos := rl.GetMousePosition()
	mouse_world_pos := rl.GetScreenToWorld2D(mouse_screen_pos, camera)
	ret := rl.Vector2{ X : float32(math.Floor(float64(mouse_world_pos.X / float32(utils.TileSize)))), Y : float32(math.Floor(float64(mouse_world_pos.Y / float32(utils.TileSize)))) }
	//fmt.Printf("test : %d %d",int(),int())
	if ret.X < 0 || int32(ret.X) >= width || ret.Y < 0 || int32(ret.Y) >= height {
		return ret,true
	}else {
		return ret,false
	}
}

func getRectangle2Pt(p1 rl.Vector2, p2 rl.Vector2) rl.Rectangle{
    width :=  float32(math.Abs(float64(p1.X-p2.X)))
    height := float32(math.Abs(float64(p1.Y-p2.Y)))
    startX := float32(math.Min(float64(p1.X),float64(p2.X)))
    startY := float32(math.Min(float64(p1.Y),float64(p2.Y)))
    return rl.NewRectangle(startX,startY,width,height)
}
