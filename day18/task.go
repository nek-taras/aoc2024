package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ItemType int

const (
	Wall      ItemType = iota
	FreeSpace          = iota
)

type Move int

type Coordinate struct {
	y int
	x int
}

const (
	Left  Move = iota
	Up         = iota
	Right      = iota
	Down       = iota
)

type LabyrintItem struct {
	item_type ItemType
	is_locked bool
	min_price int
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func printMap(room_map [][]LabyrintItem, start Coordinate, moves []Move) {

	move_map := make(map[string]Move)
	next := Coordinate{y: start.y, x: start.x}
	for i := 0; i < len(moves); i++ {
		m := moves[i]
		next = getNextCoordinate(next, m)
		coordinate := fmt.Sprintf("%d-%d", next.y, next.x)
		move_map[coordinate] = m
	}
	fmt.Println(move_map)
	for y, row := range room_map {
		for x, item := range row {
			coordinate := fmt.Sprintf("%d-%d", y, x)
			move, ok := move_map[coordinate]
			if ok {
				switch move {
				case Up:
					fmt.Printf("^")
				case Down:
					fmt.Printf("v")
				case Left:
					fmt.Printf("<")
				case Right:
					fmt.Printf(">")
				}

			} else {
				switch item.item_type {
				case Wall:
					fmt.Printf("#")
				case FreeSpace:
					fmt.Printf(".")
				}

			}

		}
		fmt.Printf("\n")
	}

}

func getNextCoordinate(pos Coordinate, move Move) Coordinate {
	y, x := pos.y, pos.x
	switch move {
	case Up:
		y--
	case Down:
		y++
	case Left:
		x--
	case Right:
		x++
	}

	return Coordinate{y: y, x: x}
}

func getAnswerPrice(answer []Move) int {
	return len(answer)
}

func getBestAnswer(answers [][]Move) []Move {
	sort.Slice(answers, func(i, j int) bool {
		first_price := getAnswerPrice(answers[i])
		second_price := getAnswerPrice(answers[j])
		return first_price < second_price
	})
	return answers[0]
}

func cleanUpScores(room_map [][]LabyrintItem) {
	width := len(room_map[0])
	height := len(room_map)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			room_item := room_map[y][x]
			if room_item.item_type == FreeSpace {
				room_map[y][x] = LabyrintItem{item_type: FreeSpace, is_locked: false, min_price: 0}
			}
		}
	}

}

func isValidItem(item LabyrintItem) bool {
	return !item.is_locked
}

func findPath(room_map [][]LabyrintItem, start_pos Coordinate, end_pos Coordinate, curr_pos Coordinate, answers_pt *[][]Move, current_path_pt *[]Move, min_price *int) {
	item := room_map[curr_pos.y][curr_pos.x]
	if item.item_type != FreeSpace {
		return
	}
	if curr_pos.x == end_pos.x && curr_pos.y == end_pos.y {
		answers := *answers_pt
		current_path := *current_path_pt
		ct_prise := getAnswerPrice(current_path)
		if *min_price == 0 || ct_prise <= *min_price {
			*min_price = ct_prise
		}
		answer_path := append([]Move{}, current_path...)
		answers = append(answers, answer_path)
		*answers_pt = answers
		return
	}
	item.is_locked = true
	room_map[curr_pos.y][curr_pos.x] = item
	moves := []Move{Left, Right, Up, Down}
	room_height := len(room_map)
	room_width := len(room_map[0])
	for _, m := range moves {
		next_pos := getNextCoordinate(curr_pos, m)
		if next_pos.y < 0 || next_pos.x < 0 || next_pos.y >= room_height || next_pos.x >= room_width {
			continue
		}
		n_item := room_map[next_pos.y][next_pos.x]
		if isValidItem(n_item) {
			current_path := *current_path_pt
			current_path = append(current_path, m)
			*current_path_pt = current_path
			c_price := getAnswerPrice(current_path)
			if (n_item.min_price == 0 || c_price < n_item.min_price) && (*min_price == 0 || c_price < *min_price) {
				n_item.min_price = c_price
				room_map[next_pos.y][next_pos.x] = n_item
				findPath(room_map, start_pos, end_pos, next_pos, answers_pt, current_path_pt, min_price)
			}
			current_path = *current_path_pt
			if len(current_path) > 0 {
				current_path = current_path[:len(current_path)-1]
				*current_path_pt = current_path
			}
		}

	}
	if item.item_type == FreeSpace {
		item.is_locked = false
		room_map[curr_pos.y][curr_pos.x] = item
	}
}

func findAnyPath(room_map [][]LabyrintItem, start_pos Coordinate, end_pos Coordinate, curr_pos Coordinate, answers_pt *[][]Move, current_path_pt *[]Move, min_price *int) {
	item := room_map[curr_pos.y][curr_pos.x]
	if item.item_type != FreeSpace {
		return
	}
	if curr_pos.x == end_pos.x && curr_pos.y == end_pos.y {
		answers := *answers_pt
		current_path := *current_path_pt
		ct_prise := getAnswerPrice(current_path)
		if *min_price == 0 || ct_prise <= *min_price {
			*min_price = ct_prise
		}
		answer_path := append([]Move{}, current_path...)
		answers = append(answers, answer_path)
		*answers_pt = answers
		return
	}
	item.is_locked = true
	room_map[curr_pos.y][curr_pos.x] = item
	moves := []Move{Left, Right, Up, Down}
	room_height := len(room_map)
	room_width := len(room_map[0])
	for _, m := range moves {
		next_pos := getNextCoordinate(curr_pos, m)
		if next_pos.y < 0 || next_pos.x < 0 || next_pos.y >= room_height || next_pos.x >= room_width {
			continue
		}
		n_item := room_map[next_pos.y][next_pos.x]
		if isValidItem(n_item) {
			current_path := *current_path_pt
			current_path = append(current_path, m)
			*current_path_pt = current_path
			c_price := getAnswerPrice(current_path)
			if (n_item.min_price == 0 || c_price < n_item.min_price) && (*min_price == 0) {
				n_item.min_price = c_price
				room_map[next_pos.y][next_pos.x] = n_item
				findPath(room_map, start_pos, end_pos, next_pos, answers_pt, current_path_pt, min_price)
			}
			current_path = *current_path_pt
			if len(current_path) > 0 {
				current_path = current_path[:len(current_path)-1]
				*current_path_pt = current_path
			}
		}

	}
	if item.item_type == FreeSpace {
		item.is_locked = false
		room_map[curr_pos.y][curr_pos.x] = item
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

	map_width := 71
	map_height := 71
	broken_bits_count := 1024
	room_map := [][]LabyrintItem{}
	broken_bits := []Coordinate{}

	start := Coordinate{y: 0, x: 0}
	end := Coordinate{y: map_height - 1, x: map_width - 1}
	for y := 0; y < map_height; y++ {
		row_items := []LabyrintItem{}
		for x := 0; x < map_width; x++ {
			l_item := LabyrintItem{item_type: FreeSpace, is_locked: false}
			row_items = append(row_items, l_item)
		}
		room_map = append(room_map, row_items)
	}
	for scanner.Scan() {
		row := scanner.Text()
		coords := strings.Split(row, ",")
		x, err := strconv.Atoi(coords[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(coords[1])
		if err != nil {
			panic(err)
		}
		broken_bits = append(broken_bits, Coordinate{y: y, x: x})
	}
	for i := 0; i < broken_bits_count; i++ {
		broken_bit := broken_bits[i]
		room_map[broken_bit.y][broken_bit.x] = LabyrintItem{item_type: Wall, is_locked: true}
	}
	answers := [][]Move{}
	current_path := []Move{}
	min_price := 0
	start_pos := Coordinate{y: 0, x: 0}
	//printMap(room_map, start, current_path)
	findPath(room_map, start, end, start_pos, &answers, &current_path, &min_price)
	best_answer := getBestAnswer(answers)
	//printMap(room_map, start, best_answer)
	best_answer_price := getAnswerPrice(best_answer)
	fmt.Printf("Best score: %d\n", best_answer_price)
	bad_bit := 0
	for i := broken_bits_count; i < len(broken_bits)-1; i++ {
		broken_bit := broken_bits[i]
		room_map[broken_bit.y][broken_bit.x] = LabyrintItem{item_type: Wall, is_locked: true}
		answers2 := [][]Move{}
		current_path2 := []Move{}
		min_price2 := 0
		start_pos2 := Coordinate{y: 0, x: 0}
		cleanUpScores(room_map)
		findAnyPath(room_map, start, end, start_pos2, &answers2, &current_path2, &min_price2)
		if len(answers2) == 0 {
			bad_bit = i
			break
		}
	}
	fmt.Printf("Bad bit: %d,%d\n", broken_bits[bad_bit].x, broken_bits[bad_bit].y)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...127528 125528
}
