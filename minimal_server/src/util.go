package main

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
