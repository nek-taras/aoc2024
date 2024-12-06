package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

type Direction struct {
	x int
	y int
}

func checkIfLoop(room_map [][]byte, guard_x int, guard_y int) bool {
	guard_x_direction := 0
	guard_y_direction := -1
	room_width := len(room_map[0])
	room_height := len(room_map)
	cell_directions := make(map[Direction][]Direction)
	for {
		next_cell_x := guard_x + guard_x_direction
		next_cell_y := guard_y + guard_y_direction
		if next_cell_x < 0 || next_cell_x >= room_width || next_cell_y < 0 || next_cell_y >= room_height {
			return false
		}
		next_cell := room_map[next_cell_y][next_cell_x]
		if next_cell == '#' {
			// rotate
			if guard_x_direction == 0 && guard_y_direction == -1 {
				guard_x_direction = 1
				guard_y_direction = 0
			} else if guard_x_direction == 1 && guard_y_direction == 0 {
				guard_x_direction = 0
				guard_y_direction = 1
			} else if guard_x_direction == 0 && guard_y_direction == 1 {
				guard_x_direction = -1
				guard_y_direction = 0
			} else if guard_x_direction == -1 && guard_y_direction == 0 {
				guard_x_direction = 0
				guard_y_direction = -1
			}
		} else {
			// set current direction
			curr_position := Direction{x: guard_x, y: guard_y}
			directions, ok := cell_directions[curr_position]
			if !ok {
				cell_directions[curr_position] = []Direction{{x: guard_x_direction, y: guard_y_direction}}
			} else {
				for _, dir := range directions {
					if dir.x == guard_x_direction && dir.y == guard_y_direction {
						return true
					}
				}
				directions = append(directions, Direction{x: guard_x_direction, y: guard_y_direction})
			}
			cell_directions[curr_position] = directions
			guard_x = next_cell_x
			guard_y = next_cell_y
		}
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
	room_map := [][]byte{}
	guard_x, guard_y := 0, 0
	y := 0
	guard_found := false
	for scanner.Scan() {
		row := []byte(scanner.Text())
		room_map = append(room_map, row)
		if !guard_found {
			for x, r := range row {
				if r == '^' {
					guard_found = true
					guard_x, guard_y = x, y
					break
				}
			}
		}
		y += 1
	}
	guard_start_x, guard_start_y := guard_x, guard_y
	room_width := len(room_map[0])
	room_height := len(room_map)
	guard_x_direction := 0
	guard_y_direction := -1
	cell_visited := 0
	for {
		next_cell_x := guard_x + guard_x_direction
		next_cell_y := guard_y + guard_y_direction
		if next_cell_x < 0 || next_cell_x >= room_width || next_cell_y < 0 || next_cell_y >= room_height {
			room_map[guard_y][guard_x] = 'X'
			break
		}
		next_cell := room_map[next_cell_y][next_cell_x]
		if next_cell == '.' || next_cell == 'X' {
			room_map[guard_y][guard_x] = 'X'
			guard_x = next_cell_x
			guard_y = next_cell_y
		} else if next_cell == '#' {
			// rotate
			if guard_x_direction == 0 && guard_y_direction == -1 {
				guard_x_direction = 1
				guard_y_direction = 0
			} else if guard_x_direction == 1 && guard_y_direction == 0 {
				guard_x_direction = 0
				guard_y_direction = 1
			} else if guard_x_direction == 0 && guard_y_direction == 1 {
				guard_x_direction = -1
				guard_y_direction = 0
			} else if guard_x_direction == -1 && guard_y_direction == 0 {
				guard_x_direction = 0
				guard_y_direction = -1
			}
		}
	}
	successful_obsticles_count := 0
	for y := range room_map {
		row := room_map[y]
		for x, cell := range row {
			if cell == 'X' {
				cell_visited += 1
				duplicate_map := make([][]byte, len(room_map))
				for i := range room_map {
					duplicate_map[i] = make([]byte, len(room_map[i]))
					copy(duplicate_map[i], room_map[i])
				}
				duplicate_map[y][x] = '#'
				if checkIfLoop(duplicate_map, guard_start_x, guard_start_y) {
					successful_obsticles_count += 1
				}
			}
		}
	}
	fmt.Printf("Cell visited: %d\n", cell_visited)
	fmt.Printf("Successful obsticle placed: %d\n", successful_obsticles_count)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
