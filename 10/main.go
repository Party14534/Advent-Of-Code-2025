package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
)


type Wire struct {
	index int
	lights []int
}

type Machine struct {
	lights []bool
	wirings []Wire
	joltage []int
	solution []Wire
	solution2 int
}

var machines []Machine
var isPart1 = false

func stateIsBad(joltage []int, state []int) bool {
	for i := range state {
		if state[i] > joltage[i] { return true}
	}

	return false
}

func stateIsIt(m *Machine, state []int) bool {
	for i := range m.joltage {
		if state[i] != m.joltage[i] { return false }
	}
	return true
}

func getKey(s []int) string {
	return fmt.Sprint(s)
}

func addWire(state *[]int, wire Wire) {
	for _, w := range wire.lights {
		(*state)[w]++
	}
}

func subWire(state *[]int, wire Wire) {
	for _, w := range wire.lights {
		(*state)[w]--
	}
}

// Lower equals better
func closeness(goal []int, current []int) int {
	closeness := 0
	
	for i := range goal {
		//fmt.Println(current[i], goal[i])
		if current[i] > goal[i] { return -1 }
		
		closeness += goal[i] - current[i]
	}

	//fmt.Println("--",closeness)
	return closeness
}

func getSolution5 (m *Machine) {
	type Choice struct {
		close int
		wire Wire
	}

	state := make([]int, len(m.joltage))
	visited := make(map[string]bool)
	visited[getKey(state)] = true

	for {
		minIndex := 0
		minVal := 1000000000
		for i, val := range m.joltage {
			if state[i] < val && val < minVal {
				minVal = val
				minIndex = i
			}
		}

		choices := make([]Choice, 0, len(m.wirings))
		for _, wire := range m.wirings {
			safe := slices.Contains(wire.lights, minIndex)
			if !safe { continue }

			newState := append([]int{}, state...)
			addWire(&newState, wire)
			val, ok := visited[getKey(newState)]
			if ok && val { 
				continue 
			}

			wireCloseness := closeness(m.joltage, newState)
			if wireCloseness == - 1 { continue }
			choices = append(choices, Choice{close: wireCloseness, wire: wire})
		}

		sort.Slice(choices, func(i, j int) bool {
			return choices[i].close < choices[j].close
		})

		if len(m.solution) == 0 && len(choices) == 0 { 
			panic("Error")
		}

		if len(choices) == 0 {
			subWire(&state, m.solution[len(m.solution) - 1])
			m.solution = m.solution[:len(m.solution) - 1]
			continue
		}

		addWire(&state, choices[0].wire)
		visited[getKey(state)] = true
		m.solution = append(m.solution, choices[0].wire)
		if choices[0].close == 0 { return }
	}

}

func getSolution3(m *Machine) {
	type QueueItem struct {
		state []int
		solution []Wire
		close int
	}

	type Choice struct {
		close int
		wire Wire
	}

	initialState := make([]int, len(m.joltage))
	queue := make([]QueueItem, 1)
	queue[0] = QueueItem{state: initialState}
	visited := make(map[string]bool)
	visited[getKey(initialState)] = true

	for len(queue) > 0 {
		current := queue[len(queue) - 1]
		queue = queue[:len(queue) - 1]
		//if len(queue) % 100 == 0 {fmt.Println(current.state)}

		choices := make([]QueueItem, 0, len(m.wirings))
		for _, wire := range m.wirings {
			newState := append([]int{}, current.state...)
			addWire(&newState, wire)
			_, ok := visited[getKey(newState)]
			if ok { 
				continue 
			}

			visited[getKey(newState)] = true
			wireCloseness := closeness(m.joltage, newState)
			if wireCloseness > 10 * 10 * 10 { continue }
			choices = append(choices, QueueItem{
				close: wireCloseness,
				solution: append(current.solution, wire),
				state: newState,
			})
		}

		sort.Slice(choices, func(i, j int) bool {
			return choices[i].close < choices[j].close
		})

		if len(current.solution) == 0 && len(choices) == 0 { 
			panic("Error")
		}

		if len(choices) == 0 {
			continue
		}

		if choices[0].close == 0 { 
			m.solution = choices[0].solution
			return 
		}

		queue = append(queue, choices...)
	}

}

func getSolution2 (m *Machine) {
	type Choice struct {
		close int
		wire Wire
	}

	state := make([]int, len(m.joltage))
	visited := make(map[string]bool)
	visited[getKey(state)] = true

	for {
		choices := make([]Choice, 0, len(m.wirings))
		for _, wire := range m.wirings {
			newState := append([]int{}, state...)
			addWire(&newState, wire)
			val, ok := visited[getKey(newState)]
			if ok && val { 
				continue 
			}

			wireCloseness := closeness(m.joltage, newState)
			if wireCloseness > 10 * 10 * 10 { continue }
			choices = append(choices, Choice{close: wireCloseness, wire: wire})
		}

		sort.Slice(choices, func(i, j int) bool {
			return choices[i].close < choices[j].close
		})

		if len(m.solution) == 0 && len(choices) == 0 { 
			panic("Error")
		}

		if len(choices) == 0 {
			subWire(&state, m.solution[len(m.solution) - 1])
			m.solution = m.solution[:len(m.solution) - 1]
			continue
		}

		addWire(&state, choices[0].wire)
		visited[getKey(state)] = true
		m.solution = append(m.solution, choices[0].wire)
		if choices[0].close == 0 { return }
	}
}

func part2() {
	var wg sync.WaitGroup

	for i := range machines {
		fmt.Printf("%v %v\n", float32(i)/float32(len(machines)) * 100, i)
		//getSolution2(&machines[i])
		//getSolution5(&machines[i])
		wg.Add(1)
		go func(m *Machine, i int) {
			defer wg.Done()
			defer fmt.Println(i)
			getSolution5(m)
		}(&machines[i], i)
	}		

	wg.Wait()
	fmt.Println("Done")

	total := 0
	for _, machine := range machines {
		fmt.Println(len(machine.solution))
		total += len(machine.solution)
	}

	fmt.Println(total)
}

func performSolution (wirings []Wire, index int) {
	lights := make([]bool, len(machines[index].lights))
	
	for _, wiring := range wirings {
		for _, wire := range wiring.lights {
			lights[wire] = !lights[wire]
		}
	}

	// Check if equal
	for i := range lights {
		if lights[i] != machines[index].lights[i] { return }
	}
	
	if len(machines[index].solution) == 0 ||
		len(machines[index].solution) > len(wirings) {
		machines[index].solution = wirings		
	}
}

func processAllCombinations(wirings []Wire, index int) {
    n := len(wirings)
	fmt.Println(n)
    for i := 0; i < (1 << n); i++ {
        var subset []Wire
        for j := 0; j < n; j++ {
            if i & (1 << j) != 0 {
                subset = append(subset, wirings[j])
            }
        }

        performSolution(subset, index)
    }
}

func part1() {
	for i, machine := range machines {
		processAllCombinations(machine.wirings, i)
	}

	total := 0
	for _, machine := range machines {
		fmt.Println(len(machine.solution))
		total += len(machine.solution)
	}

	fmt.Println(total)
}

func main() {
	f, err := os.Create("cpu.prof")
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

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

	if isPart1 {
		part1()
	} else {
		part2()
	}
}
