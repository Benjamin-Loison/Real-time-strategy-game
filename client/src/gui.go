package main

import (
	"fmt"
	"os"
	_ "image/png"
	"math"
	"time"
	"encoding/json"
	"rts/events"
	"rts/utils"
	"github.com/gen2brain/raylib-go/raylib"
	"strings"
	"sort"
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

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var map_width int32 = gmap.Width
	var map_height int32 = gmap.Height

	map_middle := rl.Vector2{X:float32(utils.TileSize)*float32(map_width)/2.0,Y: float32(utils.TileSize)*float32(map_height)/2.0 }


	camera := rl.NewCamera2D(rl.Vector2{X:screenWidth/2.0,Y:screenHeight/2.0},rl.Vector2{X: float32(utils.TileSize)*float32(map_width)/2.0, Y: float32(utils.TileSize)*float32(map_height)/2.0}, 0, 1.0)

	rl.BeginMode2D(camera)

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
				currentMessages = NewMessageItem(currentMessages, x[5:], has_send)
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
		             | Update the game |
		             +~~~~~~~~~~~~~~~~~+ */

		// Check wether the movements have to be detected
		if (currentState & StateNone > 0) {
			offsetThisFrame := cameraSpeed*rl.GetFrameTime()

			if (rl.IsKeyDown(config.Keys.Chat)) {
				utils.Logging("GUI", "Entering chat state.")
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

			if (rl.IsMouseButtonPressed(rl.MouseRightButton)) { // && len(selectedUnits)> 0){
                //fmt.Println(rl.GetScreenToWorld2D(rl.GetMousePosition(),camera))
                //flowField = utils.PathFinding(*gmap,rl.GetScreenToWorld2D(rl.GetMousePosition(),camera),ffstep)
				
				// We send the order to move units to the server
				dest := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
				units := []string{}

				for k := range (*players)[client_id].Units {
					if selectedUnits[k] {
						units = append(units, k)
					}
				}

				move := events.MoveUnits_e{Units: units, Dest: dest}
				data, err := json.Marshal(move)
				utils.Check(err)
				e := events.Event{EventType: events.MoveUnits, Data: string(data)}
				e_marsh, err := json.Marshal(e)
				utils.Check(err)
				chan_gui_link<-string(e_marsh)
				has_send++


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
				// Reset the lmessage and state
				ChatText = ""
				currentState = StateNone
			} else if (key == rl.KeyBackspace) {
				length := len(ChatText)
				if length > 2 {
					ChatText = ChatText[:(length-1)]
				}
			} else if (time.Since(lastInputTime) > 10 * time.Millisecond && isPrintable(key)) {
				ChatText += fmt.Sprintf("%c", key)
			}
			if (key > 0) {
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
				for i, player := range *players {
					for k , player_unit := range player.Units {
						_, found := selectedUnits[k]
						utils.DrawUnit(player_unit, i== client_id, found )
					}
				}
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
				for i := 0 ; i < len(current_menu.Elements) ; i ++ {
					rl.DrawText(
						fmt.Sprintf("%c) %s",
							current_menu.Elements[i].Key,
							current_menu.Elements[i].Name),
						100,
						int32(40 + (20 * i)),
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

