package main

import (
	"fmt"
	"os"
	"strings"
)


type Tile struct {
	TileType int
	Timelines int
}

const (
	Empty = iota
	Splitter
	Taken
)

func findParentSplitters(x, y int, board [][]Tile) int {
	count := 0
	leftAvailable := (x - 1) >= 0
	rightAvailable := (x + 1) < len(board[0])
	
	for i := y-1; i >= 0; i-- {
		if board[i][x].TileType == Splitter { 
			return count 
		}
		if board[i][x].TileType == Taken { return count + 1 }

		if leftAvailable && board[i][x-1].TileType == Splitter {
			count += board[i][x-1].Timelines
		}

		if rightAvailable && board[i][x+1].TileType == Splitter {
			count += board[i][x+1].Timelines
		}
	}
	
	return count
}

func part2(board [][]Tile) {
	for i := range board {
		for j := range board[i] {
			if board[i][j].TileType != Splitter { continue }	
			
			board[i][j].Timelines = findParentSplitters(j, i, board)
		}
	}

	total := 0
	i := len(board) - 1
	for j := range board[i] {
		val := findParentSplitters(j, i, board)
		total += val
		fmt.Print(val, " ")
	}
	//printBoard(board)
	fmt.Println("\n", total, "\n", total*2)
}

func printBoard(board [][]Tile) {
	for i := range board {
		for _, n := range board[i] {
			switch n.TileType {
			case Taken:
				fmt.Print("|")
				break
			case Splitter:
				fmt.Print( n.Timelines)
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
	board := make([][]Tile, len(lines) - 1)
	for i := range lines { 
		if len(lines[i]) == 0 { continue }
		board[i] = make([]Tile, len(lines[i]))
	}

	for i, line := range lines {
		if len(lines[i]) == 0 { continue }
		for j, r := range line {
			switch r {
			case '.':
				var t Tile
				t.TileType = Empty
				t.Timelines = 0
				board[i][j] = t
				break
			case 'S':
				var t Tile
				t.TileType = Taken
				t.Timelines = 0
				board[i][j] = t
				break
			case '^':

				var t Tile
				t.TileType = Splitter
				t.Timelines = 0

				board[i][j] = t
				break
			}
		}
	}

	part2(board)
}
