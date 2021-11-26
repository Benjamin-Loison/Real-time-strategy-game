package main

import (
    "github.com/gen2brain/raylib-go/raylib"
	"fmt"
	"encoding/json"
)

const (
	screenWidth = 1280
	screenHeight = 720
	zoomFactor = 1.03
	cameraSpeed = 300.0
	config_path = "../conf/conf.json"
	TileSize = 32
)
// TODO : Remove struct once integrated with client
type Keys_t struct {
	Left int32
	Right int32
	Up int32
	Down int32
	ZoomIn int32
	ZoomOut int32
	Menu int32
	ResetCamera int32
}

func writeConfig(key_codes [8]int32) {
	keys := Keys_t{}
	keys.Left = key_codes[1]
	keys.Right = key_codes[0]
	keys.Up = key_codes[2]
	keys.Down = key_codes[3]
	keys.ZoomIn = key_codes[4]
	keys.ZoomOut = key_codes[5]
	keys.Menu = key_codes[7]
	keys.ResetCamera = key_codes[6]

	data, _ := json.MarshalIndent(keys, "", "\t")
	fmt.Println(string(data))
}

func main() {
	rl.SetTraceLog(rl.LogNone)
	rl.InitWindow(screenWidth, screenHeight, "RTS")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)


	camera := rl.NewCamera2D(rl.Vector2{X:screenWidth/2.0,Y:screenHeight/2.0},
							rl.Vector2{X: 0, Y: 0}, 0, 1.0)

	keys_name := [8]string{"Camera right", "Camera left", "Camera up", "Camera Down", "Zoom in", "Zoom out", "Reset Camera", "Open menu"}
	keys_code := [8]int32{}
	index := 0
	last_key := int32(-1)

	rl.BeginMode2D(camera)


	for !rl.WindowShouldClose() {
		// Draw to screenTexture
		//----------------------------------------------------------------------------------
		rl.BeginDrawing();
			rl.ClearBackground(rl.Black);
			rl.BeginMode2D(camera);
			rl.DrawText(fmt.Sprintf("Enter key for %s", keys_name[index]), 0, 0 , 32, rl.White)
			rl.EndMode2D();
		rl.EndDrawing();

		curr_key := rl.GetKeyPressed()
		if curr_key <= 0 || curr_key == last_key {
			continue
		}
		last_key = curr_key
		keys_code[index] = curr_key
		index++
		if index == 8 {
			break
		}


	}

	writeConfig(keys_code)

}

