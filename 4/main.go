package main

import (
	"fmt"
	"os"
	"strings"
)

func numSurrounding(x, y int, grid [][]bool) int {
	count := 0

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 { continue }
			if y + i < 0 || y + i >= len(grid) { continue }
			if x + j < 0 || x + j >= len(grid[y]) { continue }

			if grid[y + i][x + j] { count++ }
		}
	}

	return count
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")

	// Build grid
	grid := make([][]bool, len(lines) - 1)
	for i, line := range lines {
		if len(line) == 0 { continue; }

		grid[i] = make([]bool, len(line))
		for j, r := range line {
			grid[i][j] = r == '@'	
		}
	}

	// print grid
	for _, line := range grid {
		for _, val := range line {
			c := '.'
			if val { c = '@'}
			fmt.Print(string(c))
		}
		fmt.Println()
	}

	fmt.Println("\n---\n")

	prevTotal := 0
	total := 0
	
	for {
		for y, line := range grid {
			for x, val := range line {
				if val && numSurrounding(x, y, grid) <= 3 { 
					total++ 
					fmt.Print("x")
					grid[y][x] = false
				} else {
					c := "."
					if val { c = "@"}
					fmt.Print(c)
				}
			}
			fmt.Println()
		}

		if total == prevTotal { break }

		prevTotal = total
	}

	fmt.Println(total)
}
