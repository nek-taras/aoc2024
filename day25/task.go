package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type Key struct {
	pins []int
}

type Lock struct {
	pins []int
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func parseKey(kt []string) Key {
	key := Key{}
	for i := 0; i < 5; i++ {
		pin_sum := 0
		for j := 1; j < 6; j++ {
			if kt[j][i] == '#' {
				pin_sum++
			}
		}
		key.pins = append(key.pins, pin_sum)
	}
	return key
}

func parseLock(lt []string) Lock {
	lock := Lock{}
	for i := 0; i < 5; i++ {
		pin_sum := 0
		for j := 1; j < 6; j++ {
			if lt[j][i] == '#' {
				pin_sum++
			}
		}
		lock.pins = append(lock.pins, pin_sum)
	}
	return lock
}

func pickLock(key Key, lock Lock) bool {
	for i := 0; i < 5; i++ {
		if key.pins[i]+lock.pins[i] > 5 {
			return false
		}
	}
	return true
}

func getKeyLockCombos(keys []Key, locks []Lock) int {
	combos := 0
	for _, k := range keys {
		for _, l := range locks {
			if pickLock(k, l) {
				combos++
			}
		}
	}
	return combos
}

func Parsetext(text [][]string) ([]Key, []Lock) {
	keys := []Key{}
	locks := []Lock{}
	for _, t := range text {
		if t[0] == "#####" {
			locks = append(locks, parseLock(t))
		} else if t[0] == "....." {
			keys = append(keys, parseKey(t))
		}
	}

	return keys, locks

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

	text := [][]string{}

	scanner := bufio.NewScanner(f)
	curr_row := []string{}
	for scanner.Scan() {
		row := scanner.Text()
		if len(row) == 0 {
			text = append(text, curr_row)
			curr_row = []string{}

		} else {
			curr_row = append(curr_row, row)
		}
	}
	text = append(text, curr_row)
	keys, locks := Parsetext(text)
	combos := getKeyLockCombos(keys, locks)
	fmt.Printf("Key lock combos: %d\n", combos)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}
