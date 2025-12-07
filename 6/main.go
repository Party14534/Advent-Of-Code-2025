package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var part1 = false

func leftPad(s *string, length int) {
	for len(*s) < length {
		*s = " " + *s
	}
}

func getTotalForColumn(nums [][]int, column int, multiply bool) int {
	total := 0

	sVals := make([]string, len(nums))
	maxLen := 0
	for i := 0; i < len(nums); i++ {
		sVals[i] = strconv.Itoa(nums[i][column])
		if len(sVals[i]) > maxLen { maxLen = len(sVals[i]) }
	}

	for i := 0; i < len(sVals); i++ {
		leftPad(&sVals[i], maxLen)
	}

	vals := make([]int, 0)	
	for index := 0; index < maxLen; index++ {
		sum := 0
		addedToSum := 0
		for i := len(sVals) - 1; i >= 0; i-- {
			s := sVals[i]
			if string(s[index]) == " " { continue } 

			val, err := strconv.Atoi(string(s[index]))
			if err != nil {
				panic("Failed")
			}

			sum += val * int(math.Pow10(addedToSum))		
			addedToSum++
		}

		if (sum != 0) { vals = append(vals, sum) }
	}

	if multiply { total = 1 }
	for _, val := range vals {
		fmt.Println(vals)
		if multiply { total *= val } else { total += val }
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

	numbers := make([][]int, 0)

	for i, line := range lines {
		if i == len(lines) - 2 { break; }
			
		numLine := strings.Split(line, " ")
		nums := make([]int, 0)
		for _, num := range numLine {
			if len(num) == 0 { continue }
			
			val, err := strconv.Atoi(num)
			if err != nil {
				panic("Unable to convert " + line)
			}

			nums = append(nums, val)
		}

		numbers = append(numbers, nums)
	}

	index := 0
	total := 0
	for _, op := range strings.Split(lines[len(lines)-2], " ") {
		if len(op) == 0 { continue }

		sum := 0
		if op == "*" {
			sum++
			for i := 0; i < len(numbers); i++ {
				if part1 {
					sum *= numbers[i][index]
				} else {
					sum *= getTotalForColumn(numbers, index, true)
				}
			}
		} else if op == "+" {
			for i := 0; i < len(numbers); i++ {
				if part1 {
					sum += numbers[i][index]
				} else {
					sum += getTotalForColumn(numbers, index, false)
				}
			}
		} else { panic("Unknown operator") }

		//fmt.Println(sum)

		total += sum
		index++
	}

	fmt.Println(total)
}
