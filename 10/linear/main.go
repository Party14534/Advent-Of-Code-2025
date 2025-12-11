package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/draffensperger/golp"
)

/*
// Tried my own gaussian solver but ended up going with a premade one
func multiplyRow(matrix *[][]int, row int, val float32) {
	for j := range (*matrix)[row] {
		(*matrix)[row][j] = int(float32((*matrix)[row][j]) * val)
	}
}

func addToRow(matrix *[][]int, rowOne, rowTwo int) {
	for j := range (*matrix)[rowOne] {
		(*matrix)[rowOne][j] += (*matrix)[rowTwo][j]
	}
}

func subFromRow(matrix *[][]int, rowOne, rowTwo int) {
	for j := range (*matrix)[rowOne] {
		(*matrix)[rowOne][j] -= (*matrix)[rowTwo][j]
	}
}

func addRows(rowOne, rowTwo []int) []int {
	for i := range rowOne {
		rowOne[i] += rowTwo[i]
	}

	return rowOne
}

func subRows(rowOne, rowTwo []int) []int {
	for i := range rowOne {
		rowOne[i] -= rowTwo[i]
	}

	return rowOne
}

func rankRow(row []int, column, target int) int {
	offset := 0.0
	for j := range column {
		offset += math.Abs(float64(row[j]))
	}

	offset += math.Abs(float64(row[column] - target))

	return int(offset)
}

func findBestRow(matrix [][]int, row, column, target int) Operation {
	bestRowRank := 10000000
	bestRowIndex := -1
	isAddition := false
	targetHigher := target > matrix[row][column]
	current := matrix[row][column]
	for i := range len(matrix) {
		if i == row { continue }

		// Temporary
		// If there are values before the target column which aren't 0
		// Don't choose the column
		safe := true
		for x := range column {
			if matrix[i][x] != 0 {
				safe = false
				break
			}
		} 
		if !safe { continue }
		// Not necessarily what we want but it'll do for now

		isPositive := matrix[i][column] > 0
		
		var result []int
		rowOne := make([]int, len(matrix[row]))
		copy(rowOne, matrix[row])

		if targetHigher == isPositive {
			result = addRows(rowOne, matrix[i])	
		} else {
			result = subRows(rowOne, matrix[i])	
		}
		
		rank := rankRow(result, column, target)	
		if rank < bestRowRank {
			bestRowRank = rank
			bestRowIndex = i
			isAddition = targetHigher == isPositive
		}
	}

	// If best operation doesn't leave it at target
	// Find multiplication to leave it at target
	var mult float32 = 0.0
	bestVal := matrix[bestRowIndex][column]
	if isAddition && current + bestVal != target {
		val := current - bestVal
		mult = float32(target) / float32(val)
	} else if !isAddition && current - bestVal != target {
		val := current - bestVal
		mult = float32(target) / float32(val)
	}

	if bestRowIndex == -1 { panic("Couldn't find row") }

	fmt.Println(row, "\n", bestRowIndex, "\n", isAddition)

	return Operation {
		rowOne: row,
		rowTwo: bestRowIndex,
		mult: mult,
		isAddition: isAddition,
	}
}

func solve(m *Machine) {
	// Create equation matrix

	// Rows correspond to joltage values
	numButtons := len(m.wirings)
	numAnswers := len(m.joltage)
	matrix := make([][]int, numAnswers)
	for i := range matrix { 
		// Columns corresbond to buttons
		matrix[i] = make([]int, numButtons + 1)
	}

	// Build matrix
	for j := range numButtons {
		for i := range numAnswers {
			contains := slices.Contains(m.wirings[j].lights, i)
			
			if contains {
				matrix[i][j] = 1
			} else {
				matrix[i][j] = 0
			}
		}
	}

	// Set target values
	for i := range numAnswers {
		matrix[i][numButtons] = m.joltage[i]
	}

	printMatrix(&matrix)
	fmt.Println("---")

	// Now that the matrix is built we do Gaussian Elimination to solve
	for i := range numAnswers {
		index := i
		fmt.Println(index)
		// First make sure there are no leading numbers	
		for j := range index {
			if matrix[i][j] == 0 { continue }
			fmt.Println("Fixing it")

			// If there is a non-zero value we use other equations to set it to zero
			o := findBestRow(matrix, i, j, 0)	
			
			if o.isAddition {
				addToRow(&matrix, o.rowOne, o.rowTwo)
				if o.mult == 0 { continue }

				fmt.Println("Multiplying: ", o.mult)
				multiplyRow(&matrix, o.rowOne, o.mult)
			} else {
				subFromRow(&matrix, o.rowOne, o.rowTwo)
				if o.mult == 0 { continue }

				fmt.Println("Multiplying: ", o.mult)
				multiplyRow(&matrix, o.rowOne, o.mult)
			}
		}

		// Then make sure number is 1
		if matrix[i][index] == 1 { continue }
		fmt.Println("Fixing number")

		// If there is a non-zero value we use other equations to set it to zero
		o := findBestRow(matrix, i, i, 1)	
		
		if o.isAddition {
			addToRow(&matrix, o.rowOne, o.rowTwo)
			if o.mult == 0 { continue }

			fmt.Println("Multiplying: ", o.mult)
			multiplyRow(&matrix, o.rowOne, o.mult)
		} else {
			subFromRow(&matrix, o.rowOne, o.rowTwo)
			if o.mult == 0 { continue }

			fmt.Println("Multiplying: ", o.mult)
			multiplyRow(&matrix, o.rowOne, o.mult)
		}
	}
	
	printMatrix(&matrix)
}

func printMatrix(matrix *[][]int) {
	for i := range *matrix {
		fmt.Print("[ ")
		for j := range (*matrix)[i] {
			fmt.Print((*matrix)[i][j]," ")
		}
		fmt.Print("]\n")
	}
}
*/

func solveGolp(m *Machine) {
	// Create equation matrix

	// Rows correspond to joltage values
	numButtons := len(m.wirings)
	numAnswers := len(m.joltage)
	matrix := make([][]float64, numAnswers)
	for i := range matrix { 
		// Columns corresbond to buttons
		matrix[i] = make([]float64, numButtons)
	}

	// Build matrix
	for j := range numButtons {
		for i := range numAnswers {
			contains := slices.Contains(m.wirings[j].lights, i)
			
			if contains {
				matrix[i][j] = 1
			} else {
				matrix[i][j] = 0
			}
		}
	}

	// Solve matrix with Golp
	lp := golp.NewLP(numAnswers, numButtons)

	// Setting the objectives to one makes Golp find the smallest possible answers
	objective := make([]float64, numButtons)
	for i := range objective { objective[i] = 1 }
	lp.SetObjFn(objective)

	// Rows of the matrix
	for i := range matrix {
		lp.AddConstraint(matrix[i], golp.EQ, float64(m.joltage[i]))
	}

	// Make Golp use integers
	for j := range numButtons {
		lp.SetInt(j, true)
	}

	// Solve the equation
	lp.Solve()
	
	total := 0

	for _, val := range lp.Variables() {
		total += int(val)
	}

	m.solution = total
}

type Operation struct {
	rowOne int
	rowTwo int
	mult float32
	isAddition bool
	isDivision bool
}

type Wire struct {
	index int
	lights []int
}

type Machine struct {
	lights []bool
	wirings []Wire
	joltage []int
	solution int
}

var machines []Machine

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")

	// Build machines
	machines = make([]Machine, 0, len(lines) - 1)
	for _, line := range lines {
		if len(line) == 0 { continue }	

		parts := strings.Split(line, " ")
		var m Machine

		// Build lights
		lights := strings.Trim(parts[0], "[]")
		m.lights = make([]bool, len(lights))
		for i, r := range lights {
			m.lights[i] = r == '#'
		}
		
		// Build wirings
		m.wirings = make([]Wire, 0, len(parts) - 2)
		for i := 1; i < len(parts) - 1; i++ {
			nums := strings.Split(strings.Trim(parts[i], "()"), ",")
			numbers := make([]int, 0, len(nums))	
			for _, num := range nums {
				val, err := strconv.Atoi(num)
				if err != nil { panic("Failed converting " + num + " " + line)}

				numbers = append(numbers, val)
			}

			m.wirings = append(m.wirings, Wire{
				index: i - 1,
				lights: numbers,
			})
		}
		
		// Joltage
		jolts := strings.Split(strings.Trim(parts[len(parts) - 1], "{}"), ",")
		m.joltage = make([]int, 0, len(jolts))

		for _, jolt := range jolts {
			val, err := strconv.Atoi(jolt)
			if err != nil { panic("Failed to jolt: " + jolt + " " + line) }
			m.joltage = append(m.joltage, val)
		}

		machines = append(machines, m)
	}


	// Solve day 2
	total := 0

	for i := range machines {
		solveGolp(&machines[i])
		total += machines[i].solution
	}

	fmt.Println(total)
}
