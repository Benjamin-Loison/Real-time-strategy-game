module teamrts/rts

go 1.17

require github.com/gen2brain/raylib-go/raylib v0.0.0-20211114111602-29ba3cc50849

require "rts/events" v0.0.1
require "rts/utils" v0.0.1

replace "rts/events" v0.0.1 => "../events"
replace "rts/utils" v0.0.1 => "../utils"
