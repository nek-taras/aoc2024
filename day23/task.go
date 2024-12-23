package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func intersections(conn_one, conn_second []string) []string {
	intersection := []string{}
	for _, el := range conn_one {
		if slices.Contains(conn_second, el) {
			intersection = append(intersection, el)
		}
	}
	return intersection
}

func validateIntersection(intersection []string, network map[string][]string) bool {
	for i, v := range intersection {
		for j, check_v := range intersection {
			if i == j {
				continue
			}
			if !slices.Contains(network[check_v], v) {
				return false
			}
		}
	}
	return true
}

func findCycles(start string, network map[string][]string, seen_connections map[string]bool) int {
	cycles := 0
	connections := network[start]
	for _, conn := range connections {
		second_conn := network[conn]
		for _, sc := range second_conn {
			if slices.Contains(connections, sc) {
				m_n := []string{start, conn, sc}
				slices.Sort(m_n)
				key := strings.Join(m_n, ",")
				_, ok := seen_connections[key]
				if !ok {
					cycles++
					seen_connections[key] = true

				}
			}
		}
	}
	return cycles
}

func findCyclesLong(start string, network map[string][]string, seen_connections map[string]int) {
	connections := network[start]
	for _, conn := range connections {
		second_conn := network[conn]
		intersections := intersections(connections, second_conn)
		intersections = append(intersections, start, conn)
		if !validateIntersection(intersections, network) {
			continue
		}
		slices.Sort(intersections)
		key := strings.Join(intersections, ",")
		_, ok := seen_connections[key]
		if !ok {
			seen_connections[key] = len(intersections)

		}
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
	network := make(map[string][]string)
	for scanner.Scan() {
		row := scanner.Text()
		computers := strings.Split(row, "-")
		connections, ok := network[computers[0]]
		if !ok {
			connections = []string{}
		}
		connections = append(connections, computers[1])
		network[computers[0]] = connections
		connections, ok = network[computers[1]]
		if !ok {
			connections = []string{}
		}
		connections = append(connections, computers[0])
		network[computers[1]] = connections
	}
	total_cycles := 0
	seen_connections := make(map[string]bool)
	long_connections := make(map[string]int)
	for k, _ := range network {
		if strings.HasPrefix(k, "t") {
			total_cycles += findCycles(k, network, seen_connections)
		}
		findCyclesLong(k, network, long_connections)
	}
	fmt.Printf("Total 3 comp networs: %d\n", total_cycles)
	max_connections := 0
	password := ""
	for k, v := range long_connections {
		if v > max_connections {
			max_connections = v
			password = k
		}
	}
	fmt.Printf("Password: %s\n", password)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}
