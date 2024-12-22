package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func mix(a, b int) int {
	return a ^ b
}

func prune(s int) int {
	return s % 16777216
}

func generateNext(secret int) int {
	out := 0
	//In particular, each buyer's secret number evolves into the next secret number in the sequence via the following process:

	//Calculate the result of multiplying the secret number by 64. Then, mix this result into the secret number. Finally, prune the secret number.
	//Calculate the result of dividing the secret number by 32. Round the result down to the nearest integer.
	//Then, mix this result into the secret number. Finally, prune the secret number.
	//Calculate the result of multiplying the secret number by 2048. Then, mix this result into the secret number. Finally, prune the secret number.
	//Each step of the above process involves mixing and pruning:

	//To mix a value into the secret number, calculate the bitwise XOR of the given value and the secret number.
	//Then, the secret number becomes the result of that operation. (If the secret number is 42 and you were to mix 15 into the secret number,
	//the secret number would become 37.)
	//To prune the secret number, calculate the value of the secret number modulo 16777216.
	//Then, the secret number becomes the result of that operation.
	//(If the secret number is 100000000 and you were to prune the secret number, the secret number would become 16113920.)
	out = secret * 64
	secret = mix(out, secret)
	secret = prune(secret)
	out = secret / 32
	secret = mix(out, secret)
	secret = prune(secret)
	out = secret * 2048
	secret = mix(out, secret)
	secret = prune(secret)
	return secret

}

type PriceChange struct {
	price    int
	sequence []int
}

func getSequenceChange(secret, cycles int, sequence_prices map[string]int) {
	first_price := secret % 10
	price_change := PriceChange{}
	buyer_sequnce_prices := make(map[string]int)
	for i := 0; i < cycles; i++ {
		secret = generateNext(secret)
		secret_price := secret % 10
		change := secret_price - first_price
		first_price = secret_price
		price_change.sequence = append(price_change.sequence, change)
		price_change.price = secret_price
		if i > 2 {
			key := ""
			for _, s := range price_change.sequence {
				key = fmt.Sprintf("%s|%d", key, s)
			}
			_, ok := buyer_sequnce_prices[key]
			if !ok {
				buyer_sequnce_prices[key] = price_change.price
				sequence_price := sequence_prices[key]
				sequence_price += price_change.price
				sequence_prices[key] = sequence_price
			}
			price_change.sequence = price_change.sequence[1:]
		}
	}
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
	secrets := []int{}
	for scanner.Scan() {
		row := scanner.Text()
		secret, err := strconv.Atoi(row)
		if err != nil {
			panic(err)
		}
		secrets = append(secrets, secret)
	}
	secret_cycles := 2000
	result := 0
	//sequences_changes := make(map[string]int)
	for _, secret := range secrets {
		for i := 0; i < secret_cycles; i++ {
			secret = generateNext(secret)
		}
		result += secret

	}
	sequences_changes := make(map[string]int)
	for _, secret := range secrets {
		getSequenceChange(secret, secret_cycles, sequences_changes)
	}
	fmt.Printf("First reuslt: %d\n", result)
	max_change_price := 0
	for _, v := range sequences_changes {
		max_change_price = max(max_change_price, v)
	}
	fmt.Printf("Max bananas: %d\n", max_change_price)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}
