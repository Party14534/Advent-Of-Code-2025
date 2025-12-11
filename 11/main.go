package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Device struct {
	key string
	connections []string
}

var devices map[string]Device

func test() {
	type State struct {
		node   string
		hasDac bool
		hasFFT bool
	}

	memo := make(map[State]int)

	var dfs func(node string, hasDac, hasFFT bool, visited map[string]bool) int

	dfs = func(node string, hasDac, hasFFT bool, visited map[string]bool) int {
		// Base case: reached destination
		if node == "out" {
			if hasDac && hasFFT {
				return 1
			}
			return 0
		}

		// Check memo
		state := State{node, hasDac, hasFFT}
		if count, ok := memo[state]; ok {
			return count
		}

		// Explore connections
		totalPaths := 0
		d := devices[node]

		for _, next := range d.connections {
			if visited[next] {
				continue
			}

			visited[next] = true
			totalPaths += dfs(
				next,
				hasDac || next == "dac",
				hasFFT || next == "fft",
				visited,
			)
			delete(visited, next)
		}

		memo[state] = totalPaths
		return totalPaths
	}

	visited := map[string]bool{"svr": true}
	result := dfs("svr", false, false, visited)
	fmt.Println(result)
}

func getKey(path []string) string {
	return fmt.Sprint(path)
}

func part2() {
	type QueueItem struct {
		path []string
		hasDac bool
		hasFft bool
	}

	count := 0

	queue := make([]QueueItem, 1)
	queue[0] = QueueItem{
		path: []string{"svr"},
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		d, ok := devices[current.path[len(current.path)-1]]
		if !ok { panic("Not found") }
		for _, connection := range d.connections {
			// Cycle detection
			skip := slices.Contains(current.path, connection)
			if skip { continue }

			if connection == "out" {
				if !current.hasDac || !current.hasFft {
					continue
				}

				count++
				continue
			}

			if connection == "dac" { current.hasDac = true }
			if connection == "fft" { current.hasFft = true }

			newPath := make([]string, len(current.path), len(current.path)+1)
			copy(newPath, current.path)
			newPath = append(newPath, connection)

			queue = append(queue, QueueItem { 
				path: newPath, 
				hasDac: current.hasDac,
				hasFft: current.hasFft,
			})
		}
	}

	fmt.Println(count)
}

func part1() {
	type QueueItem struct {
		path []string
	}

	solutions := make([][]string, 0)

	queue := make([]QueueItem, 1)
	queue[0] = QueueItem{
		path: []string{"you"},
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		d, ok := devices[current.path[len(current.path)-1]]
		if !ok { panic("Not found") }
		for _, connection := range d.connections {
			if connection == "out" {
				solutions = append(solutions, append(current.path, "out"))
				continue
			}


			newPath := make([]string, len(current.path), len(current.path)+1)
			copy(newPath, current.path)
			newPath = append(newPath, connection)

			queue = append(queue, QueueItem{ path: newPath })
		}
	}

	fmt.Println(len(solutions))
}

func main() { 
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")
	devices = make(map[string]Device)

	for _, line := range lines {
		if len(line) == 0 { continue }

		segments := strings.Split(line, ":")
		key := segments[0]
		connections := strings.Split(segments[1][1:], " ")

		devices[key] = Device{key: key, connections: connections}
	}

	test()
}

