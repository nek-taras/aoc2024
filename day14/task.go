package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Coordinate struct {
}

type Velocity struct {
}

type Robot struct {
	x  int
	y  int
	dx int
	dy int
}

func (r *Robot) Move(times int, map_wdth int, map_height int) {
	new_x := r.x + (r.dx*times)%map_wdth
	new_y := r.y + (r.dy*times)%map_height
	if new_x > map_wdth-1 {
		new_x = new_x % map_wdth
	} else if new_x < 0 {
		new_x = map_wdth + new_x

	}
	if new_y > map_height-1 {
		new_y = new_y % map_height
	} else if new_y < 0 {
		new_y = map_height + new_y

	}
	r.x = new_x
	r.y = new_y
}

func (r Robot) GetQuadrant(map_wdth int, map_height int) int {
	q_width := map_wdth / 2
	q_height := map_height / 2
	if r.x >= 0 && r.x < q_width {
		if r.y >= 0 && r.y < q_height {
			return 1
		} else if r.y > q_height {
			return 3
		}

	} else if r.x > q_width {
		if r.y >= 0 && r.y < q_height {
			return 2

		} else if r.y > q_height {
			return 4
		}

	}
	return 0
}

func getSqaureNeighbours(robots []Robot) (int, int) {
	width := 0
	height := 0

	x_s := make(map[int][]int)
	y_s := make(map[int][]int)

	for _, r := range robots {
		x, ok := x_s[r.x]
		if !ok {
			x = []int{}
		}
		x = append(x, r.y)
		x_s[r.x] = x
		y, ok := y_s[r.y]
		if !ok {
			y = []int{}
		}
		y = append(y, r.x)
		y_s[r.y] = y
	}

	cur_width := 1
	for _, v := range x_s {
		sort.Ints(v)
		r_s := v[0]
		for i := 1; i < len(v); i++ {
			r_n := v[i]
			if r_s == r_n || r_s == r_n-1 {
				cur_width++
			} else {
				width = max(width, cur_width)
				cur_width = 1
			}
			r_s = r_n
		}
		width = max(width, cur_width)
		cur_height := 1
		for _, v := range y_s {
			sort.Ints(v)
			r_s := v[0]
			for i := 1; i < len(v); i++ {
				r_n := v[i]
				if r_s == r_n || r_s == r_n-1 {
					cur_height++
				} else {
					height = max(height, cur_height)
					cur_height = 1
				}
				r_s = r_n
			}
			height = max(height, cur_height)
		}
	}
	return width, height
}

func createIage(robots []Robot, img_width int, img_height int, k int) {
	col_red := color.RGBA{255, 0, 0, 255}   // Red
	col_green := color.RGBA{0, 255, 0, 255} // Green
	img := image.NewRGBA(image.Rect(0, 0, img_width, img_height))
	f, err := os.Create(fmt.Sprintf("images/%d.png", k))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for y := 0; y < img_height; y++ {
		for x := 0; x < img_width; x++ {
			img.Set(x, y, col_green)
		}
	}

	for _, r := range robots {
		img.Set(r.x, r.y, col_red)
	}
	png.Encode(f, img)

}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
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
	robots := []Robot{}
	map_width := 101
	map_height := 103
	moves := 100
	for scanner.Scan() {
		row := scanner.Text()
		parts := strings.Split(row, " ")
		r_coords := strings.Split(strings.Split(parts[0], "=")[1], ",")
		x, err := strconv.Atoi(r_coords[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(r_coords[1])
		if err != nil {
			panic(err)
		}
		r_vel := strings.Split(strings.Split(parts[1], "=")[1], ",")
		dx, err := strconv.Atoi(r_vel[0])
		if err != nil {
			panic(err)
		}
		dy, err := strconv.Atoi(r_vel[1])
		if err != nil {
			panic(err)
		}
		robot := Robot{x: x, y: y, dx: dx, dy: dy}
		robots = append(robots, robot)
	}
	quadrant_robots := make(map[int]int)

	for _, r := range robots {
		r.Move(moves, map_width, map_height)
		r_q := r.GetQuadrant(map_width, map_height)
		if r_q != 0 {
			quadrant_robots[r_q]++
		}
	}
	total_mul := 1
	for _, v := range quadrant_robots {
		total_mul *= v
	}
	fmt.Printf("Total scode: %d\n", total_mul)
	k := 1
	for {
		for i, r := range robots {
			r.Move(1, map_width, map_height)
			robots[i] = r
		}
		r_w, r_h := getSqaureNeighbours(robots)
		if r_w > 30 && r_h > 30 {
			createIage(robots, map_width, map_height, k)
			fmt.Printf("Seconds to star: %d\n", k)
			break
		}
		k++

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
