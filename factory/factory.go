package factory

import (
    "github.com/gen2brain/raylib-go/raylib"
)


type Owner int64

const(
    Player1 Owner = 1
    Player2       = 2
    NoOne         = 0
)
type Unit struct {
	X int32 `json:"X"`
	Y int32 `json:"Y"`
	Name string `json:"Name"`
	Id int32 `json:"Id"`
	OwnerPlayer Owner `json:"OwnerPlayer"`
    Health float32 `json:"Health"`
    MaxHealth float32 `json:"MaxHealth"`
    LastAttack uint32 `json:"LastAttack"` // time of last attack
    AttackCoolDown uint32 `json:"AttackCoolDown"` // number of server ticks between atacks
    Speed float32 `json:"Speed"`
    FlowField *[][]rl.Vector2
    FlowStep  int32
    FlowTarget rl.Vector2
}

//********************************** HUMANS ***********************************

func MakeHumanPeon(x int32, y int32, id  int32, own Owner) Unit {
    return Unit{X : x, Y: y ,Name : "P",Id: id, OwnerPlayer: own, Health: 200.0, MaxHealth: 200.0, LastAttack: 0, AttackCoolDown: 150, Speed: 0.5, FlowField: nil, FlowStep: 0, FlowTarget: rl.Vector2{}}
}
