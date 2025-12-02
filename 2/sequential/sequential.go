package sequential

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/Party14534/Advent-Of-Code-2025/2/lib"
)

func RunSequential(input string, part1 bool) {
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
				if lib.HasSequence(num) {
					total += i
				}
			}
		}
	}

	fmt.Println(total)
}
