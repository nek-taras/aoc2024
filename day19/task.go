package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func findPattern(target_matches map[string]int, available_towels []string, expected_pattern string) int {
	ans, ok := target_matches[expected_pattern]
	if ok {
		return ans
	} else {
		ans = 0
	}
	if len(expected_pattern) == 0 {
		return 1
	}
	for _, towel := range available_towels {
		if strings.HasPrefix(expected_pattern, towel) {
			ans += findPattern(target_matches, available_towels, strings.TrimPrefix(expected_pattern, towel))
		}
	}
	target_matches[expected_pattern] = ans

	return ans
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
	available_towels := []string{}
	expected_patterns := []string{}
	line_counter := 0

	for scanner.Scan() {
		row := scanner.Text()
		if line_counter == 0 {
			available_towels = strings.Split(row, ", ")
			// read towels
		} else if line_counter > 1 {
			expected_patterns = append(expected_patterns, row)
			//read patterns
		}
		line_counter++
	}
	target_matches := make(map[string]int)
	found_patterns := 0
	total_answers := 0
	for _, expected_pattern := range expected_patterns {

		answers := findPattern(target_matches, available_towels, expected_pattern)
		if answers > 0 {
			found_patterns++
		}
		total_answers += answers
	}

	fmt.Printf("Found patterns: %d\n", found_patterns)
	fmt.Printf("Found total patterns: %d\n", total_answers)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...127528 125528
}
