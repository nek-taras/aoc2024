package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

type Coordinate struct {
	y int
	x int
}

type TrialHead struct {
	start Coordinate
	end   Coordinate
}

func findAPath(tpg_map [][]int, check_coordinate Coordinate, nextValue int, prev_points []Coordinate, found_pathes *[][]Coordinate) {

	if tpg_map[check_coordinate.y][check_coordinate.x] == nextValue {
		points := append(make([]Coordinate, 0, len(prev_points)), prev_points...)
		points = append(points, check_coordinate)
		// found path
		if nextValue == 9 {
			*found_pathes = append(*found_pathes, points)
		} else {
			map_width := len(tpg_map[0])
			map_height := len(tpg_map)
			cc_top := Coordinate{x: check_coordinate.x, y: check_coordinate.y - 1}
			cc_bottom := Coordinate{x: check_coordinate.x, y: check_coordinate.y + 1}
			cc_left := Coordinate{x: check_coordinate.x - 1, y: check_coordinate.y}
			cc_right := Coordinate{x: check_coordinate.x + 1, y: check_coordinate.y}
			if cc_top.x >= 0 && cc_top.x < map_width && cc_top.y >= 0 && cc_top.y < map_height {
				findAPath(tpg_map, cc_top, nextValue+1, points, found_pathes)
			}
			if cc_bottom.x >= 0 && cc_bottom.x < map_width && cc_bottom.y >= 0 && cc_bottom.y < map_height {
				findAPath(tpg_map, cc_bottom, nextValue+1, points, found_pathes)
			}
			if cc_left.x >= 0 && cc_left.x < map_width && cc_left.y >= 0 && cc_left.y < map_height {
				findAPath(tpg_map, cc_left, nextValue+1, points, found_pathes)
			}
			if cc_right.x >= 0 && cc_right.x < map_width && cc_right.y >= 0 && cc_right.y < map_height {
				findAPath(tpg_map, cc_right, nextValue+1, points, found_pathes)
			}
		}
	} else {
		return
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
	tpg_map := [][]int{}
	for scanner.Scan() {
		row := scanner.Text()
		tpg_map_row := []int{}
		for _, el := range row {
			m, err := strconv.Atoi(string(el))
			if err != nil {
				panic(err)
			}
			tpg_map_row = append(tpg_map_row, m)
		}
		tpg_map = append(tpg_map, tpg_map_row)
	}
	found_pathes := [][]Coordinate{}
	head_scores := make(map[Coordinate]int)

	for j := range tpg_map {
		for i := range tpg_map[j] {
			if tpg_map[j][i] == 0 {
				prev_points := []Coordinate{}
				findAPath(tpg_map, Coordinate{x: i, y: j}, 0, prev_points, &found_pathes)
			}
		}
	}
	ths := make(map[TrialHead]bool)
	trailheads_sum := 0
	trailheads_sum_second_part := 0
	head_scores_second_part := make(map[Coordinate]int)
	for _, path := range found_pathes {
		th := TrialHead{start: path[0], end: path[len(path)-1]}
		_, ok := ths[th]
		h_score_second, ok := head_scores_second_part[path[0]]
		if !ok {
			h_score_second = 0
		}
		h_score_second++
		head_scores_second_part[path[0]] = h_score_second
		if ok {
			continue
		}
		ths[th] = true
		h_score, ok := head_scores[path[0]]
		if !ok {
			h_score = 0
		}
		h_score++
		head_scores[path[0]] = h_score
	}
	for _, score := range head_scores {
		trailheads_sum += score
	}
	for _, score := range head_scores_second_part {
		trailheads_sum_second_part += score
	}
	fmt.Printf("Trailhead sum: %d\n", trailheads_sum)
	fmt.Printf("Trailhead sum second: %d\n", trailheads_sum_second_part)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
