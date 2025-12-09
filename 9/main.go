package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

var points []Point
var areas []int
var safe [][]bool
var isPart1 = false

func scanlineFill(startX, startY int) {
	if len(safe) == 0 || len(safe[0]) == 0 {
        return
    }

    height := len(safe)
    width := len(safe[0])
	
    if startX < 0 || startX >= width || startY < 0 || startY >= height {
        return
    }

    if safe[startY][startX] {
        return
    }

    stack := []Point{{startX, startY}}

    for len(stack) > 0 {
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        x, y := current.x, current.y

        if y < 0 || y >= height || safe[y][x] {
            continue
        }

        // Find left edge of scanline
        left := x
        for left > 0 && !safe[y][left-1] {
            left--
        }

        // Find right edge of scanline
        right := x
        for right < width-1 && !safe[y][right+1] {
            right++
        }

        // Fill the scanline
        for i := left; i <= right; i++ {
            safe[y][i] = true
        }

        // Check scanlines above and below - only add when transitioning from filled to unfilled
		spanAbove := false
		spanBelow := false
		for i := left; i <= right; i++ {
			if y > 0 && !safe[y-1][i] {
				if !spanAbove {
					stack = append(stack, Point{i, y - 1})
					spanAbove = true
				}
			} else {
				spanAbove = false
			}

			if y < height-1 && !safe[y+1][i] {
				if !spanBelow {
					stack = append(stack, Point{i, y + 1})
					spanBelow = true
				}
			} else {
				spanBelow = false
			}
		}
    }
}

func fillBoard() {
	for i := range len(points) {
		j := i+1
		if i == len(points) - 1 {
			j = 0
		}

		moveX := points[i].x != points[j].x
		modifier := 1
		pos := points[i].x
		x := points[i].x
		y := points[i].y


		if (moveX && points[i].x > points[j].x) || 
			!moveX && points[i].y > points[j].y{
			modifier = -1
		}
		
		if !moveX { pos = points[i].y }
		
		for {
			if moveX {
				safe[y][pos] = true	
				pos += modifier
				if pos == points[j].x { 
					safe[y][pos] = true
					break 
				}
			} else {
				safe[pos][x] = true	
				pos += modifier
				if pos == points[j].y { 
					safe[pos][x] = true
					break 
				}
			}	
		}
	}

	// Flood fill
	scanlineFill(points[0].x - 1, points[0].y + 1)	
}

func area(p1, p2 Point) int {
	x := p1.x - p2.x
	y := p1.y - p2.y

	if x < 0 { x *= -1 }
	if y < 0 { y *= -1 }
	x++
	y++

	return x * y
}

func shapeIsSafe(p1, p2 Point) bool {
	yModifer := 1
	xModifer := 1
	
	if p1.y > p2.y { yModifer *= -1}
	if p1.x > p2.x { xModifer *= -1}

	for i := p1.y; i != p2.y + yModifer; i += yModifer {
		for j := p1.x; j != p2.x + xModifer; j += xModifer {
			if !safe[i][j] { return false }
		}
	}

	return true
}

func printGrid() {
    for _, row := range safe {
        for _, cell := range row {
            if cell {
                fmt.Print("█ ")
            } else {
                fmt.Print("· ")
            }
        }
        fmt.Println()
    }
}

func part2() {
	fillBoard()
	fmt.Println("Finished filling board")
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			if shapeIsSafe(points[i], points[j]) {
				areas = append(areas, area(points[i], points[j]))		
			}
		}
	}

	sort.Slice(areas, func(i, j int) bool {
		return areas[i] > areas[j]
	})

	fmt.Println(areas[0])
}

func part1() {
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			areas = append(areas, area(points[i], points[j]))		
		}
	}

	sort.Slice(areas, func(i, j int) bool {
		return areas[i] > areas[j]
	})

	fmt.Println(areas[0])
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")

	points = make([]Point, 0, len(lines) - 1)
	maxX := 0
	maxY := 0

	for _, line := range lines {
		if len(line) == 0 { continue }
		nums := strings.Split(line, ",")
		
		x, err := strconv.Atoi(nums[0])
		if err != nil {
			panic("Failed to convert x: " + line)
		}

		y, err := strconv.Atoi(nums[1])
		if err != nil {
			panic("Failed to convert x: " + line)
		}

		var p Point
		p.x = x
		p.y = y

		maxX = max(p.x, maxX)
		maxY = max(p.y, maxY)

		points = append(points, p)
	}

	if isPart1 {
		part1()
	} else {
		safe = make([][]bool, maxY + 1)
		for i := range maxY+1 { safe[i] = make([]bool, maxX + 1)}
		part2()
	}
}
