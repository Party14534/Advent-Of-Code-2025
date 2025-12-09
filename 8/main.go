package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int
}

type Pair struct {
	p1 int
	p2 int
	dist int
}

type Circuit struct {
	connections []Point
}

var points []Point
var shortestDistances [][]int
var pointDistances []Pair

var depth = 1000

func distance(p1, p2 Point) int {
	d := (p1.x - p2.x) * (p1.x - p2.x) + (p1.y - p2.y) * (p1.y - p2.y) + (p1.z - p2.z) * (p1.z - p2.z)
	
	return int(math.Cbrt(float64(d)))
}

func fillShortestDistances() {
	for i := range shortestDistances {
		for j := range shortestDistances[i] {
			if i == j { continue }
			if shortestDistances[j][i] != 0 { continue }

			shortestDistances[i][j] = distance(points[i], points[j])
			var dist Pair
			dist.p1 = i
			dist.p2 = j
			dist.dist = shortestDistances[i][j]

			pointDistances = append(pointDistances, dist)
		}
	} 

	sort.Slice(pointDistances, func(i, j int) bool {
		return pointDistances[i].dist < pointDistances[j].dist
	})
}

func buildCircuits() {
	fmt.Println(pointDistances[:10])
	circuits := make([]Circuit, 0, depth)

	var lastCombined1 Point
	var lastCombined2 Point
	for i := range len(pointDistances) {
		dist := pointDistances[i]
		p1InCircuits := contains(points[dist.p1], circuits)
		p2InCircuits := contains(points[dist.p2], circuits)

		if (p1InCircuits == p2InCircuits && p1InCircuits != -1) { continue }

		fmt.Print(p1InCircuits, p2InCircuits)
		
		if (p1InCircuits == -1 && p2InCircuits == -1) {
			// Create new circuit
			var c Circuit
			c.connections = append(c.connections, points[dist.p1])
			c.connections = append(c.connections, points[dist.p2])
			circuits = append(circuits, c)

			lastCombined1 = points[dist.p1]
			lastCombined2 = points[dist.p2]
		} else if (p1InCircuits != -1 && p2InCircuits != -1) {
			// Combine circuits	
			circuits[p1InCircuits].connections = 
									append(circuits[p1InCircuits].connections,
									circuits[p2InCircuits].connections...)

			// Remove old circuit
			circuits = append(circuits[:p2InCircuits], circuits[p2InCircuits+1:]...)

			lastCombined1 = points[dist.p1]
			lastCombined2 = points[dist.p2]
		} else if (p1InCircuits == -1 && p2InCircuits != -1) {
			circuits[p2InCircuits].connections =
				append(circuits[p2InCircuits].connections, points[dist.p1])

			lastCombined1 = points[dist.p1]
			lastCombined2 = points[dist.p2]
		} else if p1InCircuits != -1 && p2InCircuits == -1 {
			circuits[p1InCircuits].connections =
				append(circuits[p1InCircuits].connections, points[dist.p2])

			lastCombined1 = points[dist.p1]
			lastCombined2 = points[dist.p2]
		}

		fmt.Println(" ", len(circuits))
	}

	fmt.Println(circuits)

	// Find longest circuits
	sort.Slice(circuits, func(i, j int) bool {
		return len(circuits[i].connections) > len(circuits[j].connections)
	})

	/*
	one := circuits[0]
	two := circuits[1]
	three := circuits[2]

	fmt.Println(len(one.connections) * len(two.connections) * len(three.connections))
	fmt.Println(len(one.connections), len(two.connections), len(three.connections))
	*/
	//fmt.Println(one, "\n--\n", two, "\n--\n", three)
	fmt.Println(lastCombined1.x * lastCombined2.x)
	fmt.Println(lastCombined1, lastCombined2)
}

/*
func findPaths(c *CircuitPoint, base *CircuitPoint, paths *[][]CircuitPoint, index int) {
	if len(c.connections) == 1 {
		// Can't return to base
		if c.p.equals(base.p) { return }

		// No cycles
		if contains(c.p, (*paths)[index]) != -1 { return }

		(*paths)[index] = append((*paths)[index], *c)
		findPaths(c.connections[0], base, paths, index)
		return
	} else if len(c.connections) > 1 {
		for _, connection := range c.connections {
			// Can't return to base
			if connection.p.equals(base.p) { continue }

			// No cycles
			if contains(connection.p, (*paths)[index]) != -1 { 
				continue
			}
			
			newIndex := len(*paths)
			(*paths) = append(*paths, (*paths)[index])
			(*paths)[newIndex] = append((*paths)[newIndex], *connection)

			fmt.Println((*paths)[index], "Here")
			findPaths(connection, base, paths, newIndex)
		}
	}
}*/

func contains(p Point, c []Circuit) int {
	for i, circuit := range c {
		for _, point := range circuit.connections {
			if point.equals(p) { return i }
		}
	}
	return -1
}

func (this *Point) equals(p Point) bool {
	return (this.x == p.x && this.y == p.y && this.z == p.z)
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Error reading file")
	}

	input := string(bytes)
	lines := strings.Split(input, "\n")

	// Build points
	points = make([]Point, len(lines) - 1)	
	for i, line := range lines {
		if len(line) == 0 { continue }

		var p Point
		nums := strings.Split(line, ",")

		p.x, err = strconv.Atoi(nums[0])
		if err != nil { panic("Failed to convert x: " + line)}
		p.y, err = strconv.Atoi(nums[1])
		if err != nil { panic("Failed to convert y: " + line)}
		p.z, err = strconv.Atoi(nums[2])
		if err != nil { panic("Failed to convert z: " + line)}

		points[i] = p
	}

	shortestDistances = make([][]int, len(points))
	for i := range shortestDistances { shortestDistances[i] = make([]int, len(points))}

	fillShortestDistances()
	buildCircuits()	

}
