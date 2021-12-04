package main

import (
	"bufio"
	_ "image/png"
	"math"
	"time"
    "encoding/json"
	"rts/events"
	"github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth = 1280
	screenHeight = 720
	zoomFactor = 1.03
	cameraSpeed = 300.0
)

type GameState int32

const (
	StateNone GameState = 0
	StateMenu GameState = 1
	StateAction GameState = 2
)

var (
	currentState = StateNone
	currentMenu = -1
	currentAction = -1
	timeMenu = time.Now()
)

func drawGrid(width int32, height int32) {
	for i := int32(0) ; i <= height ; i++ {
		rl.DrawLine(0, TileSize *i, TileSize*width ,TileSize*i,rl.Red)
	}
	for i := int32(0) ; i <= width ; i++ {
		rl.DrawLine(TileSize*i, 0, TileSize*i ,TileSize*height,rl.Red)
	}
}

func get_mouse_grid_pos(camera rl.Camera2D, width , height int32) (rl.Vector2, bool) {
	mouse_screen_pos := rl.GetMousePosition()
	mouse_world_pos := rl.GetScreenToWorld2D(mouse_screen_pos, camera)
	ret := rl.Vector2{ X : float32(math.Floor(float64(mouse_world_pos.X / float32(TileSize)))), Y : float32(math.Floor(float64(mouse_world_pos.Y / float32(TileSize)))) }
	//fmt.Printf("test : %d %d",int(),int())
	if ret.X < 0 || int32(ret.X) >= width || ret.Y < 0 || int32(ret.Y) >= height {
		return ret,true
	}else {
		return ret,false
	}
}

func RunGui(gmap *Map, players *[]Player, config Configuration_t, config_menus MenuConfiguration_t, chan_client chan string) {
	rl.SetTraceLog(rl.LogNone)
	rl.InitWindow(screenWidth, screenHeight, "RTS")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var map_width int32 = gmap.Width
	var map_height int32 = gmap.Height

	map_middle := rl.Vector2{X:float32(TileSize)*float32(map_width)/2.0,Y: float32(TileSize)*float32(map_height)/2.0 }


	camera := rl.NewCamera2D(rl.Vector2{X:screenWidth/2.0,Y:screenHeight/2.0},rl.Vector2{X: float32(TileSize)*float32(map_width)/2.0, Y: float32(TileSize)*float32(map_height)/2.0},0,1.0)

	rl.BeginMode2D(camera)

    writer := bufio.NewWriter(serv_conn)

	for !rl.WindowShouldClose() {
		// check that shouldn't quit
		select {
		case x, _ := <-chan_client:
			if x == "QUIT" {
				logging("gui","Forced to quit")
				return
			}
		default:
		}

		// Update
		//----------------------------------------------------------------------------------

		switch currentState {
			case StateMenu:
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
				//~Check the delay since last interaction
				if(time.Since(timeMenu) > time.Second) {
					for key, val := range menu_options {
						if(rl.IsKeyDown(key)) {
							switch val.Type {
								case MenuElementSubMenu:
									currentMenu = val.Ref
									timeMenu = time.Now()
								default:
									///
							}
						}
					}
				}
			case StateAction:
				// Wait for the action to be achieved.
			default:
				break
		}

		offsetThisFrame := cameraSpeed*rl.GetFrameTime()

        // test d'envoie d'event
        if (rl.IsKeyDown(rl.KeyQ)) {
            e := events.Event{EventType : events.MoveUnits,Data : "test"}
            e_marsh, err := json.Marshal(e)
            Check(err)
			logging("GUI","Trying to send event")
            _,err = writer.Write([]byte(string(e_marsh)+"\n"))
            writer.Flush()
            Check(err)
			logging("GUI","Event sent")

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
			if (currentMenu == -1 && time.Since(timeMenu) > time.Second) {
				currentState = StateMenu
				currentMenu = 0
				timeMenu = time.Now()
			} else if (time.Since(timeMenu) > time.Second) {
				currentState = StateNone
				currentMenu = -1
				timeMenu = time.Now()
			}
		}
		if (rl.IsKeyDown(config.Keys.ResetCamera)){
			camera.Zoom = 1.0
			camera.Target.X = map_middle.X
			camera.Target.Y = map_middle.Y
        }

		// Draw to screenTexture
		//----------------------------------------------------------------------------------
		rl.BeginDrawing();
			rl.ClearBackground(rl.Black);
			rl.BeginMode2D(camera);
				// DRAW MAP
				DrawMap(*gmap)
				// DRAW UNITS
				for i, player := range *players {
					for _, player_unit := range player.Units {
						DrawUnit(player_unit, i== client_id)
					}
				}

			rl.EndMode2D();
		rl.EndDrawing();

	}

}
