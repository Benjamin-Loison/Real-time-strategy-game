package utils

import (
	"github.com/gen2brain/raylib-go/raylib"
	"container/heap"
//	"fmt"
)

type Point struct {
	X int32
	Y int32
}

// From https://pkg.go.dev/container/heap
type PointAndDistance struct {
	value    Point
	distance float32    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

type PriorityQueue []*PointAndDistance

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PointAndDistance)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *PointAndDistance, value Point, distance float32) {
	item.value = value
	item.distance = distance
	heap.Fix(pq, item.index)
}

func PathFinding(mapp Map, endPos rl.Vector2, step int32) [][]rl.Vector2 {

	map_width  := mapp.Width * TileSize
	map_height := mapp.Height * TileSize

	array_width  := map_width / step + 1
	array_height := map_height / step + 1

	var cost_field = make([][]uint16, array_width)
	var integration_field = make([][]float32, array_width)
	var flow_field = make([][]rl.Vector2, array_width)
	for x, _ := range flow_field {
		flow_field[x] = make([]rl.Vector2, array_height)
		cost_field[x] = make([]uint16, array_height)
		integration_field[x] = make([]float32, array_height)
	}

	for x, line := range cost_field {
		for y, _ := range line {
			cost_field[x][y] = 1
			if (int32(x)+1)*step > map_width || (int32(y)+1)*step > map_height {
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

					if int32(matrix.X) < mapp.Width &&
                       int32(matrix.Y) < mapp.Height &&
                       mapp.Grid[int32(matrix.X)][int32(matrix.Y)].Tile_Type != None {
						center := rl.Vector2Add(curr_vect, half_vect)
						center = rl.Vector2Add(center, v)
						center = rl.Vector2Scale(center, float32(TileSize))
						if rl.CheckCollisionCircleRec(center, float32(TileSize)/2-1, rect) {
							cost_field[x][y] = ^uint16(0)
						}
					}
				}
            }
        }
    }
	endPos = rl.Vector2Scale(endPos, 1/float32(step))
	end_x := int32(endPos.X)
	end_y := int32(endPos.Y)
    if cost_field[end_x][end_y] == ^uint16(0) {
        return nil
    }
	for x, line := range integration_field {
		for y, _ := range line {
			integration_field[x][y] = float32(^uint16(0))
		}
	}

	/*var visit_order = make([][]uint16, array_width)
	for x, _ := range visit_order {
		visit_order[x] = make([]uint16, array_height)
		for y, _ := range visit_order[x] {
			visit_order[x][y] = ^uint16(0)
		}
	}
	ivisit := uint16(0)*/
	queue := make(PriorityQueue, 1)
	queue[0] = &PointAndDistance{
		value:		Point{X: end_x, Y: end_y},
		distance: 	0,
		index: 0,
	}
	to_test := [8]Point{ {X: -1, Y: -1},{X: -1, Y: 0}, {X: -1, Y: 1},{X: 0, Y: -1}, {X: 0, Y: 1}, {X: 1, Y: -1}, {X: 1, Y:0}, {X: 1, Y: 1} }

	for queue.Len() > 0 {
		point_dist := heap.Pop(&queue).(*PointAndDistance)
		p := point_dist.value
		if integration_field[p.X][p.Y] != float32(^uint16(0)) {
			continue
		}
		integration_field[p.X][p.Y] = point_dist.distance
		/*visit_order[p.X][p.Y] = ivisit
		ivisit = ivisit + 1*/
		for _, voisin := range to_test {
			v_x := p.X + voisin.X
			v_y := p.Y + voisin.Y

			if v_x < 0 || v_x >= array_width || v_y < 0 || v_y >= array_height {
				continue
			}

			if cost_field[v_x][v_y] != ^uint16(0) {
                c_p := integration_field[p.X][p.Y]
				c_v := float32(cost_field[v_x][v_y])*rl.Vector2Length(rl.Vector2{X: float32(voisin.X),Y:float32(voisin.Y)})

				if c_p + c_v < integration_field[v_x][v_y] {
					heap.Push(&queue, &PointAndDistance{
						value:		Point{v_x, v_y},
						distance:	c_p + c_v,
					})
				}
			}
		}
	}
	// Log visiting order for integration field
	/*Logging("Pathfinding", "Visiting order for integration field:")
	for _, line := range visit_order {
		Logging("Pathfinding", fmt.Sprintf("%v", line))
	}*/
	for x, line := range flow_field {
		for y, _ := range line {
			min_cost := float32(^uint16(0))

            flow_field[x][y] = rl.Vector2{X: 0.0, Y: 0.0}
            //if integration_field[x][y] == 0 || cost_field[x][y] == ^uint16(0) {
            if integration_field[x][y] == 0 {
                continue
            }

			for _, voisin := range to_test {
				v_x := int32(x) + voisin.X
				v_y := int32(y) + voisin.Y

				if v_x < 0 || v_x >= array_width || v_y < 0 || v_y >= array_height {
					continue
				}

				if integration_field[v_x][v_y] < min_cost {
					min_cost = integration_field[v_x][v_y]
                    flow_field[x][y] = rl.Vector2Normalize( rl.Vector2{X: float32(voisin.X), Y: float32(voisin.Y)})
				}
			}
		}
	}

	return flow_field
}
