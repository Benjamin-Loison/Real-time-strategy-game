package main

import "strings"

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
}

type Player struct {
    Units []Unit `json:"Units"`
}

type ServerMessage struct {
    MessageType ServerMessageType `json:"MessageType"`
    GameMap Map `json:"GameMap"`
    Units []Unit `json:"Units"`
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
			panic("Not implemented: recognition of non-ascii characters and description")
	}
}

