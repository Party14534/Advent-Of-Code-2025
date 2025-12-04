package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var isPart1 = false

func part2(banks [][]int) int {
	total := 0

	for _, bank := range banks {
		if len(bank) == 0 { break }

		maximums := make([]int, len(bank)) 
		maximumIds := make([]int, len(bank))

		for length := 11; length >= 0; length-- {
			prevId := -1
			index := 0
			if length != 11 {
				prevId = maximumIds[10 - length]
				index = 11 - length
			}

			for i := prevId + 1; i < len(bank) - length; i++ {
				if bank[i] > maximums[index] {
					maximums[index] = bank[i]
					maximumIds[index] = i
				}
			}
		}
		
		sum := 0;
		for i, val := range maximums {
			sum += val * int(math.Pow10(11 - i))
		}
		
		// Add to total
		fmt.Println(sum)
		total += sum
	}
	
	return total
}

func part1(banks [][]int) int {
	total := 0

	for _, bank := range banks {
		if len(bank) == 0 { break }

		// First get max
		maximum := -1
		maxIndex := -1
		for i := 0; i < len(bank) - 1; i++ {
			if bank[i] > maximum {
				maximum = bank[i]
				maxIndex = i
			}
		}
		
		secondMax := -1
		for i := maxIndex + 1; i < len(bank); i++ {
			if bank[i] > secondMax {
				secondMax = bank[i]
			}
		}

		// Add to total
		fmt.Println(maximum * 10 + secondMax)
		total += (maximum * 10) + secondMax
	}
	
	return total
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")

	// Set up data
	banks := make([][]int, len(lines))
	for i := range banks {
		banks[i] = make([]int, len(lines[i]))
	}

	for i, line := range lines {
		for j, num := range line {
			val, err := strconv.Atoi(string(num))
			if err != nil {
				panic("Unable to convert: " + line)
			}

			banks[i][j] = val
		}
	}
	
	if isPart1 {
		fmt.Println(part1(banks))
	} else {
		fmt.Println(part2(banks))
	}
}
