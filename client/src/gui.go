package main

import (
	"fmt"
	_ "image/png"
	"math"
	"time"
    "encoding/json"
	"rts/events"
    "rts/utils"
	"github.com/gen2brain/raylib-go/raylib"
	"strings"
)

type MessageItem_t struct {
	Ownership bool
	Message string
	Position_x int
	Position_y int
	ArrivalTime time.Time
}

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
	max_messages_nb = 15
	message_duration = 10 * time.Second
	message_font_size = 15
	message_max_len = 20
	message_color = rl.Red
	message_color_ours = rl.Blue
	message_padding = 5
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
				if (currentMenu == -1 && time.Since(lastInputTime) > time.Second) {
					currentState = StateMenu
					currentMenu = 0
					lastInputTime = time.Now()
				} else if (time.Since(lastInputTime) > time.Second) {
					currentState = StateNone
					currentMenu = -1
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

		// Check wether a menu has to be printed
		if (currentState & StateMenu > 0) {
			menu_options := make(map[int32]MenuElement_t)
			if currentMenu >= 0 {
				// Print the current menu and its elements, and check for its hotkeys:
				current_menu := FindMenuByRef(config_menus.Menus, currentMenu)
				// Menu title
				rl.DrawText(current_menu.Title, 0, 0, 40, rl.Red)
				// Menu options and keys
				for i := 0 ; i < len(current_menu.Elements) ; i ++ {
					rl.DrawText(current_menu.Elements[i].Name,
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
								currentMenu = val.Ref
								lastInputTime = time.Now()
								break
							default:
								break
						}
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
	}
}




/*                +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+
                  | Messages auxiliary functions |
                  +~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~+ */
func SplitMessage(m string, length int) []string {
	split := make([]string, 0)
	if len(m) < 2 {
		return split
	}
	m = m[1:]
	for {
		if len(m) <= length {
			return append(split, m)
		} else {
			split = append(split, m[:length])
			m = m[length:]
		}
	}
}

func FilterMessages(messages []MessageItem_t) []MessageItem_t {
	res := make([]MessageItem_t, 0)
	for i := 0 ; i < len(messages) ; i ++ {
		if (time.Since(messages[i].ArrivalTime) < message_duration) {
			res = append(res, messages[i])
		}
	}
	return res
}

func organizeMessages(messages []MessageItem_t) {
	// The messages are printed on the bottom left corner of the screen
	for i := 0 ; i < len(messages) ; i ++ {
		messages[i].Position_x = 0
		messages[i].Position_y = rl.GetScreenHeight() - (i+1) * (message_font_size + message_padding)
	}
}

func NewMessageItem(current []MessageItem_t, new_message string, has_send int) ([]MessageItem_t) {
	ownership := has_send > 0
	if len(current) == max_messages_nb {
		// On enlève le premier!
		return append(current[1:], (MessageItem_t {ownership, new_message, 0, 0, time.Now()}))
	}
	return append(current, (MessageItem_t {ownership, new_message, 0, 0, time.Now()}))
}

func isPrintable(key int32) (bool) {
	return key >= 32 && key <= 126
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

