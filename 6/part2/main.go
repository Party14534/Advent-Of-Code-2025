package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type numRange struct {
	nums []int
	mult bool
}

func main () {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")

	ranges := make([]numRange, 0)
	
	for j := 0; j < len(lines[0]); j++ {
		sum := 0
		addedToSum := 0
		for i := 0; i < len(lines) - 2; i++ {
			if string(lines[i][j]) == " " { continue }

			val, err := strconv.Atoi(string(lines[i][j]))
			if err != nil { panic("error at " + lines[i]) }

			sum *= 10
			sum += val
			addedToSum++
		}

		if sum == 0 { continue }
		
		b := string(lines[len(lines) - 2][j])
		if b == "*" {
			var r numRange
			r.nums = make([]int, 1)
			r.nums[0] = sum
			r.mult = true

			ranges = append(ranges, r)
			continue
		} else if b == "+" {
			var r numRange
			r.nums = make([]int, 1)
			r.nums[0] = sum
			r.mult = false

			ranges = append(ranges, r)
			continue
		}

		ranges[len(ranges)-1].nums = append(ranges[len(ranges)-1].nums, sum)
	}

	total := 0
	for _, r := range ranges {
		if r.mult {
			sum := 1
			for _, val := range r.nums {
				sum *= val
			}
			total += sum
		} else {
			sum := 0
			for _, val := range r.nums {
				sum += val
			}
			total += sum
		}
	}

	fmt.Println(total)
}
