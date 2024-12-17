package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
)

type ItemType int

const (
	Wall      ItemType = iota
	FreeSpace          = iota
	Start              = iota
	End                = iota
)

type Move int

const (
	Left  Move = iota
	Up         = iota
	Right      = iota
	Down       = iota
)

type LabyrintItem struct {
	item_type ItemType
	is_locked bool
	min_price map[Move]int
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func printMap(room_map [][]LabyrintItem, y int, x int, moves []Move) {

	move_map := make(map[string]Move)
	for i := 1; i < len(moves); i++ {
		m := moves[i]
		y, x = getNextCoordinate(y, x, m)
		coordinate := fmt.Sprintf("%d-%d", y, x)
		move_map[coordinate] = m
	}
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
				case Start:
					fmt.Printf("S")
				case End:
					fmt.Printf("E")
				}

			}

		}
		fmt.Printf("\n")
	}

}

func printSpotMap(room_map [][]LabyrintItem, y int, x int, spots map[string]Move) {

	for y, row := range room_map {
		for x, item := range row {
			coordinate := fmt.Sprintf("%d-%d", y, x)
			_, ok := spots[coordinate]
			if ok {
				fmt.Printf("0")

			} else {
				switch item.item_type {
				case Wall:
					fmt.Printf("#")
				case FreeSpace:
					fmt.Printf(".")
				case Start:
					fmt.Printf("S")
				case End:
					fmt.Printf("E")
				}

			}

		}
		fmt.Printf("\n")
	}

}

func getNextCoordinate(y int, x int, move Move) (int, int) {
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

	return y, x
}

func getAnswerPrice(answer []Move) int {
	price := 0
	move := answer[0]
	for i := 1; i < len(answer); i++ {
		next_move := answer[i]
		if move == next_move {
			price++
		} else {
			price += 1001
		}
		move = next_move
	}
	return price
}

func getBestAnswer(answers [][]Move) []Move {
	sort.Slice(answers, func(i, j int) bool {
		first_price := getAnswerPrice(answers[i])
		second_price := getAnswerPrice(answers[j])
		return first_price < second_price
	})
	return answers[0]
}

func isValidItem(item LabyrintItem) bool {
	return !item.is_locked
}

func findPath(room_map [][]LabyrintItem, start_y int, start_x int, y int, x int, answers_pt *[][]Move, current_path_pt *[]Move, min_price *int) {
	item := room_map[y][x]
	if item.item_type == End {
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
	room_map[y][x] = item
	moves := []Move{Up, Down, Left, Right}
	item_price := item.min_price
	for _, m := range moves {
		n_y, n_x := getNextCoordinate(y, x, m)
		n_item := room_map[n_y][n_x]
		if isValidItem(n_item) {
			current_path := *current_path_pt
			current_path = append(current_path, m)
			*current_path_pt = current_path
			c_price := getAnswerPrice(current_path)
			n_item_dir_price, ok := n_item.min_price[m]
			if !ok {
				n_item_dir_price = c_price
			}
			if c_price <= n_item_dir_price && (*min_price == 0 || c_price <= *min_price) {
				n_item.min_price[m] = c_price
				room_map[n_y][n_x] = n_item
				findPath(room_map, start_y, start_x, n_y, n_x, answers_pt, current_path_pt, min_price)
			}
			current_path = *current_path_pt
			item.min_price = item_price

			room_map[y][x] = item
			if len(current_path) > 0 {
				current_path = current_path[:len(current_path)-1]
				*current_path_pt = current_path
			}
		}

	}
	if item.item_type != Wall {
		item.is_locked = false
		room_map[y][x] = item
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
	room_map := [][]LabyrintItem{}
	start_y, start_x := 0, 0
	y := 0
	for scanner.Scan() {
		row := scanner.Text()
		row_items := []LabyrintItem{}
		for x, item_symbol := range row {
			switch item_symbol {
			case '#':
				l_item := LabyrintItem{item_type: Wall, is_locked: true, min_price: make(map[Move]int)}
				row_items = append(row_items, l_item)
			case '.':
				l_item := LabyrintItem{item_type: FreeSpace, is_locked: false, min_price: make(map[Move]int)}
				row_items = append(row_items, l_item)
			case 'S':
				l_item := LabyrintItem{item_type: Start, is_locked: false, min_price: make(map[Move]int)}
				row_items = append(row_items, l_item)
				start_y, start_x = y, x
			case 'E':
				l_item := LabyrintItem{item_type: End, is_locked: false, min_price: make(map[Move]int)}
				row_items = append(row_items, l_item)
			}
		}
		room_map = append(room_map, row_items)
		y++
	}
	answers := [][]Move{}
	current_path := []Move{Right}
	min_price := 0
	findPath(room_map, start_y, start_x, start_y, start_x, &answers, &current_path, &min_price)
	best_answer := getBestAnswer(answers)
	//printMap(room_map, start_y, start_x, best_answer)
	spot_map := make(map[string]Move)
	best_answer_price := getAnswerPrice(best_answer)
	coordinate := fmt.Sprintf("%d-%d", start_y, start_x)
	spot_map[coordinate] = Right
	for _, moves := range answers {
		y, x := start_y, start_x
		if getAnswerPrice(moves) == best_answer_price {
			for i := 1; i < len(moves); i++ {
				m := moves[i]
				y, x = getNextCoordinate(y, x, m)
				coordinate := fmt.Sprintf("%d-%d", y, x)
				spot_map[coordinate] = m
			}
		}

	}

	fmt.Printf("Best score: %d\n", best_answer_price)
	//printSpotMap(room_map, start_y, start_x, spot_map)
	fmt.Printf("Best spots: %d\n", len(spot_map))
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...127528 125528
}
