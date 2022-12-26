package main

import (
	"os"
	"strconv"
	"strings"
)

func abs(value int) int {
	if value < 0 {
		return -value
	}

	return value
}

type linkedNumber struct {
	left   *linkedNumber
	right  *linkedNumber
	number int
}

func main() {
	inputData, err := os.ReadFile("day20/input.txt")
	if err != nil {
		panic(err)
	}

	var numbers []int
	for _, line := range strings.Split(string(inputData), "\n") {
		lineNum, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}

		numbers = append(numbers, lineNum)
	}

	var linkedNumbers []*linkedNumber
	var zeroLinkedNumber *linkedNumber
	for _, number := range numbers {
		newLinkedNumber := &linkedNumber{number: number * 811589153}
		linkedNumbers = append(linkedNumbers, newLinkedNumber)

		if number == 0 {
			zeroLinkedNumber = newLinkedNumber
		}
	}

	for i, cur := range linkedNumbers {
		if i == 0 {
			cur.left = linkedNumbers[len(linkedNumbers)-1]
		} else {
			cur.left = linkedNumbers[i-1]
		}

		if i == len(linkedNumbers)-1 {
			cur.right = linkedNumbers[0]
		} else {
			cur.right = linkedNumbers[i+1]
		}
	}

	for i := 0; i < 10; i++ {
		println(i)
		for _, linkedNumber := range linkedNumbers {
			linkedNumber.left.right = linkedNumber.right
			linkedNumber.right.left = linkedNumber.left

			newLeft := linkedNumber.left
			newRight := linkedNumber.right
			for i := 0; i < (abs(linkedNumber.number) % (len(linkedNumbers) - 1)); i++ {
				if linkedNumber.number < 0 {
					newLeft = newLeft.left
					newRight = newRight.left
				} else {
					newLeft = newLeft.right
					newRight = newRight.right
				}
			}

			newLeft.right = linkedNumber
			linkedNumber.left = newLeft

			newRight.left = linkedNumber
			linkedNumber.right = newRight
		}
	}

	sum := 0
	cur := zeroLinkedNumber
	for i := 0; i < 3; i++ {
		for j := 0; j < 1000; j++ {
			cur = cur.right
		}

		sum += cur.number
	}

	println(sum)
}
