package main

import (
	"fmt"
	"os"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	var L, C, N int
	fmt.Scan(&L, &C, &N)
	fmt.Fprintf(os.Stderr, "\n limit %v", L)
	groups := []int{}
	for i := 0; i < N; i++ {
		var pi int
		fmt.Scan(&pi)
		groups = append(groups, pi)
	}
	fmt.Fprintf(os.Stderr, "\n groups %v", groups)
	totalRevenue := 0
	totalPassengers := sum(groups...)

	if totalPassengers < L {
		fmt.Fprintf(os.Stderr, "\n sum of groups(%d) greater than Limit %d", totalPassengers, L)
		totalRevenue = totalPassengers * C
	} else {
		idx := 0
		preComputed := map[int][2]int{}
		for j := 0; j < C; j++ {
			fmt.Fprintf(os.Stderr, "\n ===> Cycle: %d, start idx %d", j, idx)

			numPassengers := 0
			// [nextIdx, numPassengers]
			val, exists := preComputed[idx]
			fmt.Fprintf(os.Stderr, "\n  preComputed: %v", preComputed)
			if exists {
				fmt.Fprintf(os.Stderr, "\n exist in precomputed: %v", val)
				idx = val[0]
				numPassengers = val[1]
			} else {
				beginIdx := idx
				for groups[idx]+numPassengers <= L {
					fmt.Fprintf(os.Stderr, "\n status passengers %d", groups[idx]+numPassengers)

					fmt.Fprintf(os.Stderr, "\n idx: %d", idx)
					fmt.Fprintf(os.Stderr, "\n groups[idx]: %d", groups[idx])

					numPassengers += groups[idx]
					fmt.Fprintf(os.Stderr, "\n after passengers %d  ", numPassengers)
					idx = (idx + 1) % N
				}
				preComputed[beginIdx] = [2]int{idx, numPassengers}
			}
			totalRevenue += numPassengers
			fmt.Fprintf(os.Stderr, "\n end cycle %d <====== \n ", j)
		}
	}

	fmt.Fprintf(os.Stderr, "totalRevenue: %d \n", totalRevenue)
	fmt.Println(totalRevenue) // Write answer to stdout
}

func sum(numbers ...int) (total int) {
	for i := range numbers {
		total += numbers[i]
	}
	return total
}
