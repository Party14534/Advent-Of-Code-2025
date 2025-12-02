package parallel

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Party14534/Advent-Of-Code-2025/2/lib"
)

var totals []int

func RunParallel(input string, part1 bool) {
	ranges := strings.Split(input, ",")
	totals = make([]int, len(ranges))
	
	total := 0
	var wg sync.WaitGroup

	for index, str := range ranges {
		wg.Add(1)
		go func (str string, index int) {
			defer wg.Done()
			totals[index] = 0

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
					if firstHalf == secondHalf { totals[index] += i }
				} else {
					if lib.HasSequence(num) {
						totals[index] += i
					}
				}
			}
		}(str, index)
	}

	wg.Wait()

	for _, sum := range totals {
		total += sum	
	}
	fmt.Println(total)
}
