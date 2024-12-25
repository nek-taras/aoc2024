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

type OperationType int

const (
	AND OperationType = iota
	OR                = iota
	XOR               = iota
)

type Operation struct {
	operation_type OperationType
}

func OperationFromStr(op_type string) OperationType {
	switch op_type {
	case "AND":
		return AND
	case "OR":
		return OR
	case "XOR":
		return XOR
	default:
		panic("Wrong operation type")
	}
}

func OperationToStr(op_type OperationType) string {
	switch op_type {
	case AND:
		return "AND"
	case OR:
		return "OR"
	case XOR:
		return "XOR"
	default:
		panic("Wrong operation type")
	}
}

func (op Operation) Compute(input_one int, input_two int) int {
	switch op.operation_type {
	case AND:
		return input_one & input_two
	case OR:
		return input_one | input_two
	case XOR:
		return input_one ^ input_two
	default:
		panic("Wrong operation type")
	}
}

type Wire struct {
	value  int
	is_set bool
}

type Gate struct {
	name           string
	input_one_name string
	input_two_name string
	output_name    string
	operation      Operation
	is_computed    bool
}

func (g *Gate) Compute(gate_map map[string]Gate, wire_results map[string]Wire) {
	input_one := wire_results[g.input_one_name]
	input_two := wire_results[g.input_two_name]
	if !g.is_computed && input_one.is_set && input_two.is_set {
		res := g.operation.Compute(input_one.value, input_two.value)
		ouptut := wire_results[g.output_name]
		ouptut.value = res
		ouptut.is_set = true
		wire_results[g.output_name] = ouptut
		g.is_computed = true
		gate_map[g.name] = *g
	}
}

func (g *Gate) WaitingForCompute(wire_results map[string]Wire) bool {
	input_one := wire_results[g.input_one_name]
	input_two := wire_results[g.input_two_name]
	return !g.is_computed && input_one.is_set && input_two.is_set
}

type ResultItem struct {
	index int
	value int
}

func resetGates(gate_map map[string]Gate) {
	for k, g := range gate_map {
		g.is_computed = false
		gate_map[k] = g
	}
}

func setWires(x, y int, wire_results map[string]Wire) {
	x_str := strconv.FormatInt(int64(x), 2)
	y_str := strconv.FormatInt(int64(y), 2)
	x_dict := make(map[int]int)
	y_dict := make(map[int]int)
	for index := range x_str {
		val, err := strconv.Atoi(string(x_str[len(x_str)-index-1]))
		if err != nil {
			panic(err)
		}
		x_dict[index] = val
	}
	for index := range y_str {
		val, err := strconv.Atoi(string(y_str[len(y_str)-index-1]))
		if err != nil {
			panic(err)
		}
		y_dict[index] = val
	}
	for k, v := range wire_results {
		first_letter := k[0]
		if first_letter == 'x' {
			index, err := strconv.Atoi(k[1:])
			if err != nil {
				panic(err)
			}
			v.value = x_dict[index]

		} else if first_letter == 'y' {
			index, err := strconv.Atoi(k[1:])
			if err != nil {
				panic(err)
			}
			v.value = y_dict[index]
		} else {
			v.is_set = false
		}
		wire_results[k] = v
	}

}

func computeGates(gate_map map[string]Gate, wire_results map[string]Wire) int {

	for {
		waiting_for_compute := []Gate{}
		for _, g := range gate_map {
			if g.WaitingForCompute(wire_results) {
				waiting_for_compute = append(waiting_for_compute, g)
			}
		}
		if len(waiting_for_compute) == 0 {
			break
		}
		for _, g := range waiting_for_compute {
			g.Compute(gate_map, wire_results)
		}
	}
	return getResult(wire_results, "z")
}

func getResult(wire_results map[string]Wire, k_s string) int {
	result_items := []ResultItem{}
	for k, v := range wire_results {
		if strings.HasPrefix(k, k_s) {
			index_str := strings.TrimPrefix(k, k_s)
			index, err := strconv.Atoi(index_str)
			if err != nil {
				panic(err)
			}
			res_item := ResultItem{index: index, value: v.value}
			result_items = append(result_items, res_item)
		}
	}
	sort.Slice(result_items, func(i, j int) bool {
		ii := result_items[i]
		jj := result_items[j]
		return ii.index > jj.index
	})
	out_str := ""
	for _, el := range result_items {
		out_str = fmt.Sprintf("%s%d", out_str, el.value)
	}
	out, err := strconv.ParseInt(out_str, 2, 64)
	if err != nil {
		panic(err)
	}
	return int(out)
}

func isInputGate(g Gate) bool {
	return strings.HasPrefix(g.input_one_name, "x") || strings.HasPrefix(g.input_two_name, "x") || strings.HasPrefix(g.input_one_name, "y") || strings.HasPrefix(g.input_two_name, "y")

}

func findBadGates(gate_map map[string]Gate) []string {
	wrong_gates := []string{}
	wires_ops := make(map[string][]OperationType)
	for _, g := range gate_map {
		ops, ok := wires_ops[g.input_one_name]
		if !ok {
			ops = []OperationType{}
		}
		if !slices.Contains(ops, g.operation.operation_type) {
			ops = append(ops, g.operation.operation_type)

		}
		wires_ops[g.input_one_name] = ops
		ops, ok = wires_ops[g.input_two_name]
		if !ok {
			ops = []OperationType{}
		}
		if !slices.Contains(ops, g.operation.operation_type) {
			ops = append(ops, g.operation.operation_type)

		}
		wires_ops[g.input_two_name] = ops

	}

	for _, g := range gate_map {
		if g.output_name == "z45" {
			if isInputGate(g) || g.operation.operation_type != OR {
				if !slices.Contains(wrong_gates, g.output_name) {
					wrong_gates = append(wrong_gates, g.output_name)
				}
			}
			continue
		}

		if g.output_name == "z00" {
			ins := []string{}
			ins = append(ins, g.input_one_name)
			ins = append(ins, g.input_two_name)
			slices.Sort(ins)
			if ins[0] != "x00" || ins[1] != "y00" {
				if !slices.Contains(wrong_gates, g.output_name) {
					wrong_gates = append(wrong_gates, g.output_name)
				}
			}
			if g.operation.operation_type != XOR {
				if !slices.Contains(wrong_gates, g.output_name) {
					wrong_gates = append(wrong_gates, g.output_name)
				}

			}
			continue
		}

		if g.input_one_name == "x00" || g.input_one_name == "y00" || g.input_two_name == "x00" || g.input_two_name == "y00" {
			if (strings.HasPrefix(g.input_one_name, "x") && strings.HasPrefix(g.input_two_name, "y")) || (strings.HasPrefix(g.input_one_name, "y") && strings.HasPrefix(g.input_two_name, "x")) {
				if g.operation.operation_type == OR {
					if !slices.Contains(wrong_gates, g.output_name) {
						wrong_gates = append(wrong_gates, g.output_name)
					}
				}
			}
			continue
		}

		if g.operation.operation_type == XOR {
			if strings.HasPrefix(g.input_one_name, "x") || strings.HasPrefix(g.input_one_name, "y") {
				if !strings.HasPrefix(g.input_two_name, "x") && !strings.HasPrefix(g.input_two_name, "y") {
					if !slices.Contains(wrong_gates, g.output_name) {
						wrong_gates = append(wrong_gates, g.output_name)
					}
				}
				if strings.HasPrefix(g.output_name, "z") {
					if !slices.Contains(wrong_gates, g.output_name) {
						wrong_gates = append(wrong_gates, g.output_name)
					}
				}
				if !slices.Contains(wires_ops[g.output_name], AND) || !slices.Contains(wires_ops[g.output_name], XOR) {
					if !slices.Contains(wrong_gates, g.output_name) {
						wrong_gates = append(wrong_gates, g.output_name)
					}
				}

			} else if !strings.HasPrefix(g.output_name, "z") {
				if !slices.Contains(wrong_gates, g.output_name) {
					wrong_gates = append(wrong_gates, g.output_name)
				}
			}
		} else if g.operation.operation_type == OR {
			if isInputGate(g) || strings.HasPrefix(g.output_name, "z") {
				if !slices.Contains(wrong_gates, g.output_name) {
					wrong_gates = append(wrong_gates, g.output_name)
				}
			}
			if !slices.Contains(wires_ops[g.output_name], AND) || !slices.Contains(wires_ops[g.output_name], XOR) {
				if !slices.Contains(wrong_gates, g.output_name) {
					wrong_gates = append(wrong_gates, g.output_name)
				}
			}
		} else if g.operation.operation_type == AND {

			if strings.HasPrefix(g.input_one_name, "x") || strings.HasPrefix(g.input_one_name, "y") {
				if !strings.HasPrefix(g.input_two_name, "x") && !strings.HasPrefix(g.input_two_name, "y") {
					if !slices.Contains(wrong_gates, g.output_name) {
						wrong_gates = append(wrong_gates, g.output_name)
					}
				}
			}
			if !slices.Contains(wires_ops[g.output_name], OR) {
				if !slices.Contains(wrong_gates, g.output_name) {
					wrong_gates = append(wrong_gates, g.output_name)
				}

			}
		}
	}
	return wrong_gates
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
	gates := make(map[string]Gate)
	wire_results := make(map[string]Wire)
	for scanner.Scan() {
		row := scanner.Text()
		if strings.Contains(row, ":") {
			parts := strings.Split(row, ": ")
			wire_name := parts[0]
			wire_val, err := strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}
			wire_results[wire_name] = Wire{value: wire_val, is_set: true}
		} else if strings.Contains(row, "->") {
			parts := strings.Split(row, " -> ")
			out_wire_name := parts[1]
			_, ok := wire_results[out_wire_name]
			if !ok {
				wire_results[out_wire_name] = Wire{}
			}
			parts = strings.Split(parts[0], " ")
			input_one_name := parts[0]
			_, ok = wire_results[input_one_name]
			if !ok {
				wire_results[input_one_name] = Wire{}
			}
			operation_name := parts[1]
			input_two_name := parts[2]
			_, ok = wire_results[input_two_name]
			if !ok {
				wire_results[input_two_name] = Wire{}
			}
			operation := Operation{operation_type: OperationFromStr(operation_name)}
			gate := Gate{
				input_one_name: input_one_name,
				input_two_name: input_two_name,
				operation:      operation,
				output_name:    out_wire_name,
				name:           fmt.Sprintf("%s|%s|%s|%s", input_one_name, operation_name, input_two_name, out_wire_name),
			}
			gates[gate.name] = gate
		}
	}
	res := computeGates(gates, wire_results)
	fmt.Printf("Result: %d\n", res)
	wrong_gates := findBadGates(gates)
	slices.Sort(wrong_gates)
	fmt.Printf("Swaps: %s\n", strings.Join(wrong_gates, ","))
	//swap_key, _ := GetSwaps(gates, wire_results, 8)
	//outs := []string{}
	//pairs := strings.Split(swap_key, "-")
	//for _, p := range pairs {
	//	out := strings.Split(p, "|")
	//	outs = append(outs, out...)
	//}
	//slices.Sort(outs)
	//fmt.Printf("Swaps: %s\n", strings.Join(outs, ","))

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}
