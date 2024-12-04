package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
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

func searchVertical(data [][]byte, word []byte) int {
	found := 0
	i_n := len(data[0])
	j_n := len(data)
	for j := 0; j < j_n; j++ {
		for i := 0; i < i_n; i++ {
			next_letter_i := i
			next_letter_j := j
			word_found := true
			for _, char := range word {
				if next_letter_j == j_n || next_letter_i == i_n {
					word_found = false
					break
				}
				check_letter := data[next_letter_j][next_letter_i]
				if char == check_letter {
					next_letter_i += 1
				} else {
					word_found = false
					break
				}
			}
			if word_found {
				found += 1
			}
		}
	}
	return found
}

func searchHorizontal(data [][]byte, word []byte) int {
	found := 0
	i_n := len(data[0])
	j_n := len(data)
	for j := 0; j < j_n; j++ {
		for i := 0; i < i_n; i++ {
			next_letter_i := i
			next_letter_j := j
			word_found := true
			for _, char := range word {
				if next_letter_j == j_n || next_letter_i == i_n {
					word_found = false
					break
				}
				check_letter := data[next_letter_j][next_letter_i]
				if char == check_letter {
					next_letter_j += 1
				} else {
					word_found = false
					break
				}
			}
			if word_found {
				found += 1
			}
		}
	}
	return found
}

func searchDiagonal(data [][]byte, word []byte) int {
	found := 0
	i_n := len(data[0])
	j_n := len(data)
	for j := 0; j < j_n; j++ {
		for i := 0; i < i_n; i++ {
			next_letter_i := i
			next_letter_j := j
			word_found := true
			for _, char := range word {
				if next_letter_j == j_n || next_letter_i == i_n {
					word_found = false
					break
				}
				check_letter := data[next_letter_j][next_letter_i]
				if char == check_letter {
					next_letter_j += 1
					next_letter_i += 1
				} else {
					word_found = false
					break
				}
			}
			if word_found {
				found += 1
			}
		}
	}
	return found
}

func searchDiagonalReverse(data [][]byte, word []byte) int {
	found := 0
	i_n := len(data[0])
	j_n := len(data)
	for j := 0; j < j_n; j++ {
		for i := 0; i < i_n; i++ {
			next_letter_i := i
			next_letter_j := j
			word_found := true
			for _, char := range word {
				if next_letter_j < 0 || next_letter_i == i_n {
					word_found = false
					break
				}
				check_letter := data[next_letter_j][next_letter_i]
				if char == check_letter {
					next_letter_j -= 1
					next_letter_i += 1
				} else {
					word_found = false
					break
				}
			}
			if word_found {
				found += 1
			}
		}
	}
	return found
}

func searchCros(data [][]byte, word []byte) int {
	found := 0
	i_n := len(data[0])
	j_n := len(data)
	l_word := len(word)
	rew_word := make([]byte, len(word))

	for i, v := range word {
		rew_word[l_word-1-i] = v
	}
	for j := 0; j < j_n; j++ {
		for i := 0; i < i_n; i++ {
			if i+l_word > i_n || j+l_word > j_n {
				continue
			}
			first_plank := make([]byte, l_word)
			second_plank := make([]byte, l_word)
			for w_i := range word {
				first_plank[w_i] = data[j+w_i][i+w_i]
				second_plank[w_i] = data[j+l_word-1-w_i][i+w_i]
			}

			if (slices.Compare(first_plank, word) == 0 || slices.Compare(first_plank, rew_word) == 0) && (slices.Compare(second_plank, word) == 0 || slices.Compare(second_plank, rew_word) == 0) {
				found += 1
			}
		}
	}
	return found
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
	letters := [][]byte{}

	for scanner.Scan() {
		row := []byte(scanner.Text())
		letters = append(letters, row)
	}
	total_count := 0
	total_count += searchVertical(letters, []byte("XMAS"))
	total_count += searchVertical(letters, []byte("SAMX"))
	total_count += searchHorizontal(letters, []byte("XMAS"))
	total_count += searchHorizontal(letters, []byte("SAMX"))
	total_count += searchDiagonal(letters, []byte("XMAS"))
	total_count += searchDiagonal(letters, []byte("SAMX"))
	total_count += searchDiagonalReverse(letters, []byte("XMAS"))
	total_count += searchDiagonalReverse(letters, []byte("SAMX"))
	fmt.Printf("Count XMAS: %d\n", total_count)
	x_mas_count := searchCros(letters, []byte("MAS"))
	fmt.Printf("Count X-MAS: %d\n", x_mas_count)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
