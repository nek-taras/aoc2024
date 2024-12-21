package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ItemType int

type Move int

const (
	Left      Move = iota
	Up             = iota
	Right          = iota
	Down           = iota
	Press          = iota
	EmptyMove      = iota
)

func MoveToStr(m Move) string {
	switch m {
	case Left:
		return "<"
	case Up:
		return "^"
	case Right:
		return ">"
	case Down:
		return "v"
	case Press:
		return "A"
	}
	return ""
}

func MovesToStr(moves []Move) string {
	out := ""
	for _, m := range moves {
		out = fmt.Sprintf("%s%s", out, MoveToStr(m))
	}
	return out
}

type KeypadKey int

const (
	One   KeypadKey = 1
	Two             = 2
	Three           = 3
	Four            = 4
	Five            = 5
	Six             = 6
	Seven           = 7
	Eight           = 8
	Nine            = 9
	Zero            = 0
	A               = 10
	Empty           = 11
)

func KeypadKeyToStr(k KeypadKey) string {
	switch k {
	case A:
		return "A"
	case Empty:
		return " "
	default:
		return fmt.Sprintf("%d", k)
	}

}

func ComboStr(keys []KeypadKey) string {
	out := ""
	for _, k := range keys {
		out = fmt.Sprintf("%s%s", out, KeypadKeyToStr(k))
	}
	return out
}

type KeypadItem struct {
	key KeypadKey
	x   int
	y   int
}

// +---+---+---+
// | 7 | 8 | 9 |
// +---+---+---+
// | 4 | 5 | 6 |
// +---+---+---+
// | 1 | 2 | 3 |
// +---+---+---+
// | . | 0 | A |
// +---+---+---+
func initKeypadKeys() []KeypadItem {
	keys := []KeypadItem{}
	keys = append(keys, KeypadItem{key: Seven, y: 0, x: 0})
	keys = append(keys, KeypadItem{key: Eight, y: 0, x: 1})
	keys = append(keys, KeypadItem{key: Nine, y: 0, x: 2})
	keys = append(keys, KeypadItem{key: Four, y: 1, x: 0})
	keys = append(keys, KeypadItem{key: Five, y: 1, x: 1})
	keys = append(keys, KeypadItem{key: Six, y: 1, x: 2})
	keys = append(keys, KeypadItem{key: One, y: 2, x: 0})
	keys = append(keys, KeypadItem{key: Two, y: 2, x: 1})
	keys = append(keys, KeypadItem{key: Three, y: 2, x: 2})
	keys = append(keys, KeypadItem{key: Empty, y: 3, x: 0})
	keys = append(keys, KeypadItem{key: Zero, y: 3, x: 1})
	keys = append(keys, KeypadItem{key: A, y: 3, x: 2})
	return keys
}

func getKeyPosition(key KeypadKey, keypad []KeypadItem) (int, int) {
	for _, k := range keypad {
		if k.key == key {
			return k.y, k.x
		}
	}
	return 0, 0
}

func getPositionKey(y, x int, keypad []KeypadItem) KeypadKey {
	for _, k := range keypad {
		if k.y == y && k.x == x {
			return k.key
		}
	}
	return Empty
}

// +---+---+---+
// |   | ^ | A |
// +---+---+---+
// | < | v | > |
// +---+---+---+
func initMovepadKeys() []MovePadItem {
	keys := []MovePadItem{}
	keys = append(keys, MovePadItem{key: EmptyMove, y: 0, x: 0})
	keys = append(keys, MovePadItem{key: Up, y: 0, x: 1})
	keys = append(keys, MovePadItem{key: Press, y: 0, x: 2})
	keys = append(keys, MovePadItem{key: Left, y: 1, x: 0})
	keys = append(keys, MovePadItem{key: Down, y: 1, x: 1})
	keys = append(keys, MovePadItem{key: Right, y: 1, x: 2})
	return keys
}

func getMovePosition(move Move, keypad []MovePadItem) (int, int) {
	for _, k := range keypad {
		if k.key == move {
			return k.y, k.x
		}
	}
	return 0, 0
}

func getPositionMove(y, x int, keypad []MovePadItem) Move {
	for _, k := range keypad {
		if k.y == y && k.x == x {
			return k.key
		}
	}
	return EmptyMove
}

type MovePadItem struct {
	key Move
	x   int
	y   int
}

type KeypadRobod struct {
	y             int
	x             int
	robot_movepad *Movepad
	keypad        []KeypadItem
}

func (r *KeypadRobod) resetKeypad() {
	r.y = 3
	r.x = 2
	if r.robot_movepad != nil {
		r.robot_movepad.resetMovepad()
	}
}

func GetDistanceMoves(sy, sx, ey, ex, empty_space_y, empty_space_x int, horizontal_first bool) []Move {
	request_moves := []Move{}
	dy := ey - sy
	dx := ex - sx
	horizont_moves := []Move{}
	vertiacl_moves := []Move{}
	if dx < 0 {
		for i := 0; i < (-dx); i++ {
			horizont_moves = append(horizont_moves, Left)

		}
	} else {
		for i := 0; i < dx; i++ {
			horizont_moves = append(horizont_moves, Right)
		}
	}
	if dy < 0 {
		for i := 0; i < (-dy); i++ {
			vertiacl_moves = append(vertiacl_moves, Up)

		}
	} else {
		for i := 0; i < dy; i++ {
			vertiacl_moves = append(vertiacl_moves, Down)
		}
	}
	if sx+dx == empty_space_x && sy == empty_space_y {
		request_moves = append(request_moves, vertiacl_moves...)
		request_moves = append(request_moves, horizont_moves...)
	} else if sy+dy == empty_space_y && sx == empty_space_x {
		request_moves = append(request_moves, horizont_moves...)
		request_moves = append(request_moves, vertiacl_moves...)
	} else if horizontal_first {
		request_moves = append(request_moves, horizont_moves...)
		request_moves = append(request_moves, vertiacl_moves...)

	} else {
		request_moves = append(request_moves, vertiacl_moves...)
		request_moves = append(request_moves, horizont_moves...)

	}
	return request_moves

}

func (r *KeypadRobod) inputCombination(combination []KeypadKey) int {
	move_count := 0
	move_memory := make(map[string]int)
	for _, key := range combination {
		move_count += r.moveToKey(key, move_memory)
	}
	return move_count
}

func (r *KeypadRobod) moveToKey(k KeypadKey, move_memory map[string]int) int {
	ky, kx := getKeyPosition(k, r.keypad)
	empty_y, empty_x := getKeyPosition(Empty, r.keypad)
	request_moves := GetDistanceMoves(r.y, r.x, ky, kx, empty_y, empty_x, true)
	request_moves = append(request_moves, Press)
	horizontal_move_count := r.robot_movepad.requestMoves(request_moves, move_memory)
	request_moves = GetDistanceMoves(r.y, r.x, ky, kx, empty_y, empty_x, false)
	request_moves = append(request_moves, Press)
	vertiacl_move_count := r.robot_movepad.requestMoves(request_moves, move_memory)
	move_count := 0
	if horizontal_move_count < vertiacl_move_count {
		move_count = horizontal_move_count
	} else {
		move_count = vertiacl_move_count
	}

	r.x = kx
	r.y = ky
	return move_count
}

type Movepad struct {
	name         string
	y            int
	x            int
	movepad_keys []MovePadItem
	prev_movepad *Movepad
}

func (m *Movepad) resetMovepad() {
	m.y = 0
	m.x = 2
	if m.prev_movepad != nil {
		m.prev_movepad.resetMovepad()
	}
}

func (m *Movepad) requestMoves(request []Move, move_memory map[string]int) int {
	if m.prev_movepad != nil {
		memory_key := fmt.Sprintf("n:%s|m:%s", m.name, MovesToStr(request))
		moves, ok := move_memory[memory_key]
		if ok {
			return moves
		}
		for _, mi := range request {
			move_count := 0
			my, mx := getMovePosition(mi, m.movepad_keys)
			empty_y, empty_x := getMovePosition(EmptyMove, m.movepad_keys)
			request_moves := GetDistanceMoves(m.y, m.x, my, mx, empty_y, empty_x, true)
			request_moves = append(request_moves, Press)
			horizontal_move_count := m.prev_movepad.requestMoves(request_moves, move_memory)
			request_moves = GetDistanceMoves(m.y, m.x, my, mx, empty_y, empty_x, false)
			request_moves = append(request_moves, Press)
			vertiacl_move_count := m.prev_movepad.requestMoves(request_moves, move_memory)
			if horizontal_move_count < vertiacl_move_count {
				move_count = horizontal_move_count
			} else {
				move_count = vertiacl_move_count
			}
			m.y = my
			m.x = mx
			moves += move_count
		}
		move_memory[memory_key] = moves

		return moves
	}
	return len(request)

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

	first_robot_chain := [2]Movepad{}
	// inint movepads
	for i := 0; i < len(first_robot_chain); i++ {
		movepad := first_robot_chain[i]
		movepad.name = fmt.Sprintf("Move Pad %d", i)
		movepad.y = 0
		movepad.x = 2
		movepad.movepad_keys = initMovepadKeys()
		if i > 0 {
			movepad.prev_movepad = &first_robot_chain[i-1]
		}
		first_robot_chain[i] = movepad
	}
	my_keypad := Movepad{name: "My pad", y: 0, x: 2, movepad_keys: initMovepadKeys()}
	first_robot_chain[0].prev_movepad = &my_keypad

	keypad_robot := KeypadRobod{
		y:             3,
		x:             2,
		robot_movepad: &first_robot_chain[len(first_robot_chain)-1],
		keypad:        initKeypadKeys(),
	}

	second_robot_chain := [25]Movepad{}
	// inint movepads
	for i := 0; i < len(second_robot_chain); i++ {
		movepad := second_robot_chain[i]
		movepad.name = fmt.Sprintf("Move Pad2 %d", i)
		movepad.y = 0
		movepad.x = 2
		movepad.movepad_keys = initMovepadKeys()
		if i > 0 {
			movepad.prev_movepad = &second_robot_chain[i-1]
		}
		second_robot_chain[i] = movepad
	}
	my_second_keypad := Movepad{name: "My pad2", y: 0, x: 2, movepad_keys: initMovepadKeys()}
	second_robot_chain[0].prev_movepad = &my_second_keypad

	second_keypad_robot := KeypadRobod{
		y:             3,
		x:             2,
		robot_movepad: &second_robot_chain[len(second_robot_chain)-1],
		keypad:        initKeypadKeys(),
	}

	combinations := [][]KeypadKey{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		row := scanner.Text()
		row_items := []KeypadKey{}
		for _, item_symbol := range row {
			switch item_symbol {
			case 'A':
				row_items = append(row_items, A)
			default:
				int_repr, err := strconv.Atoi(string(item_symbol))
				if err != nil {
					panic(err)
				}
				row_items = append(row_items, KeypadKey(int_repr))
			}
		}
		combinations = append(combinations, row_items)
	}
	result := 0
	for _, combination := range combinations {
		com_moves := keypad_robot.inputCombination(combination)
		combo_str := ComboStr(combination)
		key_int, err := strconv.Atoi(strings.TrimSuffix(combo_str, "A"))
		if err != nil {
			panic(err)
		}
		result += key_int * com_moves
		keypad_robot.resetKeypad()

	}
	fmt.Printf("Result: %d\n", result)
	second_result := 0
	for _, combination := range combinations {
		com_moves := second_keypad_robot.inputCombination(combination)
		combo_str := ComboStr(combination)
		key_int, err := strconv.Atoi(strings.TrimSuffix(combo_str, "A"))
		if err != nil {
			panic(err)
		}
		second_result += key_int * com_moves

		keypad_robot.resetKeypad()

	}
	fmt.Printf("Second Result: %d\n", second_result)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}
