package events

import (
    "github.com/gen2brain/raylib-go/raylib"
    "rts/factory"
)

type Event_t int32

const (
    MoveUnits Event_t = iota
    ChatEvent
    ServerUpdate
    Attack
)

type MoveUnits_e struct {
    Units   []string        `json:"Units"`
    Dest    rl.Vector2      `json:"Dest"`
}

type ServerUpdate_e struct {
    Units   []factory.Unit        `json:"Units"`
}

type Event struct {
    EventType  Event_t `json:"EventType"`
    Data       string `json:"Data"`
}
