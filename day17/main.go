package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type jet int

const (
	jetLeft  jet = -1
	jetRight jet = 1
)

func (j jet) dx() int {
	return int(j)
}

type stone struct {
	x, y int
}

type rock []stone

func (r rock) copy() rock {
	newStones := make([]stone, len(r))
	copy(newStones, r)

	return newStones
}

func (r rock) move(dx, dy int) rock {
	movedRock := r.copy()
	for i := range movedRock {
		movedRock[i].x += dx
		movedRock[i].y += dy
	}

	return movedRock
}

var rocks = []rock{
	{
		{x: 0, y: 0},
		{x: 1, y: 0},
		{x: 2, y: 0},
		{x: 3, y: 0},
	},
	{
		{x: 1, y: 0},
		{x: 0, y: 1},
		{x: 1, y: 1},
		{x: 2, y: 1},
		{x: 1, y: 2},
	},
	{
		{x: 0, y: 0},
		{x: 1, y: 0},
		{x: 2, y: 0},
		{x: 2, y: 1},
		{x: 2, y: 2},
	},
	{
		{x: 0, y: 0},
		{x: 0, y: 1},
		{x: 0, y: 2},
		{x: 0, y: 3},
	},
	{
		{x: 0, y: 0},
		{x: 1, y: 0},
		{x: 0, y: 1},
		{x: 1, y: 1},
	},
}

func parseJets(line string) []jet {
	jets := make([]jet, 0, len(line))
	for _, c := range line {
		if c == '<' {
			jets = append(jets, jetLeft)
		} else {
			jets = append(jets, jetRight)
		}
	}

	return jets
}

type stack struct {
	layers [][7]bool
	height int
}

func (s *stack) top() int {
	return len(s.layers) - 1
}

func (s *stack) addStone(stone stone) {
	for stone.y > s.top() {
		s.layers = append([][7]bool{{}}, s.layers...)
		s.height++
	}

	s.layers[s.top()-stone.y][stone.x] = true
}

func (s *stack) addRock(rock rock) {
	for _, stone := range rock {
		s.addStone(stone)
	}
}

func (s *stack) stoneCollides(stone stone) bool {
	if stone.y < 0 {
		return true
	}

	if stone.x < 0 || stone.x >= 7 {
		return true
	}

	if stone.y > s.top() {
		return false
	}

	return s.layers[s.top()-stone.y][stone.x]
}

func (s *stack) rockCollides(rock rock) bool {
	for _, stone := range rock {
		if s.stoneCollides(stone) {
			return true
		}
	}

	return false
}

func (s *stack) attemptMove(rock rock, dx, dy int) (rock, bool) {
	movedRock := rock.move(dx, dy)
	if s.rockCollides(movedRock) {
		return rock, false
	}

	return movedRock, true
}

func (s *stack) relevantSubStack() *stack {
	for y := 0; y < s.top(); y++ {
		isLine := true
		for x := 0; x < 7; x++ {
			if !s.layers[y][x] {
				isLine = false
				break
			}
		}

		if isLine {
			subStack := stack{
				layers: s.layers[0:y],
				height: s.height,
			}
			return &subStack
		}
	}

	return s
}

func (s *stack) string() string {
	var sb strings.Builder
	for _, line := range s.layers {
		sb.WriteRune('|')
		for _, c := range line {
			if c {
				sb.WriteRune('#')
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('|')

		sb.WriteRune('\n')
	}
	sb.WriteString("+-------+")

	return sb.String()
}

func main() {
	inputData, err := os.ReadFile("day17/input.txt")
	if err != nil {
		panic(err)
	}

	jets := parseJets(string(inputData))

	type memoizedKey struct {
		layers    [][7]bool
		rockIndex int
		jetIndex  int
	}
	type memoizedValue struct {
		numRocksFallen int
		height         int
	}

	var memoizedKeys []memoizedKey
	var memoizedValues []memoizedValue

	stack := &stack{}
	curRockIndex := 0
	curJetIndex := 0
	numRocksFallen := 0
	for numRocksFallen < 1000000000000 {
		curRock := rocks[curRockIndex]
		curRock = curRock.move(2, stack.top()+4)
		curRockIndex = (curRockIndex + 1) % len(rocks)

		for true {
			// Blow the rock by the jet.
			curJet := jets[curJetIndex]
			curJetIndex = (curJetIndex + 1) % len(jets)

			curRock, _ = stack.attemptMove(curRock, curJet.dx(), 0)

			// Then try to move the rock downwards.
			var ok bool
			curRock, ok = stack.attemptMove(curRock, 0, -1)
			if !ok {
				stack.addRock(curRock)
				break
			}
		}

		numRocksFallen++

		// Check to see if we've been in this state before.
		stack = stack.relevantSubStack()

		newKey := memoizedKey{
			layers:    stack.layers,
			rockIndex: curRockIndex,
			jetIndex:  curJetIndex,
		}
		for i, memoizedKey := range memoizedKeys {
			if reflect.DeepEqual(memoizedKey, newKey) {
				memoizedValue := memoizedValues[i]
				fmt.Printf("%v\n", memoizedValue)

				rocksDiff := numRocksFallen - memoizedValue.numRocksFallen
				heightDiff := stack.height - memoizedValue.height

				repeatsToCome := (1000000000000 - numRocksFallen) / rocksDiff

				numRocksFallen += rocksDiff * repeatsToCome
				stack.height += heightDiff * repeatsToCome
			}
		}

		memoizedKeys = append(memoizedKeys, newKey)
		memoizedValues = append(memoizedValues, memoizedValue{
			numRocksFallen: numRocksFallen,
			height:         stack.height,
		})

		//if numRocksFallen%1000 == 0 {
		//	oldLength := len(stack.layers)
		//	stack = stack.relevantSubStack()
		//	println(numRocksFallen, oldLength, len(stack.layers))
		//}
	}

	println(stack.height)
}
