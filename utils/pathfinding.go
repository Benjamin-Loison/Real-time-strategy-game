package utils

import (
	"github.com/gen2brain/raylib-go/raylib"
)

type Point struct {
	X int32
	Y int32
}
func pathFinding(mapp Map, endPos rl.Vector2, step int32) [][]rl.Vector2 {

	map_width  := mapp.Width * TileSize
	map_height := mapp.Height * TileSize

	array_width  := map_width / step + 2
	array_height := map_height / step + 2

	var cost_field = make([][]uint16, array_width)
	var integration_field = make([][]uint16, array_width)
	var flow_field = make([][]rl.Vector2, array_width)
	for x, _ := range flow_field {
		flow_field[x] = make([]rl.Vector2, array_height)
		cost_field[x] = make([]uint16, array_height)
		integration_field[x] = make([]uint16, array_height)
	}


	for x, line := range cost_field {
		for y, _ := range line {
			if x == 0 || int32(x) == array_width -1 || y == 0 || int32(y) == array_height -1 {
				cost_field[x][y] = ^uint16(0)
			} else {
				x_s := int32(x) * step
				y_s := int32(y) * step
				map_x := x_s / TileSize
				map_y := y_s / TileSize
				rect := rl.NewRectangle(float32(x_s), float32(y_s), float32(step), float32(step))
				to_test := [4]rl.Vector2{ {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y:1}, {X: 1, Y: 1} } 
				curr_vect := rl.Vector2{X: float32(map_x), Y: float32(map_y)}
				half_vect := rl.Vector2{X: 0.5, Y: 0.5}
				for _, v :=  range to_test {
					matrix := rl.Vector2Add(curr_vect, v)

					cost_field[x][y] = 1
					if mapp.Grid[int32(matrix.X)][int32(matrix.Y)].Tile_Type != None {
						center := rl.Vector2Add(curr_vect, half_vect)
						center = rl.Vector2Add(center, v)
						center = rl.Vector2Scale(center, float32(TileSize))
						if rl.CheckCollisionCircleRec(center, float32(TileSize)/2, rect) {
							cost_field[x][y] = ^uint16(0)
						}
					}
				}
			}
		}
	}
	for x, line := range integration_field {
		for y, _ := range line {
			integration_field[x][y] = ^uint16(0)
		}
	}
	endPos = rl.Vector2Scale(endPos, 1/float32(step))
	end_x := int32(endPos.Y) + 1
	end_y := int32(endPos.Y) + 1

	integration_field[end_x][end_y] = 0
	stack := []Point{{X:end_x, Y:end_y}}
	to_test := [8]Point{ {X: -1, Y: -1},{X: -1, Y: 0}, {X: -1, Y: 1},{X: 0, Y: -1}, {X: 0, Y: 1}, {X: 1, Y: -1}, {X: 1, Y:0}, {X: 1, Y: 1} } 

	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, voisin := range to_test {
			v_x := p.X + voisin.X
			v_y := p.Y + voisin.Y

			if v_x < 0 || v_x >= array_width || v_y < 0 || v_y >= array_height {
				continue
			}

			if cost_field[v_x][v_y] != ^uint16(0) {
				c_p := integration_field[p.X][p.Y]
				c_v := cost_field[v_x][v_y]

				if c_p + c_v < integration_field[v_x][v_y] {
					integration_field[v_x][v_y] = c_p + c_v
					stack = append(stack, Point{v_x, v_y})
				}
			}
		}
	}
	for x, line := range flow_field {
		for y, _ := range line {
			min_cost := ^uint16(0)

			for _, voisin := range to_test {
				v_x := int32(x) + voisin.X
				v_y := int32(y) + voisin.Y

				if v_x < 0 || v_x >= array_width || v_y < 0 || v_y >= array_height {
					continue
				}

				if integration_field[v_x][v_y] < min_cost {
					min_cost = integration_field[v_x][v_y]
					flow_field[x][y] = rl.Vector2{float32(voisin.X), float32(voisin.Y)}
				}
			}
		}
	}

	return flow_field
}
