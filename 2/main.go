package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var part1 = false

func splitByLength(s string, chunkSize int) []string {
    var chunks []string

    for i := 0; i < len(s); i += chunkSize {
        end := i + chunkSize
		end = min(end, len(s))
        chunks = append(chunks, s[i:end])
    }

    return chunks
}

func hasSequence(num string) bool {
	length := len(num)
	
	for i := 1; i <= length/2; i++ {
		if length % i != 0 { continue }

		units := splitByLength(num, i)
		allEqual := true
		for i := 0; i < len(units) - 1; i++ {
			if units[i] != units[i+1] {
				allEqual = false
				break;
			} 
		}
		if (allEqual) { return true }
	}
	
	return false
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	ranges := strings.Split(input, ",")
	
	total := 0

	for _, str := range ranges {
		nums := strings.Split(str, "-")
		start, err := strconv.Atoi(nums[0])
		if err != nil {
			panic("Unable to convert starting number: " + str)
		}
		
		end, err := strconv.Atoi(strings.TrimSpace(nums[1]))
		if err != nil {
			panic("Unable to convert ending number: " + str)
		}

		for i := start; i <= end; i++ {
			num := strconv.Itoa(i)
			numLen := len(num)
			if part1 && numLen & 1 != 0 { continue }

			if part1 {
				firstHalf := num [0:numLen/2] 
				secondHalf := num [numLen/2:] 
				if firstHalf == secondHalf { total += i }
			} else {
				if hasSequence(num) {
					total += i
				}
			}
		}
	}

	fmt.Println(total)
}
