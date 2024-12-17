package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	A      int
	B      int
	C      int
	output []int
}

//Combo operands 0 through 3 represent literal values 0 through 3.
//Combo operand 4 represents the value of register A.
//Combo operand 5 represents the value of register B.
//Combo operand 6 represents the value of register C.

func (m *Machine) getCombo(op int) int {
	if op == 6 {
		return m.C
	} else if op == 5 {
		return m.B
	} else if op == 4 {
		return m.A
	}
	return op
}

func (m *Machine) runProgram(program []Command) {
	running := true
	idx := 0
	for running {
		c := program[idx]
		switch c.instruction {
		case adv:
			//The adv instruction (opcode 0) performs division. The numerator is the value in the A register.
			//The denominator is found by raising 2 to the power of the instruction's combo operand.
			//(So, an operand of 2 would divide A by 4 (2^2); an operand of 5 would divide A by 2^B.)
			//The result of the division operation is truncated to an integer and then written to the A register.
			m.A = m.A / int(math.Pow(2, float64(m.getCombo(c.operand))))
		case bxl:
			//The bxl instruction (opcode 1) calculates the bitwise XOR of register B and the instruction's literal operand,
			// then stores the result in register B.
			m.B = m.B ^ c.operand
		case bst:
			//The bst instruction (opcode 2) calculates the value of its combo operand modulo 8 (thereby keeping only its lowest 3 bits),
			//then writes that value to the B register.
			m.B = (m.getCombo(c.operand) % 8) & 0b111
		case jnz:
			//The jnz instruction (opcode 3) does nothing if the A register is 0. However, if the A register is not zero,
			//it jumps by setting the instruction pointer to the value of its literal operand; if this instruction jumps, the instruction pointer is not increased by 2 after this instruction.
			if m.A != 0 {
				idx = c.operand
				continue
			}
		case bxc:
			//The bxc instruction (opcode 4) calculates the bitwise XOR of register B and register C,
			// then stores the result in register B. (For legacy reasons, this instruction reads an operand but ignores it.)
			m.B = m.B ^ m.C
		case out:
			//The out instruction (opcode 5) calculates the value of its combo operand modulo 8, then outputs that value.
			// (If a program outputs multiple values, they are separated by commas.)
			m.output = append(m.output, m.getCombo(c.operand)%8)
		case bdv:
			//The bdv instruction (opcode 6) works exactly like the adv instruction except that the result is stored in the B register.
			//(The numerator is still read from the A register.)
			m.B = m.A / int(math.Pow(2, float64(m.getCombo(c.operand))))
		case cdv:
			//The cdv instruction (opcode 7) works exactly like the adv instruction except that the result is stored in the C register.
			//(The numerator is still read from the A register.)
			m.C = m.A / int(math.Pow(2, float64(m.getCombo(c.operand))))
		}
		idx++
		if idx >= len(program) {
			running = false
		}
	}
}

type Instruction int

const (
	adv Instruction = iota
	bxl             = iota
	bst             = iota
	jnz             = iota
	bxc             = iota
	out             = iota
	bdv             = iota
	cdv             = iota
)

type Command struct {
	instruction Instruction
	operand     int
}

func hackSolution(program []Command, expected_output []int, current_a int) (int, bool) {
	if len(expected_output) == 0 {
		return current_a, true
	} else {
		for i := 0; i < 8; i++ {
			possible_a := current_a*8 + i
			machine := Machine{A: possible_a}
			machine.runProgram(program)
			if machine.output[len(machine.output)-1] == expected_output[len(expected_output)-1] {
				sol, ok := hackSolution(program, expected_output[:len(expected_output)-1], possible_a)
				if ok {
					return sol, ok
				} else {
					continue
				}

			}

		}
	}
	return -1, false
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
	program := []Command{}
	machine := Machine{}
	expected_output := []int{}
	for scanner.Scan() {
		row := scanner.Text()
		if strings.HasPrefix(row, "Register A:") {
			parts := strings.Split(row, ": ")
			a_s := parts[1]
			a, err := strconv.Atoi(a_s)
			if err != nil {
				panic(err)
			}
			machine.A = a
		} else if strings.HasPrefix(row, "Register B:") {
			parts := strings.Split(row, ": ")
			b_s := parts[1]
			b, err := strconv.Atoi(b_s)
			if err != nil {
				panic(err)
			}
			machine.B = b
		} else if strings.HasPrefix(row, "Register C:") {
			parts := strings.Split(row, ": ")
			c_s := parts[1]
			c, err := strconv.Atoi(c_s)
			if err != nil {
				panic(err)
			}
			machine.C = c
		} else if strings.HasPrefix(row, "Program: ") {
			parts := strings.Split(row, ": ")
			i_s := strings.Split(parts[1], ",")
			for i := 1; i < len(i_s); i += 2 {
				ins, err := strconv.Atoi(i_s[i-1])
				if err != nil {
					panic(err)
				}
				expected_output = append(expected_output, ins)
				op, err := strconv.Atoi(i_s[i])
				if err != nil {
					panic(err)
				}
				program = append(program, Command{instruction: Instruction(ins), operand: op})
				expected_output = append(expected_output, op)
			}

		}
	}
	machine.runProgram(program)
	fmt.Printf("Output: ")
	for _, out := range machine.output {
		fmt.Printf("%d,", out)
	}
	fmt.Printf("\n")
	fmt.Println("Output: ", machine.output)
	program_without_loop := program[:len(program)-1]
	possible_solution, ok := hackSolution(program_without_loop, expected_output, 0)
	if ok {
		fmt.Printf("Start A: %d\n", possible_solution)
	}
	new_machine := Machine{A: possible_solution}
	new_machine.runProgram(program)
	fmt.Println("Output: ", new_machine.output)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// ...
}
