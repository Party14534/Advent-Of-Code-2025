package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	Empty = iota
	Splitter
	Taken
)

var timesSplit = 0
var isPart1 = false
var timelines = 1

func part2(board [][]int, x, y int) {
	if y + 1 >= len(board) { 
		fmt.Println(timelines)
		return 
	}	

	tile := board[y+1][x]
	switch tile {
	case Empty:
		part2(board, x, y+1)
		return
	case Splitter:
		timelines++
		if x - 1 >= 0 && board[y+1][x-1] == Empty {
			part2(board, x-1, y+1)
		}
		if x + 1 < len(board) && board[y+1][x+1] == Empty {
			part2(board, x+1, y+1)
		}
		return
	}

	fmt.Println("Here")
}

func part1(board [][]int, x, y int) {
	if y + 1 >= len(board) { return }	

	tile := board[y+1][x]
	switch tile {
	case Empty:
		board[y+1][x] = Taken
		part1(board, x, y+1)
		return
	case Splitter:
		timesSplit++
		if x - 1 >= 0 && board[y+1][x-1] == Empty {
			board[y+1][x-1] = Taken
			part1(board, x-1, y+1)
		}
		if x + 1 < len(board) && board[y+1][x+1] == Empty {
			board[y+1][x+1] = Taken
			part1(board, x+1, y+1)
		}
		return
	case Taken:
		return
	}
}

func printBoard(board [][]int) {
	for i := range board {
		for _, n := range board[i] {
			switch n {
			case Taken:
				fmt.Print("|")
				break
			case Splitter:
				fmt.Print("^")
				break
			case Empty:
				fmt.Print(".")
				break
			}
		}

		fmt.Println()
	}
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")

	// build board
	board := make([][]int, len(lines) - 1)
	for i := range lines { 
		if len(lines[i]) == 0 { continue }
		board[i] = make([]int, len(lines[i]))
	}

	startX, startY := 0, 0

	for i, line := range lines {
		if len(lines[i]) == 0 { continue }
		for j, r := range line {
			switch r {
			case '.':
				board[i][j] = Empty
				break
			case 'S':
				startX = j
				startY = i
				board[i][j] = Taken
				break
			case '^':
				board[i][j] = Splitter
				break
			}
		}
	}
	
	if isPart1 {
		part1(board, startX, startY)
	} else {
	 	part2(board, startX, startY)
	}
	
	fmt.Println(timesSplit)
	fmt.Println(timelines)
}
