package main

import (
	"os"
	"strconv"
	"strings"
)

func max(first, second int) int {
	if first > second {
		return first
	}

	return second
}

type robotBlueprint struct {
	mines        string
	oreCost      int
	clayCost     int
	obsidianCost int
}

func parseRobotBlueprint(robotStr string) robotBlueprint {
	robotStrSplit := strings.Split(robotStr, " ")

	cost := map[string]int{}
	for i := 4; i < len(robotStrSplit); i += 3 {
		price, _ := strconv.ParseInt(robotStrSplit[i], 0, 64)
		cost[robotStrSplit[i+1]] = int(price)
	}

	return robotBlueprint{
		mines:        robotStrSplit[1],
		obsidianCost: cost["obsidian"],
		clayCost:     cost["clay"],
		oreCost:      cost["ore"],
	}
}

type blueprint struct {
	robots []robotBlueprint
}

func parseBlueprint(line string) blueprint {
	lineSplit := strings.Split(line, ":")

	var robots []robotBlueprint
	for _, robotStr := range strings.Split(lineSplit[1], ".") {
		if strings.TrimSpace(robotStr) == "" {
			continue
		}

		robots = append(robots, parseRobotBlueprint(strings.TrimSpace(robotStr)))
	}

	return blueprint{robots: robots}
}

type state struct {
	geodes         int
	geodeRobots    int
	clay           int
	clayRobots     int
	ore            int
	oreRobots      int
	obsidian       int
	obsidianRobots int
	minute         int
}

func (s state) nextStates(blueprint blueprint) []state {
	if s.minute == 32 {
		return nil
	}

	var nextStates []state

	s.minute++
	s.ore += s.oreRobots
	s.clay += s.clayRobots
	s.obsidian += s.obsidianRobots
	s.geodes += s.geodeRobots

	nextStates = append(nextStates, s)

	for _, robotBlueprint := range blueprint.robots {
		newState := s
		newState.clay -= robotBlueprint.clayCost
		newState.ore -= robotBlueprint.oreCost
		newState.obsidian -= robotBlueprint.obsidianCost

		if newState.clay-s.clayRobots < 0 || newState.ore-s.oreRobots < 0 || newState.obsidian-s.obsidianRobots < 0 {
			continue
		}

		switch robotBlueprint.mines {
		case "clay":
			newState.clayRobots++
		case "ore":
			newState.oreRobots++
		case "geode":
			newState.geodeRobots++
		case "obsidian":
			newState.obsidianRobots++
		}

		nextStates = append(nextStates, newState)
	}

	return nextStates
}

func main() {
	inputData, err := os.ReadFile("day19/input.txt")
	if err != nil {
		panic(err)
	}

	var robotBlueprints []blueprint
	for _, line := range strings.Split(string(inputData), "\n") {
		robotBlueprints = append(robotBlueprints, parseBlueprint(line))
	}

	result := 1
	for i, blueprint := range robotBlueprints[:3] {
		println("starting", i)
		maxGeodes := 0
		maxGeodeRobots := 0

		curToVisit := map[state]struct{}{
			state{
				oreRobots: 1,
				minute:    0,
			}: {},
		}
		nextToVisit := map[state]struct{}{}

		minute := 0
		for len(curToVisit) > 0 {
			println(minute, len(curToVisit))
			minute++

			for state := range curToVisit {
				for _, nextState := range state.nextStates(blueprint) {
					nextToVisit[nextState] = struct{}{}
				}
			}

			for state := range nextToVisit {
				maxGeodes = max(maxGeodes, state.geodes)
				maxGeodeRobots = max(maxGeodeRobots, state.geodeRobots)
			}

			curToVisit = map[state]struct{}{}
			for state := range nextToVisit {
				if state.geodeRobots >= maxGeodeRobots-2 {
					curToVisit[state] = struct{}{}
				}
			}

			nextToVisit = map[state]struct{}{}
		}

		println("done", i+1, maxGeodes)
		result *= maxGeodes
	}

	println(result)
}
