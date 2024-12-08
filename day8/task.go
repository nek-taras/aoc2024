package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func getDirection(a, b Coordinate) Direction {
	dx := b.x - a.x
	dy := b.y - a.y
	return Direction{dx: dx, dy: dy}
}

type Coordinate struct {
	y int
	x int
}

type Direction struct {
	dy int
	dx int
}

type AntenaPath struct {
	start     Coordinate
	end       Coordinate
	direction Direction
	points    []Coordinate
}

func buildAPath(start Coordinate, points []Coordinate, direction Direction) AntenaPath {
	path_points := []Coordinate{start}
	curr_point := start
	for _, point := range points {
		if getDirection(curr_point, point) == direction {
			curr_point = point
			path_points = append(path_points, point)
		} else {
			break
		}
	}
	return AntenaPath{start: start, end: curr_point, points: path_points, direction: direction}
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
	y := 0
	antena_map := make(map[byte][]Coordinate)
	antena_pathes := make(map[byte][]AntenaPath)
	whole_map := [][]byte{}

	for scanner.Scan() {
		s := []byte(scanner.Text())
		whole_map = append(whole_map, s)
		for x, v := range s {
			if v == '.' {
				continue
			} else {
				a_m, ok := antena_map[v]
				if !ok {
					antena_map[v] = []Coordinate{{x: x, y: y}}
				} else {
					a_m = append(a_m, Coordinate{x: x, y: y})
					antena_map[v] = a_m
				}
			}
		}
		y += 1
	}
	map_height := len(whole_map)
	map_width := len(whole_map[0])
	for antena, a_points := range antena_map {
		if len(a_points) == 1 {
			continue
		}
		for i := 0; i < len(a_points); i++ {
			start_point := a_points[i]
			for j := i + 1; j < len(a_points); j++ {

				next_point := a_points[j]
				direction := getDirection(start_point, next_point)
				pathes, ok := antena_pathes[antena]
				new_path := buildAPath(start_point, a_points[j:], direction)
				if ok {
					for _, path := range pathes {
						if path.direction == new_path.direction && slices.Contains(path.points, new_path.points[0]) && slices.Contains(path.points, new_path.points[1]) {
							continue
						}
					}
					pathes = append(pathes, new_path)
				} else {
					pathes = []AntenaPath{new_path}
				}
				antena_pathes[antena] = pathes
			}
		}
	}
	unique_antinodes := 0
	for _, pathes := range antena_pathes {
		for _, path := range pathes {
			start_an_x := path.start.x - path.direction.dx
			start_an_y := path.start.y - path.direction.dy
			end_an_x := path.end.x + path.direction.dx
			end_an_y := path.end.y + path.direction.dy
			if start_an_x >= 0 && start_an_x < map_width && start_an_y >= 0 && start_an_y < map_height {
				cur_map_val := whole_map[start_an_y][start_an_x]
				if cur_map_val != '#' {
					whole_map[start_an_y][start_an_x] = '#'
					unique_antinodes += 1
				}
			}
			if end_an_x >= 0 && end_an_x < map_width && end_an_y >= 0 && end_an_y < map_height {
				cur_map_val := whole_map[end_an_y][end_an_x]
				if cur_map_val != '#' {
					whole_map[end_an_y][end_an_x] = '#'
					unique_antinodes += 1
				}
			}

		}
	}
	antinode_with_resonants := 0
	for _, pathes := range antena_pathes {
		for _, path := range pathes {
			start_an_x := path.start.x - path.direction.dx
			start_an_y := path.start.y - path.direction.dy
			end_an_x := path.end.x + path.direction.dx
			end_an_y := path.end.y + path.direction.dy
			for start_an_x >= 0 && start_an_x < map_width && start_an_y >= 0 && start_an_y < map_height {
				cur_map_val := whole_map[start_an_y][start_an_x]
				if cur_map_val != '#' {
					whole_map[start_an_y][start_an_x] = '#'
					antinode_with_resonants += 1
				}
				start_an_x = start_an_x - path.direction.dx
				start_an_y = start_an_y - path.direction.dy
			}
			for end_an_x >= 0 && end_an_x < map_width && end_an_y >= 0 && end_an_y < map_height {
				cur_map_val := whole_map[end_an_y][end_an_x]
				if cur_map_val != '#' {
					whole_map[end_an_y][end_an_x] = '#'
					antinode_with_resonants += 1
				}
				end_an_x = end_an_x + path.direction.dx
				end_an_y = end_an_y + path.direction.dy
			}
			for _, point := range path.points {
				cur_map_val := whole_map[point.y][point.x]
				if cur_map_val != '#' {
					whole_map[point.y][point.x] = '#'
					antinode_with_resonants += 1
				}

			}

		}
	}

	fmt.Printf("Unique antinodes: %d\n", unique_antinodes)
	fmt.Printf("Unique antinodes with resonance: %d\n", unique_antinodes+antinode_with_resonants)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
