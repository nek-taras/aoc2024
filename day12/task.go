package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

type Coordinate struct {
	x int
	y int
}

type Fence struct {
	x           int
	y           int
	is_vertical bool
}

type PlantInField struct {
	diff_neighbours int
	kind            string
	in_a_group      bool
	coord           Coordinate
	fences          []Fence
}

func calcualteSides(group []PlantInField) int {
	if len(group) == 1 {
		return 4
	}
	sides := 0
	// filter out ones that is inside
	xs := make(map[int][]PlantInField)
	ys := make(map[int][]PlantInField)
	for _, p := range group {
		if len(p.fences) > 0 {
			rowx, ok := xs[p.coord.x]
			if !ok {
				rowx = []PlantInField{}
			}
			rowx = append(rowx, p)
			xs[p.coord.x] = rowx
			rowy, ok := ys[p.coord.y]
			if !ok {
				rowy = []PlantInField{}
			}
			rowy = append(rowy, p)
			ys[p.coord.y] = rowy
		}
	}
	for k, v := range xs {
		left := k - 1
		right := k + 1
		left_f := []int{}
		right_f := []int{}
		for _, p := range v {
			for _, f := range p.fences {
				if !f.is_vertical {
					continue
				}
				if f.x == left {
					left_f = append(left_f, f.y)
				}
				if f.x == right {
					right_f = append(right_f, f.y)
				}
			}
		}
		sort.Ints(left_f)
		sort.Ints(right_f)
		if len(left_f) > 0 {
			sides++
			if len(left_f) > 1 {
				curr := left_f[0]
				for i := 1; i < len(left_f); i++ {
					next := left_f[i]
					if curr != next-1 {
						sides++
					}
					curr = next
				}

			}
		}
		if len(right_f) > 0 {
			sides++
			if len(right_f) > 1 {
				curr := right_f[0]
				for i := 1; i < len(right_f); i++ {
					next := right_f[i]
					if curr != next-1 {
						sides++
					}
					curr = next
				}

			}
		}
	}
	for k, v := range ys {
		top := k - 1
		bottom := k + 1
		top_f := []int{}
		bottom_f := []int{}
		for _, p := range v {
			for _, f := range p.fences {
				if f.is_vertical {
					continue
				}
				if f.y == top {
					top_f = append(top_f, f.x)
				}
				if f.y == bottom {
					bottom_f = append(bottom_f, f.x)
				}
			}
		}
		sort.Ints(top_f)
		sort.Ints(bottom_f)
		if len(top_f) > 0 {
			sides++
			if len(top_f) > 1 {
				curr := top_f[0]
				for i := 1; i < len(top_f); i++ {
					next := top_f[i]
					if curr != next-1 {
						sides++
					}
					curr = next
				}

			}
		}
		if len(bottom_f) > 0 {
			sides++
			if len(bottom_f) > 1 {
				curr := bottom_f[0]
				for i := 1; i < len(bottom_f); i++ {
					next := bottom_f[i]
					if curr != next-1 {
						sides++
					}
					curr = next
				}

			}
		}
	}
	return sides

}

func findAGroup(x int, y int, field_map [][]PlantInField) []PlantInField {
	pif := field_map[y][x]
	if pif.in_a_group {
		return []PlantInField{}
	}
	map_height := len(field_map)
	map_width := len(field_map[0])
	pif.in_a_group = true
	same_naoghbours := []PlantInField{}
	top_y, top_x := y-1, x
	bottom_y, bottom_x := y+1, x
	left_y, left_x := y, x-1
	right_y, right_x := y, x+1
	if top_y >= 0 && top_y < map_height && top_x >= 0 && top_x < map_width {
		neighbour := field_map[top_y][top_x]
		if neighbour.kind == pif.kind {
			same_naoghbours = append(same_naoghbours, neighbour)
		} else {
			pif.fences = append(pif.fences, Fence{y: top_y, x: top_x, is_vertical: false})
		}
	} else {
		pif.fences = append(pif.fences, Fence{y: top_y, x: top_x, is_vertical: false})
	}
	if bottom_y >= 0 && bottom_y < map_height && bottom_x >= 0 && bottom_x < map_width {
		neighbour := field_map[bottom_y][bottom_x]
		if neighbour.kind == pif.kind {
			same_naoghbours = append(same_naoghbours, neighbour)
		} else {
			pif.fences = append(pif.fences, Fence{y: bottom_y, x: bottom_x, is_vertical: false})
		}
	} else {
		pif.fences = append(pif.fences, Fence{y: bottom_y, x: bottom_x, is_vertical: false})
	}
	if left_y >= 0 && left_y < map_height && left_x >= 0 && left_x < map_width {
		neighbour := field_map[left_y][left_x]
		if neighbour.kind == pif.kind {
			same_naoghbours = append(same_naoghbours, neighbour)
		} else {
			pif.fences = append(pif.fences, Fence{y: left_y, x: left_x, is_vertical: true})
		}
	} else {
		pif.fences = append(pif.fences, Fence{y: left_y, x: left_x, is_vertical: true})
	}
	if right_y >= 0 && right_y < map_height && right_x >= 0 && right_x < map_width {
		neighbour := field_map[right_y][right_x]
		if neighbour.kind == pif.kind {
			same_naoghbours = append(same_naoghbours, neighbour)
		} else {
			pif.fences = append(pif.fences, Fence{y: right_y, x: right_x, is_vertical: true})
		}
	} else {
		pif.fences = append(pif.fences, Fence{y: right_y, x: right_x, is_vertical: true})
	}
	pif.diff_neighbours = 4 - len(same_naoghbours)
	field_map[y][x] = pif
	group := []PlantInField{pif}
	for _, n := range same_naoghbours {
		if !n.in_a_group {
			n_group := findAGroup(n.coord.x, n.coord.y, field_map)
			group = append(group, n_group...)
		}
	}
	return group
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Input file is missing.")
		os.Exit(1)
	}
	//fmt.Printf("opening %s\n", args[0])
	f, err := os.Open(args[0])

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	field_map := [][]PlantInField{}
	plant_groups := [][]PlantInField{}
	y := 0
	for scanner.Scan() {
		row := scanner.Text()
		field_map_row := []PlantInField{}
		for x, el := range row {
			pf := PlantInField{
				diff_neighbours: 0,
				kind:            string(el),
				in_a_group:      false,
				coord:           Coordinate{y: y, x: x},
				fences:          []Fence{},
			}
			field_map_row = append(field_map_row, pf)
		}
		y++
		field_map = append(field_map, field_map_row)
	}
	for y := range field_map {
		row := field_map[y]
		for x := range row {
			if !field_map[y][x].in_a_group {
				plant_group := findAGroup(x, y, field_map)
				plant_groups = append(plant_groups, plant_group)
			}
		}
	}
	total_fence_price := 0
	total_fence_price_second_part := 0
	for _, pg := range plant_groups {
		group_squaer := len(pg)
		group_p := 0
		for _, el := range pg {
			group_p += el.diff_neighbours
		}
		total_fence_price += group_squaer * group_p
		sides := calcualteSides(pg)
		total_fence_price_second_part += group_squaer * sides
	}
	fmt.Printf("Total price: %d\n", total_fence_price)
	fmt.Printf("Total price with sides: %d\n", total_fence_price_second_part)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
