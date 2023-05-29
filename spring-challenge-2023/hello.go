package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func normalizedAnts(amountAnt int) int {
	return amountAnt / (2 + 1 + 2)
}
func wait() string {
	return fmt.Sprintf("WAIT")
}
func showMsg(text string) string {
	return fmt.Sprintf("MESSAGE %s", text)
}
func placeBeacon(cellIdx, strength int) string {
	// Places a beacon that lasts one turn.
	return fmt.Sprintf("BEACON %d %d ", cellIdx, strength)
}
func placeLine(srcIdx, destIdx, strength int) string {
	// Places beacons along a path between the two provided cells.
	return fmt.Sprintf("LINE %d %d %d", srcIdx, destIdx, strength)
}

var biasEgg = 50.0

const biasCrystal = 100.0

type State struct {
	Cells          []Cell
	MyBaseIndices  []int
	OppBaseIndices []int
}

func (s *State) decideAct() {
	crystalCells := []Cell{}
	eggCells := []Cell{}
	myBaseCells := []Cell{}
	commands := make([]string, 0)
	commands = append(commands, placeBeacon(s.MyBaseIndices[0], 1))
	amountCrystal := 0.0
	amountEgg := 0.0
	initialAmountCrystal := 0.0
	initialAmountEgg := 0.0
	for i := range s.Cells {
		cell := s.Cells[i]
		if Contains(s.MyBaseIndices, cell.index) {
			myBaseCells = append(myBaseCells, cell)
		} else if cell.isCrystal() && cell.resources > 0 {
			crystalCells = append(crystalCells, cell)
			amountCrystal += cell.resources
			initialAmountCrystal += cell.initialResources
		} else if cell.isEgg() && cell.resources > 0 {
			eggCells = append(eggCells, cell)
			amountEgg += cell.resources
			initialAmountEgg += cell.initialResources
		}
	}
	if initialAmountCrystal > 3000 {
		biasEgg = 300
	} else if initialAmountCrystal > 1000 {
		biasEgg = 150
	}
	for j := range crystalCells {
		c := crystalCells[j]
		weight := roundInt((c.resources / amountCrystal) * biasCrystal)
		if c.resources < 1 {
			weight = 0
		} else if c.resources < 10 {
			weight = 1
		}
		for idx := range myBaseCells {
			mbc := myBaseCells[idx]
			dist := calcDist(s.Cells, mbc, c, 0, []int{mbc.index})
			fmt.Fprintf(os.Stderr, "\n from : %d to: %d distance: %d", mbc.index, c.index, dist)
			if dist < 2 {
				weight *= 4
			} else if dist < 4 {
				weight *= 2
			} else {
				weight /= 2
			}
		}
		// fmt.Fprintf(os.Stderr, "resource: %f amountCrystal: %f placeBeacon: to index: %d with weight: %d  \n", c.resources, amountCrystal, c.index, weight)
		// commands = append(commands, showMsg())
		commands = append(commands, placeLine(s.MyBaseIndices[0], c.index, weight))
	}
	// crystalRatio := (amountCrystal / initialAmountCrystal) * 100.0
	eggRatio := (amountEgg / initialAmountEgg) * 100.0
	// fmt.Fprintf(os.Stderr, "\n crystals ratio:  %f, amount: %f initialAmount: %f", crystalRatio, amountCrystal, initialAmountCrystal)
	// fmt.Fprintf(os.Stderr, "\n egg ratio:  %f, amount: %f initialAmount: %f", eggRatio, amountEgg, initialAmountEgg)
	if eggRatio > 5.0 {
		for z := range eggCells {
			weight := roundInt((eggCells[z].resources / amountEgg) * biasEgg)
			fmt.Fprintf(os.Stderr, "\n placeBeacon: to index: %d with weight: %d", eggCells[z].index, weight)
			commands = append(commands, placeLine(s.MyBaseIndices[0], eggCells[z].index, weight))
		}
	}
	fmt.Println(strings.Join(commands, ";"))
}

type Cell struct {
	// initialResources: the initial amount of eggs/crystals on this cell
	// neigh0: the index of the neighbouring cell for each direction
	index, cellType                              int
	initialResources, resources, myAnts, oppAnts float64

	neighbors []int
}

func (c *Cell) isEmpty() bool {
	return c.cellType == 0
}
func (c *Cell) isEgg() bool {
	return c.cellType == 1
}
func (c *Cell) isCrystal() bool {
	return c.cellType == 2
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	var inputs []string
	// numberOfCells: amount of hexagonal cells in this map
	var numberOfCells int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &numberOfCells)

	cells := make([]Cell, 0, numberOfCells)
	state := State{cells, make([]int, 0), make([]int, 0)}

	for i := 0; i < numberOfCells; i++ {
		var _type, neigh0, neigh1, neigh2, neigh3, neigh4, neigh5 int
		var initialResources float64
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &_type, &initialResources, &neigh0, &neigh1, &neigh2, &neigh3, &neigh4, &neigh5)
		neighbors := []int{neigh0, neigh1, neigh2, neigh3, neigh4, neigh5}
		state.Cells = append(state.Cells, Cell{i, _type, initialResources, 0, 0, 0, filter(neighbors, func(v int) bool {
			return v > -1
		})})
	}
	//  The players’ ants will start the game on these bases.
	var numberOfBases int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &numberOfBases)

	scanner.Scan()
	inputs = strings.Split(scanner.Text(), " ")
	for i := 0; i < numberOfBases; i++ {
		myBaseIndex, _ := strconv.ParseInt(inputs[i], 10, 32)
		state.MyBaseIndices = append(state.MyBaseIndices, int(myBaseIndex))
	}
	scanner.Scan()
	inputs = strings.Split(scanner.Text(), " ")
	for i := 0; i < numberOfBases; i++ {
		oppBaseIndex, _ := strconv.ParseInt(inputs[i], 10, 32)
		state.OppBaseIndices = append(state.OppBaseIndices, int(oppBaseIndex))
	}
	// fmt.Fprintf(os.Stderr, "Game State %+v\n", state)
	for {
		for i := 0; i < numberOfCells; i++ {
			// resources: the current amount of eggs/crystals on this cell
			// myAnts: the amount of your ants on this cell
			// oppAnts: the amount of opponent ants on this cell
			var resources, myAnts, oppAnts float64
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &resources, &myAnts, &oppAnts)
			state.Cells[i].resources = resources
			state.Cells[i].myAnts = myAnts
			state.Cells[i].oppAnts = oppAnts
		}
		// fmt.Fprintf(os.Stderr, "Cells:  %+v\n", state.Cells)
		state.decideAct()
	}
}

func calcDist(cells []Cell, src Cell, dest Cell, initialDist int, searchedIndices []int) int {
	// FIXME
	// 자기 자신이거나 이웃에 있을경우 거리를 반환
	// 이미 탐색된 이웃을 제외하고 해당 함수를 적용
	// 무엇도 안될경우 -1 반환
	if src.index == dest.index {
		return initialDist
	} else if Contains(src.neighbors, dest.index) {
		return initialDist + 1
	}
	for _, idx := range src.neighbors {
		neighbors := filter(cells, func(v Cell) bool {
			return v.index == idx && len(filter(searchedIndices, func(v int) bool { return v == idx })) == 0
		})
		if len(neighbors) > 0 {
			searchedIndices = append(searchedIndices, idx)
			return calcDist(cells, src, neighbors[0], initialDist+1, searchedIndices)
		}
	}
	return -1
}

func filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func roundInt(f float64) int {
	return int(f + 0.5)
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
