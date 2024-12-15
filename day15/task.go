package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type ItemType int

const (
	Wall        ItemType = iota
	Box                  = iota
	Robot                = iota
	FreeSpace            = iota
	BigBoxStart          = iota
	BigBoxEnd            = iota
)

type BigBox struct {
	y_f int
	x_f int
	y_s int
	x_s int
}

type NextItem struct {
	item_type ItemType
	y         int
	x         int
}

func getBox(b_type ItemType, b_y int, b_x int) BigBox {
	if b_type == BigBoxStart {
		return BigBox{y_f: b_y, x_f: b_x, y_s: b_y, x_s: b_x + 1}
	}
	return BigBox{y_f: b_y, x_f: b_x - 1, y_s: b_y, x_s: b_x}

}

type Move int

const (
	Left  Move = iota
	Rgiht      = iota
	Up         = iota
	Down       = iota
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func getNext(move Move, y int, x int) (int, int) {
	next_y := y
	next_x := x
	switch move {
	case Up:
		next_y--
	case Down:
		next_y++
	case Rgiht:
		next_x++
	case Left:
		next_x--
	}
	return next_y, next_x

}

func getNextInLine(room_map [][]ItemType, move Move, boxes []BigBox) []NextItem {
	next_items := []NextItem{}
	for _, box := range boxes {
		if move == Up || move == Down {
			n_y, n_x := getNext(move, box.y_f, box.x_f)
			next_item_type := room_map[n_y][n_x]
			first_item := NextItem{item_type: next_item_type, y: n_y, x: n_x}
			next_items = append(next_items, first_item)
			n_y, n_x = getNext(move, box.y_s, box.x_s)
			next_item_type = room_map[n_y][n_x]
			second_item := NextItem{item_type: next_item_type, y: n_y, x: n_x}
			next_items = append(next_items, second_item)
		} else if move == Rgiht {
			n_y, n_x := getNext(move, box.y_s, box.x_s)
			next_item_type := room_map[n_y][n_x]
			item := NextItem{item_type: next_item_type, y: n_y, x: n_x}
			next_items = append(next_items, item)
		} else if move == Left {
			n_y, n_x := getNext(move, box.y_f, box.x_f)
			next_item_type := room_map[n_y][n_x]
			item := NextItem{item_type: next_item_type, y: n_y, x: n_x}
			next_items = append(next_items, item)

		}

	}
	return next_items

}

func checkIfFree(row []NextItem) bool {
	for _, item := range row {
		if item.item_type != FreeSpace {
			return false
		}
	}
	return true
}

func checkIfWall(row []NextItem) bool {
	for _, item := range row {
		if item.item_type == Wall {
			return true
		}
	}
	return false
}

func moveBoxes(room_map [][]ItemType, move Move, boxes []BigBox) {
	for i := len(boxes) - 1; i >= 0; i-- {
		box := boxes[i]
		room_map[box.y_s][box.x_s] = FreeSpace
		room_map[box.y_f][box.x_f] = FreeSpace
		n_y, n_x := getNext(move, box.y_s, box.x_s)
		room_map[n_y][n_x] = BigBoxEnd
		n_y, n_x = getNext(move, box.y_f, box.x_f)
		room_map[n_y][n_x] = BigBoxStart
	}

}

func moveRobot(room_map [][]ItemType, move Move, r_y int, r_x int) (int, int) {
	next_y, next_x := getNext(move, r_y, r_x)
	next_item := room_map[next_y][next_x]
	if next_item == Wall {
		return r_y, r_x
	} else if next_item == FreeSpace {
		room_map[r_y][r_x] = FreeSpace
		room_map[next_y][next_x] = Robot
		return next_y, next_x
	} else if next_item == Box {
		b_y, b_x := getNext(move, next_y, next_x)
		b_item := room_map[b_y][b_x]
		for b_item == Box {
			b_y, b_x = getNext(move, b_y, b_x)
			b_item = room_map[b_y][b_x]
		}
		if b_item == Wall {
			return r_y, r_x
		} else if b_item == FreeSpace {
			room_map[r_y][r_x] = FreeSpace
			room_map[next_y][next_x] = Robot
			room_map[b_y][b_x] = Box
			return next_y, next_x

		}
	} else if next_item == BigBoxStart || next_item == BigBoxEnd {
		box := getBox(next_item, next_y, next_x)
		boxes_in_a_way := []BigBox{box}
		next_in_line := getNextInLine(room_map, move, boxes_in_a_way)
		is_all_free := checkIfFree(next_in_line)
		is_a_wall := checkIfWall(next_in_line)
		for !is_a_wall && !is_all_free {
			next_boxes_in_line := []BigBox{}
			for _, el := range next_in_line {
				if el.item_type == BigBoxStart || el.item_type == BigBoxEnd {
					box := getBox(el.item_type, el.y, el.x)
					if !slices.Contains(next_boxes_in_line, box) {
						next_boxes_in_line = append(next_boxes_in_line, box)
					}
				}
			}
			boxes_in_a_way = append(boxes_in_a_way, next_boxes_in_line...)
			next_in_line = getNextInLine(room_map, move, next_boxes_in_line)
			is_all_free = checkIfFree(next_in_line)
			is_a_wall = checkIfWall(next_in_line)
		}
		if is_a_wall {
			return r_y, r_x
		} else if is_all_free {
			// move boxes
			moveBoxes(room_map, move, boxes_in_a_way)
			room_map[next_y][next_x] = Robot
			room_map[r_y][r_x] = FreeSpace
			return next_y, next_x
		}

	}
	return r_y, r_x
}

func printMap(room_map [][]ItemType) {
	for _, row := range room_map {
		for _, item := range row {
			switch item {
			case Wall:
				fmt.Printf("#")
			case Box:
				fmt.Printf("O")
			case Robot:
				fmt.Printf("@")
			case FreeSpace:
				fmt.Printf(".")
			case BigBoxStart:
				fmt.Printf("[")
			case BigBoxEnd:
				fmt.Printf("]")
			}
		}
		fmt.Printf("\n")
	}

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
	room_map := [][]ItemType{}
	big_room_map := [][]ItemType{}
	moves := []Move{}
	y := 0
	r_x, r_y := 0, 0
	rb_x, rb_y := 0, 0
	for scanner.Scan() {
		row := scanner.Text()
		if strings.HasPrefix(row, "#") {
			row_items := []ItemType{}
			big_row_items := []ItemType{}
			for x, item_symbol := range row {
				var item ItemType
				switch item_symbol {
				case '#':
					item = Wall
					big_row_items = append(big_row_items, Wall, Wall)
				case 'O':
					item = Box
					big_row_items = append(big_row_items, BigBoxStart, BigBoxEnd)
				case '.':
					item = FreeSpace
					big_row_items = append(big_row_items, FreeSpace, FreeSpace)
				case '@':
					item = Robot
					big_row_items = append(big_row_items, Robot, FreeSpace)
					r_x, r_y = x, y
					rb_x, rb_y = x*2, y
				}
				row_items = append(row_items, item)
			}
			room_map = append(room_map, row_items)
			big_room_map = append(big_room_map, big_row_items)
			y++
		} else if strings.HasPrefix(row, "<") || strings.HasPrefix(row, ">") || strings.HasPrefix(row, "^") || strings.HasPrefix(row, "v") {
			for _, move_symbol := range row {
				var move_item Move
				switch move_symbol {
				case '<':
					move_item = Left
				case '>':
					move_item = Rgiht
				case '^':
					move_item = Up
				case 'v':
					move_item = Down
				}
				moves = append(moves, move_item)
			}
		}
	}

	for _, move := range moves {
		r_y, r_x = moveRobot(room_map, move, r_y, r_x)
	}
	for _, move := range moves {
		rb_y, rb_x = moveRobot(big_room_map, move, rb_y, rb_x)
	}
	coordinate_score := 0
	for y, row := range room_map {
		for x, item := range row {
			if item == Box {
				coordinate_score += y*100 + x
			}
		}
	}
	big_coordinate_score := 0
	for y, row := range big_room_map {
		for x, item := range row {
			if item == BigBoxStart {
				big_coordinate_score += y*100 + x
			}
		}
	}
	printMap(room_map)
	fmt.Printf("Coordinate score:%d\n", coordinate_score)
	printMap(big_room_map)
	fmt.Printf("Coordinate score big room:%d\n", big_coordinate_score)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
