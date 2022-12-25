package main

import (
	"os"
	"strconv"
	"strings"
)

type droplet struct {
	x, y, z int
}

func parseDroplet(line string) droplet {
	splitLine := strings.Split(line, ",")

	x, _ := strconv.ParseInt(splitLine[0], 10, 0)
	y, _ := strconv.ParseInt(splitLine[1], 10, 0)
	z, _ := strconv.ParseInt(splitLine[2], 10, 0)

	return droplet{
		x: int(x),
		y: int(y),
		z: int(z),
	}
}

func min(left, right int) int {
	if left > right {
		return right
	}

	return left
}

func max(left, right int) int {
	if left > right {
		return left
	}

	return right
}

func main() {
	inputData, err := os.ReadFile("day18/input.txt")
	if err != nil {
		panic(err)
	}

	var droplets []droplet
	for _, line := range strings.Split(string(inputData), "\n") {
		droplets = append(droplets, parseDroplet(line))
	}

	minX := droplets[0].x
	maxX := droplets[0].x
	minY := droplets[0].y
	maxY := droplets[0].y
	minZ := droplets[0].z
	maxZ := droplets[0].z

	for _, droplet := range droplets[1:] {
		minX = min(droplet.x, minX)
		maxX = max(droplet.x, maxX)
		minY = min(droplet.y, minY)
		maxY = max(droplet.y, maxY)
		minZ = min(droplet.z, minZ)
		maxZ = max(droplet.z, maxZ)
	}

	width := maxX - minX + 3
	depth := maxZ - minZ + 3
	height := maxY - minY + 3

	// [Y/height][Z/depth][X/width]
	var space [][][]bool
	for y := 0; y < height; y++ {
		space = append(space, [][]bool{})
		for z := 0; z < depth; z++ {
			space[y] = append(space[y], []bool{})
			for x := 0; x < width; x++ {
				space[y][z] = append(space[y][z], false)
			}
		}
	}

	for _, droplet := range droplets {
		space[droplet.y-minY+1][droplet.z-minZ+1][droplet.x-minX+1] = true
	}

	// [Y/height][Z/depth][X/width]
	var exteriorSpace [][][]bool
	for y := 0; y < height; y++ {
		exteriorSpace = append(exteriorSpace, [][]bool{})
		for z := 0; z < depth; z++ {
			exteriorSpace[y] = append(exteriorSpace[y], []bool{})
			for x := 0; x < width; x++ {
				exteriorSpace[y][z] = append(exteriorSpace[y][z], false)
			}
		}
	}

	// [Y/height][Z/depth][X/width]
	var notVisited [][][]bool
	for y := 0; y < height; y++ {
		notVisited = append(notVisited, [][]bool{})
		for z := 0; z < depth; z++ {
			notVisited[y] = append(notVisited[y], []bool{})
			for x := 0; x < width; x++ {
				notVisited[y][z] = append(notVisited[y][z], true)
			}
		}
	}

	isEmpty := func(space [][][]bool, x, y, z int) bool {
		if x < 0 || y < 0 || z < 0 {
			return true
		}

		if x >= width || y >= height || z >= depth {
			return true
		}

		return !space[y][z][x]
	}

	deltas := []struct{ x, y, z int }{
		{x: 1},
		{x: -1},
		{y: 1},
		{y: -1},
		{z: 1},
		{z: -1},
	}

	type point struct{ x, y, z int }
	var toVisit []point
	toVisit = append(toVisit, point{})
	for len(toVisit) > 0 {
		cur := toVisit[0]
		toVisit = toVisit[1:]

		if isEmpty(notVisited, cur.x, cur.y, cur.z) {
			continue
		}
		notVisited[cur.y][cur.z][cur.x] = false

		if isEmpty(space, cur.x, cur.y, cur.z) {
			exteriorSpace[cur.y][cur.z][cur.x] = true
		} else {
			continue
		}

		for _, delta := range deltas {
			toVisit = append(toVisit, point{
				x: cur.x + delta.x,
				y: cur.y + delta.y,
				z: cur.z + delta.z,
			})
		}
	}

	// [Y/height][Z/depth][X/width]
	var reversedExteriorSpace [][][]bool
	for y := 0; y < height; y++ {
		reversedExteriorSpace = append(reversedExteriorSpace, [][]bool{})
		for z := 0; z < depth; z++ {
			reversedExteriorSpace[y] = append(reversedExteriorSpace[y], []bool{})
			for x := 0; x < width; x++ {
				reversedExteriorSpace[y][z] = append(reversedExteriorSpace[y][z], !exteriorSpace[y][z][x])
			}
		}
	}

	surfaceArea := 0
	for y := 0; y < height; y++ {
		for z := 0; z < depth; z++ {
			for x := 0; x < width; x++ {
				if isEmpty(reversedExteriorSpace, x, y, z) {
					continue
				}

				for _, delta := range deltas {
					if isEmpty(reversedExteriorSpace, x+delta.x, y+delta.y, z+delta.z) {
						surfaceArea++
					}
				}
			}
		}
	}

	println(surfaceArea)
}
