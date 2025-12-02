package main

import (
	"os"
	"github.com/Party14534/Advent-Of-Code-2025/2/parallel"
	"github.com/Party14534/Advent-Of-Code-2025/2/sequential"
)

var part1 = false
var isParallel = false

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	if (len(os.Args) != 1) { 
		isParallel = os.Args[1] == "-p"
	}

	input := string(bytes)

	if isParallel {
		parallel.RunParallel(input, part1)
	} else {
		sequential.RunSequential(input, part1)
	}
}
