package main

import (
	"fmt"
    "os"
	_ "image/png"
    "time"
    "github.com/gen2brain/raylib-go/raylib"
    "math"
    "strconv"
)

const (
    screenWidth = 1280
    screenHeight = 720
    zoomFactor = 1.03
    cameraSpeed = 300.0
)

var (
)

func logging(src string, message string) {
	fmt.Println(time.Now().Format(time.ANSIC) + "[" + src + "] " + message)
}

func drawGrid(width int32, height int32) {
    for i := int32(0) ; i <= height ; i++ {
        rl.DrawLine(0,tileSize*i, tileSize*width ,tileSize*i,rl.Red)
    }
    for i := int32(0) ; i <= width ; i++ {
        rl.DrawLine(tileSize*i, 0, tileSize*i ,tileSize*height,rl.Red)
    }
}

func get_mouse_grid_pos(camera rl.Camera2D, width , height int32) (rl.Vector2, bool) {
    mouse_screen_pos := rl.GetMousePosition()
    mouse_world_pos := rl.GetScreenToWorld2D(mouse_screen_pos, camera)
    ret := rl.Vector2{ X : float32(math.Floor(float64(mouse_world_pos.X / float32(tileSize)))), Y : float32(math.Floor(float64(mouse_world_pos.Y / float32(tileSize)))) }
    //fmt.Printf("test : %d %d",int(),int())
    if ret.X < 0 || int32(ret.X) >= width || ret.Y < 0 || int32(ret.Y) >= height {
        return ret,true
    }else {
        return ret,false
    }
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "MAP BUILDER")

	rl.SetTargetFPS(60)

    var map_width int32
    var map_height int32
    var game_map Map

    if len(os.Args) == 3 {
        tmp1,e1 := strconv.Atoi(os.Args[1])
        check(e1)
        tmp2,e2 := strconv.Atoi(os.Args[2])
        check(e2)
        map_width = int32(tmp1)
        map_height = int32(tmp2)
        game_map = makeMap(map_width, map_height)
    }else if len(os.Args) == 2 {
        path := os.Args[1]
        game_map = loadMap(path)
        map_width = int32(game_map.Width)
        map_height = int32(game_map.Height)
        printMap(game_map)
    }else {
        fmt.Printf("\n\nUSAGE : map_builder width height | map_builder path_to_map\n\n")
        os.Exit(-1)
    }


    map_middle := rl.Vector2{X:float32(tileSize)*float32(map_width)/2.0,Y: float32(tileSize)*float32(map_height)/2.0 }


    camera := rl.NewCamera2D(rl.Vector2{X:screenWidth/2.0,Y:screenHeight/2.0},rl.Vector2{X: float32(tileSize)*float32(map_width)/2.0, Y: float32(tileSize)*float32(map_height)/2.0},0,1.0)

    rl.BeginMode2D(camera)


	for !rl.WindowShouldClose() {

        // Update
        //----------------------------------------------------------------------------------

        offsetThisFrame := cameraSpeed*rl.GetFrameTime()

        if (rl.IsKeyDown(rl.KeyRight)){
            //camera.Offset.X -= 2.0
            camera.Target.X += offsetThisFrame
        }
        if (rl.IsKeyDown(rl.KeyLeft)){
            //camera.Offset.X += 2.0
            camera.Target.X -= offsetThisFrame
        }
        if (rl.IsKeyDown(rl.KeyUp)){
            //camera.Offset.Y += 2.0
            camera.Target.Y -= offsetThisFrame
        }
        if (rl.IsKeyDown(rl.KeyDown)){
            //camera.Offset.Y -= 2.0
            camera.Target.Y += offsetThisFrame
        }
        if (rl.IsKeyDown(rl.KeyP)){
            camera.Zoom *= zoomFactor
        }
        if (rl.IsKeyDown(rl.KeyO)){
            camera.Zoom /= zoomFactor
        }
        if (rl.IsKeyDown(rl.KeyS)){
            saveMap(game_map)
        }
        if (rl.IsKeyDown(rl.KeySpace)){
            camera.Zoom = 1.0
            camera.Target.X = map_middle.X
            camera.Target.Y = map_middle.Y
        }
        mouse_grid_pos, err := get_mouse_grid_pos(camera, map_width, map_height)

        // Handle map modifications
        if !err {
            mx,my := int(mouse_grid_pos.X) , int(mouse_grid_pos.Y)
            if (rl.IsKeyDown(rl.KeyR)){
                game_map.Grid[mx][my].Tile_Type = Rock
                game_map.Grid[mx][my].Startpoint= NoOne
            }
            if (rl.IsKeyDown(rl.KeyD)){
                game_map.Grid[mx][my].Tile_Type = None
                game_map.Grid[mx][my].Startpoint= NoOne
            }
            if (rl.IsKeyDown(rl.KeyT)){
                game_map.Grid[mx][my].Tile_Type = Tree
                game_map.Grid[mx][my].Startpoint= NoOne
            }
            if (rl.IsKeyDown(rl.KeyG)){
                game_map.Grid[mx][my].Tile_Type = Gold
                game_map.Grid[mx][my].Startpoint= NoOne
            }
            if (rl.IsKeyDown(rl.KeyOne)){
                game_map.Grid[mx][my].Startpoint = Player1
                game_map.Grid[mx][my].Tile_Type = None
            }
            if (rl.IsKeyDown(rl.KeyTwo)){
                game_map.Grid[mx][my].Startpoint = Player2
                game_map.Grid[mx][my].Tile_Type = None
            }
        }
        // Draw to screenTexture
        //----------------------------------------------------------------------------------
        rl.BeginDrawing();
            rl.ClearBackground(rl.Black);
            rl.BeginMode2D(camera);

                drawMap(game_map)

                drawGrid(map_width,map_height)

            rl.EndMode2D();
        rl.EndDrawing();

    }

	rl.CloseWindow()
}
