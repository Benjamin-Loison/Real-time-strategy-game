package main

import (
	_ "image/png"
    "github.com/gen2brain/raylib-go/raylib"
    "math"
	"time"
)

const (
    screenWidth = 1280
    screenHeight = 720
    zoomFactor = 1.03
    cameraSpeed = 300.0
)

var (
	currentMenu = ""
	timeMenu = time.Now()
)

func drawGrid(width int32, height int32) {
    for i := int32(0) ; i <= height ; i++ {
        rl.DrawLine(0,TileSize *i, TileSize*width ,TileSize*i,rl.Red)
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

		if len(currentMenu) > 0 {
			// Print the current menu and its elements, and check for its hotkeys:
			current_menu := FindMenuByRef(config_menus.Menus, currentMenu)
			rl.DrawText(current_menu.Title, 0, 0, 40, rl.Red)
			for i := 0 ; i < len(current_menu.Elements) ; i ++ {
				rl.DrawText(current_menu.Elements[i].Name, 100, int32(40 + (20 * i)), 15, rl.Blue)
				if(rl.IsKeyDown(current_menu.Elements[i].Key) && time.Since(timeMenu) > time.Second) {
					if current_menu.Elements[i].Type == MenuElementSubMenu {
						currentMenu = current_menu.Elements[i].Name
					}
					timeMenu = time.Now()
				}
			}
		}

        offsetThisFrame := cameraSpeed*rl.GetFrameTime()

        if (rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(config.Keys.Right)){
            //camera.Offset.X -= 2.0
            camera.Target.X += offsetThisFrame
        }
        if (rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(config.Keys.Left)){
            //camera.Offset.X += 2.0
            camera.Target.X -= offsetThisFrame
        }
        if (rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(config.Keys.Up)){
            //camera.Offset.Y += 2.0
            camera.Target.Y -= offsetThisFrame
        }
        if (rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(config.Keys.Up)){
            //camera.Offset.Y -= 2.0
            camera.Target.Y += offsetThisFrame
        }
        if (rl.IsKeyDown(rl.KeyP)){
            camera.Zoom *= zoomFactor
        }
        if (rl.IsKeyDown(rl.KeyO)){
            camera.Zoom /= zoomFactor
        }
        if (rl.IsKeyDown(rl.KeyM)){
			if (currentMenu == "" && time.Since(timeMenu) > time.Second) {
				currentMenu = "Main"
				timeMenu = time.Now()
			} else if (time.Since(timeMenu) > time.Second) {
				currentMenu = ""
				timeMenu = time.Now()
			}
        }
        if (rl.IsKeyDown(rl.KeySpace)){
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
