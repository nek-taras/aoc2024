package main

import (
	"bufio"
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

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getDirectionAndDistance(x, y int) (string, int) {
	distance := x - y
	direction := "inc"
	if distance > 0 {
		direction = "dec"
	}
	return direction, AbsInt(x - y)
}

func checkReport(reports []int) bool {
	first := reports[0]
	next := reports[1]
	el_range := len(reports)
	initial_direction, distance := getDirectionAndDistance(first, next)

	if distance == 0 || distance > 3 {
		return false
	}

	for i := 1; i < el_range-1; i++ {
		first := reports[i]
		next := reports[i+1]
		direction, distance := getDirectionAndDistance(first, next)
		if direction != initial_direction || (distance == 0 || distance > 3) {
			return false
		}
	}
	return true

}

func removeIndex(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
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
	good_reports := 0
	good_with_tolerance := 0

	fmt.Printf("reading data\n")
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		reports := []int{}
		for i := range s {
			e, err := strconv.Atoi(s[i])
			if err != nil {
				panic(err)
			}
			reports = append(reports, e)
		}
		ok := checkReport(reports)
		if ok {
			good_reports += 1
			good_with_tolerance += 1
		} else {
			for i := 0; i < len(reports); i++ {
				fixed_reports := removeIndex(reports, i)
				ok_with_t := checkReport(fixed_reports)
				if ok_with_t {
					good_with_tolerance += 1
					break
				}
			}
		}
	}
	fmt.Printf("Good reports:%d\n", good_reports)
	fmt.Printf("Good reports with tolerance:%d\n", good_with_tolerance)
}
