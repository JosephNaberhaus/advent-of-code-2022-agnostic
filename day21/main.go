package main

import (
	"os"
	"strconv"
	"strings"
)

type operation int

const (
	operationAdd operation = iota
	operationSubtract
	operationMultiply
	operationDivide
	operationShout
)

type listener struct {
	monkeyName string

	monkey *monkey
	isLeft bool
}

type monkey struct {
	name string

	op     operation
	number int

	leftMonkeyName  string
	leftMemory      *int
	leftMonkey      *monkey
	rightMonkeyName string
	rightMemory     *int
	rightMonkey     *monkey

	listeners []listener
}

func (m *monkey) copy() *monkey {
	newMonkey := *m
	return &newMonkey
}

func parseMonkey(line string) *monkey {
	lineSplit := strings.Split(line, ": ")

	monkey := new(monkey)
	monkey.name = lineSplit[0]

	equationSplit := strings.Split(lineSplit[1], " ")
	if len(equationSplit) == 1 {
		monkey.op = operationShout
		monkey.number, _ = strconv.Atoi(equationSplit[0])
	} else {
		monkey.leftMonkeyName = equationSplit[0]

		switch equationSplit[1] {
		case "+":
			monkey.op = operationAdd
		case "-":
			monkey.op = operationSubtract
		case "/":
			monkey.op = operationDivide
		case "*":
			monkey.op = operationMultiply
		}

		monkey.rightMonkeyName = equationSplit[2]
	}

	return monkey
}

func (m *monkey) equation() string {
	if m.name == "root" {
		return m.leftMonkey.equation() + " = " + m.rightMonkey.equation()
	} else if m.name == "humn" {
		return "X"
	}

	if m.op == operationShout {
		return strconv.Itoa(m.number)
	}

	switch m.op {
	case operationAdd:
		return "(" + m.leftMonkey.equation() + " + " + m.rightMonkey.equation() + ")"
	case operationSubtract:
		return "(" + m.leftMonkey.equation() + " - " + m.rightMonkey.equation() + ")"
	case operationMultiply:
		return "(" + m.leftMonkey.equation() + " * " + m.rightMonkey.equation() + ")"
	case operationDivide:
		return "(" + m.leftMonkey.equation() + " / " + m.rightMonkey.equation() + ")"
	}

	panic("unkown")
}

func main() {
	inputData, err := os.ReadFile("day21/input.txt")
	if err != nil {
		panic(err)
	}

	var monkeys []*monkey
	for _, line := range strings.Split(string(inputData), "\n") {
		monkeys = append(monkeys, parseMonkey(line))
	}

	monkeysByName := map[string]*monkey{}
	for _, monkey := range monkeys {
		monkeysByName[monkey.name] = monkey
	}

	for _, monkey := range monkeys {
		if monkey.op != operationShout {
			leftSourceMonkey := monkeysByName[monkey.leftMonkeyName]
			leftSourceMonkey.listeners = append(leftSourceMonkey.listeners, listener{
				monkeyName: monkey.name,
				monkey:     monkey,
				isLeft:     true,
			})
			monkey.leftMonkey = leftSourceMonkey

			rightSourceMonkey := monkeysByName[monkey.rightMonkeyName]
			rightSourceMonkey.listeners = append(rightSourceMonkey.listeners, listener{
				monkeyName: monkey.name,
				monkey:     monkey,
				isLeft:     false,
			})
			monkey.rightMonkey = rightSourceMonkey
		}
	}

	var shoutingMonkeys []*monkey
	for _, monkey := range monkeys {
		if monkey.op == operationShout {
			shoutingMonkeys = append(shoutingMonkeys, monkey)
		}
	}

	println(monkeysByName["root"].equation())
}
