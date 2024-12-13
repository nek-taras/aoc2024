package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Coordinate struct {
	x int
	y int
}

type Game struct {
	buttonA Coordinate
	buttonB Coordinate
	reward  Coordinate
}

func move(cur_coordinate Coordinate, game Game, a_count int, b_count int) Coordinate {
	cur_coordinate.x += game.buttonA.x*a_count + game.buttonB.x*b_count
	cur_coordinate.y += game.buttonA.y*a_count + game.buttonB.y*b_count
	return cur_coordinate
}

func getPrice(a_count int, b_count int) int {
	return a_count*3 + b_count
}

type Result int

const (
	Lower  Result = iota
	Equal         = iota
	Higher        = iota
)

func checkSolution(cur_coordinte Coordinate, reward Coordinate) Result {
	if cur_coordinte.x == reward.x && cur_coordinte.y == reward.y {
		return Equal
	}
	if cur_coordinte.x > reward.x {
		return Higher
	}
	if cur_coordinte.y > reward.y {
		return Higher
	}
	return Lower
}

func playGame(game Game, is_limited bool, reward_add int) (bool, int) {
	game.reward.x += reward_add
	game.reward.y += reward_add
	b_f := (game.reward.y*game.buttonA.x - game.reward.x*game.buttonA.y) / (game.buttonB.y*game.buttonA.x - game.buttonB.x*game.buttonA.y)
	a_f := (game.reward.y*game.buttonB.x - game.reward.x*game.buttonB.y) / (game.buttonA.y*game.buttonB.x - game.buttonA.x*game.buttonB.y)
	b := int(b_f)
	a := (game.reward.x - b*game.buttonB.x) / game.buttonA.x
	if is_limited && (a > 100 || b > 100) {
		return false, 0
	}
	res := move(Coordinate{x: 0, y: 0}, game, a, b)
	if checkSolution(res, game.reward) == Equal {
		return true, getPrice(a, b)
	}
	a = int(a_f)
	b = (game.reward.x - a*game.buttonA.x) / game.buttonB.x
	res = move(Coordinate{x: 0, y: 0}, game, a, b)
	if checkSolution(res, game.reward) == Equal {
		return true, getPrice(a, b)
	}

	return false, 0
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
	games := []Game{}
	game_elements := 0
	butonA := Coordinate{}
	butonB := Coordinate{}
	reward := Coordinate{}
	//Button A: X+94, Y+34
	//Button B: X+22, Y+67
	//Prize: X=8400, Y=5400
	for scanner.Scan() {
		row := scanner.Text()
		if strings.HasPrefix(row, "Button A:") {
			parts := strings.Split(row, " ")
			x_s := strings.TrimRight(strings.Split(parts[2], "+")[1], ",")
			x, err := strconv.Atoi(x_s)
			if err != nil {
				panic(err)
			}
			y_s := strings.Split(parts[3], "+")[1]
			y, err := strconv.Atoi(y_s)
			if err != nil {
				panic(err)
			}
			butonA = Coordinate{x: x, y: y}
			game_elements++

		} else if strings.HasPrefix(row, "Button B:") {
			parts := strings.Split(row, " ")
			x_s := strings.TrimRight(strings.Split(parts[2], "+")[1], ",")
			x, err := strconv.Atoi(x_s)
			if err != nil {
				panic(err)
			}
			y_s := strings.Split(parts[3], "+")[1]
			y, err := strconv.Atoi(y_s)
			if err != nil {
				panic(err)
			}
			butonB = Coordinate{x: x, y: y}
			game_elements++

		} else if strings.HasPrefix(row, "Prize:") {
			parts := strings.Split(row, " ")
			x_s := strings.TrimRight(strings.Split(parts[1], "=")[1], ",")
			x, err := strconv.Atoi(x_s)
			if err != nil {
				panic(err)
			}
			y_s := strings.Split(parts[2], "=")[1]
			y, err := strconv.Atoi(y_s)
			if err != nil {
				panic(err)
			}
			reward = Coordinate{x: x, y: y}
			game_elements++
		}
		if game_elements == 3 {
			game := Game{buttonA: butonA, buttonB: butonB, reward: reward}
			games = append(games, game)
			game_elements = 0
		}
	}
	total_price := 0
	total_price_second_round := 0
	for _, game := range games {
		win, price := playGame(game, true, 0)
		if win {
			total_price += price
		}
		win, price = playGame(game, false, 10000000000000)
		if win {
			total_price_second_round += price
		}
	}
	fmt.Printf("Total price: %d\n", total_price)
	fmt.Printf("Total price second: %d\n", total_price_second_round)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
