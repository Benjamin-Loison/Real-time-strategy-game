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
)

var (
	currentState = StateNone
	currentMenu = -1
	currentAction = -1
	lastInputTime = time.Now()
	max_messages_nb = 15
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
	currentMessages := make([]MessageItem_t, 0)

    var selectionStart rl.Vector2
    var inSelection = false

	rl.SetTraceLog(rl.LogNone)
	rl.InitWindow(screenWidth, screenHeight, "RTS")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var map_width int32 = gmap.Width
	var map_height int32 = gmap.Height

	map_middle := rl.Vector2{X:float32(utils.TileSize)*float32(map_width)/2.0,Y: float32(utils.TileSize)*float32(map_height)/2.0 }


	camera := rl.NewCamera2D(rl.Vector2{X:screenWidth/2.0,Y:screenHeight/2.0},rl.Vector2{X: float32(utils.TileSize)*float32(map_width)/2.0, Y: float32(utils.TileSize)*float32(map_height)/2.0},0,1.0)

	rl.BeginMode2D(camera)

	chan_link_gui<- "Haha, je suis là!\n"

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
				utils.Logging("gui", " Hé mais j'ai un message quoi !")
				currentMessages = NewMessageItem(currentMessages, x[5:])
			} else {
				utils.Logging("gui", fmt.Sprintf("Je comprendspas %s", x))
			}
		default:
		}

		/*           +~~~~~~~~~~~~~~~~~+
		             | Update the game |
		             +~~~~~~~~~~~~~~~~~+ */

		// Check wether the movements have to be detected
		if (currentState & StateNone > 0) {
			offsetThisFrame := cameraSpeed*rl.GetFrameTime()

			// test d'envoie d'event
			if (rl.IsKeyDown(rl.KeyQ)) {
				e := events.Event{EventType : events.MoveUnits,Data : "test"}
				e_marsh, err := json.Marshal(e)
				utils.Check(err)
				chan_gui_link<- string(e_marsh)
				utils.Logging("GUI","Event sent")
			}

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

			// Print chat messages
			for i := 0 ; i < len(currentMessages) ; i ++ {
				rl.DrawText(currentMessages[i].Message,
					int32(currentMessages[i].Position_x),
					int32(currentMessages[i].Position_y),
					40,
					rl.Red)
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
							default:
								///
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
				if len(ChatText) > 0 {
					e := events.Event{
						EventType : events.ChatEvent,
						Data: fmt.Sprintf("%s: %s", config.Pseudo, ChatText[1:])}
					e_marsh, err := json.Marshal(e)
					utils.Check(err)
					chan_gui_link<- string(e_marsh)

					// Reset the lmessage and state
					ChatText = ""
					currentState = StateNone
				}
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

			rl.EndMode2D();
            if( inSelection ){
                selRect := getRectangle2Pt(selectionStart, rl.GetMousePosition() )
                rl.DrawRectangleLines( selRect.ToInt32().X , selRect.ToInt32().Y, selRect.ToInt32().Width, selRect.ToInt32().Height, rl.Magenta)
            }
		rl.EndDrawing();
	}

}

func drawGrid(width int32, height int32) {
	for i := int32(0) ; i <= height ; i++ {
		rl.DrawLine(0, utils.TileSize *i, utils.TileSize*width ,utils.TileSize*i,rl.Red)
	}
	for i := int32(0) ; i <= width ; i++ {
		rl.DrawLine(utils.TileSize*i, 0, utils.TileSize*i ,utils.TileSize*height,rl.Red)
	}
}

func isPrintable(key int32) (bool) {
	return key >= 32 && key <= 126
}

func NewMessageItem(current []MessageItem_t, new_message string) ([]MessageItem_t) {
	if len(current) == max_messages_nb {
		// On enlève le premier!
		return append(current[1:], (MessageItem_t {new_message, 0, 0, time.Now()}))
	}
	return append(current, (MessageItem_t {new_message, 0, 0, time.Now()}))
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
