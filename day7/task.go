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

func checkResult(res int, parts []int) bool {
	if len(parts) == 2 {
		return (parts[0]+parts[1]) == res || (parts[0]*parts[1]) == res
	} else {
		first_mul := []int{parts[0] * parts[1]}
		first_add := []int{parts[0] + parts[1]}
		first_mul = append(first_mul, parts[2:]...)
		first_add = append(first_add, parts[2:]...)
		return checkResult(res, first_mul) || checkResult(res, first_add)
	}
}

func checkResultWithConcrete(res int, parts []int) bool {
	if len(parts) == 2 {
		add_res := parts[0] + parts[1]
		mul_res := parts[0] * parts[1]
		concrete_res, err := strconv.Atoi(fmt.Sprintf("%d%d", parts[0], parts[1]))
		if err != nil {
			panic(err)
		}
		return add_res == res || mul_res == res || concrete_res == res
	} else {
		first_mul := []int{parts[0] * parts[1]}
		first_add := []int{parts[0] + parts[1]}
		first_mul = append(first_mul, parts[2:]...)
		first_add = append(first_add, parts[2:]...)
		concrete, err := strconv.Atoi(fmt.Sprintf("%d%d", parts[0], parts[1]))
		if err != nil {
			panic(err)
		}
		first_concrete := []int{concrete}
		first_concrete = append(first_concrete, parts[2:]...)
		return checkResultWithConcrete(res, first_mul) || checkResultWithConcrete(res, first_add) || checkResultWithConcrete(res, first_concrete)
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
	good_test_sum := 0
	good_test_with_concrete_sum := 0
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ": ")
		result, err := strconv.Atoi(s[0])
		if err != nil {
			panic(err)
		}
		parts := strings.Split(s[1], " ")
		nums := []int{}
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				panic(err)
			}
			nums = append(nums, num)
		}
		if checkResult(result, nums) {
			good_test_sum += result
		}
		if checkResultWithConcrete(result, nums) {
			good_test_with_concrete_sum += result
		}
	}
	fmt.Printf("Correct tests: %d\n", good_test_sum)
	fmt.Printf("Correct tests wiith concrete: %d\n", good_test_with_concrete_sum)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
