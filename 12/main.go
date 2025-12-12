package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Present struct {
	shape [][]bool
}

type Region struct {
	area [][]bool
	test [][]string
	neededIds []int
}

type PresentPos struct {
	id int
	shape [][]bool
	x, y int
}

var presents []Present
var regions []Region

func shapeToString(shape [][]bool) string {
	return fmt.Sprint(shape)
}

func copy2DSlice[T any](original [][]T) [][]T {
    if original == nil {
        return nil
    }

    copied := make([][]T, len(original))
    for i := range original {
        copied[i] = make([]T, len(original[i]))
        copy(copied[i], original[i])
    }
    return copied
}

func rotate90(matrix [][]bool) [][]bool {
    n := len(matrix)

    // Transpose
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
        }
    }

    // Reverse each row
    for i := 0; i < n; i++ {
        for j := 0; j < n/2; j++ {
            matrix[i][j], matrix[i][n-1-j] = matrix[i][n-1-j], matrix[i][j]
        }
    }

	return matrix
}

var rotationCache = make(map[string][][][]bool)

func getRotations(shape [][]bool) [][][]bool {
    key := shapeToString(shape) // implement a hash function

    if cached, ok := rotationCache[key]; ok {
        return cached
    }

    rotations := computeRotations(shape)
    rotationCache[key] = rotations
    return rotations
}

func computeRotations(shape [][]bool) [][][]bool {
	rotations := make([][][]bool, 0, 4)
	rotations = append(rotations, shape)
	oneShape := rotate90(copy2DSlice(shape))
	rotations = append(rotations, oneShape)
	twoShape := rotate90(copy2DSlice(oneShape))
	rotations = append(rotations, twoShape)
	threeShape := rotate90(copy2DSlice(twoShape))
	rotations = append(rotations, threeShape)
	return rotations
}

func size (area [][]bool) int {
	count := 0
	for i := range area {
		for j := range area[i] {
			if area[i][j] { count++ }
		}
	}
	return count
}

// canPlaceShape checks if shape can be placed at (x, y)
func canPlaceShape(grid [][]bool, shape [][]bool, x, y int) bool {
    for i := 0; i < len(shape); i++ {
        for j := 0; j < len(shape[i]); j++ {
            if shape[i][j] {
                newY := y + i
                newX := x + j

                if newY >= len(grid) || newX >= len(grid[0]) {
                    return false
                }

                if grid[newY][newX] {
                    return false
                }
            }
        }
    }
    return true
}

// placeShape marks/unmarks the shape on the grid
func placeShape(grid [][]bool, shape [][]bool, x, y int, place bool) {
    for i := 0; i < len(shape); i++ {
        for j := 0; j < len(shape[i]); j++ {
            if shape[i][j] {
                grid[y+i][x+j] = place
            }
        }
    }
}

func placeGreedy(grid [][]bool, shapes [][][]bool, shapeId int) bool {
    if shapeId >= len(shapes) {
        return true
    }

    shape := shapes[shapeId]
    rotations := getRotations(shape)

    // Get candidate positions (only near existing shapes)
    positions := getCandidatePositions(grid, shapeId == 0)

    for _, rotated := range rotations {
        for _, pos := range positions {
            if canPlaceShape(grid, rotated, pos.x, pos.y) {
                placeShape(grid, rotated, pos.x, pos.y, true)

                if placeGreedy(grid, shapes, shapeId+1) {
                    return true
                }

                placeShape(grid, rotated, pos.x, pos.y, false)
            }
        }
    }

    return false
}


type Position struct{ x, y int }


func getCandidatePositions(grid [][]bool, isFirst bool) []Position {
    if isFirst {
        // First shape: only try top-left corner
        return []Position{{0, 0}}
    }

    // Only try positions adjacent to occupied cells
    positions := make([]Position, 0)
    seen := make(map[Position]bool)

    for y := 0; y < len(grid); y++ {
        for x := 0; x < len(grid[0]); x++ {
            if grid[y][x] {
                // Check all adjacent positions
                for dy := -1; dy <= 1; dy++ {
                    for dx := -1; dx <= 1; dx++ {
                        ny, nx := y+dy, x+dx
						
                        if ny >= 0 && ny < len(grid) && 
                           nx >= 0 && nx < len(grid[0]) && 
                           !grid[ny][nx] {

                            pos := Position{nx, ny}
								
                            if !seen[pos] {
                                positions = append(positions, pos)
                                seen[pos] = true
                            }
                        }
                    }
                }
            }
        }
    }

    return positions
}

/*
func placeGreedy(grid [][]bool, shapes [][][]bool, shapeId int) bool {
	if shapeId >= len(shapes) {
		return true
	}

	shape := shapes[shapeId]
	rotations := getRotations(shape)

	// Try each rotation

    for _, rotated := range rotations {
        // Try each position (left-to-right, top-to-bottom)
        for y := 0; y < len(grid); y++ {
            for x := 0; x < len(grid[0]); x++ {
                if canPlaceShape(grid, rotated, x, y) {
                    placeShape(grid, rotated, x, y, true)

                    if placeGreedy(grid, shapes, shapeId+1) {
                        return true
                    }

                    // Backtrack
                    placeShape(grid, rotated, x, y, false)
                }
            }
        }
    }

    return false
}*/

func canFitGreedy (r *Region) bool {
	numShapes := 0
	for _, val := range r.neededIds { numShapes += val }

	shapes := make([][][]bool, 0, numShapes)
	for id, val := range r.neededIds {
		for i := 0; i < val; i++ {
			shapes = append(shapes, copy2DSlice(presents[id].shape))
		}
	}

	// Sort by size
	sort.Slice(shapes, func(i, j int) bool {
		return size(shapes[i]) > size(shapes[j])
	})

	// Greedy algorithm
	return placeGreedy(r.area, shapes, 0) 
}

func printRegion(r [][]bool) {
	for i := range r {
		for j := range r[i] {
			if r[i][j] == false {
				fmt.Print(".")
			} else {
				fmt.Print("X")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func part1() {
	count := make([]bool, len(regions))
	for _, p := range presents {
		getRotations(p.shape)
	}

	var wg sync.WaitGroup

	for i, region := range regions {
		wg.Add(1)	
		go func (r *Region, i int) {
			defer wg.Done()
			val := canFitGreedy(r)
			count[i] = val
			if val {
				fmt.Println(i)
			}
		}(&region, i)
	}

	wg.Wait()

	total := 0
	for _, val := range count {
		if val { total++ }
	}

	fmt.Println(total)
}

func main () {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")

	presents = make([]Present, 0)
	regions = make([]Region, 0)

	startingPresent := true
	lastIndex := 0
	for i, line := range lines {
		if startingPresent {
			if len(line) > 3 {
				lastIndex = i
				break
			}
			presents = append(presents, Present{shape: make([][]bool, 0)})
			startingPresent = false
			continue
		} else {
			if len(line) == 0 { 
				startingPresent = true
				continue
			}

			index := len(presents) - 1
			shape := make([]bool, 3)
			for j := range 3 {
				shape[j] = line[j] == '#'
			}
			presents[index].shape = append(presents[index].shape, shape)
		}
	}

	// Build regions
	for i := lastIndex; i < len(lines); i++ {
		line := lines[i]
		if len(line) == 0 { continue }

		parts := strings.Split(line, ":")
		sizes := strings.Split(parts[0], "x")
		
		w, err := strconv.Atoi(sizes[0])
		if err != nil {panic("Failed w" + line)}
		h, err := strconv.Atoi(sizes[1])
		if err != nil {panic("Failed h" + line)}

		area := make([][]bool, h)	
		test := make([][]string, h)
		for i := range area {
			area[i] = make([]bool, w)
			test[i] = make([]string, w)
		}

		// numbers
		counts := make([]int, 0)
		nums := strings.Split(parts[1][1:], " ")
		for _, num := range nums {
			val, err := strconv.Atoi(num)
			if err != nil { panic("Help" + line + num) }
			counts = append(counts, val)
		}

		regions = append(regions, Region{
			area: area,
			neededIds: counts,
			test: test,
		})
	}

	part1()
}
