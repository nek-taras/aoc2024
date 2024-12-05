package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

type Rule struct {
	lower  []int
	higher []int
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

	manuals := [][]int{}
	rules := make(map[int]Rule)

	for scanner.Scan() {
		row := scanner.Text()
		if strings.Contains(row, "|") {
			rules_s := strings.Split(row, "|")
			lower, err := strconv.Atoi(rules_s[0])
			if err != nil {
				panic(err)
			}
			higher, err := strconv.Atoi(rules_s[1])
			if err != nil {
				panic(err)
			}
			r, ok := rules[lower]
			if !ok {
				r = Rule{lower: []int{}, higher: []int{}}
			}
			r.higher = append(r.higher, higher)
			rules[lower] = r

			r, ok = rules[higher]
			if !ok {
				r = Rule{lower: []int{}, higher: []int{}}
			}
			r.lower = append(r.lower, lower)
			rules[higher] = r
		} else if strings.Contains(row, ",") {
			// parse manuals
			man_s := strings.Split(row, ",")
			manual := []int{}
			for _, v := range man_s {
				m, err := strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
				manual = append(manual, m)
			}
			manuals = append(manuals, manual)
		}
	}
	middle_sum := 0
	middle_fixed_sum := 0

	for _, man := range manuals {
		correct := true
		for i := range man {
			el := man[i]
			lowest := man[:i]
			highest := man[i+1:]
			r, ok := rules[el]
			if ok {
				for _, l := range lowest {
					if slices.Contains(r.higher, l) {
						correct = false
						break
					}
				}
				for _, h := range highest {
					if slices.Contains(r.lower, h) {
						correct = false
						break
					}
				}
			}
		}
		if correct {
			m_l := len(man)
			idx := int(m_l / 2)
			middle_sum += man[idx]
		} else {
			sort.Slice(man, func(i, j int) bool {
				prev := man[i]
				cur := man[j]
				prev_r, ok := rules[prev]
				if ok {
					if slices.Contains(prev_r.lower, cur) {
						return false
					}
				}
				cur_r, ok := rules[cur]
				if ok {
					if slices.Contains(cur_r.higher, prev) {
						return false
					}

				}
				return true
			})
			m_l := len(man)
			idx := int(m_l / 2)
			middle_fixed_sum += man[idx]
		}
	}
	fmt.Printf("Middle sum:%d\n", middle_sum)
	fmt.Printf("Middle fixed sum:%d\n", middle_fixed_sum)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
