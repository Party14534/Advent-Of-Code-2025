package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	Min int
	Max int
}

var ranges []Range

func printRanges() {
	for _, r := range ranges {
		fmt.Printf("%v-%v\n", r.Min, r.Max)
	}
}

func countRanges () int {
	total := 0
	for _, r := range ranges {
		total += (r.Max - r.Min) + 1
	}

	return total
}

func inRanges (x int) bool {
	for _, r := range ranges {
		if x >= r.Min && x <= r.Max { return true }	
	}

	return false
}

func addRange(r Range) {
	for i := 0; i < len(ranges); i++ {
		// rm Rm rM RM
		if (r.Min >= ranges[i].Min && r.Min <= ranges[i].Max &&
			r.Max >= ranges[i].Max) {
			ranges[i].Max = r.Max
			return
		}

		// Rm rm RM rM
		if (r.Min <= ranges[i].Min && r.Max <= ranges[i].Max &&
			r.Max >= ranges[i].Min) {
			ranges[i].Min = r.Min
			return
		}

		// Rm rm rM RM
		if (r.Min <= ranges[i].Min && r.Max >= ranges[i].Max) {
			ranges[i].Min = r.Min
			ranges[i].Max = r.Max
			return
		}

		// Rm RM rm rM
		if (r.Max < ranges[i].Min) {
			ranges = append(ranges[:i], append([]Range{r}, ranges[i:]...)... )
			return
		}
	}

	ranges = append(ranges, r)
}

func shrinkRanges() {
	for i := 0; i < len(ranges) - 1; i++ {
		if ranges[i].Max >= ranges[i+1].Max {
			ranges = append(ranges[:i+1], ranges[i+2:]...)
			i--
		} else if ranges[i].Max >= ranges[i+1].Min {
			ranges[i].Max = ranges[i+1].Max
			ranges = append(ranges[:i+1], ranges[i+2:]...)
			i--
		}
	}
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")


	// Build ranges
	lastIndex := 0
	for i, line := range lines {
		if len(line) == 0 { 
			lastIndex = i
			break
		}

		var r Range
		nums := strings.Split(line, "-")

		r.Min, err = strconv.Atoi(nums[0])
		if err != nil {
			panic("Unable to convert " + line)
		}

		r.Max, err = strconv.Atoi(nums[1])
		if err != nil {
			panic("Unable to convert " + line)
		}

		addRange(r)
	}
	printRanges()
	fmt.Println("HHHH")
	shrinkRanges()

	// Count expired
	total := 0
	for i := lastIndex + 1; i < len(lines); i++ {
		line := lines[i]
		if len(line) == 0 { break }

		num, err := strconv.Atoi(line)
		if err != nil {
			panic("Unable to convert id " + line)
		}

		if inRanges(num) { total += 1 }
	}

	printRanges()
	fmt.Println(total)
	fmt.Println(countRanges())
}
