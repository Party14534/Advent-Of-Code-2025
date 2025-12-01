package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func mod (val int, base int) int {
	return ((val % base) + base) % base
}

func abs (val int) int {
	if (val < 0) { return -val }
	return val
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")

	dial := 50
	password := 0

	for _, line := range lines {
		if (len(line) == 0) { break }

		modifier := -1
		if line[0] == 'R' {
			modifier = 1
		}

		val, err := strconv.Atoi(line[1:])
		if err != nil {
			panic("Error converting to string")
		}

		val *= modifier
		
		prev := dial
		dial = mod(dial + val, 100)
		if (abs(val) > 100) {
			password += abs(val) / 100
			val = mod(val, 100)
		}

		if dial == 0 || 
		(modifier == 1 && prev > dial && prev != 0) ||
		(modifier == -1 && prev < dial && prev != 0) {
			password++
		}
	}

	fmt.Println(password)
}
