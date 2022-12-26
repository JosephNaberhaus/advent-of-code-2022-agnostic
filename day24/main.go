package main

import (
	"os"
	"strings"
)

const maxBlizzards = 3500

type direction int8

const (
	up direction = iota
	down
	left
	right
)

type pos struct {
	x, y int
}

type blizzard struct {
	pos
	dir direction
}

type state struct {
	you          pos
	numBlizzards int
	blizzards    [maxBlizzards]blizzard
}

func (s state) nextStates(width, height int, goal pos) []state {
	var newBlizzards [maxBlizzards]blizzard
	newBlizzardPositions := map[pos]struct{}{}
	for i, b := range s.blizzards[:s.numBlizzards] {
		switch b.dir {
		case up:
			b.y--
			if b.y == 0 {
				b.y = height - 2
			}

			newBlizzards[i] = b
			newBlizzardPositions[b.pos] = struct{}{}
		case down:
			b.y++
			if b.y == height-1 {
				b.y = 1
			}

			newBlizzards[i] = b
			newBlizzardPositions[b.pos] = struct{}{}
		case left:
			b.x--
			if b.x == 0 {
				b.x = width - 2
			}

			newBlizzards[i] = b
			newBlizzardPositions[b.pos] = struct{}{}
		case right:
			b.x++
			if b.x == width-1 {
				b.x = 1
			}

			newBlizzards[i] = b
			newBlizzardPositions[b.pos] = struct{}{}
		}
	}

	var nextStates []state

	if _, ok := newBlizzardPositions[s.you]; !ok {
		nextStates = append(nextStates, state{
			you:          s.you,
			blizzards:    newBlizzards,
			numBlizzards: s.numBlizzards,
		})
	}

	deltas := []struct{ dx, dy int }{
		{dx: 1},
		{dx: -1},
		{dy: 1},
		{dy: -1},
	}
	for _, delta := range deltas {
		newPos := pos{
			x: s.you.x + delta.dx,
			y: s.you.y + delta.dy,
		}

		if newPos.x <= 0 || newPos.x >= width-1 || newPos.y <= 0 || newPos.y >= height-1 {
			if !(newPos == goal) {
				continue
			}
		}

		if _, ok := newBlizzardPositions[newPos]; ok {
			continue
		}

		nextStates = append(nextStates, state{
			you:          newPos,
			blizzards:    newBlizzards,
			numBlizzards: s.numBlizzards,
		})
	}

	return nextStates
}

func main() {
	inputData, err := os.ReadFile("day24/input.txt")
	if err != nil {
		panic(err)
	}

	numBlizzards := 0
	var blizzards [maxBlizzards]blizzard
	var valley [][]bool
	for y, line := range strings.Split(string(inputData), "\n") {
		var curLine []bool

		for x, c := range line {
			switch c {
			case '#':
				curLine = append(curLine, true)
			case '.':
				curLine = append(curLine, false)
			case '<':
				curLine = append(curLine, false)
				blizzards[numBlizzards] = blizzard{
					pos: pos{x: x, y: y},
					dir: left,
				}
				numBlizzards++
			case '>':
				curLine = append(curLine, false)
				blizzards[numBlizzards] = blizzard{
					pos: pos{x: x, y: y},
					dir: right,
				}
				numBlizzards++
			case '^':
				curLine = append(curLine, false)
				blizzards[numBlizzards] = blizzard{
					pos: pos{x: x, y: y},
					dir: up,
				}
				numBlizzards++
			case 'v':
				curLine = append(curLine, false)
				blizzards[numBlizzards] = blizzard{
					pos: pos{x: x, y: y},
					dir: down,
				}
				numBlizzards++
			}
		}

		valley = append(valley, curLine)
	}

	width := len(valley[0])
	height := len(valley)

	curGoalIndex := 0
	goals := []pos{
		{x: width - 2, y: height - 1},
		{x: 1, y: 0},
		{x: width - 2, y: height - 1},
	}

	minute := 0
	curStates := map[state]struct{}{
		state{
			you:          pos{x: 1, y: 0},
			blizzards:    blizzards,
			numBlizzards: numBlizzards,
		}: {},
	}
	nextStates := map[state]struct{}{}
	for len(curStates) != 0 {
		println(minute)

		for s := range curStates {
			if s.you == goals[curGoalIndex] {
				if curGoalIndex == len(goals)-1 {
					println(minute)
					os.Exit(42)
				}

				curGoalIndex++

				nextStates = map[state]struct{}{}
				for _, nextState := range s.nextStates(width, height, goals[curGoalIndex]) {
					nextStates[nextState] = struct{}{}
				}

				break
			}

			for _, nextState := range s.nextStates(width, height, goals[curGoalIndex]) {
				nextStates[nextState] = struct{}{}
			}
		}

		minute++
		curStates = nextStates
		nextStates = map[state]struct{}{}
	}
}
