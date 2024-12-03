package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func getClosestDo(does []int, pos int) int {
	if pos < does[0] {
		return 0
	}
	prev := 0
	for i := range does {
		if prev < pos && pos < does[i] {
			return prev
		}
		prev = does[i]
	}
	return prev
}

func getClosestDont(donts []int, pos int, max int) int {
	if pos < donts[0] {
		return -1
	}
	for i := 1; i < len(donts); i++ {
		prev := donts[i-1]
		next := donts[i]
		if prev < pos && pos < next {
			return prev
		}
	}
	return max
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Input file is missing.")
		os.Exit(1)
	}
	f_data, err := os.ReadFile(args[0])
	data := string(f_data)

	if err != nil {
		fmt.Println(err)
	}

	sum := 0
	sum_with_instructions := 0
	do_indexes := []int{}
	dont_indexes := []int{}
	r, _ := regexp.Compile(`mul\([0-9]+,[0-9]+\)`)
	do_r, _ := regexp.Compile(`do\(\)`)
	dont_r, _ := regexp.Compile(`don\'t\(\)`)
	does := do_r.FindAllStringIndex(data, -1)
	donts := dont_r.FindAllStringIndex(data, -1)

	for _, e := range does {
		do_indexes = append(do_indexes, e[0])
	}
	for _, e := range donts {
		dont_indexes = append(dont_indexes, e[0])
	}

	muls := r.FindAllStringIndex(data, -1)
	for _, e := range muls {
		pp := data[e[0]:e[1]]
		parts := strings.Split(pp, ",")
		first, err := strconv.Atoi(strings.Replace(parts[0], "mul(", "", 1))
		if err != nil {
			panic(err)
		}
		second, err := strconv.Atoi(strings.Replace(parts[1], ")", "", 1))
		if err != nil {
			panic(err)
		}
		closestDo := getClosestDo(do_indexes, e[0])
		closestDont := getClosestDont(dont_indexes, e[0], len(data))
		sum += first * second

		if closestDont < closestDo {
			sum_with_instructions += first * second
		}
	}
	fmt.Printf("Sum:%d\n", sum)
	fmt.Printf("Sum with instructions:%d\n", sum_with_instructions)
}
