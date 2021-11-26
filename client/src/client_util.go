package main

import (
	"strings"
	"github.com/gen2brain/raylib-go/raylib"
	)

type ServerMessageType int32

const (
	MapInfo ServerMessageType = 0
	StartingUnits = 1
	Update = 2
)

type Unit struct {
	X int32 `json:"X"`
	Y int32 `json:"Y"`
	Name string `json:"Name"`
	Id int32 `json:"Id"`
	OwnerPlayer Owner `json:"OwnerPlayer"`
}

type Player struct {
	Units map[string]Unit `json:"Units"`
	Seed int `json:"Seed"`
}

type ServerMessage struct {
	MessageType ServerMessageType `json:"MessageType"`
	GameMap Map `json:"GameMap"`
	Players []Player `json:"Players"`
	Id int `json:"Id"`
}

//Raylib represents keys as int32 values
func keyOfString(s string)(int32) {
	switch l := len(s); l {
		case 0:
			panic("An empty string cannot represent a key in the configuration file.")
		case 1:
			return int32([]rune(strings.ToUpper(s))[0])
		default:
			switch s {
				case "SPACE":
					return rl.KeySpace
				case "RIGHT":
					return rl.KeyRight
				case "LEFT":
					return rl.KeyLeft
				case "DOWN":
					return rl.KeyDown
				case "UP":
					return rl.KeyUp
				default:
					panic("Not implemented: recognition of non-ascii characters and description")
			}
	}
}

