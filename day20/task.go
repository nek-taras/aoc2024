package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
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
	item_type        ItemType
	price_from_start int
	y                int
	x                int
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func printMap(room_map [][]LabyrintItem) {

	for _, row := range room_map {
		for _, item := range row {
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

func isValidPathItem(item LabyrintItem) bool {
	return item.item_type == FreeSpace || item.item_type == End
}
func getNeighbours(room_map [][]LabyrintItem, y int, x int) []LabyrintItem {
	moves := []Move{Up, Down, Left, Right}
	height := len(room_map)
	width := len(room_map[0])
	neighbours := []LabyrintItem{}
	for _, m := range moves {
		n_y, n_x := getNextCoordinate(y, x, m)
		if n_y < 0 || n_y >= height || n_x < 0 || n_x >= width {
			continue
		}
		n_item := room_map[n_y][n_x]
		neighbours = append(neighbours, n_item)
	}
	return neighbours
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func reprPath(fp LabyrintItem, lp LabyrintItem) string {
	fp_distance := math.Sqrt(math.Pow(float64(fp.y), 2) + math.Pow(float64(fp.x), 2))
	lp_distance := math.Sqrt(math.Pow(float64(lp.y), 2) + math.Pow(float64(lp.x), 2))
	if fp_distance >= lp_distance {
		return fmt.Sprintf("sx%dsy%dex%dey%d", lp.x, lp.y, fp.x, fp.y)

	}
	return fmt.Sprintf("sx%dsy%dex%dey%d", fp.x, fp.y, lp.x, lp.y)
}

func getShortcutsMore(first_path []LabyrintItem, chets_available int) map[string]int {
	cheets := make(map[string]int)
	for i, item := range first_path {
		for j := i + 1; j < len(first_path); j++ {
			s_item := first_path[j]
			distance := AbsInt(item.y-s_item.y) + AbsInt(item.x-s_item.x)
			price_cut := AbsInt(item.price_from_start-s_item.price_from_start) - distance
			if price_cut > 0 && distance <= chets_available {
				cut_key := reprPath(item, s_item)
				_, ok := cheets[cut_key]
				if !ok {
					ticks := distance
					price := AbsInt(item.price_from_start-s_item.price_from_start) - ticks
					cheets[cut_key] = price
				}

			}

		}
	}
	return cheets
}

func getShortcuts(room_map [][]LabyrintItem) []int {
	shortcuts := []int{}
	for y, row := range room_map {
		for x, item := range row {
			if item.item_type == Wall {
				neighbours := getNeighbours(room_map, y, x)
				//
				if len(neighbours) < 4 {
					continue
				}
				race_naighbours := []LabyrintItem{}
				for _, n := range neighbours {
					if n.item_type != Wall {
						race_naighbours = append(race_naighbours, n)
					}
				}
				if len(race_naighbours) < 2 {
					continue
				}
				cut_ways := []int{}
				for i := 0; i < len(race_naighbours); i++ {
					for j := 0; j < len(race_naighbours); j++ {
						if i == j {
							continue
						}
						cut_ways = append(cut_ways, AbsInt(race_naighbours[i].price_from_start-race_naighbours[j].price_from_start)-2)
					}
				}
				if len(cut_ways) > 0 {
					cur_max := cut_ways[0]
					for i := 1; i < len(cut_ways); i++ {
						cur_max = max(cur_max, cut_ways[i])
					}
					if cur_max > 0 {
						shortcuts = append(shortcuts, cur_max)
					}
				}

			}
		}
	}
	return shortcuts
}

func getFirstPath(room_map [][]LabyrintItem) []LabyrintItem {
	first_path := []LabyrintItem{}
	for _, row := range room_map {
		for _, item := range row {
			if item.item_type != Wall {
				first_path = append(first_path, item)
			}
		}
	}
	sort.Slice(first_path, func(i, j int) bool {
		f := first_path[i]
		n := first_path[j]
		return f.price_from_start < n.price_from_start
	})
	return first_path
}

func findPath(room_map [][]LabyrintItem, start_y int, start_x int, y int, x int, prev_price int) int {
	item := room_map[y][x]
	if item.item_type == End {
		return prev_price
	}
	moves := []Move{Up, Down, Left, Right}
	for _, m := range moves {
		n_y, n_x := getNextCoordinate(y, x, m)
		n_item := room_map[n_y][n_x]
		if isValidPathItem(n_item) && n_item.price_from_start == 0 {
			n_item.price_from_start = prev_price + 1
			room_map[n_y][n_x] = n_item
			return findPath(room_map, start_y, start_x, n_y, n_x, n_item.price_from_start)

		}

	}
	return item.price_from_start
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
				l_item := LabyrintItem{item_type: Wall, y: y, x: x}
				row_items = append(row_items, l_item)
			case '.':
				l_item := LabyrintItem{item_type: FreeSpace, y: y, x: x}
				row_items = append(row_items, l_item)
			case 'S':
				l_item := LabyrintItem{item_type: Start, y: y, x: x}
				row_items = append(row_items, l_item)
				start_y, start_x = y, x
			case 'E':
				l_item := LabyrintItem{item_type: End, y: y, x: x}
				row_items = append(row_items, l_item)
			}
		}
		room_map = append(room_map, row_items)
		y++
	}
	original_time := findPath(room_map, start_y, start_x, start_y, start_x, 0)
	//printMap(room_map)

	fmt.Printf("Original time: %d\n", original_time)
	shortcuts := getShortcuts(room_map)
	count_gte_100 := 0
	for _, s := range shortcuts {
		if s >= 100 {
			count_gte_100++
		}
	}
	fmt.Printf("Shortcuts gte 100: %d\n", count_gte_100)
	first_path := getFirstPath(room_map)

	cheets := getShortcutsMore(first_path, 20)

	second_count := 0
	for _, v := range cheets {
		if v >= 100 {
			second_count += 1
		}
	}
	fmt.Printf("Shortcuts 20 gte 100: %d\n", second_count)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...127528 125528 48573
}
