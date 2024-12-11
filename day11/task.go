package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func splitStone(stone string) (string, string) {
	s_half := int(len(stone) / 2)
	stone_bytes := []byte(stone)
	first_s := strings.TrimLeft(string(stone_bytes[:s_half]), "0")
	if first_s == "" {
		first_s = "0"
	}
	second_s := strings.TrimLeft(string(stone_bytes[s_half:]), "0")
	if second_s == "" {
		second_s = "0"
	}
	return first_s, second_s

}

func countStoneIterations(stone string, iterations int, existing_counts map[string]map[int]int) int {
	existing_count, ok := existing_counts[stone]
	if !ok {
		existing_count = make(map[int]int)
	}
	count, ok := existing_count[iterations]
	if ok {
		return count
	}
	stone_count := 0
	if iterations == 1 {
		if len(stone)%2 == 0 {
			return 2
		} else {
			return 1
		}
	} else {
		stone_int, err := strconv.Atoi(stone)
		if err != nil {
			panic(err)
		}
		if stone_int == 0 {
			stone_count += countStoneIterations("1", iterations-1, existing_counts)
		} else if len(stone)%2 == 0 {
			f_n, s_n := splitStone(stone)
			stone_count += countStoneIterations(f_n, iterations-1, existing_counts)
			stone_count += countStoneIterations(s_n, iterations-1, existing_counts)
		} else {
			stone_count += countStoneIterations(fmt.Sprint(stone_int*2024), iterations-1, existing_counts)
		}
	}
	existing_count[iterations] = stone_count
	existing_counts[stone] = existing_count
	return stone_count
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Input file is missing.")
		os.Exit(1)
	}

	content, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Print(err)
	}
	stones := strings.Split(string(content), " ")
	iterations := 75
	total_stones_count := 0
	existing_counts := make(map[string]map[int]int)
	for _, stone := range stones {
		total_stones_count += countStoneIterations(stone, iterations, existing_counts)
	}
	fmt.Printf("Stones: %d\n", total_stones_count)

	// ...
}
