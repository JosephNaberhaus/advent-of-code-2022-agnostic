package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func max(first, second int) int {
	if first > second {
		return first
	}

	return second
}

func min(first, second int) int {
	if first > second {
		return second
	}

	return first
}

type pos struct {
	x, y int
}

type elf struct{}

type delta struct {
	dx, dy int
}

var deltaToIndex = map[delta]int{
	{dy: -1}: 1,
	{dy: 1}:  2,
	{dx: -1}: 3,
	{dx: 1}:  4,
}

func moveToBack(arr [4]int, value int) [4]int {
	found := false
	for i := 0; i < 3; i++ {
		if arr[i] == value {
			found = true
		}

		if found {
			temp := arr[i]
			arr[i] = arr[i+1]
			arr[i+1] = temp
		}
	}

	return arr
}

func main() {
	inputData, err := os.ReadFile("day23/input.txt")
	if err != nil {
		panic(err)
	}

	var hasElf [][]bool
	for _, line := range strings.Split(string(inputData), "\n") {
		var curLine []bool

		for _, c := range line {
			curLine = append(curLine, c == '#')
		}

		hasElf = append(hasElf, curLine)
	}

	width := len(hasElf[0])
	height := len(hasElf)

	elves := map[pos]elf{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if hasElf[y][x] {
				elves[pos{x: x, y: y}] = elf{}
			}
		}
	}

	clearIn := func(p pos, deltas ...delta) bool {
		for _, delta := range deltas {
			if _, ok := elves[pos{x: p.x + delta.dx, y: p.y + delta.dy}]; ok {
				return false
			}
		}

		return true
	}

	movePreferenceOrder := [4]int{1, 2, 3, 4}
	roundNum := 1
	for true {
		proposedMoves := map[pos]delta{}

		for p := range elves {
			if clearIn(p, delta{dx: -1, dy: -1}, delta{dy: -1}, delta{dx: 1, dy: -1}, delta{dx: 1}, delta{dx: 1, dy: 1}, delta{dy: 1}, delta{dx: -1, dy: 1}, delta{dx: -1}) {
				continue
			}

			for _, movePreference := range movePreferenceOrder {
				if movePreference == 1 && clearIn(p, delta{dx: -1, dy: -1}, delta{dy: -1}, delta{dx: 1, dy: -1}) {
					proposedMoves[p] = delta{dy: -1}
					break
				}

				if movePreference == 2 && clearIn(p, delta{dx: -1, dy: 1}, delta{dy: 1}, delta{dx: 1, dy: 1}) {
					proposedMoves[p] = delta{dy: 1}
					break
				}

				if movePreference == 3 && clearIn(p, delta{dx: -1, dy: -1}, delta{dx: -1}, delta{dx: -1, dy: 1}) {
					proposedMoves[p] = delta{dx: -1}
					break
				}

				if movePreference == 4 && clearIn(p, delta{dx: 1, dy: -1}, delta{dx: 1}, delta{dx: 1, dy: 1}) {
					proposedMoves[p] = delta{dx: 1}
					break
				}
			}
		}

		resultLocationsCount := map[pos]int{}
		for pos, move := range proposedMoves {
			pos.x += move.dx
			pos.y += move.dy
			resultLocationsCount[pos]++
		}

		actualMoves := map[pos]delta{}
		for pos, move := range proposedMoves {
			hax := pos
			pos.x += move.dx
			pos.y += move.dy
			if resultLocationsCount[pos] <= 1 {
				actualMoves[hax] = move
			}
		}

		if len(actualMoves) == 0 {
			println(roundNum)
			os.Exit(42)
		}

		for p, move := range actualMoves {
			delete(elves, p)
			p.x += move.dx
			p.y += move.dy
			elves[p] = elf{}
		}

		movePreferenceOrder = moveToBack(movePreferenceOrder, movePreferenceOrder[0])
		fmt.Printf("%v\n", movePreferenceOrder)

		roundNum++
	}

	minX := math.MaxInt
	minY := math.MaxInt
	maxX := math.MinInt
	maxY := math.MinInt
	for elf := range elves {
		minX = min(minX, elf.x)
		minY = min(minY, elf.y)
		maxX = max(maxX, elf.x)
		maxY = max(maxY, elf.y)
	}

	println(minX, minY, maxX, maxY)

	count := 0
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if _, ok := elves[pos{x: x, y: y}]; !ok {
				count++
			}
		}
	}

	println(count)
}
