package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Input file is missing.")
		os.Exit(1)
	}
	fmt.Printf("opening %s\n", args[0])
	f, err := os.Open(args[0])

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	left := []int{}
	right := []int{}

	fmt.Printf("reading data\n")
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "   ")
		l, err := strconv.Atoi(s[0])
		if err != nil {
			panic(err)
		}
		left = append(left, l)
		r, err := strconv.Atoi(s[1])
		if err != nil {
			panic(err)
		}
		right = append(right, r)
	}
	fmt.Printf("Calculatin\n")
	sort.Ints(left)
	sort.Ints(right)
	distance := 0
	for i := range left {
		l := left[i]
		r := right[i]
		distance += AbsInt(l - r)

	}
	fmt.Printf("Distance sum is:%d\n", distance)
	r_length := len(right)
	l_min := 0
	similarity := 0
	for i := range left {
		l := left[i]
		ocurrences := 0
		for j := l_min; i < r_length; j++ {
			r := right[j]
			if r < l {
				l_min = j
			} else if r == l {
				ocurrences += 1
			} else {
				break
			}
		}
		similarity += l * ocurrences
	}

	fmt.Printf("Similarity is:%d\n", similarity)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
