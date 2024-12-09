package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

type DisckBlock struct {
	index  int
	isFile bool
}

type File struct {
	start int
	size  int
}

func removeIndex(s [][]int, index int) [][]int {
	ret := make([][]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func getFreeSpaceGroups(space []int) [][]int {
	last_v := space[0]
	free_space := []int{last_v}
	free_spaces := [][]int{}
	el_added := 0
	for i := 1; i < len(space); i++ {
		if space[i]-1 == last_v {
			free_space = append(free_space, space[i])
		} else {
			free_spaces = append(free_spaces, free_space)
			el_added += len(free_space)
			free_space = []int{space[i]}
		}
		last_v = space[i]
	}
	if el_added < len(space) {
		free_spaces = append(free_spaces, free_space)
	}
	return free_spaces
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Input file is missing.")
		os.Exit(1)
	}
	disk_blocks := []DisckBlock{}
	disk_blocks_for_second_part := []DisckBlock{}
	files := []File{}
	content, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Print(err)
	}
	is_file := true
	file_count := 0
	prev_file_count := -1
	pos := 0
	free_space_indexes := []int{}
	for _, b := range content {

		b_size, err := strconv.Atoi(string(b))
		for i_b := 0; i_b < b_size; i_b++ {
			disk_blocks = append(disk_blocks, DisckBlock{index: file_count, isFile: is_file})
			disk_blocks_for_second_part = append(disk_blocks_for_second_part, DisckBlock{index: file_count, isFile: is_file})
			if is_file && prev_file_count < file_count {
				files = append(files, File{start: pos, size: b_size})
				prev_file_count = file_count
			}
		}
		pos += b_size
		if err != nil {
			panic(err)
		}
		is_file = !is_file
		if !is_file {
			file_count++
		}
	}
	for i, db := range disk_blocks {
		if !db.isFile {
			free_space_indexes = append(free_space_indexes, i)
		}

	}
	free_space_indexes_for_second_part := getFreeSpaceGroups(free_space_indexes)
	for i := len(disk_blocks) - 1; i >= 0; i-- {
		block := disk_blocks[i]
		// skip empty space
		if !block.isFile {
			continue
		}
		// no vacant free space
		if len(free_space_indexes) == 0 {
			break
		}
		j := 0
		j, free_space_indexes = free_space_indexes[0], free_space_indexes[1:]
		// do not move into free space that is far from here
		if j > i {
			break
		}
		disk_blocks[j], disk_blocks[i] = disk_blocks[i], disk_blocks[j]

	}
	checksum := 0
	for i, db := range disk_blocks {
		if db.isFile {
			checksum += i * db.index
		}
	}
	fmt.Printf("checksum :%d\n", checksum)
	for f_i := len(files) - 1; f_i >= 0; f_i-- {
		file := files[f_i]
		// no vacant free space
		if len(free_space_indexes_for_second_part) == 0 {
			break
		}
		free_space_found := false
		free_space_index := 0
		if free_space_indexes_for_second_part[0][0] > file.start {
			// all free space is behind
			break
		}
		for fi, fs := range free_space_indexes_for_second_part {
			if len(fs) >= file.size {
				free_space_found = true
				free_space_index = fi
				break
			}
		}
		if free_space_found {
			free_space := free_space_indexes_for_second_part[free_space_index]
			// free space behind file
			if free_space[0] > file.start {
				continue
			}
			for i := file.start; i < file.start+file.size; i++ {
				j := 0
				j, free_space = free_space[0], free_space[1:]
				disk_blocks_for_second_part[j], disk_blocks_for_second_part[i] = disk_blocks_for_second_part[i], disk_blocks_for_second_part[j]
			}
			if len(free_space) == 0 {
				free_space_indexes_for_second_part = removeIndex(free_space_indexes_for_second_part, free_space_index)
			} else {
				free_space_indexes_for_second_part[free_space_index] = free_space
			}
		}

	}
	checksum_second_part := 0
	for i, db := range disk_blocks_for_second_part {
		if db.isFile {
			checksum_second_part += i * db.index
		}
	}
	fmt.Printf("checksum second part:%d\n", checksum_second_part)

	// ...
}
