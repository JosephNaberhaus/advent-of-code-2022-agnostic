package main

import (
	"os"
	"strconv"
	"strings"
)

type tile int8

const (
	empty tile = iota
	wall
	air
)

type direction int8

const (
	up    direction = 3
	down            = 1
	left            = 2
	right           = 0
)

type state struct {
	board         [][]tile
	width, height int
	xPos, yPos    int
	direction     direction
}

func (s state) curTile() tile {
	if s.xPos < 0 || s.xPos >= s.width {
		return air
	}

	if s.yPos < 0 || s.yPos >= s.height {
		return air
	}

	return s.board[s.yPos][s.xPos]
}

func (s state) left() state {
	switch s.direction {
	case left:
		s.direction = down
		return s
	case right:
		s.direction = up
		return s
	case up:
		s.direction = left
		return s
	case down:
		s.direction = right
		return s
	}

	panic("oh no")
}

func (s state) right() state {
	newState := s.left()
	newState = newState.left()
	newState = newState.left()
	return newState
}

func (s state) moveOne() (state, bool) {
	switch s.direction {
	case left:
		return s.moveX(-1)
	case right:
		return s.moveX(1)
	case up:
		return s.moveY(-1)
	case down:
		return s.moveY(1)
	}

	panic("oh no")
}

func (s state) moveX(dx int) (state, bool) {
	newState := s
	newState.xPos += dx

	if newState.curTile() == air {
		switch {
		case newState.xPos == 49 && newState.yPos >= 0 && newState.yPos <= 49:
			newState.xPos = 0
			newState.yPos = 100 + (49 - newState.yPos)
			newState.direction = right
		case newState.xPos == 150 && newState.yPos >= 0 && newState.yPos <= 49:
			newState.xPos = 99
			newState.yPos = 100 + (49 - newState.yPos)
			newState.direction = left
		case newState.xPos == 49 && newState.yPos >= 50 && newState.yPos <= 99:
			newState.xPos = newState.yPos - 50
			newState.yPos = 100
			newState.direction = down
		case newState.xPos == 100 && newState.yPos >= 50 && newState.yPos <= 99:
			newState.xPos = 50 + newState.yPos
			newState.yPos = 49
			newState.direction = up
		case newState.xPos == -1 && newState.yPos >= 100 && newState.yPos <= 149:
			newState.xPos = 50
			newState.yPos = 49 - (newState.yPos - 100)
			newState.direction = right
		case newState.xPos == 100 && newState.yPos >= 100 && newState.yPos <= 149:
			newState.xPos = 149
			newState.yPos = 49 - (newState.yPos - 100)
			newState.direction = left
		case newState.xPos == -1 && newState.yPos >= 150 && newState.yPos <= 199:
			newState.xPos = 50 + (newState.yPos - 150)
			newState.yPos = 0
			newState.direction = down
		case newState.xPos == 50 && newState.yPos >= 150 && newState.yPos <= 199:
			newState.xPos = 50 + (newState.yPos - 150)
			newState.yPos = 149
			newState.direction = up
		default:
			panic("oh no")
		}

		if newState.curTile() == wall {
			return s, false
		}

		return newState, true
	} else if newState.curTile() == wall {
		return s, false
	}

	return newState, true
}

func (s state) moveY(dy int) (state, bool) {
	newState := s
	newState.yPos += dy

	if newState.curTile() == air {
		switch {
		case newState.xPos >= 0 && newState.xPos <= 49 && newState.yPos == 99:
			newState.yPos = 50 + newState.xPos
			newState.xPos = 50
			newState.direction = right
		case newState.xPos >= 0 && newState.xPos <= 49 && newState.yPos == 200:
			newState.yPos = 0
			newState.xPos = 100 + newState.xPos
			newState.direction = down
		case newState.xPos >= 50 && newState.xPos <= 99 && newState.yPos == -1:
			newState.yPos = 150 + (newState.xPos - 50)
			newState.xPos = 0
			newState.direction = right
		case newState.xPos >= 50 && newState.xPos <= 99 && newState.yPos == 150:
			newState.yPos = 150 + (newState.xPos - 50)
			newState.xPos = 49
			newState.direction = left
		case newState.xPos >= 100 && newState.xPos <= 149 && newState.yPos == -1:
			newState.yPos = 199
			newState.xPos = newState.xPos - 100
			newState.direction = up
		case newState.xPos >= 100 && newState.xPos <= 149 && newState.yPos == 50:
			newState.yPos = 50 + (newState.xPos - 100)
			newState.xPos = 99
			newState.direction = left
		default:
			panic("oh no")
		}

		if newState.curTile() == wall {
			return s, false
		}

		return newState, true
	} else if newState.curTile() == wall {
		return s, false
	}

	return newState, true
}

func (s state) string() string {
	output := ""
	for y, line := range s.board {
		for x, tile := range line {
			if s.xPos == x && s.yPos == y {
				switch s.direction {
				case up:
					output += "^"
				case down:
					output += "v"
				case left:
					output += "<"
				case right:
					output += ">"
				}

				continue
			}

			switch tile {
			case air:
				output += " "
			case wall:
				output += "#"
			case empty:
				output += "."
			}
		}

		output += "\n"
	}

	return output
}

func max(first, second int) int {
	if first > second {
		return first
	}

	return second
}

func main() {
	inputData, err := os.ReadFile("day22/input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(inputData), "\n")

	maxWidth := 0
	instructionsIndex := -1
	var board [][]tile
	for i, line := range lines {
		var curLine []tile

		if line == "" {
			instructionsIndex = i + 1
		}

		for _, c := range line {
			switch c {
			case '.':
				curLine = append(curLine, empty)
			case '#':
				curLine = append(curLine, wall)
			default:
				curLine = append(curLine, air)
			}
		}

		maxWidth = max(maxWidth, len(curLine))
		board = append(board, curLine)
	}

	for i, line := range board {
		for j := 0; j < maxWidth-len(line); j++ {
			board[i] = append(board[i], air)
		}
	}

	startX := -1
	for x := 0; x < maxWidth; x++ {
		if board[0][x] != air {
			startX = x
			break
		}
	}

	state := state{
		board:     board,
		width:     maxWidth,
		height:    len(board),
		xPos:      startX,
		yPos:      0,
		direction: right,
	}
	curNumber := ""
	for i, c := range lines[instructionsIndex] + "!" {
		println(i, curNumber)
		if c == 'R' || c == 'L' || c == '!' {
			if curNumber != "" {
				number, err := strconv.Atoi(curNumber)
				curNumber = ""
				if err != nil {
					panic(err)
				}

				for i := 0; i < number; i++ {
					var ok bool
					state, ok = state.moveOne()
					if !ok {
						break
					}
				}
			}

			if c == 'R' {
				state = state.right()
			} else if c == 'L' {
				state = state.left()
			} else if c == '!' {
				break
			} else {
				panic("oh no")
			}
		} else {
			curNumber += string(c)
		}
	}

	println(4*(state.xPos+1) + 1000*(state.yPos+1) + int(state.direction))
}
